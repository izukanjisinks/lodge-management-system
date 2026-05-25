package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{db: database.DB}
}

func (r *OrganizationRepository) Create(org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, logo_url, street_address, city, country, location, phone, email, parking, restaurant, check_in_time, check_out_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	org.ID = uuid.New()
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now
	_, err := r.db.Exec(query,
		org.ID, org.Name, org.LogoURL, org.StreetAddress, org.City, org.Country, org.Location, org.Phone, org.Email, org.Parking, org.Restaurant, org.CheckInTime, org.CheckOutTime, org.IsActive, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *OrganizationRepository) CreateTx(tx *sql.Tx, org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, logo_url, street_address, city, country, location, phone, email, parking, restaurant, check_in_time, check_out_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	org.ID = uuid.New()
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now
	_, err := tx.Exec(query,
		org.ID, org.Name, org.LogoURL, org.StreetAddress, org.City, org.Country, org.Location, org.Phone, org.Email, org.Parking, org.Restaurant, org.CheckInTime, org.CheckOutTime, true, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *OrganizationRepository) GetByID(id uuid.UUID) (*models.Organization, error) {
	query := `
		SELECT id, name, logo_url, street_address, city, country, location, phone, email, parking, restaurant, check_in_time, check_out_time, is_active, created_at, updated_at
		FROM organizations WHERE id = $1`
	return r.scanOrganization(r.db.QueryRow(query, id))
}

func (r *OrganizationRepository) List() ([]models.Organization, error) {
	query := `
		SELECT id, name, logo_url, street_address, city, country, location, phone, email, parking, restaurant, check_in_time, check_out_time, is_active, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		org, err := r.scanOrganization(rows)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, *org)
	}
	if orgs == nil {
		orgs = []models.Organization{}
	}
	return orgs, rows.Err()
}

// ListIDs returns just the IDs of all active organizations — used by background jobs.
func (r *OrganizationRepository) ListIDs() ([]uuid.UUID, error) {
	rows, err := r.db.Query(`SELECT id FROM organizations WHERE is_active = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *OrganizationRepository) Update(org *models.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, logo_url = $2, street_address = $3, city = $4, country = $5, location = $6, phone = $7, email = $8, parking = $9, restaurant = $10, check_in_time = $11, check_out_time = $12, is_active = $13, updated_at = $14
		WHERE id = $15`
	_, err := r.db.Exec(query,
		org.Name, org.LogoURL, org.StreetAddress, org.City, org.Country, org.Location, org.Phone, org.Email, org.Parking, org.Restaurant, org.CheckInTime, org.CheckOutTime, org.IsActive, time.Now(), org.ID,
	)
	return err
}

func (r *OrganizationRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM organizations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

// ListPublic returns active organizations for guest-facing lodge discovery, paginated.
func (r *OrganizationRepository) ListPublic(page, pageSize int) ([]models.Organization, int, error) {
	offset := (page - 1) * pageSize

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM organizations WHERE is_active = TRUE`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, name, logo_url, street_address, city, country, location, phone, email, parking, restaurant, check_in_time, check_out_time, is_active, created_at, updated_at
		FROM organizations
		WHERE is_active = TRUE
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		org, err := r.scanOrganization(rows)
		if err != nil {
			return nil, 0, err
		}
		orgs = append(orgs, *org)
	}
	if orgs == nil {
		orgs = []models.Organization{}
	}
	return orgs, total, rows.Err()
}

func (r *OrganizationRepository) scanOrganization(row rowScanner) (*models.Organization, error) {
	var org models.Organization
	var logoURL, streetAddress, city, country, location, phone, email, checkInTime, checkOutTime sql.NullString
	err := row.Scan(
		&org.ID, &org.Name, &logoURL, &streetAddress, &city, &country, &location, &phone, &email, &org.Parking, &org.Restaurant, &checkInTime, &checkOutTime, &org.IsActive, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	org.LogoURL = logoURL.String
	org.StreetAddress = streetAddress.String
	org.City = city.String
	org.Country = country.String
	org.Location = location.String
	org.Phone = phone.String
	org.Email = email.String
	if checkInTime.Valid {
		org.CheckInTime = &checkInTime.String
	}
	if checkOutTime.Valid {
		org.CheckOutTime = &checkOutTime.String
	}
	return &org, nil
}
