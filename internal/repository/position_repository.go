package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/interfaces"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type PositionRepository struct {
	db *sql.DB
}

func NewPositionRepository() *PositionRepository {
	return &PositionRepository{db: database.DB}
}

func (r *PositionRepository) Create(pos *models.Position) error {
	pos.ID = uuid.New()
	now := time.Now()
	pos.CreatedAt = now
	pos.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO positions (id, title, code, department_id, role_id, grade_level, base_salary, housing_allowance, transport_allowance, medical_allowance, income_tax, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		pos.ID, pos.Title, pos.Code, pos.DepartmentID, pos.RoleID, pos.GradeLevel,
		pos.BaseSalary, pos.HousingAllowance, pos.TransportAllowance, pos.MedicalAllowance, pos.IncomeTax,
		pos.Description, pos.IsActive, pos.CreatedAt, pos.UpdatedAt,
	)
	return err
}

func (r *PositionRepository) GetByID(id uuid.UUID) (*models.Position, error) {
	var p models.Position
	err := r.db.QueryRow(`
		SELECT id, title, code, department_id, role_id, grade_level, base_salary, housing_allowance, transport_allowance, medical_allowance, income_tax, description, is_active, created_at, updated_at, deleted_at
		FROM positions WHERE id=$1 AND deleted_at IS NULL`, id,
	).Scan(&p.ID, &p.Title, &p.Code, &p.DepartmentID, &p.RoleID, &p.GradeLevel,
		&p.BaseSalary, &p.HousingAllowance, &p.TransportAllowance, &p.MedicalAllowance, &p.IncomeTax,
		&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PositionRepository) List(filter interfaces.PositionFilter, page, pageSize int) ([]models.Position, int, error) {
	args := []interface{}{}
	where := []string{"p.deleted_at IS NULL"}
	i := 1

	if filter.DepartmentID != nil {
		where = append(where, fmt.Sprintf("p.department_id=$%d", i))
		args = append(args, *filter.DepartmentID)
		i++
	}
	if filter.GradeLevel != "" {
		where = append(where, fmt.Sprintf("p.grade_level=$%d", i))
		args = append(args, filter.GradeLevel)
		i++
	}
	if filter.IsActive != nil {
		where = append(where, fmt.Sprintf("p.is_active=$%d", i))
		args = append(args, *filter.IsActive)
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM positions p WHERE %s`, whereStr), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT p.id, p.title, p.code, p.department_id, p.role_id, p.grade_level, p.base_salary, p.housing_allowance, p.transport_allowance, p.medical_allowance, p.income_tax, p.description, p.is_active, p.created_at, p.updated_at, p.deleted_at,
		       COALESCE(d.name, '') AS department_name
		FROM positions p
		LEFT JOIN departments d ON p.department_id = d.id
		WHERE %s ORDER BY p.title LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var positions []models.Position
	for rows.Next() {
		var p models.Position
		if err := rows.Scan(&p.ID, &p.Title, &p.Code, &p.DepartmentID, &p.RoleID, &p.GradeLevel,
			&p.BaseSalary, &p.HousingAllowance, &p.TransportAllowance, &p.MedicalAllowance, &p.IncomeTax,
			&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
			&p.DepartmentName); err != nil {
			return nil, 0, err
		}
		positions = append(positions, p)
	}
	return positions, total, rows.Err()
}

func (r *PositionRepository) Update(pos *models.Position) error {
	pos.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE positions SET title=$1, code=$2, department_id=$3, role_id=$4, grade_level=$5, base_salary=$6,
		housing_allowance=$7, transport_allowance=$8, medical_allowance=$9, income_tax=$10,
		description=$11, is_active=$12, updated_at=$13 WHERE id=$14 AND deleted_at IS NULL`,
		pos.Title, pos.Code, pos.DepartmentID, pos.RoleID, pos.GradeLevel, pos.BaseSalary,
		pos.HousingAllowance, pos.TransportAllowance, pos.MedicalAllowance, pos.IncomeTax,
		pos.Description, pos.IsActive, pos.UpdatedAt, pos.ID,
	)
	return err
}

func (r *PositionRepository) SoftDelete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE positions SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL`, time.Now(), id)
	return err
}

func (r *PositionRepository) CodeExists(code string, excludeID *uuid.UUID) (bool, error) {
	var count int
	if excludeID != nil {
		err := r.db.QueryRow(`SELECT COUNT(1) FROM positions WHERE code=$1 AND id!=$2 AND deleted_at IS NULL`, code, excludeID).Scan(&count)
		return count > 0, err
	}
	err := r.db.QueryRow(`SELECT COUNT(1) FROM positions WHERE code=$1 AND deleted_at IS NULL`, code).Scan(&count)
	return count > 0, err
}

func (r *PositionRepository) ActiveEmployeeCount(id uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM employees WHERE position_id=$1 AND deleted_at IS NULL AND employment_status='active'`, id).Scan(&count)
	return count, err
}
