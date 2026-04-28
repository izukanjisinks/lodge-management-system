package services

import (
	"errors"
	"fmt"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	repo    *repository.OrderRepository
	invoice *repository.InvoiceRepository
	booking *repository.BookingRepository
}

func NewOrderService(
	repo *repository.OrderRepository,
	invoice *repository.InvoiceRepository,
	booking *repository.BookingRepository,
) *OrderService {
	return &OrderService{repo: repo, invoice: invoice, booking: booking}
}

// PlaceOrder creates a new in-house order tied to a confirmed/checked-in booking
// and immediately appends a line item to the booking's invoice.
func (s *OrderService) PlaceOrder(orgID uuid.UUID, req *models.PlaceOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	b, err := s.booking.GetByID(req.BookingID, orgID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if b.Status != models.BookingStatusConfirmed && b.Status != models.BookingStatusCheckedIn {
		return nil, fmt.Errorf("orders can only be placed on confirmed or checked-in bookings (current status: %s)", b.Status)
	}

	o := &models.Order{
		BookingID: &req.BookingID,
		Type:      models.OrderTypeInHouse,
		Notes:     req.Notes,
	}
	order, err := s.repo.Create(o, req.Items, orgID)
	if err != nil {
		return nil, err
	}

	s.appendToInvoice(order, orgID)
	return order, nil
}

// PlaceWalkInOrder creates a new walk-in order with no booking — no invoice entry.
func (s *OrderService) PlaceWalkInOrder(orgID uuid.UUID, req *models.PlaceWalkInOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}
	o := &models.Order{
		Type:  models.OrderTypeWalkIn,
		Notes: req.Notes,
	}
	return s.repo.Create(o, req.Items, orgID)
}

// AddItems appends more items to an existing order and updates the invoice immediately.
func (s *OrderService) AddItems(orderID uuid.UUID, orgID uuid.UUID, req *models.AddOrderItemsRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}
	order, err := s.repo.AddItems(orderID, orgID, req.Items)
	if err != nil {
		return nil, err
	}

	s.appendToInvoice(order, orgID)
	return order, nil
}

func (s *OrderService) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Order, error) {
	o, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (s *OrderService) List(orgID uuid.UUID, orderType string, bookingID *uuid.UUID, page, pageSize int) ([]models.Order, int, error) {
	return s.repo.List(orgID, orderType, bookingID, page, pageSize)
}

// appendToInvoice writes a line item for this order onto the booking's invoice.
// Non-fatal — logs on failure so the order itself is never rolled back.
func (s *OrderService) appendToInvoice(order *models.Order, orgID uuid.UUID) {
	if order.BookingID == nil {
		return
	}

	total, err := s.repo.GetItemsTotal(order.ID)
	if err != nil || total == 0 {
		return
	}

	inv, err := s.invoice.GetByBookingID(*order.BookingID, orgID)
	if err != nil {
		// Invoice doesn't exist yet (booking not yet confirmed) — skip silently
		return
	}

	description := fmt.Sprintf("Food & Beverage — Order %s (%d item(s))", order.OrderNumber, len(order.Items))
	orderID := order.ID
	lineItem := &models.InvoiceLineItem{
		OrderID:     &orderID,
		Description: description,
		Quantity:    1,
		UnitPrice:   total,
		Total:       total,
	}
	if err := s.invoice.AppendOrderLineItem(inv.ID, orgID, lineItem); err != nil {
		fmt.Printf("warning: failed to append order %s to invoice %s: %v\n", order.ID, inv.ID, err)
	}
}
