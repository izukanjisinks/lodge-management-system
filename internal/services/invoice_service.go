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
	repo    *repository.InvoiceRepository
	booking *repository.BookingRepository
	room    *repository.RoomRepository
}

func NewInvoiceService(
	repo *repository.InvoiceRepository,
	booking *repository.BookingRepository,
	room *repository.RoomRepository,
) *InvoiceService {
	return &InvoiceService{repo: repo, booking: booking, room: room}
}

// GenerateForBooking auto-creates an invoice when a booking is confirmed.
// It is called inside the booking status transition — not exposed as an HTTP endpoint.
func (s *InvoiceService) GenerateForBooking(bookingID uuid.UUID, orgID uuid.UUID) error {
	// Idempotent — don't create a second invoice if one already exists
	existing, _ := s.repo.GetByBookingID(bookingID, orgID)
	if existing != nil {
		return nil
	}

	b, err := s.booking.GetByID(bookingID, orgID)
	if err != nil {
		return errors.New("booking not found")
	}

	room, err := s.room.GetByID(b.RoomID, orgID)
	if err != nil {
		return errors.New("room not found")
	}

	nights := int(math.Ceil(b.CheckOut.Sub(b.CheckIn).Hours() / 24))
	if nights < 1 {
		nights = 1
	}

	roomTotal := float64(nights) * room.PricePerNight
	lineItems := []models.InvoiceLineItem{
		{
			Description: fmt.Sprintf("%s — %d night(s) @ %.2f/night", room.Name, nights, room.PricePerNight),
			Quantity:    nights,
			UnitPrice:   room.PricePerNight,
			Total:       roomTotal,
		},
	}

	subtotal := roomTotal
	taxAmount := math.Round((subtotal*defaultTaxRate/100)*100) / 100
	total := math.Round((subtotal+taxAmount)*100) / 100

	invoiceNumber, err := s.repo.GenerateInvoiceNumber()
	if err != nil {
		return err
	}

	now := time.Now()
	dueDate := b.CheckOut

	inv := &models.Invoice{
		InvoiceNumber: invoiceNumber,
		BookingID:     &bookingID,
		ClientID:      b.ClientID,
		ClientType:    b.ClientType,
		ClientName:    b.ClientName,
		BranchID:      b.BranchID,
		LineItems:     lineItems,
		Subtotal:      subtotal,
		TaxRate:       defaultTaxRate,
		TaxAmount:     taxAmount,
		Total:         total,
		Status:        models.InvoiceStatusDraft,
		IssuedDate:    &now,
		DueDate:       &dueDate,
	}

	return s.repo.Create(inv, orgID)
}

// GenerateCorporateInvoice creates one consolidated invoice for a corporate booking.
// bookingIDs contains all guest bookings that belong to this corporate transaction.
// corporateClientID is the company being billed.
func (s *InvoiceService) GenerateCorporateInvoice(corporateClientID uuid.UUID, orgID uuid.UUID, bookingIDs []uuid.UUID) error {
	// Idempotent — don't create a second invoice
	existing, _ := s.repo.GetByCorporateClientID(corporateClientID, orgID)
	if existing != nil {
		return nil
	}

	var lineItems []models.InvoiceLineItem
	var subtotal float64
	var branchID *uuid.UUID

	for _, bookingID := range bookingIDs {
		b, err := s.booking.GetByID(bookingID, orgID)
		if err != nil {
			continue
		}
		if branchID == nil {
			branchID = b.BranchID
		}
		room, err := s.room.GetByID(b.RoomID, orgID)
		if err != nil {
			continue
		}
		nights := int(math.Ceil(b.CheckOut.Sub(b.CheckIn).Hours() / 24))
		if nights < 1 {
			nights = 1
		}
		roomTotal := float64(nights) * room.PricePerNight
		subtotal += roomTotal

		bID := b.ID
		lineItems = append(lineItems, models.InvoiceLineItem{
			BookingID:   &bID,
			Description: fmt.Sprintf("%s — %s — %d night(s) @ %.2f/night", b.ClientName, room.Name, nights, room.PricePerNight),
			Quantity:    nights,
			UnitPrice:   room.PricePerNight,
			Total:       roomTotal,
		})
	}

	if len(lineItems) == 0 {
		return errors.New("no valid bookings to invoice")
	}

	taxAmount := math.Round((subtotal*defaultTaxRate/100)*100) / 100
	total := math.Round((subtotal+taxAmount)*100) / 100

	invoiceNumber, err := s.repo.GenerateInvoiceNumber()
	if err != nil {
		return err
	}

	now := time.Now()
	inv := &models.Invoice{
		InvoiceNumber:     invoiceNumber,
		CorporateClientID: &corporateClientID,
		ClientID:          corporateClientID,
		ClientType:        models.BookingClientTypeCorporate,
		BranchID:          branchID,
		LineItems:         lineItems,
		Subtotal:          subtotal,
		TaxRate:           defaultTaxRate,
		TaxAmount:         taxAmount,
		Total:             total,
		Status:            models.InvoiceStatusDraft,
		IssuedDate:        &now,
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

func (s *InvoiceService) List(orgID uuid.UUID, branchID *uuid.UUID, status string, page, pageSize int) ([]models.Invoice, int, error) {
	return s.repo.List(orgID, branchID, status, page, pageSize)
}

func (s *InvoiceService) UpdateDueDate(bookingID uuid.UUID, orgID uuid.UUID, dueDate time.Time) error {
	return s.repo.UpdateDueDate(bookingID, orgID, dueDate)
}

// RecalculateRoomCharge updates the room line item and invoice totals to reflect
// a changed check_in or check_out date.
func (s *InvoiceService) RecalculateRoomCharge(bookingID uuid.UUID, orgID uuid.UUID) error {
	b, err := s.booking.GetByID(bookingID, orgID)
	if err != nil {
		return nil // invoice may not exist yet (pending booking), skip silently
	}
	room, err := s.room.GetByID(b.RoomID, orgID)
	if err != nil {
		return nil
	}
	nights := int(math.Ceil(b.CheckOut.Sub(b.CheckIn).Hours() / 24))
	if nights < 1 {
		nights = 1
	}
	return s.repo.UpdateRoomLineItem(bookingID, orgID, nights, room.PricePerNight, room.Name)
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
