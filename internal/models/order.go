package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderTypeInHouse = "in_house"
	OrderTypeWalkIn  = "walk_in"
)

type Order struct {
	ID            uuid.UUID   `json:"id"`
	OrgID         uuid.UUID   `json:"org_id"`
	BookingID     *uuid.UUID  `json:"booking_id,omitempty"`
	BookingNumber string      `json:"booking_number,omitempty"`
	RoomName      string      `json:"room_name,omitempty"`
	ClientName    string      `json:"client_name,omitempty"`
	OrderNumber   string      `json:"order_number"`
	Type          string      `json:"type"`
	Notes         string      `json:"notes,omitempty"`
	Total         float64     `json:"total"`
	Items         []OrderItem `json:"items,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	MenuItemID uuid.UUID `json:"menu_item_id"`
	ItemName   string    `json:"item_name,omitempty"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	Subtotal   float64   `json:"subtotal"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// PlaceOrderRequest creates a new in-house order tied to a booking.
type PlaceOrderRequest struct {
	BookingID uuid.UUID              `json:"booking_id"`
	Notes     string                 `json:"notes,omitempty"`
	Items     []PlaceOrderItemRequest `json:"items"`
}

// PlaceWalkInOrderRequest creates a new walk-in order with no booking.
type PlaceWalkInOrderRequest struct {
	Notes string                 `json:"notes,omitempty"`
	Items []PlaceOrderItemRequest `json:"items"`
}

type PlaceOrderItemRequest struct {
	MenuItemID uuid.UUID `json:"menu_item_id"`
	Quantity   int       `json:"quantity"`
	Notes      string    `json:"notes,omitempty"`
}

// AddOrderItemsRequest appends items to an existing order.
type AddOrderItemsRequest struct {
	Items []PlaceOrderItemRequest `json:"items"`
}
