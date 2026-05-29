package services

import (
	"errors"
	"fmt"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type GuestBookingService struct {
	bookingRepo    *repository.BookingRepository
	roomRepo       *repository.RoomRepository
	guestAuth      *GuestAuthService
	workflow       *WorkflowService
	bookingService *BookingService
}

func NewGuestBookingService(
	bookingRepo *repository.BookingRepository,
	roomRepo *repository.RoomRepository,
	guestAuth *GuestAuthService,
) *GuestBookingService {
	return &GuestBookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		guestAuth:   guestAuth,
	}
}

func (s *GuestBookingService) SetBookingService(svc *BookingService) {
	s.bookingService = svc
}

// SetWorkflowService injects the workflow service after construction to avoid a circular dependency.
func (s *GuestBookingService) SetWorkflowService(workflow *WorkflowService) {
	s.workflow = workflow
}

// Create makes a booking on behalf of the logged-in guest.
// client_id and client_type are resolved automatically from the JWT user.
func (s *GuestBookingService) Create(userID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error) {
	profile, err := s.guestAuth.GetProfileByGuestID(userID)
	if err != nil {
		return nil, errors.New("guest profile not found — please complete your registration")
	}

	if req.IDPassportNumber != "" && profile.IDPassportNumber != req.IDPassportNumber {
		if err := s.guestAuth.UpdateProfileIDPassport(profile.ID, req.IDPassportNumber); err != nil {
			return nil, fmt.Errorf("failed to save ID/passport number: %w", err)
		}
		profile.IDPassportNumber = req.IDPassportNumber
	}

	// Look up the room without an org filter — guests have no org in their JWT.
	// orgID is derived from the room itself and used for all subsequent scoped calls.
	room, err := s.roomRepo.GetByIDUnscoped(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	orgID := uuid.Nil
	if room.OrgID != nil {
		orgID = *room.OrgID
	}

	if req.Guests > room.Capacity {
		return nil, fmt.Errorf("room capacity is %d, requested %d guests", room.Capacity, req.Guests)
	}
	if req.CheckOut.Before(req.CheckIn.Time) || req.CheckOut.Equal(req.CheckIn.Time) {
		return nil, errors.New("check_out must be after check_in")
	}

	available, err := s.bookingRepo.IsRoomAvailable(req.RoomID, req.CheckIn.Time, req.CheckOut.Time, nil)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for the selected dates")
	}

	b := &models.Booking{
		RoomID:          req.RoomID,
		BranchID:        room.BranchID,
		ClientID:        profile.ID,
		ClientType:      models.BookingClientTypeIndividual,
		CheckIn:         req.CheckIn.Time,
		CheckOut:        req.CheckOut.Time,
		Guests:          req.Guests,
		Status:          models.BookingStatusPending,
		SpecialRequests: req.SpecialRequests,
	}

	if err := s.bookingRepo.Create(b, orgID); err != nil {
		return nil, err
	}

	if err := s.guestAuth.UpdateProfileOrg(userID, orgID); err != nil {
		fmt.Printf("warning: failed to stamp org on guest profile for guest %s: %v\n", userID, err)
	}

	created, err := s.bookingRepo.GetByIDUnscoped(b.ID)
	if err != nil {
		return nil, err
	}

	if s.workflow != nil {
		go func() {
			taskDetails := models.TaskDetails{
				TaskID:          created.ID.String(),
				TaskRef:         created.BookingNumber,
				TaskType:        "booking",
				TaskDescription: fmt.Sprintf("Review booking request for %s — %s to %s (%d guest(s))", created.ClientName, created.CheckIn.Format("2006-01-02"), created.CheckOut.Format("2006-01-02"), created.Guests),
				SenderDetails: models.SenderDetails{
					SenderID:   created.ClientID.String(),
					SenderName: created.ClientName,
					Position:   created.ClientType,
					Department: "Guest",
				},
			}
			if _, err := s.workflow.InitiateWorkflow(
				models.WorkflowTypeBookingApproval,
				taskDetails,
				userID.String(),
				"medium",
				nil,
				orgID.String(),
			); err != nil {
				fmt.Printf("warning: failed to initiate booking approval workflow for booking %s: %v\n", created.ID, err)
			}
		}()
	}

	return created, nil
}

// ListForGuest returns all bookings belonging to the logged-in guest.
func (s *GuestBookingService) ListForGuest(userID uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	profile, err := s.guestAuth.GetProfileByGuestID(userID)
	if err != nil {
		return nil, 0, errors.New("guest profile not found")
	}

	return s.bookingRepo.List(uuid.Nil, nil, "", models.BookingClientTypeIndividual, &profile.ID, page, pageSize)
}

// GetByID returns a single booking, scoped to the guest — 403 if it doesn't belong to them.
func (s *GuestBookingService) GetByID(userID uuid.UUID, bookingID uuid.UUID) (*models.Booking, error) {
	profile, err := s.guestAuth.GetProfileByGuestID(userID)
	if err != nil {
		return nil, errors.New("guest profile not found")
	}

	b, err := s.bookingRepo.GetByIDUnscoped(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if b.ClientID != profile.ID {
		return nil, errors.New("forbidden")
	}
	return b, nil
}

// CreateCorporate makes a corporate booking on behalf of a guest.
// orgID is derived from the first guest's room so no org is needed in the JWT.
func (s *GuestBookingService) CreateCorporate(req *models.CreateCorporateBookingRequest) (*models.CorporateBookingResponse, error) {
	if s.bookingService == nil {
		return nil, errors.New("corporate bookings are not available")
	}
	if len(req.Guests) == 0 {
		return nil, errors.New("at least one guest is required")
	}

	// Derive orgID from the first room — all rooms should belong to the same org.
	room, err := s.roomRepo.GetByIDUnscoped(req.Guests[0].RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	orgID := uuid.Nil
	if room.OrgID != nil {
		orgID = *room.OrgID
	}

	return s.bookingService.CreateCorporate(orgID, req)
}

// Cancel transitions a guest's booking to cancelled — only allowed from pending or confirmed.
func (s *GuestBookingService) Cancel(userID uuid.UUID, bookingID uuid.UUID) error {
	b, err := s.GetByID(userID, bookingID)
	if err != nil {
		return err
	}

	allowed := models.ValidBookingTransitions[b.Status]
	canCancel := false
	for _, s := range allowed {
		if s == models.BookingStatusCancelled {
			canCancel = true
			break
		}
	}
	if !canCancel {
		return fmt.Errorf("cannot cancel a booking with status %q", b.Status)
	}

	// Derive orgID from the room so the status update is scoped correctly.
	orgID := uuid.Nil
	if room, err := s.roomRepo.GetByIDUnscoped(b.RoomID); err == nil && room.OrgID != nil {
		orgID = *room.OrgID
	}
	return s.bookingRepo.UpdateStatusTx(bookingID, orgID, models.BookingStatusCancelled)
}
