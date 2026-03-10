package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.DB}
}

func (r *RoleRepository) Create(role *models.Role) error {
	role.RoleID = uuid.New()
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO roles (role_id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		role.RoleID, role.Name, role.Description, role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *RoleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRow(
		`SELECT role_id, name, description, created_at, updated_at FROM roles WHERE role_id = $1`,
		id,
	).Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRow(
		`SELECT role_id, name, description, created_at, updated_at FROM roles WHERE name = $1`,
		name,
	).Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	rows, err := r.db.Query(
		`SELECT role_id, name, description, created_at, updated_at FROM roles ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *RoleRepository) Update(role *models.Role) error {
	_, err := r.db.Exec(
		`UPDATE roles SET name = $1, description = $2, updated_at = $3 WHERE role_id = $4`,
		role.Name, role.Description, time.Now(), role.RoleID,
	)
	return err
}

func (r *RoleRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM roles WHERE role_id = $1`, id)
	return err
}
