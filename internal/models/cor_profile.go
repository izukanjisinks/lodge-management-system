package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CorProfile struct {
	ID         uuid.UUID  `json:"id"`
	OrgID      uuid.UUID  `json:"org_id"`
	CompanyID  uuid.UUID  `json:"company_id"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	JobTitle   string     `json:"job_title,omitempty"`
	Department string     `json:"department,omitempty"`
	Status     string     `json:"status"`
	MetaData   json.RawMessage   `json:"meta_data,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CorporateGuest struct {
	ID                 uuid.UUID `json:"id"`
	CorporateProfileID uuid.UUID `json:"corporate_profile_id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Phone              string    `json:"phone,omitempty"`
	Email              string    `json:"email,omitempty"`
	IdentificationCard string    `json:"identification_card"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateCorProfileRequest struct {
	CompanyID  uuid.UUID  `json:"company_id"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	JobTitle   string     `json:"job_title,omitempty"`
	Department string     `json:"department,omitempty"`
}

type UpdateCorProfileRequest struct {
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	FirstName  *string    `json:"first_name,omitempty"`
	LastName   *string    `json:"last_name,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	JobTitle   *string    `json:"job_title,omitempty"`
	Department *string    `json:"department,omitempty"`
	Status     *string    `json:"status,omitempty"`
}

type CreateCorporateGuestRequest struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Phone              string `json:"phone,omitempty"`
	Email              string `json:"email,omitempty"`
	IdentificationCard string `json:"identification_card"`
}

type UpdateCorporateGuestRequest struct {
	FirstName          *string `json:"first_name,omitempty"`
	LastName           *string `json:"last_name,omitempty"`
	Phone              *string `json:"phone,omitempty"`
	Email              *string `json:"email,omitempty"`
	IdentificationCard *string `json:"identification_card,omitempty"`
}
