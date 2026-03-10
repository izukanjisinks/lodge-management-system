package services

import (
	"errors"

	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

var allowedMimeTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

const maxFileSize = 10 * 1024 * 1024 // 10 MB

type EmployeeDocumentService struct {
	repo     *repository.EmployeeDocumentRepository
	empRepo  *repository.EmployeeRepository
}

func NewEmployeeDocumentService(repo *repository.EmployeeDocumentRepository, empRepo *repository.EmployeeRepository) *EmployeeDocumentService {
	return &EmployeeDocumentService{repo: repo, empRepo: empRepo}
}

func (s *EmployeeDocumentService) Create(doc *models.EmployeeDocument) error {
	if _, err := s.empRepo.GetByID(doc.EmployeeID); err != nil {
		return errors.New("employee not found")
	}
	if doc.FileURL == "" {
		return errors.New("file_url is required")
	}
	if doc.FileSize > maxFileSize {
		return errors.New("file size exceeds 10MB limit")
	}
	if doc.MimeType != "" && !allowedMimeTypes[doc.MimeType] {
		return errors.New("file type not allowed; use PDF, JPEG, PNG, or DOCX")
	}
	return s.repo.Create(doc)
}

func (s *EmployeeDocumentService) GetByID(id uuid.UUID) (*models.EmployeeDocument, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeDocumentService) ListByEmployee(employeeID uuid.UUID, docType string) ([]models.EmployeeDocument, error) {
	return s.repo.ListByEmployee(employeeID, docType)
}

func (s *EmployeeDocumentService) Verify(id uuid.UUID, verifiedBy uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("document not found")
	}
	return s.repo.Verify(id, verifiedBy)
}

func (s *EmployeeDocumentService) SoftDelete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("document not found")
	}
	return s.repo.SoftDelete(id)
}
