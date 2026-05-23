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
		INSERT INTO branches (id, org_id, name, branch_code, address, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		b.ID, b.OrgID, b.Name, b.BranchCode, b.Address, b.Location, b.CreatedAt, b.UpdatedAt,
	)
	return err
}

func (r *BranchRepository) GetByID(id, orgID uuid.UUID) (*models.Branch, error) {
	row := r.db.QueryRow(`
		SELECT id, org_id, name, branch_code, address, location, created_at, updated_at
		FROM branches WHERE id = $1 AND org_id = $2`, id, orgID)
	return scanBranch(row)
}

func (r *BranchRepository) List(orgID uuid.UUID) ([]models.Branch, error) {
	rows, err := r.db.Query(`
		SELECT id, org_id, name, branch_code, address, location, created_at, updated_at
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
		UPDATE branches SET name=$1, address=$2, location=$3, updated_at=$4
		WHERE id=$5 AND org_id=$6`,
		b.Name, b.Address, b.Location, b.UpdatedAt, b.ID, b.OrgID,
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
	var address, location sql.NullString
	err := row.Scan(&b.ID, &b.OrgID, &b.Name, &b.BranchCode, &address, &location, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if address.Valid {
		b.Address = address.String
	}
	if location.Valid {
		b.Location = location.String
	}
	return &b, nil
}
