package services

import (
	"errors"
	"strings"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type PositionService struct {
	repo       *repository.PositionRepository
	deptRepo   *repository.DepartmentRepository
}

func NewPositionService(repo *repository.PositionRepository, deptRepo *repository.DepartmentRepository) *PositionService {
	return &PositionService{repo: repo, deptRepo: deptRepo}
}

func (s *PositionService) Create(pos *models.Position) error {
	pos.Code = strings.ToUpper(strings.TrimSpace(pos.Code))
	if pos.Code == "" {
		return errors.New("position code is required")
	}
	if pos.Title == "" {
		return errors.New("position title is required")
	}

	// Verify department exists
	dept, err := s.deptRepo.GetByID(pos.DepartmentID)
	if err != nil || !dept.IsActive {
		return errors.New("department not found or inactive")
	}

	exists, err := s.repo.CodeExists(pos.Code, nil)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("position code already in use")
	}

	pos.IsActive = true
	pos.CalculateSalaryComponents()
	return s.repo.Create(pos)
}

func (s *PositionService) GetByID(id uuid.UUID) (*models.Position, error) {
	return s.repo.GetByID(id)
}

func (s *PositionService) List(filter interfaces.PositionFilter, page, pageSize int) ([]models.Position, int, error) {
	return s.repo.List(filter, page, pageSize)
}

func (s *PositionService) Update(pos *models.Position) error {
	pos.Code = strings.ToUpper(strings.TrimSpace(pos.Code))

	existing, err := s.repo.GetByID(pos.ID)
	if err != nil {
		return errors.New("position not found")
	}

	if pos.Code != existing.Code {
		exists, err := s.repo.CodeExists(pos.Code, &pos.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("position code already in use")
		}
	}

	pos.CalculateSalaryComponents()
	return s.repo.Update(pos)
}

func (s *PositionService) SoftDelete(id uuid.UUID) error {
	count, err := s.repo.ActiveEmployeeCount(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete position with active employees")
	}
	return s.repo.SoftDelete(id)
}
