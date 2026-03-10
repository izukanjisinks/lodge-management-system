package repositories

import (
	"database/sql"
	"fmt"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type PasswordPolicyRepository struct {
	db *sql.DB
}

func NewPasswordPolicyRepository() *PasswordPolicyRepository {
	return &PasswordPolicyRepository{db: database.DB}
}

// Get retrieves the password policy
func (r *PasswordPolicyRepository) Get() (*models.PasswordPolicy, error) {
	var policy models.PasswordPolicy
	query := `
		SELECT id, min_length, require_uppercase, require_lowercase,
		       require_numbers, require_special_chars, max_failed_attempts, lockout_duration_mins,
		       password_expiry_days, otp_length, otp_expiry_mins, session_timeout_mins,
		       created_at, updated_at
		FROM password_policies
		LIMIT 1
	`
	err := r.db.QueryRow(query).Scan(
		&policy.ID,
		&policy.MinLength,
		&policy.RequireUppercase,
		&policy.RequireLowercase,
		&policy.RequireNumbers,
		&policy.RequireSpecialChars,
		&policy.MaxFailedAttempts,
		&policy.LockoutDurationMins,
		&policy.PasswordExpiryDays,
		&policy.OTPLength,
		&policy.OTPExpiryMins,
		&policy.SessionTimeoutMins,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no password policy found")
		}
		return nil, err
	}
	return &policy, nil
}

// Create creates a new password policy
func (r *PasswordPolicyRepository) Create(policy *models.PasswordPolicy) error {
	query := `
		INSERT INTO password_policies (
			id, min_length, require_uppercase, require_lowercase,
			require_numbers, require_special_chars, max_failed_attempts, lockout_duration_mins,
			password_expiry_days, otp_length, otp_expiry_mins, session_timeout_mins,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`
	_, err := r.db.Exec(query,
		policy.ID,
		policy.MinLength,
		policy.RequireUppercase,
		policy.RequireLowercase,
		policy.RequireNumbers,
		policy.RequireSpecialChars,
		policy.MaxFailedAttempts,
		policy.LockoutDurationMins,
		policy.PasswordExpiryDays,
		policy.OTPLength,
		policy.OTPExpiryMins,
		policy.SessionTimeoutMins,
		policy.CreatedAt,
		policy.UpdatedAt,
	)
	return err
}

// Update updates an existing password policy
func (r *PasswordPolicyRepository) Update(policy *models.PasswordPolicy) error {
	query := `
		UPDATE password_policies
		SET min_length = $2,
		    require_uppercase = $3,
		    require_lowercase = $4,
		    require_numbers = $5,
		    require_special_chars = $6,
		    max_failed_attempts = $7,
		    lockout_duration_mins = $8,
		    password_expiry_days = $9,
		    otp_length = $10,
		    otp_expiry_mins = $11,
		    session_timeout_mins = $12,
		    updated_at = $13
		WHERE id = $1
	`
	_, err := r.db.Exec(query,
		policy.ID,
		policy.MinLength,
		policy.RequireUppercase,
		policy.RequireLowercase,
		policy.RequireNumbers,
		policy.RequireSpecialChars,
		policy.MaxFailedAttempts,
		policy.LockoutDurationMins,
		policy.PasswordExpiryDays,
		policy.OTPLength,
		policy.OTPExpiryMins,
		policy.SessionTimeoutMins,
		policy.UpdatedAt,
	)
	return err
}

// Upsert creates or updates the password policy
func (r *PasswordPolicyRepository) Upsert(policy *models.PasswordPolicy) error {
	// Check if policy exists
	existing, err := r.Get()

	if err != nil {
		// Doesn't exist, create new
		policy.ID = uuid.New()
		return r.Create(policy)
	}

	// Exists, update
	policy.ID = existing.ID
	return r.Update(policy)
}
