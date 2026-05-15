package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type GuestRepository struct {
	db *sql.DB
}

func NewGuestRepository() *GuestRepository {
	return &GuestRepository{db: database.DB}
}

func (r *GuestRepository) Create(g *models.Guest) error {
	query := `
		INSERT INTO guests (id, full_name, email, password, phone, is_active, change_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	g.ID = uuid.New()
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	_, err := r.db.Exec(query,
		g.ID, g.FullName, g.Email, g.Password, g.Phone, g.IsActive, g.ChangePassword, g.CreatedAt, g.UpdatedAt,
	)
	return err
}

func (r *GuestRepository) GetByID(id uuid.UUID) (*models.Guest, error) {
	query := `
		SELECT id, full_name, email, password, phone, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM guests WHERE id = $1`
	return r.scanGuest(r.db.QueryRow(query, id))
}

func (r *GuestRepository) GetByEmail(email string) (*models.Guest, error) {
	query := `
		SELECT id, full_name, email, password, phone, is_active, change_password,
		       failed_login_attempts, is_locked, locked_until, last_login_at, created_at, updated_at
		FROM guests WHERE email = $1`
	return r.scanGuest(r.db.QueryRow(query, email))
}

func (r *GuestRepository) Update(g *models.Guest) error {
	query := `
		UPDATE guests
		SET full_name = $1, email = $2, phone = $3, is_active = $4, change_password = $5,
		    failed_login_attempts = $6, is_locked = $7, locked_until = $8, last_login_at = $9, updated_at = $10
		WHERE id = $11`
	_, err := r.db.Exec(query,
		g.FullName, g.Email, g.Phone, g.IsActive, g.ChangePassword,
		g.FailedLoginAttempts, g.IsLocked, g.LockedUntil, g.LastLoginAt, time.Now(), g.ID,
	)
	return err
}

func (r *GuestRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM guests WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}

// CreateIndividualProfile creates an individual_profiles record linked to a guest account.
// org_id is NULL since guests are cross-org.
func (r *GuestRepository) CreateIndividualProfile(guestID uuid.UUID, g *models.Guest) error {
	_, err := r.db.Exec(`
		INSERT INTO individual_profiles (id, guest_id, full_name, email, phone, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'active', $6, $6)`,
		uuid.New(), guestID, g.FullName, g.Email, g.Phone, time.Now(),
	)
	return err
}

func (r *GuestRepository) UpdateIndividualProfileOrg(guestID uuid.UUID, orgID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE individual_profiles SET org_id=$1, updated_at=$2 WHERE guest_id=$3`,
		orgID, time.Now(), guestID,
	)
	return err
}

func (r *GuestRepository) UpdateIndividualProfileIDPassport(profileID uuid.UUID, idPassport string) error {
	_, err := r.db.Exec(`
		UPDATE individual_profiles SET id_passport_number=$1, updated_at=$2 WHERE id=$3`,
		idPassport, time.Now(), profileID,
	)
	return err
}

func (r *GuestRepository) GetIndividualProfileByGuestID(guestID uuid.UUID) (*models.IndividualClient, error) {
	query := `
		SELECT id, full_name, email, phone, id_passport_number, nationality, status, notes, created_at, updated_at
		FROM individual_profiles WHERE guest_id = $1`
	var c models.IndividualClient
	var idPassport, nationality, notes sql.NullString
	err := r.db.QueryRow(query, guestID).Scan(
		&c.ID, &c.FullName, &c.Email, &c.Phone, &idPassport,
		&nationality, &c.Status, &notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("individual profile not found for guest")
	}
	c.IDPassportNumber = idPassport.String
	c.Nationality = nationality.String
	c.Notes = notes.String
	return &c, nil
}

func (r *GuestRepository) scanGuest(row rowScanner) (*models.Guest, error) {
	var g models.Guest
	var phone sql.NullString
	var lockedUntil, lastLoginAt sql.NullTime
	err := row.Scan(
		&g.ID, &g.FullName, &g.Email, &g.Password, &phone, &g.IsActive, &g.ChangePassword,
		&g.FailedLoginAttempts, &g.IsLocked, &lockedUntil, &lastLoginAt, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	g.Phone = phone.String
	if lockedUntil.Valid {
		g.LockedUntil = &lockedUntil.Time
	}
	if lastLoginAt.Valid {
		g.LastLoginAt = &lastLoginAt.Time
	}
	return &g, nil
}
