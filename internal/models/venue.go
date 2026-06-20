package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	VenueTypeConferenceHall = "conference_hall"
	VenueTypeEventSpace     = "event_space"
	VenueTypeBoardroom      = "boardroom"
	VenueTypeOutdoor        = "outdoor"
	VenueTypeDining         = "dining"

	VenueRateHourly = "hourly"
	VenueRateDaily  = "daily"
)

var ValidVenueTypes = map[string]bool{
	VenueTypeConferenceHall: true,
	VenueTypeEventSpace:     true,
	VenueTypeBoardroom:      true,
	VenueTypeOutdoor:        true,
	VenueTypeDining:         true,
}

var ValidVenueRateTypes = map[string]bool{
	VenueRateHourly: true,
	VenueRateDaily:  true,
}

type Venue struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       *uuid.UUID `json:"org_id,omitempty"`
	BranchID    *uuid.UUID `json:"branch_id,omitempty"`
	Name        string     `json:"name"`
	VenueType   string     `json:"venue_type"`
	Capacity    int        `json:"capacity"`
	AreaSqm     float64    `json:"area_sqm,omitempty"`
	Floor       string     `json:"floor,omitempty"`
	BaseRate    float64    `json:"base_rate"`
	RateType    string     `json:"rate_type"`
	Amenities   []string   `json:"amenities"`
	Images      []string   `json:"images"`
	IsAvailable bool       `json:"is_available"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
