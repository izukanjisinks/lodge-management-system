package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type LeaveTypeRepository struct {
	db *sql.DB
}

func NewLeaveTypeRepository() *LeaveTypeRepository {
	return &LeaveTypeRepository{db: database.DB}
}

func (r *LeaveTypeRepository) Create(lt *models.LeaveType) error {
	lt.ID = uuid.New()
	now := time.Now()
	lt.CreatedAt = now
	lt.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO leave_types
		(id, name, code, description, default_days_per_year, is_paid, is_carry_forward_allowed,
		 max_carry_forward_days, requires_approval, requires_document, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		lt.ID, lt.Name, lt.Code, lt.Description, lt.DefaultDaysPerYear, lt.IsPaid,
		lt.IsCarryForwardAllowed, lt.MaxCarryForwardDays, lt.RequiresApproval, lt.RequiresDocument,
		lt.IsActive, lt.CreatedAt, lt.UpdatedAt,
	)
	return err
}

func (r *LeaveTypeRepository) GetByID(id uuid.UUID) (*models.LeaveType, error) {
	return r.scan(r.db.QueryRow(`
		SELECT id, name, code, description, default_days_per_year, is_paid, is_carry_forward_allowed,
		       max_carry_forward_days, requires_approval, requires_document, is_active, created_at, updated_at
		FROM leave_types WHERE id=$1`, id))
}

func (r *LeaveTypeRepository) GetByCode(code string) (*models.LeaveType, error) {
	return r.scan(r.db.QueryRow(`
		SELECT id, name, code, description, default_days_per_year, is_paid, is_carry_forward_allowed,
		       max_carry_forward_days, requires_approval, requires_document, is_active, created_at, updated_at
		FROM leave_types WHERE code=$1`, code))
}

func (r *LeaveTypeRepository) List(activeOnly bool) ([]models.LeaveType, error) {
	q := `SELECT id, name, code, description, default_days_per_year, is_paid, is_carry_forward_allowed,
		         max_carry_forward_days, requires_approval, requires_document, is_active, created_at, updated_at
		  FROM leave_types`
	if activeOnly {
		q += " WHERE is_active=TRUE"
	}
	q += " ORDER BY name"
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.LeaveType
	for rows.Next() {
		lt, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *lt)
	}
	return out, rows.Err()
}

func (r *LeaveTypeRepository) Update(lt *models.LeaveType) error {
	lt.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE leave_types SET name=$1, code=$2, description=$3, default_days_per_year=$4, is_paid=$5,
		is_carry_forward_allowed=$6, max_carry_forward_days=$7, requires_approval=$8,
		requires_document=$9, is_active=$10, updated_at=$11 WHERE id=$12`,
		lt.Name, lt.Code, lt.Description, lt.DefaultDaysPerYear, lt.IsPaid,
		lt.IsCarryForwardAllowed, lt.MaxCarryForwardDays, lt.RequiresApproval,
		lt.RequiresDocument, lt.IsActive, lt.UpdatedAt, lt.ID,
	)
	return err
}

func (r *LeaveTypeRepository) CodeExists(code string, excludeID *uuid.UUID) (bool, error) {
	var count int
	if excludeID != nil {
		err := r.db.QueryRow(`SELECT COUNT(1) FROM leave_types WHERE code=$1 AND id!=$2`, code, excludeID).Scan(&count)
		return count > 0, err
	}
	err := r.db.QueryRow(`SELECT COUNT(1) FROM leave_types WHERE code=$1`, code).Scan(&count)
	return count > 0, err
}

func (r *LeaveTypeRepository) scan(row rowScanner) (*models.LeaveType, error) {
	var lt models.LeaveType
	err := row.Scan(&lt.ID, &lt.Name, &lt.Code, &lt.Description, &lt.DefaultDaysPerYear, &lt.IsPaid,
		&lt.IsCarryForwardAllowed, &lt.MaxCarryForwardDays, &lt.RequiresApproval, &lt.RequiresDocument,
		&lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &lt, nil
}
