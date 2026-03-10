package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type LeaveTypeInterface interface {
	Create(lt *models.LeaveType) error
	GetByID(id uuid.UUID) (*models.LeaveType, error)
	GetByCode(code string) (*models.LeaveType, error)
	List(activeOnly bool) ([]models.LeaveType, error)
	Update(lt *models.LeaveType) error
	SeedDefaults() error
}
