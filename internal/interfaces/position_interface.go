package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type PositionFilter struct {
	DepartmentID *uuid.UUID
	GradeLevel   string
	IsActive     *bool
}

type PositionInterface interface {
	Create(pos *models.Position) error
	GetByID(id uuid.UUID) (*models.Position, error)
	List(filter PositionFilter, page, pageSize int) ([]models.Position, int, error)
	Update(pos *models.Position) error
	SoftDelete(id uuid.UUID) error
}
