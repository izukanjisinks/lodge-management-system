package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingDocument struct {
	ID                uuid.UUID `json:"id"`
	CorporateClientID uuid.UUID `json:"corporate_client_id"`
	OrgID             uuid.UUID `json:"org_id"`
	URLs              []string  `json:"urls"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
