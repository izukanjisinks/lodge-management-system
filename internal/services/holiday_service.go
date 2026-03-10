package services

import (
	"errors"

	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type HolidayService struct {
	repo *repository.HolidayRepository
}

func NewHolidayService(repo *repository.HolidayRepository) *HolidayService {
	return &HolidayService{repo: repo}
}

func (s *HolidayService) Create(h *models.Holiday) error {
	if h.Name == "" {
		return errors.New("holiday name is required")
	}
	if h.Date.IsZero() {
		return errors.New("holiday date is required")
	}
	h.IsActive = true
	return s.repo.Create(h)
}

func (s *HolidayService) GetByID(id uuid.UUID) (*models.Holiday, error) {
	return s.repo.GetByID(id)
}

func (s *HolidayService) List(year int, location string) ([]models.Holiday, error) {
	return s.repo.List(year, location)
}

func (s *HolidayService) Update(h *models.Holiday) error {
	_, err := s.repo.GetByID(h.ID)
	if err != nil {
		return errors.New("holiday not found")
	}
	return s.repo.Update(h)
}

func (s *HolidayService) Delete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("holiday not found")
	}
	return s.repo.Delete(id)
}
