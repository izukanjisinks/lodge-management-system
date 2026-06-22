package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type IndividualBookingRequestService struct {
	requestRepo    *repository.IndividualBookingRequestRepository
	roomRepo       *repository.RoomRepository
	bookingService *BookingService
	workflow       *WorkflowService
}

func NewIndividualBookingRequestService(
	requestRepo *repository.IndividualBookingRequestRepository,
	roomRepo *repository.RoomRepository,
	bookingService *BookingService,
) *IndividualBookingRequestService {
	return &IndividualBookingRequestService{
		requestRepo:    requestRepo,
		roomRepo:       roomRepo,
		bookingService: bookingService,
	}
}

func (s *IndividualBookingRequestService) SetWorkflowService(svc *WorkflowService) {
	s.workflow = svc
}

// ─── Submission ───────────────────────────────────────────────────────────────

func (s *IndividualBookingRequestService) Submit(guestID uuid.UUID, req *models.SubmitIndividualBookingRequest) (*models.IndividualBookingRequest, error) {
	if req.BookedBy.Name == "" || req.BookedBy.Email == "" {
		return nil, errors.New("booked_by name and email are required")
	}
	if req.Accommodation == nil {
		return nil, errors.New("accommodation block is required")
	}
	if req.Accommodation.CheckIn == "" || req.Accommodation.CheckOut == "" {
		return nil, errors.New("check_in and check_out are required")
	}
	if len(req.Accommodation.Rooms) == 0 {
		return nil, errors.New("at least one room must be selected")
	}

	// Resolve org: use org_id from body, or fall back to the first room's org
	orgID := req.OrgID
	if orgID == uuid.Nil {
		firstRoom, err := s.roomRepo.GetByIDUnscoped(req.Accommodation.Rooms[0].RoomID)
		if err != nil {
			return nil, errors.New("room not found")
		}
		if firstRoom.OrgID != nil {
			orgID = *firstRoom.OrgID
		}
	}
	if orgID == uuid.Nil {
		return nil, errors.New("could not resolve org_id")
	}

	// Store the entire envelope as-is in JSONB
	payloadBytes, _ := json.Marshal(req)

	r := &models.IndividualBookingRequest{
		OrgID:       orgID,
		WebUserID:   &guestID,
		BookerName:  req.BookedBy.Name,
		BookerEmail: req.BookedBy.Email,
		BookerPhone: req.BookedBy.Phone,
		BookingType: models.BookingTypeRoom,
		Status:      models.IndividualBookingStatusPending,
		Notes:       req.Accommodation.Notes,
		Documents:   req.Documents,
		Payload:     json.RawMessage(payloadBytes),
	}
	if r.Documents == nil {
		r.Documents = []string{}
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorkflow(r, orgID)
	return r, nil
}

// SubmitEvent stores a standalone individual event request (Flow B). It reuses the
// shared event envelope (no company/approver for individuals), persists the whole
// payload as JSONB with booking_type='event', and starts the approval workflow.
func (s *IndividualBookingRequestService) SubmitEvent(guestID uuid.UUID, req *models.SubmitEventBookingRequest) (*models.IndividualBookingRequest, error) {
	if req.BookedBy.Name == "" || req.BookedBy.Email == "" {
		return nil, errors.New("booked_by name and email are required")
	}
	if req.Event == nil || len(req.Event.Sessions) == 0 {
		return nil, errors.New("at least one event session is required")
	}
	for i, sess := range req.Event.Sessions {
		if sess.VenueID == "" {
			return nil, fmt.Errorf("session %d is missing a venue", i+1)
		}
		if _, err := uuid.Parse(sess.VenueID); err != nil {
			return nil, fmt.Errorf("session %d has an invalid venue", i+1)
		}
	}

	orgID := req.OrgID
	if orgID == uuid.Nil {
		return nil, errors.New("org_id is required")
	}

	payloadBytes, _ := json.Marshal(req)

	r := &models.IndividualBookingRequest{
		OrgID:       orgID,
		WebUserID:   &guestID,
		BookerName:  req.BookedBy.Name,
		BookerEmail: req.BookedBy.Email,
		BookerPhone: req.BookedBy.Phone,
		BookingType: models.BookingTypeEvent,
		Status:      models.IndividualBookingStatusPending,
		Notes:       req.Event.Notes,
		Documents:   req.Documents,
		Payload:     json.RawMessage(payloadBytes),
	}
	if r.Documents == nil {
		r.Documents = []string{}
	}

	if err := s.requestRepo.Create(r); err != nil {
		return nil, err
	}

	s.startWorkflow(r, orgID)
	return r, nil
}

// ─── Web user ─────────────────────────────────────────────────────────────────

func (s *IndividualBookingRequestService) GetForWebUser(id, webUserID uuid.UUID) (*models.IndividualBookingRequest, error) {
	return s.requestRepo.GetByIDForWebUser(id, webUserID)
}

func (s *IndividualBookingRequestService) ListForWebUser(webUserID uuid.UUID, page, pageSize int) ([]models.IndividualBookingRequest, int, error) {
	return s.requestRepo.ListByWebUser(webUserID, page, pageSize)
}

func (s *IndividualBookingRequestService) CancelForWebUser(id, webUserID uuid.UUID) error {
	req, err := s.requestRepo.GetByIDForWebUser(id, webUserID)
	if err != nil {
		return err
	}
	if req.Status != models.IndividualBookingStatusPending {
		return fmt.Errorf("only pending requests can be cancelled")
	}
	return s.requestRepo.UpdateStatus(id, req.OrgID, models.IndividualBookingStatusCancelled)
}

// ─── Backoffice ───────────────────────────────────────────────────────────────

func (s *IndividualBookingRequestService) GetByID(id, orgID uuid.UUID) (*models.IndividualBookingRequest, error) {
	return s.requestRepo.GetByID(id, orgID)
}

func (s *IndividualBookingRequestService) List(orgID uuid.UUID, status string, page, pageSize int) ([]models.IndividualBookingRequest, int, error) {
	return s.requestRepo.List(orgID, status, page, pageSize)
}

// ApproveFromWorkflow adapts Approve to the workflow's BookingRequestApprover
// interface — it materialises the booking but discards the returned record, since
// the workflow only needs the success/failure signal.
func (s *IndividualBookingRequestService) ApproveFromWorkflow(id, orgID uuid.UUID) error {
	_, err := s.Approve(id, orgID)
	return err
}

// RejectFromWorkflow adapts Reject to the workflow's BookingRequestApprover interface.
func (s *IndividualBookingRequestService) RejectFromWorkflow(id, orgID uuid.UUID) error {
	return s.Reject(id, orgID)
}

func (s *IndividualBookingRequestService) Approve(id, orgID uuid.UUID) (*models.Booking, error) {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}
	if req.Status != models.IndividualBookingStatusPending {
		return nil, errors.New("only pending requests can be approved")
	}

	// Event requests materialise via the shared multi-session engine (one
	// booking_events row per session), not the room-assignment path below.
	if req.BookingType == models.BookingTypeEvent {
		var envelope models.SubmitEventBookingRequest
		if jsonErr := json.Unmarshal(req.Payload, &envelope); jsonErr != nil || envelope.Event == nil {
			return nil, errors.New("invalid event request payload")
		}
		booking, err := s.bookingService.CreateIndividualEvent(orgID, req.WebUserID, &envelope)
		if err != nil {
			return nil, err
		}
		if err := s.requestRepo.UpdateStatus(id, orgID, models.IndividualBookingStatusApproved); err != nil {
			return nil, err
		}
		return booking, nil
	}

	// Decode the stored payload — try the new unified envelope first,
	// fall back to the old flat shape for requests submitted before this refactor.
	var roomID uuid.UUID
	var checkInStr, checkOutStr string

	var newPayload models.SubmitIndividualBookingRequest
	if jsonErr := json.Unmarshal(req.Payload, &newPayload); jsonErr == nil && newPayload.Accommodation != nil {
		// New envelope shape
		if len(newPayload.Accommodation.Rooms) == 0 {
			return nil, errors.New("no rooms in accommodation block")
		}
		roomID = newPayload.Accommodation.Rooms[0].RoomID
		checkInStr = newPayload.Accommodation.CheckIn
		checkOutStr = newPayload.Accommodation.CheckOut
	} else {
		// Legacy flat shape
		var legacy models.IndividualBookingPayload
		if jsonErr := json.Unmarshal(req.Payload, &legacy); jsonErr != nil {
			return nil, errors.New("invalid request payload")
		}
		roomID = legacy.RoomID
		checkInStr = legacy.CheckIn
		checkOutStr = legacy.CheckOut
	}

	checkIn := models.DateOnly{}
	if err := checkIn.UnmarshalJSON([]byte(`"` + checkInStr + `"`)); err != nil {
		return nil, fmt.Errorf("invalid check_in date: %w", err)
	}
	checkOut := models.DateOnly{}
	if err := checkOut.UnmarshalJSON([]byte(`"` + checkOutStr + `"`)); err != nil {
		return nil, fmt.Errorf("invalid check_out date: %w", err)
	}

	room, err := s.roomRepo.GetByIDUnscoped(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	bookingReq := &models.CreateIndividualBookingRequest{
		WebUserID:   req.WebUserID,
		BookerName:  req.BookerName,
		BookerEmail: req.BookerEmail,
		BookerPhone: req.BookerPhone,
		RoomID:      roomID,
		CheckIn:     checkIn,
		CheckOut:    checkOut,
	}

	booking, err := s.bookingService.CreateIndividual(orgID, room.BranchID, bookingReq)
	if err != nil {
		return nil, err
	}

	if err := s.requestRepo.UpdateStatus(id, orgID, models.IndividualBookingStatusApproved); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *IndividualBookingRequestService) Reject(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status != models.IndividualBookingStatusPending {
		return errors.New("only pending requests can be rejected")
	}
	return s.requestRepo.UpdateStatus(id, orgID, models.IndividualBookingStatusRejected)
}

func (s *IndividualBookingRequestService) Cancel(id, orgID uuid.UUID) error {
	req, err := s.requestRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	if req.Status == models.IndividualBookingStatusApproved {
		return errors.New("approved requests cannot be cancelled")
	}
	return s.requestRepo.UpdateStatus(id, orgID, models.IndividualBookingStatusCancelled)
}

// ─── Workflow ─────────────────────────────────────────────────────────────────

func (s *IndividualBookingRequestService) startWorkflow(r *models.IndividualBookingRequest, orgID uuid.UUID) {
	if s.workflow == nil {
		return
	}
	go func() {
		taskDetails := models.TaskDetails{
			TaskID:          r.ID.String(),
			TaskRef:         r.ID.String()[:8],
			TaskType:        "individual_booking",
			TaskDescription: fmt.Sprintf("Room booking request from %s", r.BookerName),
			SenderDetails: models.SenderDetails{
				SenderID:   r.ID.String(),
				SenderName: r.BookerName,
				Position:   "Guest",
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
