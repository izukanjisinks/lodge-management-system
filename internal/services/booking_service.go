package services

import (
	"errors"
	"fmt"
	"log"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type BookingService struct {
	repo    *repository.BookingRepository
	room    *repository.RoomRepository
	client  *repository.ClientRepository
	invoice *InvoiceService
}

func NewBookingService(repo *repository.BookingRepository, room *repository.RoomRepository, client *repository.ClientRepository) *BookingService {
	return &BookingService{repo: repo, room: room, client: client}
}

// SetInvoiceService injects the invoice service after construction to avoid a circular dependency.
func (s *BookingService) SetInvoiceService(invoice *InvoiceService) {
	s.invoice = invoice
}


func (s *BookingService) CreateIndividual(orgID uuid.UUID, req *models.CreateIndividualBookingRequest) (*models.Booking, error) {
	// Resolve client — look up existing or create on the fly
	var clientID uuid.UUID
	if req.ClientID != nil {
		if _, err := s.client.GetIndividualByID(*req.ClientID, orgID); err != nil {
			return nil, errors.New("client not found")
		}
		clientID = *req.ClientID
	} else {
		if req.Client == nil || req.Client.FullName == "" {
			return nil, errors.New("client.full_name is required when client_id is not provided")
		}
		c := &models.IndividualClient{
			FullName:         req.Client.FullName,
			Email:            req.Client.Email,
			Phone:            req.Client.Phone,
			IDPassportNumber: req.Client.IDPassportNumber,
			Status:           models.ClientStatusActive,
		}
		if err := s.client.CreateIndividual(c, orgID); err != nil {
			return nil, formatConstraintError(err)
		}
		clientID = c.ID
	}

	if req.RoomID == uuid.Nil {
		return nil, errors.New("room_id is required")
	}
	if req.CheckIn.IsZero() || req.CheckOut.IsZero() {
		return nil, errors.New("check_in and check_out are required")
	}
	if !req.CheckOut.After(req.CheckIn.Time) {
		return nil, errors.New("check_out must be after check_in")
	}
	if req.Guests <= 0 {
		return nil, errors.New("guests must be greater than 0")
	}

	room, err := s.room.GetByID(req.RoomID, orgID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if req.Guests > room.Capacity {
		return nil, fmt.Errorf("guests (%d) exceed room capacity (%d)", req.Guests, room.Capacity)
	}

	available, err := s.repo.IsRoomAvailable(req.RoomID, req.CheckIn.Time, req.CheckOut.Time, nil)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for the selected dates")
	}

	duplicate, err := s.repo.HasActiveBookingForClient(clientID, req.RoomID, req.CheckIn.Time, req.CheckOut.Time)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("client already has an active booking for this room on the selected dates")
	}

	b := &models.Booking{
		RoomID:          req.RoomID,
		ClientID:        clientID,
		ClientType:      models.BookingClientTypeIndividual,
		CheckIn:         req.CheckIn.Time,
		CheckOut:        req.CheckOut.Time,
		Guests:          req.Guests,
		Status:          models.BookingStatusPending,
		SpecialRequests: req.SpecialRequests,
	}
	if err := s.repo.Create(b, orgID); err != nil {
		return nil, err
	}

	return s.repo.GetByID(b.ID, orgID)
}

func (s *BookingService) CreateCorporate(orgID uuid.UUID, req *models.CreateCorporateBookingRequest) (*models.CorporateBookingResponse, error) {
	if len(req.Guests) == 0 {
		return nil, errors.New("at least one guest is required for a corporate booking")
	}

	// Validate all guests up front before touching the DB
	for i, g := range req.Guests {
		if g.RoomID == uuid.Nil {
			return nil, fmt.Errorf("guest %d: room_id is required", i+1)
		}
		if g.CheckIn.IsZero() || g.CheckOut.IsZero() {
			return nil, fmt.Errorf("guest %d: check_in and check_out are required", i+1)
		}
		if !g.CheckOut.After(g.CheckIn.Time) {
			return nil, fmt.Errorf("guest %d: check_out must be after check_in", i+1)
		}
		if g.ClientID == nil && g.FullName == "" {
			return nil, fmt.Errorf("guest %d: full_name is required when client_id is not provided", i+1)
		}
	}

	// Validate rooms and availability before starting the transaction
	for i, g := range req.Guests {
		room, err := s.room.GetByID(g.RoomID, orgID)
		if err != nil {
			return nil, fmt.Errorf("guest %d: room not found", i+1)
		}
		if room.Capacity < 1 {
			return nil, fmt.Errorf("guest %d: room has no capacity", i+1)
		}
		available, err := s.repo.IsRoomAvailable(g.RoomID, g.CheckIn.Time, g.CheckOut.Time, nil)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, fmt.Errorf("guest %d: room is not available for the selected dates", i+1)
		}
	}

	tx, err := s.repo.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Resolve or create the corporate client
	var corp *models.CorporateClient
	if req.ClientID != nil {
		corp, err = s.client.GetCorporateByID(*req.ClientID, orgID)
		if err != nil {
			return nil, errors.New("corporate client not found")
		}
	} else {
		if req.Client == nil || req.Client.CompanyName == "" {
			return nil, errors.New("client.company_name is required when client_id is not provided")
		}
		corp = &models.CorporateClient{
			CompanyName:      req.Client.CompanyName,
			ContactPerson:    req.Client.ContactPerson,
			Email:            req.Client.Email,
			Phone:            req.Client.Phone,
			CompanyRegNumber: req.Client.CompanyRegNumber,
			Industry:         req.Client.Industry,
			Status:           models.ClientStatusActive,
		}
		if err = s.client.CreateCorporateInTx(tx, corp, orgID); err != nil {
			log.Printf("[booking] failed to create corporate client (org %s): %v", orgID, err)
			return nil, fmt.Errorf("failed to create corporate client: %w", err)
		}
	}

	// Create each guest's individual record and booking
	var bookingIDs []uuid.UUID
	for i, g := range req.Guests {
		var guestClientID uuid.UUID
		if g.ClientID != nil {
			individual, lookupErr := s.client.GetIndividualByID(*g.ClientID, orgID)
			if lookupErr != nil {
				err = fmt.Errorf("guest %d: client not found", i+1)
				return nil, err
			}
			guestClientID = individual.ID
		} else {
			individual := &models.IndividualClient{
				FullName:         g.FullName,
				Email:            g.Email,
				Phone:            g.Phone,
				IDPassportNumber: g.IDNumber,
				Status:           models.ClientStatusActive,
			}
			if err = s.client.CreateIndividualInTx(tx, individual, orgID); err != nil {
				log.Printf("[booking] guest %d: failed to create individual client (org %s, name %q): %v", i+1, orgID, individual.FullName, err)
				err = fmt.Errorf("guest %d: %w", i+1, formatConstraintError(err))
				return nil, err
			}
			guestClientID = individual.ID
		}

		b := &models.Booking{
			RoomID:            g.RoomID,
			ClientID:          guestClientID,
			ClientType:        models.BookingClientTypeIndividual,
			CorporateClientID: &corp.ID,
			CheckIn:           g.CheckIn.Time,
			CheckOut:          g.CheckOut.Time,
			Guests:            1,
			Status:            models.BookingStatusPending,
			SpecialRequests:   g.SpecialRequests,
		}
		if err = s.repo.CreateInTx(tx, b, orgID); err != nil {
			log.Printf("[booking] guest %d: failed to save booking (org %s, room %s, client %s): %v", i+1, orgID, b.RoomID, b.ClientID, err)
			err = fmt.Errorf("guest %d: failed to create booking: %w", i+1, err)
			return nil, err
		}
		bookingIDs = append(bookingIDs, b.ID)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[booking] failed to commit corporate booking transaction (org %s, corporate client %s): %v", orgID, corp.ID, err)
		return nil, err
	}

	// Generate consolidated invoice for the corporate client
	if s.invoice != nil {
		if invErr := s.invoice.GenerateCorporateInvoice(corp.ID, orgID, bookingIDs); invErr != nil {
			log.Printf("[booking] warning: failed to generate corporate invoice for corp %s: %v", corp.ID, invErr)
		}
	}

	// Fetch the created bookings for the response
	var bookings []models.Booking
	var totalAmount float64
	for _, id := range bookingIDs {
		b, fetchErr := s.repo.GetByID(id, orgID)
		if fetchErr != nil {
			continue
		}
		bookings = append(bookings, *b)
		totalAmount += b.TotalAmount
	}

	return &models.CorporateBookingResponse{
		CorporateClientID: corp.ID,
		CompanyName:       corp.CompanyName,
		Bookings:          bookings,
		TotalAmount:       totalAmount,
	}, nil
}

func (s *BookingService) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Booking, error) {
	return s.repo.GetByID(id, orgID)
}

func (s *BookingService) List(orgID uuid.UUID, status, clientType string, clientID *uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	return s.repo.List(orgID, status, clientType, clientID, page, pageSize)
}

func (s *BookingService) Update(id uuid.UUID, orgID uuid.UUID, req *models.UpdateBookingRequest) (*models.Booking, error) {
	b, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Only pending bookings can be edited
	if b.Status != models.BookingStatusCheckedIn && b.Status != models.BookingStatusConfirmed && b.Status != models.BookingStatusPending {
		return nil, errors.New("only pending, checked-in, or confirmed bookings can be updated")
	}

	if req.CheckIn != nil {
		b.CheckIn = req.CheckIn.Time
	}
	if req.CheckOut != nil {
		b.CheckOut = req.CheckOut.Time
	}
	if req.Guests != nil {
		b.Guests = *req.Guests
	}
	if req.SpecialRequests != nil {
		b.SpecialRequests = *req.SpecialRequests
	}

	if !b.CheckOut.After(b.CheckIn) {
		return nil, errors.New("check_out must be after check_in")
	}
	if b.Guests <= 0 {
		return nil, errors.New("guests must be greater than 0")
	}

	// Re-check room availability excluding this booking
	available, err := s.repo.IsRoomAvailable(b.RoomID, b.CheckIn, b.CheckOut, &id)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for the selected dates")
	}

	if err := s.repo.Update(b, orgID); err != nil {
		return nil, err
	}

	// Recalculate room charge and due date if either date changed
	if (req.CheckIn != nil || req.CheckOut != nil) && s.invoice != nil {
		if err := s.invoice.RecalculateRoomCharge(id, orgID); err != nil {
			fmt.Printf("warning: failed to recalculate room charge for booking %s: %v\n", id, err)
		}
	}

	return s.repo.GetByID(id, orgID)
}

func (s *BookingService) UpdateStatus(id uuid.UUID, orgID uuid.UUID, newStatus string) (*models.Booking, error) {
	b, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Validate transition
	allowed := models.ValidBookingTransitions[b.Status]
	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("cannot transition booking from '%s' to '%s'", b.Status, newStatus)
	}

	if err := s.repo.UpdateStatusTx(id, orgID, newStatus); err != nil {
		return nil, err
	}

	// Auto-generate invoice when booking is confirmed.
	// Corporate guest bookings are covered by the consolidated corporate invoice — skip them.
	if newStatus == models.BookingStatusConfirmed && s.invoice != nil && b.CorporateClientID == nil {
		if err := s.invoice.GenerateForBooking(id, orgID); err != nil {
			fmt.Printf("warning: failed to generate invoice for booking %s: %v\n", id, err)
		}
	}

	return s.repo.GetByID(id, orgID)
}

// ClearOverstayed manually resolves the overstayed flag set by the nightly job.
func (s *BookingService) ClearOverstayed(id uuid.UUID, orgID uuid.UUID) error {
	_, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return errors.New("booking not found")
	}
	return s.repo.ClearOverstayed(id, orgID)
}

func (s *BookingService) Delete(id uuid.UUID, orgID uuid.UUID) error {
	b, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return errors.New("booking not found")
	}
	if b.Status == models.BookingStatusCheckedIn {
		return errors.New("cannot delete a booking that is currently checked in")
	}
	return s.repo.Delete(id, orgID)
}
