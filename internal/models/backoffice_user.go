package models

import (
	"time"

	"github.com/google/uuid"
)

type BackofficeUser struct {
	ID                   uuid.UUID  `json:"id"`
	FullName             string     `json:"full_name"`
	Email                string     `json:"email"`
	Password             string     `json:"-"`
	IsActive             bool       `json:"is_active"`
	ChangePassword       bool       `json:"change_password"`
	FailedLoginAttempts  int        `json:"-"`
	IsLocked             bool       `json:"is_locked"`
	LockedUntil          *time.Time `json:"locked_until,omitempty"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type CreateBackofficeUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UpdateBackofficeUserRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
}
