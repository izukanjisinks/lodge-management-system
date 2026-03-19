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

type Room struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Capacity      int       `json:"capacity"`
	PricePerNight float64   `json:"price_per_night"`
	Amenities     []string  `json:"amenities"`
	IsAvailable   bool      `json:"is_available"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
