package models

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"org_id"`
	Name          string     `json:"name"`
	BranchCode    string     `json:"branch_code"`
	StreetAddress *string    `json:"street_address"`
	City          *string    `json:"city"`
	Country       *string    `json:"country"`
	Location      *string    `json:"location"`
	Phone         *string    `json:"phone"`
	Email         *string    `json:"email"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateBranchRequest struct {
	Name          string `json:"name"`
	BranchCode    string `json:"branch_code"`
	StreetAddress string `json:"street_address,omitempty"`
	City          string `json:"city,omitempty"`
	Country       string `json:"country,omitempty"`
	Location      string `json:"location,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
}

type UpdateBranchRequest struct {
	Name          *string `json:"name,omitempty"`
	BranchCode    *string `json:"branch_code,omitempty"`
	StreetAddress *string `json:"street_address,omitempty"`
	City          *string `json:"city,omitempty"`
	Country       *string `json:"country,omitempty"`
	Location      *string `json:"location,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Email         *string `json:"email,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}
