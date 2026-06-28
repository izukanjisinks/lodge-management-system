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
	eventRepo      *repository.BookingEventRepository
	orderRepo      *repository.OrderRepository
}

func NewInvoiceService(
	repo *repository.InvoiceRepository,
	booking *repository.BookingRepository,
	room *repository.RoomRepository,
	assignmentRepo *repository.BookingRoomAssignmentRepository,
	eventRepo *repository.BookingEventRepository,
	orderRepo *repository.OrderRepository,
) *InvoiceService {
	return &InvoiceService{repo: repo, booking: booking, room: room, assignmentRepo: assignmentRepo, eventRepo: eventRepo, orderRepo: orderRepo}
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

	var lineItems []models.InvoiceLineItem
	subtotal := 0.0
	var latestCheckOut time.Time

	switch b.BookingType {
	case models.BookingTypeEvent:
		lineItems, subtotal, latestCheckOut, err = s.eventLineItems(bookingID)
	case models.BookingTypeMeals:
		lineItems, subtotal, latestCheckOut, err = s.mealLineItems(bookingID, orgID)
	default:
		lineItems, subtotal, latestCheckOut, err = s.roomLineItems(bookingID)
	}
	if err != nil {
		return err
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
		Metadata:      b.Metadata,
	}

	return s.repo.Create(inv, orgID)
}

// roomLineItems builds invoice lines from a booking's room assignments (room stays).
func (s *InvoiceService) roomLineItems(bookingID uuid.UUID) ([]models.InvoiceLineItem, float64, time.Time, error) {
	assignments, err := s.assignmentRepo.GetAssignmentsForInvoice(bookingID)
	if err != nil || len(assignments) == 0 {
		return nil, 0, time.Time{}, errors.New("no room assignments found for booking")
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
	return lineItems, subtotal, latestCheckOut, nil
}

// eventLineItems builds invoice lines from a booking's venue reservations
// (conference/event). Each event charges price × days.
func (s *InvoiceService) eventLineItems(bookingID uuid.UUID) ([]models.InvoiceLineItem, float64, time.Time, error) {
	events, err := s.eventRepo.ListByBookingID(bookingID)
	if err != nil || len(events) == 0 {
		return nil, 0, time.Time{}, errors.New("no venue reservation found for booking")
	}

	var lineItems []models.InvoiceLineItem
	subtotal := 0.0
	var latestCheckOut time.Time

	for _, e := range events {
		days := e.Days
		if days < 1 {
			days = 1
		}
		eventTotal := e.Price * float64(days)
		subtotal += eventTotal
		if e.EndDate.After(latestCheckOut) {
			latestCheckOut = e.EndDate
		}
		venueLabel := e.VenueName
		if venueLabel == "" {
			venueLabel = "Venue"
		}
		bID := bookingID
		lineItems = append(lineItems, models.InvoiceLineItem{
			BookingID:   &bID,
			Description: fmt.Sprintf("%s — %s (%d day(s) @ %.2f/day)", venueLabel, e.EventType, days, e.Price),
			Quantity:    days,
			UnitPrice:   e.Price,
			Total:       eventTotal,
		})
	}
	return lineItems, subtotal, latestCheckOut, nil
}

// mealLineItems builds invoice lines from a meals booking's orders. Every order
// item (per-guest selection or buffet) becomes one line at its snapshotted price.
func (s *InvoiceService) mealLineItems(bookingID, orgID uuid.UUID) ([]models.InvoiceLineItem, float64, time.Time, error) {
	if s.orderRepo == nil {
		return nil, 0, time.Time{}, errors.New("orders are not configured; cannot invoice meals booking")
	}
	items, attendeeNames, err := s.orderRepo.ListItemsByBookingID(bookingID, orgID)
	if err != nil {
		return nil, 0, time.Time{}, err
	}
	if len(items) == 0 {
		return nil, 0, time.Time{}, errors.New("no meal orders found for booking")
	}

	var lineItems []models.InvoiceLineItem
	subtotal := 0.0
	for idx, it := range items {
		subtotal += it.Subtotal
		desc := fmt.Sprintf("%s — %d @ %.2f", it.ItemName, it.Quantity, it.UnitPrice)
		if who := attendeeNames[idx]; who != "" {
			desc = fmt.Sprintf("%s (%s)", desc, who)
		}
		bID := bookingID
		lineItems = append(lineItems, models.InvoiceLineItem{
			BookingID:   &bID,
			Description: desc,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Total:       it.Subtotal,
		})
	}
	// Meals have no checkout date; due date defaults to today (caller's now()).
	return lineItems, subtotal, time.Now(), nil
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
