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
