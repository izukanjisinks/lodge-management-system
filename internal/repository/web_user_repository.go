package repository

import (
	"database/sql"
	"errors"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type WebUserRepository struct {
	db *sql.DB
}

func NewWebUserRepository() *WebUserRepository {
	return &WebUserRepository{db: database.DB}
}

func (r *WebUserRepository) Create(u *models.WebUser, password string) error {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return r.db.QueryRow(`
		INSERT INTO web_users (id, email, password, full_name, phone, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,TRUE,$6,$7)
		RETURNING id`,
		u.ID, u.Email, password, u.FullName, u.Phone, now, now,
	).Scan(&u.ID)
}

func (r *WebUserRepository) GetByID(id uuid.UUID) (*models.WebUser, error) {
	u := &models.WebUser{}
	err := r.db.QueryRow(`
		SELECT id, email, full_name, phone, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM web_users WHERE id = $1`, id,
	).Scan(
		&u.ID, &u.Email, &u.FullName, &u.Phone, &u.IsActive, &u.ChangePassword,
		&u.FailedLoginAttempts, &u.IsLocked, &u.LockedUntil, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return u, err
}

func (r *WebUserRepository) GetByEmail(email string) (*models.WebUser, string, error) {
	u := &models.WebUser{}
	var password string
	err := r.db.QueryRow(`
		SELECT id, email, password, full_name, phone, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM web_users WHERE email = $1`, email,
	).Scan(
		&u.ID, &u.Email, &password, &u.FullName, &u.Phone, &u.IsActive, &u.ChangePassword,
		&u.FailedLoginAttempts, &u.IsLocked, &u.LockedUntil, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", errors.New("user not found")
	}
	return u, password, err
}

func (r *WebUserRepository) Update(id uuid.UUID, req *models.WebUserUpdateRequest) error {
	_, err := r.db.Exec(`
		UPDATE web_users SET full_name = COALESCE($1, full_name), phone = COALESCE($2, phone), updated_at = NOW()
		WHERE id = $3`,
		req.FullName, req.Phone, id,
	)
	return err
}

func (r *WebUserRepository) UpdatePassword(id uuid.UUID, hashed string) error {
	_, err := r.db.Exec(`
		UPDATE web_users SET password = $1, change_password = FALSE, updated_at = NOW()
		WHERE id = $2`, hashed, id,
	)
	return err
}

func (r *WebUserRepository) RecordLogin(id uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE web_users SET last_login_at = NOW(), failed_login_attempts = 0, updated_at = NOW()
		WHERE id = $1`, id,
	)
	return err
}

func (r *WebUserRepository) RecordFailedLogin(id uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE web_users SET failed_login_attempts = failed_login_attempts + 1, updated_at = NOW()
		WHERE id = $1`, id,
	)
	return err
}

func (r *WebUserRepository) SetLocked(id uuid.UUID, until *time.Time) error {
	_, err := r.db.Exec(`
		UPDATE web_users SET is_locked = $1, locked_until = $2, updated_at = NOW()
		WHERE id = $3`, until != nil, until, id,
	)
	return err
}

func (r *WebUserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM web_users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}
