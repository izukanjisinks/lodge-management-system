package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InvoiceStatusDraft     = "draft"
	InvoiceStatusIssued    = "issued"
	InvoiceStatusPaid      = "paid"
	InvoiceStatusOverdue   = "overdue"
	InvoiceStatusCancelled = "cancelled"
)

var ValidInvoiceTransitions = map[string][]string{
	InvoiceStatusDraft:     {InvoiceStatusIssued, InvoiceStatusCancelled},
	InvoiceStatusIssued:    {InvoiceStatusPaid, InvoiceStatusOverdue, InvoiceStatusCancelled},
	InvoiceStatusOverdue:   {InvoiceStatusPaid, InvoiceStatusCancelled},
	InvoiceStatusPaid:      {},
	InvoiceStatusCancelled: {},
}

type InvoiceLineItem struct {
	ID          uuid.UUID  `json:"id"`
	InvoiceID   uuid.UUID  `json:"invoice_id"`
	OrderID     *uuid.UUID `json:"order_id,omitempty"`
	Description string     `json:"description"`
	Quantity    int        `json:"quantity"`
	UnitPrice   float64    `json:"unit_price"`
	Total       float64    `json:"total"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Invoice struct {
	ID            uuid.UUID         `json:"id"`
	InvoiceNumber string            `json:"invoice_number"`
	BookingID     uuid.UUID         `json:"booking_id"`
	ClientID      uuid.UUID         `json:"client_id"`
	ClientName    string            `json:"client_name"`
	ClientType    string            `json:"client_type"`
	ClientEmail   string            `json:"client_email,omitempty"`
	LineItems     []InvoiceLineItem `json:"line_items"`
	Subtotal      float64           `json:"subtotal"`
	TaxRate       float64           `json:"tax_rate"`
	TaxAmount     float64           `json:"tax_amount"`
	Total         float64           `json:"total_amount"`
	Status        string            `json:"status"`
	IssuedDate    *time.Time        `json:"issued_date,omitempty"`
	DueDate       *time.Time        `json:"due_date,omitempty"`
	PaidDate      *time.Time        `json:"paid_date,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type UpdateInvoiceStatusRequest struct {
	Status   string     `json:"status"`
	PaidDate *time.Time `json:"paid_date,omitempty"`
	Notes    *string    `json:"notes,omitempty"`
}
