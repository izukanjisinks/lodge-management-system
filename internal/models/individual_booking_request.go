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

// SubmitIndividualBookingRequest is the unified envelope from the frontend
// (accommodationBooking.js). It carries the booker, participant mode, attendants,
// and the accommodation block with pre-selected rooms.
type SubmitIndividualBookingRequest struct {
	OrgID    uuid.UUID  `json:"org_id"`
	BranchID *uuid.UUID `json:"branch_id,omitempty"`

	BookingType    string `json:"booking_type"`    // "accommodation"
	Source         string `json:"source"`
	Currency       string `json:"currency"`
	BookingContext string `json:"booking_context"` // "individual"

	ParticipantMode  string `json:"participant_mode"`  // "headcount" | "detailed"
	ParticipantCount *int   `json:"participant_count,omitempty"`

	BookedBy  IndivBookedBy  `json:"booked_by"`
	Attendants []IndivAttendant `json:"attendants,omitempty"`

	Accommodation *IndivAccommodation `json:"accommodation,omitempty"`

	Documents []string `json:"documents,omitempty"`
}

type IndivBookedBy struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	JobTitle string `json:"job_title,omitempty"`
}

type IndivAttendant struct {
	FullName      string `json:"full_name"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	IDNumber      string `json:"id_number,omitempty"`
	DietaryNotes  string `json:"dietary_notes,omitempty"`
	Company       string `json:"company,omitempty"`
	IsLeadContact bool   `json:"is_lead_contact"`
}

type IndivAccommodation struct {
	CheckIn  string          `json:"check_in"`
	CheckOut string          `json:"check_out"`
	Notes    string          `json:"notes,omitempty"`
	Rooms    []IndivRoomSlot `json:"rooms,omitempty"`
}

type IndivRoomSlot struct {
	SlotIndex    int       `json:"slot_index"`
	AttendantIdx int       `json:"attendant_idx"`
	RoomID       uuid.UUID `json:"room_id"`
	RoomName     string    `json:"room_name,omitempty"`
	RoomType     string    `json:"room_type,omitempty"`
	RatePerNight float64   `json:"rate_per_night,omitempty"`
}

// IndividualBookingPayload is what gets stored in the JSONB payload column.
// Kept for backward compat with existing requests stored in the old shape.
type IndividualBookingPayload struct {
	RoomID          uuid.UUID `json:"room_id"`
	CheckIn         string    `json:"check_in"`
	CheckOut        string    `json:"check_out"`
	SpecialRequests string    `json:"special_requests,omitempty"`
}
