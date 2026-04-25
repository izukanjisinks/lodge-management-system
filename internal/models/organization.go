package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	LogoURL   string    `json:"logo_url,omitempty"`
	Address   string    `json:"address,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrgDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}

type AdminDetails struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ProvisionOrgRequest struct {
	Organization OrgDetails   `json:"organization"`
	Admin        AdminDetails `json:"admin"`
}
