package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	IndividualBookingStatusPending   = "pending"
	IndividualBookingStatusApproved  = "approved"
	IndividualBookingStatusRejected  = "rejected"
	IndividualBookingStatusCancelled = "cancelled"
)

type IndividualBookingRequest struct {
	ID          uuid.UUID       `json:"id"`
	OrgID       uuid.UUID       `json:"org_id"`
	WebUserID   *uuid.UUID      `json:"web_user_id,omitempty"`
	BookerName  string          `json:"booker_name"`
	BookerEmail string          `json:"booker_email,omitempty"`
	BookerPhone string          `json:"booker_phone,omitempty"`
	BookingType string          `json:"booking_type"`
	Status      string          `json:"status"`
	Notes       string          `json:"notes,omitempty"`
	Documents   []string        `json:"documents,omitempty"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Joined fields
	RoomName string `json:"room_name,omitempty"`
}

// SubmitIndividualBookingRequest is the payload sent by the web user.
type SubmitIndividualBookingRequest struct {
	// Booker identity — populated from JWT in the handler, not from body
	WebUserID   *uuid.UUID `json:"web_user_id,omitempty"`
	BookerName  string     `json:"booker_name"`
	BookerEmail string     `json:"booker_email,omitempty"`
	BookerPhone string     `json:"booker_phone,omitempty"`

	// Room stay
	RoomID          uuid.UUID `json:"room_id"`
	CheckIn         DateOnly  `json:"check_in"`
	CheckOut        DateOnly  `json:"check_out"`
	SpecialRequests string    `json:"special_requests,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	Documents       []string  `json:"documents,omitempty"`
}

// IndividualBookingPayload is what gets stored in the JSONB payload column.
type IndividualBookingPayload struct {
	RoomID          uuid.UUID `json:"room_id"`
	CheckIn         string    `json:"check_in"`
	CheckOut        string    `json:"check_out"`
	SpecialRequests string    `json:"special_requests,omitempty"`
}
