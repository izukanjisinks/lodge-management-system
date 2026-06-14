package repository

import (
	"database/sql"
	"errors"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type CorCompanyRepository struct {
	db *sql.DB
}

func NewCorCompanyRepository() *CorCompanyRepository {
	return &CorCompanyRepository{db: database.DB}
}

// GetOrCreate inserts a new company if (org_id, reg_number, tpin) not seen before,
// otherwise returns the existing record untouched.
func (r *CorCompanyRepository) GetOrCreate(orgID uuid.UUID, req models.CorBookingCompanyInput) (*models.CorCompanyDetails, error) {
	c := &models.CorCompanyDetails{}
	err := r.db.QueryRow(`
		INSERT INTO cor_company_details (org_id, company_name, tpin, reg_number, industry, country)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (org_id, reg_number, tpin) DO NOTHING
		RETURNING id, org_id, company_name, tpin, reg_number, industry, country, status, created_at, updated_at`,
		orgID, req.CompanyName, req.TPIN, req.RegNumber, req.Industry, req.Country,
	).Scan(
		&c.ID, &c.OrgID, &c.CompanyName, &c.TPIN, &c.RegNumber,
		&c.Industry, &c.Country, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return r.GetByRegNumber(orgID, req.RegNumber)
	}
	return c, err
}

func (r *CorCompanyRepository) GetByID(id, orgID uuid.UUID) (*models.CorCompanyDetails, error) {
	c := &models.CorCompanyDetails{}
	err := r.db.QueryRow(`
		SELECT id, org_id, company_name, tpin, reg_number, industry, country, status, created_at, updated_at
		FROM cor_company_details WHERE id = $1 AND org_id = $2`, id, orgID,
	).Scan(
		&c.ID, &c.OrgID, &c.CompanyName, &c.TPIN, &c.RegNumber,
		&c.Industry, &c.Country, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("company not found")
	}
	return c, err
}

func (r *CorCompanyRepository) GetByRegNumber(orgID uuid.UUID, regNumber string) (*models.CorCompanyDetails, error) {
	c := &models.CorCompanyDetails{}
	err := r.db.QueryRow(`
		SELECT id, org_id, company_name, tpin, reg_number, industry, country, status, created_at, updated_at
		FROM cor_company_details WHERE org_id = $1 AND reg_number = $2`, orgID, regNumber,
	).Scan(
		&c.ID, &c.OrgID, &c.CompanyName, &c.TPIN, &c.RegNumber,
		&c.Industry, &c.Country, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("company not found")
	}
	return c, err
}

func (r *CorCompanyRepository) List(orgID uuid.UUID, page, pageSize int) ([]models.CorCompanyDetails, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM cor_company_details WHERE org_id = $1`, orgID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, org_id, company_name, tpin, reg_number, industry, country, status, created_at, updated_at
		FROM cor_company_details WHERE org_id = $1
		ORDER BY company_name ASC LIMIT $2 OFFSET $3`, orgID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []models.CorCompanyDetails
	for rows.Next() {
		var c models.CorCompanyDetails
		if err := rows.Scan(
			&c.ID, &c.OrgID, &c.CompanyName, &c.TPIN, &c.RegNumber,
			&c.Industry, &c.Country, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, c)
	}
	return companies, total, nil
}

func (r *CorCompanyRepository) Update(id, orgID uuid.UUID, req *models.UpdateCorCompanyRequest) (*models.CorCompanyDetails, error) {
	c := &models.CorCompanyDetails{}
	err := r.db.QueryRow(`
		UPDATE cor_company_details SET
			company_name = COALESCE($1, company_name),
			tpin         = COALESCE($2, tpin),
			industry     = COALESCE($3, industry),
			country      = COALESCE($4, country),
			status       = COALESCE($5, status),
			updated_at   = NOW()
		WHERE id = $6 AND org_id = $7
		RETURNING id, org_id, company_name, tpin, reg_number, industry, country, status, created_at, updated_at`,
		req.CompanyName, req.TPIN, req.Industry, req.Country, req.Status, id, orgID,
	).Scan(
		&c.ID, &c.OrgID, &c.CompanyName, &c.TPIN, &c.RegNumber,
		&c.Industry, &c.Country, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("company not found")
	}
	return c, err
}
