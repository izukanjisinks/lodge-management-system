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
	repo     *repository.InvoiceRepository
	booking  *repository.BookingRepository
	room     *repository.RoomRepository
	mealPlan *repository.MealPlanRepository
}

func NewInvoiceService(
	repo *repository.InvoiceRepository,
	booking *repository.BookingRepository,
	room *repository.RoomRepository,
	mealPlan *repository.MealPlanRepository,
) *InvoiceService {
	return &InvoiceService{repo: repo, booking: booking, room: room, mealPlan: mealPlan}
}

// GenerateForBooking auto-creates an invoice when a booking is confirmed.
// It is called inside the booking status transition — not exposed as an HTTP endpoint.
func (s *InvoiceService) GenerateForBooking(bookingID uuid.UUID, orgID uuid.UUID) error {
	// Idempotent — don't create a second invoice if one already exists
	existing, _ := s.repo.GetByBookingID(bookingID)
	if existing != nil {
		return nil
	}

	b, err := s.booking.GetByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	room, err := s.room.GetByID(b.RoomID)
	if err != nil {
		return errors.New("room not found")
	}

	nights := int(math.Ceil(b.CheckOut.Sub(b.CheckIn).Hours() / 24))
	if nights < 1 {
		nights = 1
	}

	var lineItems []models.InvoiceLineItem

	// Line item: room cost
	roomTotal := float64(nights) * room.PricePerNight
	lineItems = append(lineItems, models.InvoiceLineItem{
		Description: fmt.Sprintf("%s — %d night(s) @ %.2f/night", room.Name, nights, room.PricePerNight),
		Quantity:    nights,
		UnitPrice:   room.PricePerNight,
		Total:       roomTotal,
	})

	// Line item: meal plan cost (if attached)
	if b.MealPlanID != nil {
		mp, err := s.mealPlan.GetByID(*b.MealPlanID)
		if err == nil {
			mealTotal := float64(nights) * float64(b.Guests) * mp.PricePerPersonPerNight
			lineItems = append(lineItems, models.InvoiceLineItem{
				Description: fmt.Sprintf("%s — %d night(s) × %d guest(s) @ %.2f/person/night", mp.Name, nights, b.Guests, mp.PricePerPersonPerNight),
				Quantity:    nights * b.Guests,
				UnitPrice:   mp.PricePerPersonPerNight,
				Total:       mealTotal,
			})
		}
	}

	// Calculate totals
	subtotal := 0.0
	for _, item := range lineItems {
		subtotal += item.Total
	}
	taxAmount := math.Round((subtotal*defaultTaxRate/100)*100) / 100
	total := math.Round((subtotal+taxAmount)*100) / 100

	invoiceNumber, err := s.repo.GenerateInvoiceNumber()
	if err != nil {
		return err
	}

	now := time.Now()
	dueDate := now.AddDate(0, 0, 30) // due in 30 days

	inv := &models.Invoice{
		InvoiceNumber: invoiceNumber,
		BookingID:     bookingID,
		ClientID:      b.ClientID,
		ClientType:    b.ClientType,
		ClientName:    b.ClientName,
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

func (s *InvoiceService) GetByID(id uuid.UUID) (*models.Invoice, error) {
	inv, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}
	return inv, nil
}

func (s *InvoiceService) GetByBookingID(bookingID uuid.UUID) (*models.Invoice, error) {
	inv, err := s.repo.GetByBookingID(bookingID)
	if err != nil {
		return nil, errors.New("invoice not found for this booking")
	}
	return inv, nil
}

func (s *InvoiceService) List(orgID uuid.UUID, status string, page, pageSize int) ([]models.Invoice, int, error) {
	return s.repo.List(orgID, status, page, pageSize)
}

func (s *InvoiceService) UpdateStatus(id uuid.UUID, req *models.UpdateInvoiceStatusRequest) (*models.Invoice, error) {
	inv, err := s.repo.GetByID(id)
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

	if err := s.repo.UpdateStatus(id, req.Status, req.PaidDate, req.Notes); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
