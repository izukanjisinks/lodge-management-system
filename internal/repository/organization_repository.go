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
		INSERT INTO organizations (id, name, logo_url, address, phone, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	org.ID = uuid.New()
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now
	_, err := r.db.Exec(query,
		org.ID, org.Name, org.LogoURL, org.Address, org.Phone, org.Email, org.IsActive, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *OrganizationRepository) CreateTx(tx *sql.Tx, org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, logo_url, address, phone, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	org.ID = uuid.New()
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now
	_, err := tx.Exec(query,
		org.ID, org.Name, org.LogoURL, org.Address, org.Phone, org.Email, true, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *OrganizationRepository) GetByID(id uuid.UUID) (*models.Organization, error) {
	query := `
		SELECT id, name, logo_url, address, phone, email, is_active, created_at, updated_at
		FROM organizations WHERE id = $1`
	return r.scanOrganization(r.db.QueryRow(query, id))
}

func (r *OrganizationRepository) List() ([]models.Organization, error) {
	query := `
		SELECT id, name, logo_url, address, phone, email, is_active, created_at, updated_at
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

func (r *OrganizationRepository) Update(org *models.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, logo_url = $2, address = $3, phone = $4, email = $5, is_active = $6, updated_at = $7
		WHERE id = $8`
	_, err := r.db.Exec(query,
		org.Name, org.LogoURL, org.Address, org.Phone, org.Email, org.IsActive, time.Now(), org.ID,
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

func (r *OrganizationRepository) scanOrganization(row rowScanner) (*models.Organization, error) {
	var org models.Organization
	var logoURL, address, phone, email sql.NullString
	err := row.Scan(
		&org.ID, &org.Name, &logoURL, &address, &phone, &email, &org.IsActive, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	org.LogoURL = logoURL.String
	org.Address = address.String
	org.Phone = phone.String
	org.Email = email.String
	return &org, nil
}
