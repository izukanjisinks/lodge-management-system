package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ClientTypeIndividual = "individual"
	ClientTypeCorporate  = "corporate"
	ClientStatusActive   = "active"
	ClientStatusInactive = "inactive"
)

type IndividualClient struct {
	ID                uuid.UUID `json:"id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	IDPassportNumber  string    `json:"id_passport_number"`
	Nationality       string    `json:"nationality,omitempty"`
	Status            string    `json:"status"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CorporateClient struct {
	ID                uuid.UUID `json:"id"`
	CompanyName       string    `json:"company_name"`
	ContactPerson     string    `json:"contact_person"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	CompanyRegNumber  string    `json:"company_reg_number"`
	Industry          string    `json:"industry,omitempty"`
	Status            string    `json:"status"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CorporateBookingGuest struct {
	BookingID     string `json:"booking_id"`
	BookingNumber string `json:"booking_number"`
	ClientName    string `json:"client_name"`
	RoomName      string `json:"room_name"`
	CheckIn       string `json:"check_in"`
	CheckOut      string `json:"check_out"`
	Guests        int    `json:"guests"`
	Status        string `json:"status"`
}

type CorporateClientWithBookings struct {
	*CorporateClient
	Documents []string                `json:"documents"`
	Guests    []CorporateBookingGuest `json:"guests"`
}
