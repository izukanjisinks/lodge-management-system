package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type CorporateBookingRequestService struct {
	requestRepo *repository.CorporateBookingRequestRepository
	guestRepo   *repository.CorporateGuestRepository
	corProfile  *CorProfileService
	workflow    *WorkflowService
	venueRepo   *repository.VenueRepository
	menuRepo    *repository.MenuRepository
	booking     *BookingService
}

func NewCorporateBookingRequestService(
	requestRepo *repository.CorporateBookingRequestRepository,
	guestRepo *repository.CorporateGuestRepository,
	corProfile *CorProfileService,
) *CorporateBookingRequestService {
	return &CorporateBookingRequestService{
		requestRepo: requestRepo,
		guestRepo:   guestRepo,
		corProfile:  corProfile,
	}
}

func (s *CorporateBookingRequestService) SetWorkflowService(svc *WorkflowService) {
	s.workflow = svc
}

func (s *CorporateBookingRequestService) SetVenueRepository(repo *repository.VenueRepository) {
	s.venueRepo = repo
}

func (s *CorporateBookingRequestService) SetMenuRepository(repo *repository.MenuRepository) {
	s.menuRepo = repo
}

// SetBookingService wires the booking service so approving an event/conference
// request auto-creates the booking from the guest's chosen venue.
func (s *CorporateBookingRequestService) SetBookingService(svc *BookingService) {
	s.booking = svc
}

// validateVenue ensures a guest-chosen venue exists and is in service.
func (s *CorporateBookingRequestService) validateVenue(orgID, venueID uuid.UUID) error {
	if venueID == uuid.Nil {
		return errors.New("venue_id is required")
	}
	if s.venueRepo == nil {
		return nil
	}
	venue, err := s.venueRepo.GetByID(venueID, orgID)
	if err != nil {
		return errors.New("selected venue not found")
	}
	if !venue.IsAvailable {
		return errors.New("selected venue is not available")
	}
	return nil
}

// venueBranchID resolves the lodge branch a venue physically belongs to. This is
// the branch an event booking inherits — distinct from the request's branch_id,
// which references the corporate client's branch (cor_branch_details).
func (s *CorporateBookingRequestService) venueBranchID(orgID, venueID uuid.UUID) *uuid.UUID {
	if s.venueRepo == nil || venueID == uuid.Nil {
		return nil
	}
	venue, err := s.venueRepo.GetByID(venueID, orgID)
	if err != nil {
		return nil
	}
	return venue.BranchID
}

// ─── Submission ───────────────────────────────────────────────────────────────

func (s *CorporateBookingRequestService) SubmitAccommodation(orgID uuid.UUID, webUserID uuid.UUID, req *models.SubmitAccommodationRequest) (*models.CorporateBookingRequest, error) {
	// Validate required fields
	if req.BookedBy.Email == "" || req.BookedBy.Name == "" {
		return nil, errors.New("booked_by name and email are required")
	}
	if req.Company == nil || req.Company.Name == "" {
		return nil, errors.New("company name is required")
	}
	if req.Accommodation == nil {
		return nil, errors.New("accommodation block is required")
	}
	if req.Accommodation.RoomCount < 1 {
		return nil, errors.New("room_count must be at least 1")
	}
	if req.Accommodation.CheckIn == "" || req.Accommodation.CheckOut == "" {
		return nil, errors.New("check_in and check_out are required")
	}

	// Validate attendants if detailed mode
	if req.ParticipantMode == "detailed" && len(req.Attendants) == 0 {
		return nil, errors.New("attendants are required in detailed mode")
	}

	// Map the nested company/branch/profile into the inputs ResolveChain expects,
	// then get-or-create the company → branch → profile chain.
	company := models.CorBookingCompanyInput{
		CompanyName: req.Company.Name,
		TPIN:        req.Company.TPIN,
		Industry:    req.Company.Industry,
	}

	var branch *models.CorBookingBranchInput
	if req.Company.BranchName != "" {
		branch = &models.CorBookingBranchInput{
			Name: req.Company.BranchName,
		}
	}

	// Split full name into first/last for profile lookup.
	firstName, lastName := splitName(req.BookedBy.Name)
	profile := models.CorBookingProfileInput{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      req.BookedBy.Email,
		Phone:      req.BookedBy.Phone,
		JobTitle:   req.BookedBy.JobTitle,
		ManNumber:  req.BookedBy.ManNumber,
		Department: req.Company.DepartmentName,
	}

	companyID, branchIDPtr, corProfileID, err := s.corProfile.ResolveChain(orgID, company, branch, profile)
	if err != nil {
		return nil, err
	}

	// Store entire payload as-is
	payloadBytes, _ := json.Marshal(req)
	payload := json.RawMessage(payloadBytes)

	var webUserIDPtr *uuid.UUID
	if webUserID != uuid.Nil {
		webUserIDPtr = &webUserID
	}
	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchIDPtr,
		CorProfileID:     &corProfileID,
		CompanyID:        &companyID,
		WebUserID:        webUserIDPtr,
		BookingType:      models.CorporateBookingTypeAccommodation,
		Status:           models.CorporateBookingStatusPending,
		Notes:            req.Accommodation.Notes,
		Documents:        req.Documents,
		Payload:          payload,
		ProfileName:      req.BookedBy.Name,
		CompanyName:      req.Company.Name,
	}
	// Map approver fields from the nested approver object
	if req.Approver != nil {
		r.AuthoriserName = req.Approver.Name
		r.AuthoriserEmail = req.Approver.Email
		r.AuthoriserPhone = req.Approver.Phone
		r.AuthoriserTitle = req.Approver.Title
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

// SubmitEventBooking handles the standalone event envelope (Flow B) from
// eventBooking.js for corporate bookers. It maps the nested company/approver/
// booked_by into the ResolveChain inputs, stores the whole envelope as JSONB, and
// starts the approval workflow.
func (s *CorporateBookingRequestService) SubmitEventBooking(orgID uuid.UUID, webUserID uuid.UUID, req *models.SubmitEventBookingRequest) (*models.CorporateBookingRequest, error) {
	if req.BookedBy.Email == "" || req.BookedBy.Name == "" {
		return nil, errors.New("booked_by name and email are required")
	}
	if req.Company == nil || req.Company.Name == "" {
		return nil, errors.New("company name is required")
	}
	if req.Event == nil || len(req.Event.Sessions) == 0 {
		return nil, errors.New("at least one event session is required")
	}
	for i, sess := range req.Event.Sessions {
		if sess.VenueID == "" {
			return nil, fmt.Errorf("session %d is missing a venue", i+1)
		}
		venueID, err := uuid.Parse(sess.VenueID)
		if err != nil {
			return nil, fmt.Errorf("session %d has an invalid venue", i+1)
		}
		if err := s.validateVenue(orgID, venueID); err != nil {
			return nil, err
		}
	}

	company := models.CorBookingCompanyInput{
		CompanyName: req.Company.Name,
		TPIN:        req.Company.TPIN,
		Industry:    req.Company.Industry,
	}
	var branch *models.CorBookingBranchInput
	if req.Company.BranchName != "" {
		branch = &models.CorBookingBranchInput{Name: req.Company.BranchName}
	}
	firstName, lastName := splitName(req.BookedBy.Name)
	profile := models.CorBookingProfileInput{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      req.BookedBy.Email,
		Phone:      req.BookedBy.Phone,
		JobTitle:   req.BookedBy.JobTitle,
		ManNumber:  req.BookedBy.ManNumber,
		Department: req.Company.DepartmentName,
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, company, branch, profile)
	if err != nil {
		return nil, err
	}

	payloadBytes, _ := json.Marshal(req)

	var webUserIDPtr *uuid.UUID
	if webUserID != uuid.Nil {
		webUserIDPtr = &webUserID
	}
	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		WebUserID:        webUserIDPtr,
		BookingType:      models.CorporateBookingTypeEvent,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.Event.ReasonForBooking,
		Notes:            req.Event.Notes,
		Documents:        req.Documents,
		Payload:          json.RawMessage(payloadBytes),
		ProfileName:      req.BookedBy.Name,
		CompanyName:      req.Company.Name,
	}
	if req.Approver != nil {
		r.AuthoriserName = req.Approver.Name
		r.AuthoriserEmail = req.Approver.Email
		r.AuthoriserPhone = req.Approver.Phone
		r.AuthoriserTitle = req.Approver.Title
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

// SubmitMealBooking handles the standalone meal envelope (Flow B) from
// mealBooking.js for corporate bookers. Stores the whole envelope as JSONB and
// starts the approval workflow.
func (s *CorporateBookingRequestService) SubmitMealBooking(orgID uuid.UUID, webUserID uuid.UUID, req *models.SubmitMealBookingRequest) (*models.CorporateBookingRequest, error) {
	if req.BookedBy.Email == "" || req.BookedBy.Name == "" {
		return nil, errors.New("booked_by name and email are required")
	}
	if req.Company == nil || req.Company.Name == "" {
		return nil, errors.New("company name is required")
	}
	if req.Meal == nil || len(req.Meal.Sessions) == 0 {
		return nil, errors.New("at least one meal session is required")
	}

	company := models.CorBookingCompanyInput{
		CompanyName: req.Company.Name,
		TPIN:        req.Company.TPIN,
		Industry:    req.Company.Industry,
	}
	var branch *models.CorBookingBranchInput
	if req.Company.BranchName != "" {
		branch = &models.CorBookingBranchInput{Name: req.Company.BranchName}
	}
	firstName, lastName := splitName(req.BookedBy.Name)
	profile := models.CorBookingProfileInput{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      req.BookedBy.Email,
		Phone:      req.BookedBy.Phone,
		JobTitle:   req.BookedBy.JobTitle,
		ManNumber:  req.BookedBy.ManNumber,
		Department: req.Company.DepartmentName,
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, company, branch, profile)
	if err != nil {
		return nil, err
	}

	payloadBytes, _ := json.Marshal(req)

	var webUserIDPtr *uuid.UUID
	if webUserID != uuid.Nil {
		webUserIDPtr = &webUserID
	}
	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		WebUserID:        webUserIDPtr,
		BookingType:      models.CorporateBookingTypeMeals,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.Meal.ReasonForBooking,
		Notes:            req.Meal.Notes,
		Documents:        req.Documents,
		Payload:          json.RawMessage(payloadBytes),
		ProfileName:      req.BookedBy.Name,
		CompanyName:      req.Company.Name,
	}
	if req.Approver != nil {
		r.AuthoriserName = req.Approver.Name
		r.AuthoriserEmail = req.Approver.Email
		r.AuthoriserPhone = req.Approver.Phone
		r.AuthoriserTitle = req.Approver.Title
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

// ─── Backoffice ───────────────────────────────────────────────────────────────

func (s *CorporateBookingRequestService) GetByID(id, orgID uuid.UUID) (*models.CorporateBookingRequest, error) {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}
	if req.BookingType == models.CorporateBookingTypeMeals {
		req.MealsSummary = s.buildMealsSummary(orgID, req.Payload)
	}
	return req, nil
}

// buildMealsSummary resolves each menu_item_id in a meals payload to its current
// name + price (from menu_items, the same source billing uses) and computes
// per-line, per-guest, and grand totals for back-office display. Items whose menu
// item no longer exists are shown with a fallback name and zero price.
func (s *CorporateBookingRequestService) buildMealsSummary(orgID uuid.UUID, payload json.RawMessage) *models.MealsRequestSummary {
	if s.menuRepo == nil || len(payload) == 0 {
		return nil
	}
	var p models.SubmitMealsRequest
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil
	}

	// Cache lookups so a shared item is only fetched once.
	cache := map[uuid.UUID]*models.MenuItem{}
	resolve := func(in models.CorMealItemInput) models.MealsSummaryItem {
		mi, ok := cache[in.MenuItemID]
		if !ok {
			mi, _ = s.menuRepo.GetMenuItemByID(in.MenuItemID, orgID)
			cache[in.MenuItemID] = mi
		}
		item := models.MealsSummaryItem{
			MenuItemID: in.MenuItemID,
			Quantity:   in.Quantity,
			Notes:      in.Notes,
			Name:       "(unknown item)",
		}
		if mi != nil {
			item.Name = mi.Name
			item.UnitPrice = mi.Price
		}
		item.Subtotal = item.UnitPrice * float64(in.Quantity)
		return item
	}

	summary := &models.MealsRequestSummary{
		From:         p.From,
		To:           p.To,
		Headcount:    p.Headcount,
		DietaryNotes: p.DietaryNotes,
	}

	for _, g := range p.Guests {
		if len(g.Items) == 0 {
			continue
		}
		sg := models.MealsSummaryGuest{
			Name:               strings.TrimSpace(g.FirstName + " " + g.LastName),
			IdentificationCard: g.IdentificationCard,
		}
		for _, in := range g.Items {
			line := resolve(in)
			sg.Items = append(sg.Items, line)
			sg.Subtotal += line.Subtotal
		}
		summary.EstimatedTotal += sg.Subtotal
		summary.Guests = append(summary.Guests, sg)
	}

	for _, in := range p.Items {
		line := resolve(in)
		summary.BuffetItems = append(summary.BuffetItems, line)
		summary.EstimatedTotal += line.Subtotal
	}

	return summary
}

func (s *CorporateBookingRequestService) List(orgID uuid.UUID, bookingType, status string, page, pageSize int) ([]models.CorporateBookingRequest, int, error) {
	return s.requestRepo.List(orgID, bookingType, status, page, pageSize)
}

// ApproveFromWorkflow / RejectFromWorkflow satisfy the workflow's
// BookingRequestApprover interface, delegating to the same Approve/Reject the
// request endpoints use.
func (s *CorporateBookingRequestService) ApproveFromWorkflow(id, orgID uuid.UUID) error {
	return s.Approve(id, orgID)
}

func (s *CorporateBookingRequestService) RejectFromWorkflow(id, orgID uuid.UUID) error {
	return s.Reject(id, orgID)
}

func (s *CorporateBookingRequestService) Approve(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status != models.CorporateBookingStatusPending {
		return errors.New("only pending requests can be approved")
	}
	if err := s.requestRepo.UpdateStatus(id, orgID, models.CorporateBookingStatusApproved); err != nil {
		return err
	}

	// Event and meals requests are self-contained (venue/menu chosen at submission),
	// so approval materialises the booking automatically — no separate staff step.
	// Accommodation still needs staff to assign rooms, so it stays approved-only.
	if s.booking != nil && req.BookingType == models.CorporateBookingTypeEvent {
		var lodgeBranchID *uuid.UUID
		var envelope models.SubmitEventBookingRequest
		if json.Unmarshal(req.Payload, &envelope) == nil && envelope.Event != nil && len(envelope.Event.Sessions) > 0 {
			if vID, err := uuid.Parse(envelope.Event.Sessions[0].VenueID); err == nil {
				lodgeBranchID = s.venueBranchID(orgID, vID)
			}
		} else {
			var legacy models.SubmitEventRequest
			_ = json.Unmarshal(req.Payload, &legacy)
			lodgeBranchID = s.venueBranchID(orgID, legacy.VenueID)
		}

		if _, err := s.booking.CreateFromRequest(orgID, lodgeBranchID, id, nil, req.WebUserID); err != nil {
			return fmt.Errorf("request approved but booking creation failed: %w", err)
		}
	}
	if s.booking != nil && req.BookingType == models.CorporateBookingTypeMeals {
		var envelope models.SubmitMealBookingRequest
		if json.Unmarshal(req.Payload, &envelope) == nil && envelope.Meal != nil && len(envelope.Meal.Sessions) > 0 {
			if _, err := s.booking.CreateCorporateMeal(orgID, req.CorProfileID, req.CompanyID, req.WebUserID, &envelope); err != nil {
				return fmt.Errorf("request approved but meals booking creation failed: %w", err)
			}
		} else {
			if _, err := s.booking.CreateFromRequest(orgID, nil, id, nil, req.WebUserID); err != nil {
				return fmt.Errorf("request approved but booking creation failed: %w", err)
			}
		}
	}
	return nil
}

func (s *CorporateBookingRequestService) Reject(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status != models.CorporateBookingStatusPending {
		return errors.New("only pending requests can be rejected")
	}
	return s.requestRepo.UpdateStatus(id, orgID, models.CorporateBookingStatusRejected)
}

func (s *CorporateBookingRequestService) Cancel(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status == models.CorporateBookingStatusApproved {
		return errors.New("approved requests cannot be cancelled")
	}
	return s.requestRepo.UpdateStatus(id, orgID, models.CorporateBookingStatusCancelled)
}

// ─── Workflow ─────────────────────────────────────────────────────────────────

func (s *CorporateBookingRequestService) startWorflow(r *models.CorporateBookingRequest, orgID uuid.UUID) {
	if s.workflow == nil {
		return
	}
	go func() {
		taskDetails := models.TaskDetails{
			TaskID:          r.ID.String(),
			TaskRef:         r.ID.String()[:8],
			TaskType:        "corporate_booking",
			TaskDescription: "Corporate " + r.BookingType + " booking request from " + r.CompanyName,
			SenderDetails: models.SenderDetails{
				SenderID:   r.ID.String(),
				SenderName: r.CompanyName,
				Position:   r.BookingType,
				Department: "Guest",
			},
		}
		if _, err := s.workflow.InitiateWorkflow(
			models.WorkflowTypeBookingApproval,
			taskDetails,
			r.ID.String(),
			"medium",
			nil,
			orgID.String(),
		); err != nil {
			_ = err
		}
	}()
}

// splitName splits a full name string into first and last name.
// Everything after the first space is treated as the last name.
// If there is no space, the whole string is the first name.
func splitName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	i := strings.Index(full, " ")
	if i < 0 {
		return full, ""
	}
	return strings.TrimSpace(full[:i]), strings.TrimSpace(full[i+1:])
}
