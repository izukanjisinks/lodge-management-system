package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	AuditActorSystem = "system"
	AuditActorUser   = "user"
	AuditActorGuest  = "guest"

	AuditEntityBooking = "booking"
	AuditEntityInvoice = "invoice"
	AuditEntityOrder   = "order"
	AuditEntityRoom    = "room"

	AuditActionBookingOverstayed     = "booking.overstayed"
	AuditActionBookingStatusChanged  = "booking.status_changed"
	AuditActionBookingOverstayCleared = "booking.overstay_cleared"
)

type AuditLog struct {
	ID         uuid.UUID       `json:"id"`
	OrgID      uuid.UUID       `json:"org_id"`
	ActorType  string          `json:"actor_type"`
	ActorID    *uuid.UUID      `json:"actor_id,omitempty"`
	ActorName  string          `json:"actor_name,omitempty"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   uuid.UUID       `json:"entity_id"`
	Payload    json.RawMessage `json:"payload"`
	CreatedAt  time.Time       `json:"created_at"`
}

// OverstayedPayload is the payload written when the nightly job marks a booking as overstayed.
type OverstayedPayload struct {
	BookingNumber string `json:"booking_number"`
	RoomName      string `json:"room_name"`
	ClientName    string `json:"client_name"`
	OriginalCheckOut string `json:"original_check_out"`
	ExtendedTo    string `json:"extended_to"`
}
