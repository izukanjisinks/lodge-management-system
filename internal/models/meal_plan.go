package models

import (
	"time"

	"github.com/google/uuid"
)

type MealPlan struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	PricePerPersonPerNight float64  `json:"price_per_person_per_night"`
	Includes              []string  `json:"includes"`
	Description           string    `json:"description,omitempty"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CreateMealPlanRequest struct {
	Name                  string   `json:"name"`
	PricePerPersonPerNight float64 `json:"price_per_person_per_night"`
	Includes              []string `json:"includes"`
	Description           string   `json:"description,omitempty"`
	IsActive              bool     `json:"is_active"`
}

type UpdateMealPlanRequest struct {
	Name                  *string  `json:"name,omitempty"`
	PricePerPersonPerNight *float64 `json:"price_per_person_per_night,omitempty"`
	Includes              []string `json:"includes,omitempty"`
	Description           *string  `json:"description,omitempty"`
	IsActive              *bool    `json:"is_active,omitempty"`
}
