package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

const defaultTaxRate = 16.0 // 16% VAT — adjust as needed

type InvoiceService struct {
	repo           *repository.InvoiceRepository
	booking        *repository.BookingRepository
	room           *repository.RoomRepository
	assignmentRepo *repository.BookingRoomAssignmentRepository
}

func NewInvoiceService(
	repo *repository.InvoiceRepository,
	booking *repository.BookingRepository,
	room *repository.RoomRepository,
	assignmentRepo *repository.BookingRoomAssignmentRepository,
) *InvoiceService {
	return &InvoiceService{repo: repo, booking: booking, room: room, assignmentRepo: assignmentRepo}
}

// GenerateForBooking auto-creates an invoice when a booking is confirmed.
// It derives line items from booking_room_assignments.
func (s *InvoiceService) GenerateForBooking(bookingID uuid.UUID, orgID uuid.UUID) error {
	existing, _ := s.repo.GetByBookingID(bookingID, orgID)
	if existing != nil {
		return nil
	}

	b, err := s.booking.GetByID(bookingID, orgID)
	if err != nil {
		return errors.New("booking not found")
	}

	// Build line items from room assignments
	assignments, err := s.assignmentRepo.GetAssignmentsForInvoice(bookingID)
	if err != nil || len(assignments) == 0 {
		return errors.New("no room assignments found for booking")
	}

	var lineItems []models.InvoiceLineItem
	subtotal := 0.0
	var latestCheckOut time.Time

	for _, a := range assignments {
		nights := int(math.Ceil(a.CheckOut.Sub(a.CheckIn).Hours() / 24))
		if nights < 1 {
			nights = 1
		}
		roomTotal := float64(nights) * a.PricePerNight
		subtotal += roomTotal
		if a.CheckOut.After(latestCheckOut) {
			latestCheckOut = a.CheckOut
		}
		bID := bookingID
		lineItems = append(lineItems, models.InvoiceLineItem{
			BookingID:   &bID,
			Description: fmt.Sprintf("%s (%s) — %d night(s) @ %.2f/night", a.RoomName, a.AttendeeName, nights, a.PricePerNight),
			Quantity:    nights,
			UnitPrice:   a.PricePerNight,
			Total:       roomTotal,
		})
	}

	taxAmount := math.Round((subtotal*defaultTaxRate/100)*100) / 100
	total := math.Round((subtotal+taxAmount)*100) / 100

	invoiceNumber, err := s.repo.GenerateInvoiceNumber()
	if err != nil {
		return err
	}

	now := time.Now()
	inv := &models.Invoice{
		InvoiceNumber: invoiceNumber,
		BookingID:     &bookingID,
		ClientType:    b.BookerType,
		ClientName:    b.BookerName,
		ClientEmail:   b.BookerEmail,
		BranchID:      b.BranchID,
		LineItems:     lineItems,
		Subtotal:      subtotal,
		TaxRate:       defaultTaxRate,
		TaxAmount:     taxAmount,
		Total:         total,
		Status:        models.InvoiceStatusDraft,
		IssuedDate:    &now,
		DueDate:       &latestCheckOut,
	}

	return s.repo.Create(inv, orgID)
}

func (s *InvoiceService) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Invoice, error) {
	inv, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}
	return inv, nil
}

func (s *InvoiceService) GetByBookingID(bookingID uuid.UUID, orgID uuid.UUID) (*models.Invoice, error) {
	inv, err := s.repo.GetByBookingID(bookingID, orgID)
	if err != nil {
		return nil, errors.New("invoice not found for this booking")
	}
	return inv, nil
}

func (s *InvoiceService) List(orgID uuid.UUID, branchID *uuid.UUID, status, clientType string, page, pageSize int) ([]models.Invoice, int, error) {
	return s.repo.List(orgID, branchID, status, clientType, page, pageSize)
}

func (s *InvoiceService) UpdateDueDate(bookingID uuid.UUID, orgID uuid.UUID, dueDate time.Time) error {
	return s.repo.UpdateDueDate(bookingID, orgID, dueDate)
}

// RecalculateRoomCharge regenerates the invoice for a booking based on current assignments.
func (s *InvoiceService) RecalculateRoomCharge(bookingID uuid.UUID, orgID uuid.UUID) error {
	return s.GenerateForBooking(bookingID, orgID)
}

func (s *InvoiceService) UpdateStatus(id uuid.UUID, orgID uuid.UUID, req *models.UpdateInvoiceStatusRequest) (*models.Invoice, error) {
	inv, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	allowed := models.ValidInvoiceTransitions[inv.Status]
	valid := false
	for _, a := range allowed {
		if a == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("cannot transition invoice from '%s' to '%s'", inv.Status, req.Status)
	}

	if err := s.repo.UpdateStatus(id, orgID, req.Status, req.PaidDate, req.Notes); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id, orgID)
}
