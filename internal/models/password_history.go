package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordHistory tracks previously used passwords to prevent reuse
type PasswordHistory struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	PasswordHash string    `db:"password_hash" json:"-"` // Never expose in JSON
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
