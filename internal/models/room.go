package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoomTypeSingle     = "single"
	RoomTypeDouble     = "double"
	RoomTypeSuite      = "suite"
	RoomTypeCabin      = "cabin"
	RoomTypeConference = "conference"
)

var ValidRoomTypes = map[string]bool{
	RoomTypeSingle:     true,
	RoomTypeDouble:     true,
	RoomTypeSuite:      true,
	RoomTypeCabin:      true,
	RoomTypeConference: true,
}

type BookedDate struct {
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	Status   string `json:"status"`
}

type RoomOrganization struct {
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	LogoURL string `json:"logo_url,omitempty"`
}

type Room struct {
	ID            uuid.UUID        `json:"id"`
	OrgID         *uuid.UUID       `json:"org_id,omitempty"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	Capacity      int              `json:"capacity"`
	PricePerNight float64          `json:"price_per_night"`
	Amenities     []string         `json:"amenities"`
	Images        []string         `json:"images"`
	IsAvailable   bool             `json:"is_available"`
	Description   string           `json:"description,omitempty"`
	Organization  *RoomOrganization `json:"organization,omitempty"`
	BookedDates   []BookedDate      `json:"booked_dates,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
