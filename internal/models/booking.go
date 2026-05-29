package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DateOnly is a time.Time that marshals/unmarshals as "YYYY-MM-DD".
type DateOnly struct{ time.Time }

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	// Accept full ISO 8601 timestamps as well as plain YYYY-MM-DD
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05.999Z07:00"} {
		if t, err := time.Parse(layout, s); err == nil {
			d.Time = t.UTC().Truncate(24 * time.Hour)
			return nil
		}
	}
	return fmt.Errorf("date must be YYYY-MM-DD or ISO 8601, got %q", s)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

const (
	BookingStatusPending    = "pending"
	BookingStatusConfirmed  = "confirmed"
	BookingStatusCheckedIn  = "checked_in"
	BookingStatusCheckedOut = "checked_out"
	BookingStatusCancelled  = "cancelled"

	// Client type constants for bookings — mirrors ClientTypeIndividual / ClientTypeCorporate
	// but kept here so the booking package is self-describing.
	BookingClientTypeIndividual = ClientTypeIndividual
	BookingClientTypeCorporate  = ClientTypeCorporate
)

// ValidBookingTransitions defines the allowed next states for each booking status.
var ValidBookingTransitions = map[string][]string{
	BookingStatusPending:    {BookingStatusConfirmed, BookingStatusCancelled},
	BookingStatusConfirmed:  {BookingStatusCheckedIn, BookingStatusCancelled},
	BookingStatusCheckedIn:  {BookingStatusCheckedOut},
	BookingStatusCheckedOut: {},
	BookingStatusCancelled:  {},
}

type Booking struct {
	ID                uuid.UUID  `json:"id"`
	BookingNumber     string     `json:"booking_number"`
	RoomID            uuid.UUID  `json:"room_id"`
	RoomName          string     `json:"room_name"`
	OrgID             uuid.UUID  `json:"org_id"`
	OrgName           string     `json:"org_name,omitempty"`
	BranchID          *uuid.UUID `json:"branch_id,omitempty"`
	ClientID          uuid.UUID  `json:"client_id"`
	ClientType        string     `json:"client_type"`
	ClientName        string     `json:"client_name"`
	CorporateClientID   *uuid.UUID `json:"corporate_client_id,omitempty"`
	CorporateClientName string     `json:"corporate_client_name,omitempty"`
	CheckIn           time.Time  `json:"check_in"`
	CheckOut          time.Time  `json:"check_out"`
	Guests            int        `json:"guests"`
	Nights            int        `json:"nights"`
	RoomCost          float64    `json:"room_cost"`
	TotalAmount       float64    `json:"total_amount"`
	Status            string     `json:"status"`
	Overstayed        bool       `json:"overstayed"`
	SpecialRequests   string     `json:"special_requests,omitempty"`
	Documents         []string   `json:"documents"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ─── Individual booking request ───────────────────────────────────────────────

type NewIndividualClientDetails struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	IDPassportNumber string `json:"id_passport_number"`
}

type CreateIndividualBookingRequest struct {
	ClientID        *uuid.UUID                  `json:"client_id,omitempty"`
	Client          *NewIndividualClientDetails  `json:"client,omitempty"`
	RoomID          uuid.UUID                   `json:"room_id"`
	CheckIn         DateOnly                    `json:"check_in"`
	CheckOut        DateOnly                    `json:"check_out"`
	Guests          int                         `json:"guests"`
	SpecialRequests string                      `json:"special_requests,omitempty"`
}

// ─── Corporate booking request ────────────────────────────────────────────────

type NewCorporateClientDetails struct {
	CompanyName      string `json:"company_name"`
	ContactPerson    string `json:"contact_person"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	CompanyRegNumber string `json:"company_reg_number"`
	Industry         string `json:"industry"`
}

type CorporateGuestRequest struct {
	ClientID        *uuid.UUID `json:"client_id,omitempty"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	IDNumber        string     `json:"id_number"`
	RoomID          uuid.UUID  `json:"room_id"`
	CheckIn         DateOnly   `json:"check_in"`
	CheckOut        DateOnly   `json:"check_out"`
	SpecialRequests string     `json:"special_requests,omitempty"`
}

type CreateCorporateBookingRequest struct {
	ClientID  *uuid.UUID                 `json:"client_id,omitempty"`
	Client    *NewCorporateClientDetails `json:"client,omitempty"`
	Guests    []CorporateGuestRequest    `json:"guests"`
	Documents []string                   `json:"documents,omitempty"`
}

// ─── Guest self-service booking request ──────────────────────────────────────

type CreateBookingRequest struct {
	RoomID           uuid.UUID `json:"room_id"`
	CheckIn          DateOnly  `json:"check_in"`
	CheckOut         DateOnly  `json:"check_out"`
	Guests           int       `json:"guests"`
	SpecialRequests  string    `json:"special_requests,omitempty"`
	IDPassportNumber string    `json:"id_passport_number,omitempty"`
}

// ─── Corporate booking response ───────────────────────────────────────────────

type CorporateBookingResponse struct {
	CorporateClientID uuid.UUID `json:"corporate_client_id"`
	CompanyName       string    `json:"company_name"`
	Bookings          []Booking `json:"bookings"`
	TotalAmount       float64   `json:"total_amount"`
}

// ─── Update requests (unchanged) ─────────────────────────────────────────────

type UpdateBookingRequest struct {
	CheckIn         *DateOnly `json:"check_in,omitempty"`
	CheckOut        *DateOnly `json:"check_out,omitempty"`
	Guests          *int      `json:"guests,omitempty"`
	SpecialRequests *string   `json:"special_requests,omitempty"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status"`
}
