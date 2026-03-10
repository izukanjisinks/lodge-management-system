package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type EmergencyContactInterface interface {
	Create(contact *models.EmergencyContact) error
	GetByID(id uuid.UUID) (*models.EmergencyContact, error)
	ListByEmployee(employeeID uuid.UUID) ([]models.EmergencyContact, error)
	Update(contact *models.EmergencyContact) error
	Delete(id uuid.UUID) error
}
