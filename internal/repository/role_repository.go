package repository

import (
	"database/sql"

	"lodge-system/internal/database"
	"lodge-system/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.DB}
}

func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRow(
		`SELECT role_id, name, description, created_at, updated_at FROM roles WHERE name = $1`, name,
	).Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
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
