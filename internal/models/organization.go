package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	LogoURL       string    `json:"logo_url,omitempty"`
	StreetAddress string    `json:"street_address,omitempty"`
	City          string    `json:"city,omitempty"`
	Country       string    `json:"country,omitempty"`
	Location      string    `json:"location,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Email         string    `json:"email,omitempty"`
	Parking       bool      `json:"parking"`
	Restaurant    bool      `json:"restaurant"`
	CheckInTime   *string   `json:"check_in_time"`
	CheckOutTime  *string   `json:"check_out_time"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrgDetails struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	Location      string `json:"location"`
	LogoURL       string `json:"logo_url"`
	Parking       *bool   `json:"parking,omitempty"`
	Restaurant    *bool   `json:"restaurant,omitempty"`
	CheckInTime   *string `json:"check_in_time,omitempty"`
	CheckOutTime  *string `json:"check_out_time,omitempty"`
}

type AdminDetails struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ProvisionOrgRequest struct {
	Organization OrgDetails   `json:"organization"`
	Admin        AdminDetails `json:"admin"`
}
