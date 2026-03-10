package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordPolicy defines the password requirements and security settings
type PasswordPolicy struct {
	ID                  uuid.UUID `db:"id" json:"id"`
	MinLength           int       `db:"min_length" json:"min_length"`
	RequireUppercase    bool      `db:"require_uppercase" json:"require_uppercase"`
	RequireLowercase    bool      `db:"require_lowercase" json:"require_lowercase"`
	RequireNumbers      bool      `db:"require_numbers" json:"require_numbers"`
	RequireSpecialChars bool      `db:"require_special_chars" json:"require_special_chars"`
	MaxFailedAttempts   int       `db:"max_failed_attempts" json:"max_failed_attempts"`
	LockoutDurationMins int       `db:"lockout_duration_mins" json:"lockout_duration_mins"`
	PasswordExpiryDays  *int      `db:"password_expiry_days" json:"password_expiry_days,omitempty"` // NULL = never expires
	OTPLength           int       `db:"otp_length" json:"otp_length"`
	OTPExpiryMins       int       `db:"otp_expiry_mins" json:"otp_expiry_mins"`
	SessionTimeoutMins  int       `db:"session_timeout_mins" json:"session_timeout_mins"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// DefaultPasswordPolicy returns the default password policy settings
func DefaultPasswordPolicy() *PasswordPolicy {
	expiryDays := 90
	return &PasswordPolicy{
		ID:                  uuid.New(),
		MinLength:           8,
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireNumbers:      true,
		RequireSpecialChars: true,
		MaxFailedAttempts:   5,
		LockoutDurationMins: 30,
		PasswordExpiryDays:  &expiryDays,
		OTPLength:           6,
		OTPExpiryMins:       5,
		SessionTimeoutMins:  30,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// CreatePasswordPolicyRequest is used for creating or updating a password policy
type CreatePasswordPolicyRequest struct {
	MinLength           *int  `json:"min_length,omitempty"`
	RequireUppercase    *bool `json:"require_uppercase,omitempty"`
	RequireLowercase    *bool `json:"require_lowercase,omitempty"`
	RequireNumbers      *bool `json:"require_numbers,omitempty"`
	RequireSpecialChars *bool `json:"require_special_chars,omitempty"`
	MaxFailedAttempts   *int  `json:"max_failed_attempts,omitempty"`
	LockoutDurationMins *int  `json:"lockout_duration_mins,omitempty"`
	PasswordExpiryDays  *int  `json:"password_expiry_days,omitempty"`
	OTPLength           *int  `json:"otp_length,omitempty"`
	OTPExpiryMins       *int  `json:"otp_expiry_mins,omitempty"`
	SessionTimeoutMins  *int  `json:"session_timeout_mins,omitempty"`
}

// PasswordPolicyResponse includes additional computed fields for API responses
type PasswordPolicyResponse struct {
	PasswordPolicy
}
