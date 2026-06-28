package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"
	"lodge-system/internal/utils/email"

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
	emailService   *email.EmailService
	orgRepo        *repository.OrganizationRepository
}

// SetEmailService injects the email service used for sending invoice emails.
func (s *InvoiceService) SetEmailService(emailService *email.EmailService) {
	s.emailService = emailService
}

// SetOrganizationRepository injects the org repo so invoice emails can be
// branded with the issuing lodge's name.
func (s *InvoiceService) SetOrganizationRepository(orgRepo *repository.OrganizationRepository) {
	s.orgRepo = orgRepo
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

// SendInvoiceEmail emails the invoice (as a PDF attachment) to the client's
// billing address. The PDF is rendered by the frontend and passed in as bytes.
func (s *InvoiceService) SendInvoiceEmail(id, orgID uuid.UUID, pdf []byte) error {
	if s.emailService == nil {
		return errors.New("email service is not configured")
	}

	inv, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return errors.New("invoice not found")
	}

	recipient := inv.ClientEmail
	if recipient == "" {
		return errors.New("this invoice has no billing email address")
	}

	clientName := inv.ClientName
	if clientName == "" {
		clientName = "Customer"
	}

	// Brand the email with the issuing lodge's name when available.
	orgName := ""
	if s.orgRepo != nil {
		if org, err := s.orgRepo.GetByID(orgID); err == nil && org != nil {
			orgName = org.Name
		}
	}

	issueDate := "—"
	if inv.IssuedDate != nil {
		issueDate = inv.IssuedDate.Format("02 January 2006")
	}
	dueDate := "—"
	if inv.DueDate != nil {
		dueDate = inv.DueDate.Format("02 January 2006")
	}
	totalDue := fmt.Sprintf("ZMW %.2f", inv.Total)

	// Corporate invoices include accounting references in the summary.
	var accountingRows []string
	if inv.ClientType == "corporate" {
		if inv.GLCode != "" {
			accountingRows = append(accountingRows, email.InvoiceInfoRow("GL Code:", inv.GLCode))
		}
		if inv.CostCenterType == "internal_order" && inv.InternalOrder != "" {
			accountingRows = append(accountingRows, email.InvoiceInfoRow("Internal Order:", inv.InternalOrder))
		} else if inv.CostCenter != "" {
			accountingRows = append(accountingRows, email.InvoiceInfoRow("Cost Center:", inv.CostCenter))
		}
		if inv.ClientDepartment != "" {
			accountingRows = append(accountingRows, email.InvoiceInfoRow("Department:", inv.ClientDepartment))
		}
	}

	htmlBody := email.InvoiceEmailTemplate(orgName, clientName, inv.InvoiceNumber, issueDate, dueDate, totalDue, accountingRows...)
	subjectOrg := orgName
	if subjectOrg == "" {
		subjectOrg = "Lodge Management"
	}
	subject := fmt.Sprintf("Invoice %s from %s", inv.InvoiceNumber, subjectOrg)

	attachment := email.Attachment{
		Filename:    fmt.Sprintf("Invoice-%s.pdf", inv.InvoiceNumber),
		ContentType: "application/pdf",
		Data:        pdf,
	}

	return s.emailService.SendEmailWithAttachment([]string{recipient}, subject, htmlBody, attachment)
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
