package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type EmployeeDocumentInterface interface {
	Create(doc *models.EmployeeDocument) error
	GetByID(id uuid.UUID) (*models.EmployeeDocument, error)
	ListByEmployee(employeeID uuid.UUID, docType string) ([]models.EmployeeDocument, error)
	Verify(id uuid.UUID, verifiedBy uuid.UUID) error
	SoftDelete(id uuid.UUID) error
}
