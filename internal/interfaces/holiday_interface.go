package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type HolidayInterface interface {
	Create(h *models.Holiday) error
	GetByID(id uuid.UUID) (*models.Holiday, error)
	List(year int, location string) ([]models.Holiday, error)
	Update(h *models.Holiday) error
	Delete(id uuid.UUID) error
	IsHoliday(date string, location string) (bool, error)
}
