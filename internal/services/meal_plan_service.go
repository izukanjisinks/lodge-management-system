package services

import (
	"errors"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type MealPlanService struct {
	repo *repository.MealPlanRepository
}

func NewMealPlanService(repo *repository.MealPlanRepository) *MealPlanService {
	return &MealPlanService{repo: repo}
}

func (s *MealPlanService) Create(orgID uuid.UUID, req *models.CreateMealPlanRequest) (*models.MealPlan, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.PricePerPersonPerNight < 0 {
		return nil, errors.New("price_per_person_per_night cannot be negative")
	}
	if len(req.Includes) == 0 {
		return nil, errors.New("includes must have at least one item")
	}

	m := &models.MealPlan{
		Name:                   req.Name,
		PricePerPersonPerNight: req.PricePerPersonPerNight,
		Includes:               req.Includes,
		Description:            req.Description,
		IsActive:               req.IsActive,
	}
	if err := s.repo.Create(m, orgID); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MealPlanService) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.MealPlan, error) {
	return s.repo.GetByID(id, orgID)
}

func (s *MealPlanService) GetByIDUnscoped(id uuid.UUID) (*models.MealPlan, error) {
	return s.repo.GetByIDUnscoped(id)
}

func (s *MealPlanService) GuestList(orgID *uuid.UUID, isActive *bool, page, pageSize int) ([]models.MealPlan, int, error) {
	return s.repo.GuestList(orgID, isActive, page, pageSize)
}

func (s *MealPlanService) List(orgID uuid.UUID, isActive *bool, page, pageSize int) ([]models.MealPlan, int, error) {
	return s.repo.List(orgID, isActive, page, pageSize)
}

func (s *MealPlanService) Update(id uuid.UUID, orgID uuid.UUID, req *models.UpdateMealPlanRequest) (*models.MealPlan, error) {
	m, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("meal plan not found")
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		m.Name = *req.Name
	}
	if req.PricePerPersonPerNight != nil {
		if *req.PricePerPersonPerNight < 0 {
			return nil, errors.New("price_per_person_per_night cannot be negative")
		}
		m.PricePerPersonPerNight = *req.PricePerPersonPerNight
	}
	if req.Includes != nil {
		if len(req.Includes) == 0 {
			return nil, errors.New("includes must have at least one item")
		}
		m.Includes = req.Includes
	}
	if req.Description != nil {
		m.Description = *req.Description
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	}

	if err := s.repo.Update(m, orgID); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id, orgID)
}

func (s *MealPlanService) Delete(id uuid.UUID, orgID uuid.UUID) error {
	_, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return errors.New("meal plan not found")
	}
	return s.repo.Delete(id, orgID)
}
