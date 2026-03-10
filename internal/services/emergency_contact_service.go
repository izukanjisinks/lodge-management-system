package services

import (
	"errors"

	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type EmergencyContactService struct {
	repo    *repository.EmergencyContactRepository
	empRepo *repository.EmployeeRepository
}

func NewEmergencyContactService(repo *repository.EmergencyContactRepository, empRepo *repository.EmployeeRepository) *EmergencyContactService {
	return &EmergencyContactService{repo: repo, empRepo: empRepo}
}

func (s *EmergencyContactService) Create(c *models.EmergencyContact) error {
	if _, err := s.empRepo.GetByID(c.EmployeeID); err != nil {
		return errors.New("employee not found")
	}
	if c.Name == "" || c.Phone == "" {
		return errors.New("name and phone are required")
	}
	return s.repo.Create(c)
}

func (s *EmergencyContactService) GetByID(id uuid.UUID) (*models.EmergencyContact, error) {
	return s.repo.GetByID(id)
}

func (s *EmergencyContactService) ListByEmployee(employeeID uuid.UUID) ([]models.EmergencyContact, error) {
	return s.repo.ListByEmployee(employeeID)
}

func (s *EmergencyContactService) Update(c *models.EmergencyContact) error {
	existing, err := s.repo.GetByID(c.ID)
	if err != nil {
		return errors.New("emergency contact not found")
	}
	c.EmployeeID = existing.EmployeeID // prevent changing employee
	return s.repo.Update(c)
}

func (s *EmergencyContactService) Delete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("emergency contact not found")
	}
	return s.repo.Delete(id)
}
