package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type EmployeeFilter struct {
	Search           string
	DepartmentID     *uuid.UUID
	PositionID       *uuid.UUID
	EmploymentStatus string
	EmploymentType   string
	IncludeDeleted   bool
}

type EmployeeInterface interface {
	Create(emp *models.Employee) error
	GetByID(id uuid.UUID) (*models.Employee, error)
	GetByEmployeeNumber(number string) (*models.Employee, error)
	List(filter EmployeeFilter, page, pageSize int) ([]models.Employee, int, error)
	Update(emp *models.Employee) error
	SoftDelete(id uuid.UUID) error
	GetDirectReports(managerID uuid.UUID) ([]models.Employee, error)
	GetOrgSubtree(rootID uuid.UUID) ([]*models.Employee, error)
}
