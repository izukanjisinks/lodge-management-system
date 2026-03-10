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

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{db: database.DB}
}

func (r *DepartmentRepository) Create(dept *models.Department) error {
	dept.ID = uuid.New()
	now := time.Now()
	dept.CreatedAt = now
	dept.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO departments (id, name, code, description, parent_department_id, manager_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		dept.ID, dept.Name, dept.Code, dept.Description,
		dept.ParentDepartmentID, dept.ManagerID, dept.IsActive, dept.CreatedAt, dept.UpdatedAt,
	)
	return err
}

func (r *DepartmentRepository) GetByID(id uuid.UUID) (*models.Department, error) {
	var d models.Department
	var parentID, managerID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, name, code, description, parent_department_id, manager_id, is_active, created_at, updated_at, deleted_at
		FROM departments WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&d.ID, &d.Name, &d.Code, &d.Description, &parentID, &managerID,
		&d.IsActive, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		p, _ := uuid.Parse(parentID.String)
		d.ParentDepartmentID = &p
	}
	if managerID.Valid {
		m, _ := uuid.Parse(managerID.String)
		d.ManagerID = &m
	}
	return &d, nil
}

func (r *DepartmentRepository) List(filter interfaces.DepartmentFilter, page, pageSize int) ([]models.Department, int, error) {
	args := []interface{}{}
	where := []string{"d.deleted_at IS NULL"}
	i := 1

	if filter.IsActive != nil {
		where = append(where, fmt.Sprintf("d.is_active = $%d", i))
		args = append(args, *filter.IsActive)
		i++
	}
	if filter.Search != "" {
		where = append(where, fmt.Sprintf("(d.name ILIKE $%d OR d.code ILIKE $%d)", i, i))
		args = append(args, "%"+filter.Search+"%")
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM departments d WHERE %s`, whereStr), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT d.id, d.name, d.code, d.description, d.parent_department_id, d.manager_id, d.is_active, d.created_at, d.updated_at, d.deleted_at,
		       COALESCE(pd.name, '') AS parent_department_name,
		       COALESCE(e.first_name || ' ' || e.last_name, '') AS manager_name
		FROM departments d
		LEFT JOIN departments pd ON d.parent_department_id = pd.id
		LEFT JOIN employees e ON d.manager_id = e.id
		WHERE %s ORDER BY d.name LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		var parentID, managerID sql.NullString
		if err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.Description, &parentID, &managerID,
			&d.IsActive, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
			&d.ParentDepartmentName, &d.ManagerName); err != nil {
			return nil, 0, err
		}
		if parentID.Valid {
			p, _ := uuid.Parse(parentID.String)
			d.ParentDepartmentID = &p
		}
		if managerID.Valid {
			m, _ := uuid.Parse(managerID.String)
			d.ManagerID = &m
		}
		depts = append(depts, d)
	}
	return depts, total, rows.Err()
}

func (r *DepartmentRepository) Update(dept *models.Department) error {
	dept.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE departments SET name=$1, code=$2, description=$3, parent_department_id=$4,
		manager_id=$5, is_active=$6, updated_at=$7 WHERE id=$8 AND deleted_at IS NULL`,
		dept.Name, dept.Code, dept.Description, dept.ParentDepartmentID,
		dept.ManagerID, dept.IsActive, dept.UpdatedAt, dept.ID,
	)
	return err
}

func (r *DepartmentRepository) SoftDelete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE departments SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL`, time.Now(), id)
	return err
}

func (r *DepartmentRepository) GetAll() ([]models.Department, error) {
	rows, err := r.db.Query(`
		SELECT id, name, code, description, parent_department_id, manager_id, is_active, created_at, updated_at, deleted_at
		FROM departments WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		var parentID, managerID sql.NullString
		if err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.Description, &parentID, &managerID,
			&d.IsActive, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			p, _ := uuid.Parse(parentID.String)
			d.ParentDepartmentID = &p
		}
		if managerID.Valid {
			m, _ := uuid.Parse(managerID.String)
			d.ManagerID = &m
		}
		depts = append(depts, d)
	}
	return depts, rows.Err()
}

func (r *DepartmentRepository) GetEmployeeCount(id uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM employees WHERE department_id=$1 AND deleted_at IS NULL AND employment_status='active'`,
		id,
	).Scan(&count)
	return count, err
}

func (r *DepartmentRepository) CodeExists(code string, excludeID *uuid.UUID) (bool, error) {
	var count int
	if excludeID != nil {
		err := r.db.QueryRow(`SELECT COUNT(1) FROM departments WHERE code=$1 AND id!=$2 AND deleted_at IS NULL`, code, excludeID).Scan(&count)
		return count > 0, err
	}
	err := r.db.QueryRow(`SELECT COUNT(1) FROM departments WHERE code=$1 AND deleted_at IS NULL`, code).Scan(&count)
	return count > 0, err
}
