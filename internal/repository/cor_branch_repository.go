package repository

import (
	"database/sql"
	"errors"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type CorBranchRepository struct {
	db *sql.DB
}

func NewCorBranchRepository() *CorBranchRepository {
	return &CorBranchRepository{db: database.DB}
}

// GetOrCreate inserts a new branch if (company_id, name) not seen before,
// otherwise returns the existing record untouched.
func (r *CorBranchRepository) GetOrCreate(companyID uuid.UUID, req *models.CorBookingBranchInput) (*models.CorBranchDetails, error) {
	if req == nil || req.Name == "" {
		return nil, nil
	}
	b := &models.CorBranchDetails{}
	err := r.db.QueryRow(`
		INSERT INTO cor_branch_details (company_id, name, address, phone)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (company_id, name) DO NOTHING
		RETURNING id, company_id, name, address, phone, created_at, updated_at`,
		companyID, req.Name, req.Address, req.Phone,
	).Scan(&b.ID, &b.CompanyID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return r.GetByName(companyID, req.Name)
	}
	return b, err
}

func (r *CorBranchRepository) GetByID(id, companyID uuid.UUID) (*models.CorBranchDetails, error) {
	b := &models.CorBranchDetails{}
	err := r.db.QueryRow(`
		SELECT id, company_id, name, address, phone, created_at, updated_at
		FROM cor_branch_details WHERE id = $1 AND company_id = $2`, id, companyID,
	).Scan(&b.ID, &b.CompanyID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("branch not found")
	}
	return b, err
}

func (r *CorBranchRepository) GetByName(companyID uuid.UUID, name string) (*models.CorBranchDetails, error) {
	b := &models.CorBranchDetails{}
	err := r.db.QueryRow(`
		SELECT id, company_id, name, address, phone, created_at, updated_at
		FROM cor_branch_details WHERE company_id = $1 AND name = $2`, companyID, name,
	).Scan(&b.ID, &b.CompanyID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("branch not found")
	}
	return b, err
}

func (r *CorBranchRepository) List(companyID uuid.UUID) ([]models.CorBranchDetails, error) {
	rows, err := r.db.Query(`
		SELECT id, company_id, name, address, phone, created_at, updated_at
		FROM cor_branch_details WHERE company_id = $1 ORDER BY name ASC`, companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []models.CorBranchDetails
	for rows.Next() {
		var b models.CorBranchDetails
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (r *CorBranchRepository) Update(id, companyID uuid.UUID, req *models.UpdateCorBranchRequest) (*models.CorBranchDetails, error) {
	b := &models.CorBranchDetails{}
	err := r.db.QueryRow(`
		UPDATE cor_branch_details SET
			name       = COALESCE($1, name),
			address    = COALESCE($2, address),
			phone      = COALESCE($3, phone),
			updated_at = NOW()
		WHERE id = $4 AND company_id = $5
		RETURNING id, company_id, name, address, phone, created_at, updated_at`,
		req.Name, req.Address, req.Phone, id, companyID,
	).Scan(&b.ID, &b.CompanyID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("branch not found")
	}
	return b, err
}
