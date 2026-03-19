package models

import (
	"time"

	"github.com/google/uuid"
)

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
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	RoomID          uuid.UUID  `json:"room_id"`
	ClientID        uuid.UUID  `json:"client_id"`
	ClientType      string     `json:"client_type"`
	ClientName      string     `json:"client_name"`
	MealPlanID      *uuid.UUID `json:"meal_plan_id,omitempty"`
	MealPlanName    string     `json:"meal_plan_name,omitempty"`
	CheckIn         time.Time  `json:"check_in"`
	CheckOut        time.Time  `json:"check_out"`
	Guests          int        `json:"guests"`
	Status          string     `json:"status"`
	SpecialRequests string     `json:"special_requests,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateBookingRequest struct {
	RoomID          uuid.UUID  `json:"room_id"`
	ClientID        uuid.UUID  `json:"client_id"`
	ClientType      string     `json:"client_type"`
	MealPlanID      *uuid.UUID `json:"meal_plan_id,omitempty"`
	CheckIn         time.Time  `json:"check_in"`
	CheckOut        time.Time  `json:"check_out"`
	Guests          int        `json:"guests"`
	SpecialRequests string     `json:"special_requests,omitempty"`
}

type UpdateBookingRequest struct {
	CheckIn         *time.Time `json:"check_in,omitempty"`
	CheckOut        *time.Time `json:"check_out,omitempty"`
	Guests          *int       `json:"guests,omitempty"`
	MealPlanID      *uuid.UUID `json:"meal_plan_id,omitempty"`
	SpecialRequests *string    `json:"special_requests,omitempty"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status"`
}
