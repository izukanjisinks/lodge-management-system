package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PendingBookingInput is the minimal envelope used at customer submission to create
// a parent-only booking in the 'pending' state. The full submission payload is stored
// verbatim in Metadata and later decoded at workflow approval to materialise children.
type PendingBookingInput struct {
	OrgID        uuid.UUID
	BranchID     *uuid.UUID
	WebUserID    *uuid.UUID
	CorProfileID *uuid.UUID
	CompanyID    *uuid.UUID
	BookerType   string // individual | corporate
	BookerName   string
	BookerEmail  string
	BookerPhone  string
	BookingType  string // accommodation | event | meals
	Documents    []string
	Metadata     json.RawMessage
}

// ReservationItem is the normalised shape returned by GET /api/v1/web/my-reservations.
// record_type tells the frontend which table the row came from so it can render appropriately.
type ReservationItem struct {
	ID          uuid.UUID       `json:"id"`
	RecordType  string          `json:"record_type"`  // booking | individual_request | corporate_request
	BookingType string          `json:"booking_type"` // accommodation | event | meals
	BookerType  string          `json:"booker_type"`  // individual | corporate
	Status      string          `json:"status"`
	CompanyName string          `json:"company_name,omitempty"`
	ProfileName string          `json:"profile_name,omitempty"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	TotalAmount float64         `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at"`
}
