package models

import (
	"time"

	"github.com/google/uuid"
)

type Holiday struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	IsRecurring bool      `json:"is_recurring"`
	Location    string    `json:"location"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
