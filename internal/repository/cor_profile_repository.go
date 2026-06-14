package repository

import (
	"database/sql"
	"errors"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type CorProfileRepository struct {
	db *sql.DB
}

func NewCorProfileRepository() *CorProfileRepository {
	return &CorProfileRepository{db: database.DB}
}

// GetOrCreate inserts a new profile if (org_id, email) not seen before,
// otherwise returns the existing record untouched.
func (r *CorProfileRepository) GetOrCreate(orgID, companyID uuid.UUID, branchID *uuid.UUID, req models.CorBookingProfileInput) (*models.CorProfile, error) {
	p := &models.CorProfile{}
	err := r.db.QueryRow(`
		INSERT INTO cor_profiles (org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (org_id, email) DO NOTHING
		RETURNING id, org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department, status, created_at, updated_at`,
		orgID, companyID, branchID, req.FirstName, req.LastName, req.Email, req.Phone, req.JobTitle, req.Department,
	).Scan(
		&p.ID, &p.OrgID, &p.CompanyID, &p.BranchID,
		&p.FirstName, &p.LastName, &p.Email, &p.Phone,
		&p.JobTitle, &p.Department, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return r.GetByEmail(orgID, req.Email)
	}
	return p, err
}

func (r *CorProfileRepository) GetByID(id, orgID uuid.UUID) (*models.CorProfile, error) {
	p := &models.CorProfile{}
	err := r.db.QueryRow(`
		SELECT id, org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department, status, created_at, updated_at
		FROM cor_profiles WHERE id = $1 AND org_id = $2`, id, orgID,
	).Scan(
		&p.ID, &p.OrgID, &p.CompanyID, &p.BranchID,
		&p.FirstName, &p.LastName, &p.Email, &p.Phone,
		&p.JobTitle, &p.Department, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("profile not found")
	}
	return p, err
}

func (r *CorProfileRepository) GetByEmail(orgID uuid.UUID, email string) (*models.CorProfile, error) {
	p := &models.CorProfile{}
	err := r.db.QueryRow(`
		SELECT id, org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department, status, created_at, updated_at
		FROM cor_profiles WHERE org_id = $1 AND email = $2`, orgID, email,
	).Scan(
		&p.ID, &p.OrgID, &p.CompanyID, &p.BranchID,
		&p.FirstName, &p.LastName, &p.Email, &p.Phone,
		&p.JobTitle, &p.Department, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("profile not found")
	}
	return p, err
}

func (r *CorProfileRepository) List(companyID uuid.UUID, page, pageSize int) ([]models.CorProfile, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM cor_profiles WHERE company_id = $1`, companyID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department, status, created_at, updated_at
		FROM cor_profiles WHERE company_id = $1
		ORDER BY last_name, first_name ASC LIMIT $2 OFFSET $3`, companyID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var profiles []models.CorProfile
	for rows.Next() {
		var p models.CorProfile
		if err := rows.Scan(
			&p.ID, &p.OrgID, &p.CompanyID, &p.BranchID,
			&p.FirstName, &p.LastName, &p.Email, &p.Phone,
			&p.JobTitle, &p.Department, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, p)
	}
	return profiles, total, nil
}

func (r *CorProfileRepository) Update(id, orgID uuid.UUID, req *models.UpdateCorProfileRequest) (*models.CorProfile, error) {
	p := &models.CorProfile{}
	err := r.db.QueryRow(`
		UPDATE cor_profiles SET
			branch_id  = COALESCE($1, branch_id),
			first_name = COALESCE($2, first_name),
			last_name  = COALESCE($3, last_name),
			phone      = COALESCE($4, phone),
			job_title  = COALESCE($5, job_title),
			department = COALESCE($6, department),
			status     = COALESCE($7, status),
			updated_at = NOW()
		WHERE id = $8 AND org_id = $9
		RETURNING id, org_id, company_id, branch_id, first_name, last_name, email, phone, job_title, department, status, created_at, updated_at`,
		req.BranchID, req.FirstName, req.LastName, req.Phone,
		req.JobTitle, req.Department, req.Status, id, orgID,
	).Scan(
		&p.ID, &p.OrgID, &p.CompanyID, &p.BranchID,
		&p.FirstName, &p.LastName, &p.Email, &p.Phone,
		&p.JobTitle, &p.Department, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("profile not found")
	}
	return p, err
}
