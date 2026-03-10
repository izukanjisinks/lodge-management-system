package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type RoleInterface interface {
	CreateRole(role *models.Role) error
	GetRoleByID(id uuid.UUID) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uuid.UUID) error
	InitializePredefinedRoles() error
}
