package services

import (
	"errors"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

// ─── Individual ───────────────────────────────────────────────────────────────

func (s *ClientService) CreateIndividual(c *models.IndividualClient) error {
	if c.FullName == "" {
		return errors.New("full_name is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.Phone == "" {
		return errors.New("phone is required")
	}
	if c.IDPassportNumber == "" {
		return errors.New("id_passport_number is required")
	}
	if c.Status == "" {
		c.Status = models.ClientStatusActive
	}
	return s.repo.CreateIndividual(c)
}

func (s *ClientService) GetIndividualByID(id uuid.UUID) (*models.IndividualClient, error) {
	return s.repo.GetIndividualByID(id)
}

func (s *ClientService) ListIndividual(search, status string, page, pageSize int) ([]models.IndividualClient, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.repo.ListIndividual(search, status, page, pageSize)
}

func (s *ClientService) UpdateIndividual(id uuid.UUID, updates *models.IndividualClient) (*models.IndividualClient, error) {
	existing, err := s.repo.GetIndividualByID(id)
	if err != nil {
		return nil, errors.New("individual client not found")
	}
	if updates.FullName != "" {
		existing.FullName = updates.FullName
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Phone != "" {
		existing.Phone = updates.Phone
	}
	if updates.IDPassportNumber != "" {
		existing.IDPassportNumber = updates.IDPassportNumber
	}
	if updates.Nationality != "" {
		existing.Nationality = updates.Nationality
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}
	existing.Notes = updates.Notes
	if err := s.repo.UpdateIndividual(existing); err != nil {
		return nil, err
	}
	return s.repo.GetIndividualByID(id)
}

func (s *ClientService) DeleteIndividual(id uuid.UUID) error {
	return s.repo.DeleteIndividual(id)
}

// ─── Corporate ────────────────────────────────────────────────────────────────

func (s *ClientService) CreateCorporate(c *models.CorporateClient) error {
	if c.CompanyName == "" {
		return errors.New("company_name is required")
	}
	if c.ContactPerson == "" {
		return errors.New("contact_person is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.Phone == "" {
		return errors.New("phone is required")
	}
	if c.CompanyRegNumber == "" {
		return errors.New("company_reg_number is required")
	}
	if c.Status == "" {
		c.Status = models.ClientStatusActive
	}
	return s.repo.CreateCorporate(c)
}

func (s *ClientService) GetCorporateByID(id uuid.UUID) (*models.CorporateClient, error) {
	return s.repo.GetCorporateByID(id)
}

func (s *ClientService) ListCorporate(search, status string, page, pageSize int) ([]models.CorporateClient, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.repo.ListCorporate(search, status, page, pageSize)
}

func (s *ClientService) UpdateCorporate(id uuid.UUID, updates *models.CorporateClient) (*models.CorporateClient, error) {
	existing, err := s.repo.GetCorporateByID(id)
	if err != nil {
		return nil, errors.New("corporate client not found")
	}
	if updates.CompanyName != "" {
		existing.CompanyName = updates.CompanyName
	}
	if updates.ContactPerson != "" {
		existing.ContactPerson = updates.ContactPerson
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Phone != "" {
		existing.Phone = updates.Phone
	}
	if updates.CompanyRegNumber != "" {
		existing.CompanyRegNumber = updates.CompanyRegNumber
	}
	if updates.Industry != "" {
		existing.Industry = updates.Industry
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}
	existing.Notes = updates.Notes
	if err := s.repo.UpdateCorporate(existing); err != nil {
		return nil, err
	}
	return s.repo.GetCorporateByID(id)
}

func (s *ClientService) DeleteCorporate(id uuid.UUID) error {
	return s.repo.DeleteCorporate(id)
}
