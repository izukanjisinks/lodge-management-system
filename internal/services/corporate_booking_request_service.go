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

// ─── Submission ───────────────────────────────────────────────────────────────

func (s *CorporateBookingRequestService) SubmitAccommodation(orgID uuid.UUID, req *models.SubmitAccommodationRequest) (*models.CorporateBookingRequest, error) {
	if len(req.Guests) == 0 {
		return nil, errors.New("at least one guest is required")
	}
	if req.Profile.FirstName == "" || req.Profile.Email == "" {
		return nil, errors.New("booked_by first_name and email are required")
	}

	// Each guest must have a unique identification card. The roster is keyed on
	// (profile, id_card), so blanks or in-request duplicates would silently collapse
	// two delegates into one record.
	seen := make(map[string]bool, len(req.Guests))
	for i, g := range req.Guests {
		if strings.TrimSpace(g.IdentificationCard) == "" {
			return nil, fmt.Errorf("guest %d (%s %s) is missing an identification card", i+1, g.FirstName, g.LastName)
		}
		if seen[g.IdentificationCard] {
			return nil, fmt.Errorf("guest %d has a duplicate identification card (%s) within this request", i+1, g.IdentificationCard)
		}
		seen[g.IdentificationCard] = true
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, req.Company, req.Branch, req.Profile)
	if err != nil {
		return nil, err
	}

	// Create corporate_guests rows for each guest in the request
	if _, err := s.guestRepo.CreateMany(profileID, req.Guests); err != nil {
		return nil, errors.New("failed to save guests: " + err.Error())
	}

	payloadBytes, _ := json.Marshal(req)
	payload := json.RawMessage(payloadBytes)

	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         req.BranchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		BookingType:      models.CorporateBookingTypeAccommodation,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.ReasonForBooking,
		Notes:            req.Notes,
		Documents:        req.Documents,
		Payload:          payload,
	}
	if branchID != nil {
		r.BranchID = branchID
	}
	if req.Authoriser != nil {
		r.AuthoriserName = req.Authoriser.Name
		r.AuthoriserEmail = req.Authoriser.Email
		r.AuthoriserPhone = req.Authoriser.Phone
		r.AuthoriserTitle = req.Authoriser.Title
		r.AuthoriserDepartment = req.Authoriser.Department
		r.AuthoriserGLCode = req.Authoriser.GLCode
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

func (s *CorporateBookingRequestService) SubmitMeals(orgID uuid.UUID, req *models.SubmitMealsRequest) (*models.CorporateBookingRequest, error) {
	if req.Profile.FirstName == "" || req.Profile.Email == "" {
		return nil, errors.New("booked_by first_name and email are required")
	}
	if req.PlanType == "" || req.From == "" || req.To == "" {
		return nil, errors.New("plan_type, from, and to are required")
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, req.Company, req.Branch, req.Profile)
	if err != nil {
		return nil, err
	}

	payloadBytes, _ := json.Marshal(req)
	payload := json.RawMessage(payloadBytes)

	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		BookingType:      models.CorporateBookingTypeMeals,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.ReasonForBooking,
		Documents:        req.Documents,
		Payload:          payload,
	}
	if req.Authoriser != nil {
		r.AuthoriserName = req.Authoriser.Name
		r.AuthoriserEmail = req.Authoriser.Email
		r.AuthoriserPhone = req.Authoriser.Phone
		r.AuthoriserTitle = req.Authoriser.Title
		r.AuthoriserDepartment = req.Authoriser.Department
		r.AuthoriserGLCode = req.Authoriser.GLCode
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

func (s *CorporateBookingRequestService) SubmitConference(orgID uuid.UUID, req *models.SubmitConferenceRequest) (*models.CorporateBookingRequest, error) {
	if req.Profile.FirstName == "" || req.Profile.Email == "" {
		return nil, errors.New("booked_by first_name and email are required")
	}
	if req.StartDate == "" || req.StartTime == "" {
		return nil, errors.New("start_date and start_time are required")
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, req.Company, req.Branch, req.Profile)
	if err != nil {
		return nil, err
	}

	payloadBytes, _ := json.Marshal(req)
	payload := json.RawMessage(payloadBytes)

	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		BookingType:      models.CorporateBookingTypeConference,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.ReasonForBooking,
		Notes:            req.Notes,
		Documents:        req.Documents,
		Payload:          payload,
	}
	if req.Authoriser != nil {
		r.AuthoriserName = req.Authoriser.Name
		r.AuthoriserEmail = req.Authoriser.Email
		r.AuthoriserPhone = req.Authoriser.Phone
		r.AuthoriserTitle = req.Authoriser.Title
		r.AuthoriserDepartment = req.Authoriser.Department
		r.AuthoriserGLCode = req.Authoriser.GLCode
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

func (s *CorporateBookingRequestService) SubmitEvent(orgID uuid.UUID, req *models.SubmitEventRequest) (*models.CorporateBookingRequest, error) {
	if req.Profile.FirstName == "" || req.Profile.Email == "" {
		return nil, errors.New("booked_by first_name and email are required")
	}
	if req.EventType == "" || req.StartDate == "" {
		return nil, errors.New("event_type and start_date are required")
	}

	companyID, branchID, profileID, err := s.corProfile.ResolveChain(orgID, req.Company, req.Branch, req.Profile)
	if err != nil {
		return nil, err
	}

	payloadBytes, _ := json.Marshal(req)
	payload := json.RawMessage(payloadBytes)

	r := &models.CorporateBookingRequest{
		OrgID:            orgID,
		BranchID:         branchID,
		CorProfileID:     &profileID,
		CompanyID:        &companyID,
		BookingType:      models.CorporateBookingTypeEvent,
		Status:           models.CorporateBookingStatusPending,
		ReasonForBooking: req.ReasonForBooking,
		Notes:            req.Notes,
		Documents:        req.Documents,
		Payload:          payload,
	}
	if req.Authoriser != nil {
		r.AuthoriserName = req.Authoriser.Name
		r.AuthoriserEmail = req.Authoriser.Email
		r.AuthoriserPhone = req.Authoriser.Phone
		r.AuthoriserTitle = req.Authoriser.Title
		r.AuthoriserDepartment = req.Authoriser.Department
		r.AuthoriserGLCode = req.Authoriser.GLCode
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorflow(r, orgID)
	return r, nil
}

// ─── Backoffice ───────────────────────────────────────────────────────────────

func (s *CorporateBookingRequestService) GetByID(id, orgID uuid.UUID) (*models.CorporateBookingRequest, error) {
	return s.requestRepo.GetByID(id, orgID)
}

func (s *CorporateBookingRequestService) List(orgID uuid.UUID, bookingType, status string, page, pageSize int) ([]models.CorporateBookingRequest, int, error) {
	return s.requestRepo.List(orgID, bookingType, status, page, pageSize)
}

func (s *CorporateBookingRequestService) Approve(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status != models.CorporateBookingStatusPending {
		return errors.New("only pending requests can be approved")
	}
	return s.requestRepo.UpdateStatus(id, orgID, models.CorporateBookingStatusApproved)
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
