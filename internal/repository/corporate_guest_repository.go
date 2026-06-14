package repository

import (
	"database/sql"
	"errors"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type CorporateGuestRepository struct {
	db *sql.DB
}

func NewCorporateGuestRepository() *CorporateGuestRepository {
	return &CorporateGuestRepository{db: database.DB}
}

func (r *CorporateGuestRepository) Create(profileID uuid.UUID, req models.CorBookingGuestInput) (*models.CorporateGuest, error) {
	g := &models.CorporateGuest{}
	err := r.db.QueryRow(`
		INSERT INTO corporate_guests (corporate_profile_id, first_name, last_name, phone, email, identification_card)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, corporate_profile_id, first_name, last_name, phone, email, identification_card, created_at, updated_at`,
		profileID, req.FirstName, req.LastName, req.Phone, req.Email, req.IdentificationCard,
	).Scan(
		&g.ID, &g.CorporateProfileID, &g.FirstName, &g.LastName,
		&g.Phone, &g.Email, &g.IdentificationCard, &g.CreatedAt, &g.UpdatedAt,
	)
	return g, err
}

// CreateMany registers each guest on the corporate roster. The roster is permanent
// and shared across bookings, so a returning delegate (same profile + ID card) is
// updated rather than re-inserted — otherwise the unique constraint would reject the
// whole request the second time a company books the same person.
func (r *CorporateGuestRepository) CreateMany(profileID uuid.UUID, guests []models.CorBookingGuestInput) ([]models.CorporateGuest, error) {
	var result []models.CorporateGuest
	for _, g := range guests {
		upserted, err := r.Upsert(profileID, g)
		if err != nil {
			return nil, err
		}
		result = append(result, *upserted)
	}
	return result, nil
}

func (r *CorporateGuestRepository) Upsert(profileID uuid.UUID, g models.CorBookingGuestInput) (*models.CorporateGuest, error) {
	result := &models.CorporateGuest{}
	err := r.db.QueryRow(`
		INSERT INTO corporate_guests (corporate_profile_id, first_name, last_name, phone, email, identification_card)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (corporate_profile_id, identification_card)
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name  = EXCLUDED.last_name,
			phone      = COALESCE(EXCLUDED.phone, corporate_guests.phone),
			email      = COALESCE(EXCLUDED.email, corporate_guests.email),
			updated_at = NOW()
		RETURNING id, corporate_profile_id, first_name, last_name, phone, email, identification_card, created_at, updated_at`,
		profileID, g.FirstName, g.LastName, g.Phone, g.Email, g.IdentificationCard,
	).Scan(
		&result.ID, &result.CorporateProfileID, &result.FirstName, &result.LastName,
		&result.Phone, &result.Email, &result.IdentificationCard, &result.CreatedAt, &result.UpdatedAt,
	)
	return result, err
}

func (r *CorporateGuestRepository) GetByID(id, profileID uuid.UUID) (*models.CorporateGuest, error) {
	g := &models.CorporateGuest{}
	err := r.db.QueryRow(`
		SELECT id, corporate_profile_id, first_name, last_name, phone, email, identification_card, created_at, updated_at
		FROM corporate_guests WHERE id = $1 AND corporate_profile_id = $2`, id, profileID,
	).Scan(
		&g.ID, &g.CorporateProfileID, &g.FirstName, &g.LastName,
		&g.Phone, &g.Email, &g.IdentificationCard, &g.CreatedAt, &g.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("guest not found")
	}
	return g, err
}

func (r *CorporateGuestRepository) List(profileID uuid.UUID) ([]models.CorporateGuest, error) {
	rows, err := r.db.Query(`
		SELECT id, corporate_profile_id, first_name, last_name, phone, email, identification_card, created_at, updated_at
		FROM corporate_guests WHERE corporate_profile_id = $1
		ORDER BY last_name, first_name ASC`, profileID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []models.CorporateGuest
	for rows.Next() {
		var g models.CorporateGuest
		if err := rows.Scan(
			&g.ID, &g.CorporateProfileID, &g.FirstName, &g.LastName,
			&g.Phone, &g.Email, &g.IdentificationCard, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}
	return guests, nil
}

func (r *CorporateGuestRepository) Update(id, profileID uuid.UUID, req *models.UpdateCorporateGuestRequest) (*models.CorporateGuest, error) {
	g := &models.CorporateGuest{}
	err := r.db.QueryRow(`
		UPDATE corporate_guests SET
			first_name          = COALESCE($1, first_name),
			last_name           = COALESCE($2, last_name),
			phone               = COALESCE($3, phone),
			email               = COALESCE($4, email),
			identification_card = COALESCE($5, identification_card),
			updated_at          = NOW()
		WHERE id = $6 AND corporate_profile_id = $7
		RETURNING id, corporate_profile_id, first_name, last_name, phone, email, identification_card, created_at, updated_at`,
		req.FirstName, req.LastName, req.Phone, req.Email, req.IdentificationCard, id, profileID,
	).Scan(
		&g.ID, &g.CorporateProfileID, &g.FirstName, &g.LastName,
		&g.Phone, &g.Email, &g.IdentificationCard, &g.CreatedAt, &g.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("guest not found")
	}
	return g, err
}

func (r *CorporateGuestRepository) Delete(id, profileID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM corporate_guests WHERE id = $1 AND corporate_profile_id = $2`, id, profileID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("guest not found")
	}
	return nil
}
