package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository() *BranchRepository {
	return &BranchRepository{db: database.DB}
}

func (r *BranchRepository) Create(b *models.Branch) error {
	b.ID = uuid.New()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO branches (id, org_id, name, branch_code, street_address, city, country, location, phone, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		b.ID, b.OrgID, b.Name, b.BranchCode, b.StreetAddress, b.City, b.Country, b.Location, b.Phone, b.Email, b.IsActive, b.CreatedAt, b.UpdatedAt,
	)
	return err
}

func (r *BranchRepository) GetByID(id, orgID uuid.UUID) (*models.Branch, error) {
	row := r.db.QueryRow(`
		SELECT id, org_id, name, branch_code, street_address, city, country, location, phone, email, is_active, created_at, updated_at
		FROM branches WHERE id = $1 AND org_id = $2`, id, orgID)
	return scanBranch(row)
}

func (r *BranchRepository) List(orgID uuid.UUID) ([]models.Branch, error) {
	rows, err := r.db.Query(`
		SELECT id, org_id, name, branch_code, street_address, city, country, location, phone, email, is_active, created_at, updated_at
		FROM branches WHERE org_id = $1
		ORDER BY name ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		b, err := scanBranch(rows)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *b)
	}
	if branches == nil {
		branches = []models.Branch{}
	}
	return branches, rows.Err()
}

func (r *BranchRepository) Update(b *models.Branch) error {
	b.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE branches SET name=$1, branch_code=$2, street_address=$3, city=$4, country=$5, location=$6, phone=$7, email=$8, is_active=$9, updated_at=$10
		WHERE id=$11 AND org_id=$12`,
		b.Name, b.BranchCode, b.StreetAddress, b.City, b.Country, b.Location, b.Phone, b.Email, b.IsActive, b.UpdatedAt, b.ID, b.OrgID,
	)
	return err
}

func (r *BranchRepository) Delete(id, orgID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM branches WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("branch not found")
	}
	return nil
}

type branchScanner interface {
	Scan(dest ...interface{}) error
}

func scanBranch(row branchScanner) (*models.Branch, error) {
	var b models.Branch
	var streetAddress, city, country, location, phone, email sql.NullString
	err := row.Scan(&b.ID, &b.OrgID, &b.Name, &b.BranchCode, &streetAddress, &city, &country, &location, &phone, &email, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if streetAddress.Valid {
		b.StreetAddress = &streetAddress.String
	}
	if city.Valid {
		b.City = &city.String
	}
	if country.Valid {
		b.Country = &country.String
	}
	if location.Valid {
		b.Location = &location.String
	}
	if phone.Valid {
		b.Phone = &phone.String
	}
	if email.Valid {
		b.Email = &email.String
	}
	return &b, nil
}
