package models

import (
	"time"

	"github.com/google/uuid"
)

type WebUser struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	FullName            string     `json:"full_name"`
	Phone               string     `json:"phone,omitempty"`
	IsActive            bool       `json:"is_active"`
	ChangePassword      bool       `json:"change_password"`
	FailedLoginAttempts int        `json:"-"`
	IsLocked            bool       `json:"is_locked"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type WebUserRegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password"`
}

type WebUserUpdateRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type WebUserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
