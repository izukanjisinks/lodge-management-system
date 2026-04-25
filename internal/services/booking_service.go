package services

import (
	"errors"
	"fmt"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type BookingService struct {
	repo     *repository.BookingRepository
	room     *repository.RoomRepository
	mealPlan *repository.MealPlanRepository
	invoice  *InvoiceService
}

func NewBookingService(repo *repository.BookingRepository, room *repository.RoomRepository, mealPlan *repository.MealPlanRepository) *BookingService {
	return &BookingService{repo: repo, room: room, mealPlan: mealPlan}
}

// SetInvoiceService injects the invoice service after construction to avoid a circular dependency.
func (s *BookingService) SetInvoiceService(invoice *InvoiceService) {
	s.invoice = invoice
}


func (s *BookingService) Create(userID uuid.UUID, orgID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error) {
	if req.RoomID == uuid.Nil {
		return nil, errors.New("room_id is required")
	}
	if req.ClientID == uuid.Nil {
		return nil, errors.New("client_id is required")
	}
	if req.ClientType != models.BookingClientTypeIndividual && req.ClientType != models.BookingClientTypeCorporate {
		return nil, errors.New("client_type must be 'individual' or 'corporate'")
	}
	if req.CheckIn.IsZero() || req.CheckOut.IsZero() {
		return nil, errors.New("check_in and check_out are required")
	}
	if !req.CheckOut.After(req.CheckIn) {
		return nil, errors.New("check_out must be after check_in")
	}
	if req.Guests <= 0 {
		return nil, errors.New("guests must be greater than 0")
	}

	// Validate room exists and guests fit capacity
	room, err := s.room.GetByID(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if req.Guests > room.Capacity {
		return nil, fmt.Errorf("guests (%d) exceed room capacity (%d)", req.Guests, room.Capacity)
	}

	// Check for date conflicts with existing pending/confirmed/checked_in bookings
	available, err := s.repo.IsRoomAvailable(req.RoomID, req.CheckIn, req.CheckOut, nil)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for the selected dates")
	}

	// Prevent the same client from booking the same room twice for overlapping dates
	duplicate, err := s.repo.HasActiveBookingForClient(req.ClientID, req.RoomID, req.CheckIn, req.CheckOut)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("client already has an active booking for this room on the selected dates")
	}

	// Validate meal plan exists if provided
	if req.MealPlanID != nil {
		_, err := s.mealPlan.GetByID(*req.MealPlanID)
		if err != nil {
			return nil, errors.New("meal plan not found")
		}
	}

	b := &models.Booking{
		UserID:          userID,
		RoomID:          req.RoomID,
		ClientID:        req.ClientID,
		ClientType:      req.ClientType,
		MealPlanID:      req.MealPlanID,
		CheckIn:         req.CheckIn,
		CheckOut:        req.CheckOut,
		Guests:          req.Guests,
		Status:          models.BookingStatusPending,
		SpecialRequests: req.SpecialRequests,
	}
	if err := s.repo.Create(b, orgID); err != nil {
		return nil, err
	}

	// Fetch back to get client_name and meal_plan_name resolved via JOINs
	return s.repo.GetByID(b.ID)
}

func (s *BookingService) GetByID(id uuid.UUID) (*models.Booking, error) {
	return s.repo.GetByID(id)
}

func (s *BookingService) List(orgID uuid.UUID, status, clientType string, clientID *uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	return s.repo.List(orgID, status, clientType, clientID, page, pageSize)
}

func (s *BookingService) Update(id uuid.UUID, req *models.UpdateBookingRequest) (*models.Booking, error) {
	b, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Only pending bookings can be edited
	if b.Status != models.BookingStatusPending {
		return nil, errors.New("only pending bookings can be updated")
	}

	if req.CheckIn != nil {
		b.CheckIn = *req.CheckIn
	}
	if req.CheckOut != nil {
		b.CheckOut = *req.CheckOut
	}
	if req.Guests != nil {
		b.Guests = *req.Guests
	}
	if req.SpecialRequests != nil {
		b.SpecialRequests = *req.SpecialRequests
	}
	if req.MealPlanID != nil {
		if *req.MealPlanID == uuid.Nil {
			b.MealPlanID = nil // explicitly removing the meal plan
		} else {
			_, err := s.mealPlan.GetByID(*req.MealPlanID)
			if err != nil {
				return nil, errors.New("meal plan not found")
			}
			b.MealPlanID = req.MealPlanID
		}
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

	if err := s.repo.Update(b); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *BookingService) UpdateStatus(id uuid.UUID, orgID uuid.UUID, newStatus string) (*models.Booking, error) {
	b, err := s.repo.GetByID(id)
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

	if err := s.repo.UpdateStatusTx(id, newStatus); err != nil {
		return nil, err
	}

	// Auto-generate invoice when booking is confirmed
	if newStatus == models.BookingStatusConfirmed && s.invoice != nil {
		if err := s.invoice.GenerateForBooking(id, orgID); err != nil {
			// Log but don't fail the status update — invoice can be regenerated
			fmt.Printf("warning: failed to generate invoice for booking %s: %v\n", id, err)
		}
	}

	return s.repo.GetByID(id)
}

func (s *BookingService) Delete(id uuid.UUID) error {
	b, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("booking not found")
	}
	if b.Status == models.BookingStatusCheckedIn {
		return errors.New("cannot delete a booking that is currently checked in")
	}
	return s.repo.Delete(id)
}
