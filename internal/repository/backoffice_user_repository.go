package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type BackofficeUserRepository struct {
	db *sql.DB
}

func NewBackofficeUserRepository() *BackofficeUserRepository {
	return &BackofficeUserRepository{db: database.DB}
}

func (r *BackofficeUserRepository) Create(u *models.BackofficeUser) error {
	query := `
		INSERT INTO backoffice_users (id, full_name, email, password, is_active, change_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	_, err := r.db.Exec(query,
		u.ID, u.FullName, u.Email, u.Password, u.IsActive, u.ChangePassword, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *BackofficeUserRepository) GetByID(id uuid.UUID) (*models.BackofficeUser, error) {
	query := `
		SELECT id, full_name, email, password, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM backoffice_users WHERE id = $1`
	return r.scanUser(r.db.QueryRow(query, id))
}

func (r *BackofficeUserRepository) GetByEmail(email string) (*models.BackofficeUser, error) {
	query := `
		SELECT id, full_name, email, password, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM backoffice_users WHERE email = $1`
	return r.scanUser(r.db.QueryRow(query, email))
}

func (r *BackofficeUserRepository) List() ([]models.BackofficeUser, error) {
	query := `
		SELECT id, full_name, email, password, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM backoffice_users
		ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.BackofficeUser
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	if users == nil {
		users = []models.BackofficeUser{}
	}
	return users, rows.Err()
}

func (r *BackofficeUserRepository) Update(u *models.BackofficeUser) error {
	query := `
		UPDATE backoffice_users
		SET full_name = $1, email = $2, password = $3, is_active = $4, change_password = $5,
		    failed_login_attempts = $6, is_locked = $7, locked_until = $8, last_login_at = $9, updated_at = $10
		WHERE id = $11`
	_, err := r.db.Exec(query,
		u.FullName, u.Email, u.Password, u.IsActive, u.ChangePassword,
		u.FailedLoginAttempts, u.IsLocked, u.LockedUntil, u.LastLoginAt, time.Now(), u.ID,
	)
	return err
}

func (r *BackofficeUserRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM backoffice_users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("backoffice user not found")
	}
	return nil
}

func (r *BackofficeUserRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM backoffice_users WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}

func (r *BackofficeUserRepository) scanUser(row rowScanner) (*models.BackofficeUser, error) {
	var u models.BackofficeUser
	var lockedUntil, lastLoginAt sql.NullTime
	err := row.Scan(
		&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsActive, &u.ChangePassword,
		&u.FailedLoginAttempts, &u.IsLocked, &lockedUntil, &lastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		u.LockedUntil = &lockedUntil.Time
	}
	if lastLoginAt.Valid {
		u.LastLoginAt = &lastLoginAt.Time
	}
	return &u, nil
}
