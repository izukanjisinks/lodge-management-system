package repositories

import (
	"database/sql"
	"fmt"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type PasswordPolicyRepository struct {
	db *sql.DB
}

func NewPasswordPolicyRepository() *PasswordPolicyRepository {
	return &PasswordPolicyRepository{db: database.DB}
}

const policySelectCols = `
	id, organization_id, min_length, require_uppercase, require_lowercase,
	require_numbers, require_special_chars, max_failed_attempts, lockout_duration_mins,
	password_expiry_days, otp_length, otp_expiry_mins, session_timeout_mins,
	created_at, updated_at`

func scanPolicy(row *sql.Row) (*models.PasswordPolicy, error) {
	var p models.PasswordPolicy
	var orgID uuid.NullUUID
	err := row.Scan(
		&p.ID, &orgID,
		&p.MinLength, &p.RequireUppercase, &p.RequireLowercase,
		&p.RequireNumbers, &p.RequireSpecialChars, &p.MaxFailedAttempts, &p.LockoutDurationMins,
		&p.PasswordExpiryDays, &p.OTPLength, &p.OTPExpiryMins, &p.SessionTimeoutMins,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if orgID.Valid {
		p.OrgID = &orgID.UUID
	}
	return &p, nil
}

// GetByOrg returns the org-specific policy if it exists, otherwise falls back to the global default.
func (r *PasswordPolicyRepository) GetByOrg(orgID uuid.UUID) (*models.PasswordPolicy, error) {
	row := r.db.QueryRow(
		`SELECT`+policySelectCols+`FROM password_policies WHERE organization_id = $1`, orgID)
	p, err := scanPolicy(row)
	if err == nil {
		return p, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	// Fall back to global default
	return r.GetGlobal()
}

// GetGlobal returns the global (NULL org) policy.
func (r *PasswordPolicyRepository) GetGlobal() (*models.PasswordPolicy, error) {
	row := r.db.QueryRow(
		`SELECT` + policySelectCols + `FROM password_policies WHERE organization_id IS NULL LIMIT 1`)
	p, err := scanPolicy(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no password policy found")
		}
		return nil, err
	}
	return p, nil
}

// UpsertForOrg inserts or updates the org-specific policy row.
func (r *PasswordPolicyRepository) UpsertForOrg(orgID uuid.UUID, policy *models.PasswordPolicy) error {
	policy.OrgID = &orgID
	_, err := r.db.Exec(`
		INSERT INTO password_policies (
			id, organization_id, min_length, require_uppercase, require_lowercase,
			require_numbers, require_special_chars, max_failed_attempts, lockout_duration_mins,
			password_expiry_days, otp_length, otp_expiry_mins, session_timeout_mins,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		ON CONFLICT (organization_id) WHERE organization_id IS NOT NULL
		DO UPDATE SET
			min_length           = EXCLUDED.min_length,
			require_uppercase    = EXCLUDED.require_uppercase,
			require_lowercase    = EXCLUDED.require_lowercase,
			require_numbers      = EXCLUDED.require_numbers,
			require_special_chars= EXCLUDED.require_special_chars,
			max_failed_attempts  = EXCLUDED.max_failed_attempts,
			lockout_duration_mins= EXCLUDED.lockout_duration_mins,
			password_expiry_days = EXCLUDED.password_expiry_days,
			otp_length           = EXCLUDED.otp_length,
			otp_expiry_mins      = EXCLUDED.otp_expiry_mins,
			session_timeout_mins = EXCLUDED.session_timeout_mins,
			updated_at           = EXCLUDED.updated_at`,
		policy.ID, orgID,
		policy.MinLength, policy.RequireUppercase, policy.RequireLowercase,
		policy.RequireNumbers, policy.RequireSpecialChars, policy.MaxFailedAttempts, policy.LockoutDurationMins,
		policy.PasswordExpiryDays, policy.OTPLength, policy.OTPExpiryMins, policy.SessionTimeoutMins,
		policy.CreatedAt, policy.UpdatedAt,
	)
	return err
}

// Get — kept for backward compatibility (service startup loads global default).
func (r *PasswordPolicyRepository) Get() (*models.PasswordPolicy, error) {
	return r.GetGlobal()
}
