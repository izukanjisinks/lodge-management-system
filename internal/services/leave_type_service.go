package services

import (
	"errors"
	"strings"

	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type LeaveTypeService struct {
	repo *repository.LeaveTypeRepository
}

func NewLeaveTypeService(repo *repository.LeaveTypeRepository) *LeaveTypeService {
	return &LeaveTypeService{repo: repo}
}

func (s *LeaveTypeService) SeedDefaults() error {
	for _, lt := range models.DefaultLeaveTypes() {
		_, err := s.repo.GetByCode(lt.Code)
		if err != nil {
			lt.IsActive = true
			if err := s.repo.Create(&lt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *LeaveTypeService) Create(lt *models.LeaveType) error {
	lt.Code = strings.ToUpper(strings.TrimSpace(lt.Code))
	if lt.Code == "" || lt.Name == "" {
		return errors.New("code and name are required")
	}
	exists, err := s.repo.CodeExists(lt.Code, nil)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("leave type code already exists")
	}
	lt.IsActive = true
	return s.repo.Create(lt)
}

func (s *LeaveTypeService) GetByID(id uuid.UUID) (*models.LeaveType, error) {
	return s.repo.GetByID(id)
}

func (s *LeaveTypeService) List(activeOnly bool) ([]models.LeaveType, error) {
	return s.repo.List(activeOnly)
}

func (s *LeaveTypeService) Update(lt *models.LeaveType) error {
	lt.Code = strings.ToUpper(strings.TrimSpace(lt.Code))
	existing, err := s.repo.GetByID(lt.ID)
	if err != nil {
		return errors.New("leave type not found")
	}
	if lt.Code != existing.Code {
		exists, err := s.repo.CodeExists(lt.Code, &lt.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("leave type code already exists")
		}
	}
	return s.repo.Update(lt)
}
