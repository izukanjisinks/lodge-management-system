package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type DepartmentFilter struct {
	IsActive *bool
	Search   string
}

type DepartmentInterface interface {
	Create(dept *models.Department) error
	GetByID(id uuid.UUID) (*models.Department, error)
	List(filter DepartmentFilter, page, pageSize int) ([]models.Department, int, error)
	Update(dept *models.Department) error
	SoftDelete(id uuid.UUID) error
	GetTree() ([]*models.Department, error)
	GetEmployeeCount(id uuid.UUID) (int, error)
}
