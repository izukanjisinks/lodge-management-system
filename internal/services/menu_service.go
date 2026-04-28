package services

import (
	"errors"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type MenuService struct {
	repo *repository.MenuRepository
}

func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

// ── Menus ─────────────────────────────────────────────────────────────────────

func (s *MenuService) CreateMenu(orgID uuid.UUID, req *models.CreateMenuRequest) (*models.Menu, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	m := &models.Menu{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}
	if err := s.repo.CreateMenu(m, orgID); err != nil {
		return nil, err
	}
	return s.repo.GetMenuByID(m.ID, orgID)
}

func (s *MenuService) GetMenuByID(id uuid.UUID, orgID uuid.UUID) (*models.Menu, error) {
	m, err := s.repo.GetMenuByID(id, orgID)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	return m, nil
}

func (s *MenuService) ListMenus(orgID uuid.UUID, page, pageSize int) ([]models.Menu, int, error) {
	return s.repo.ListMenus(orgID, page, pageSize)
}

func (s *MenuService) UpdateMenu(id uuid.UUID, orgID uuid.UUID, req *models.UpdateMenuRequest) (*models.Menu, error) {
	m, err := s.repo.UpdateMenu(id, orgID, req)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	return m, nil
}

func (s *MenuService) DeleteMenu(id uuid.UUID, orgID uuid.UUID) error {
	if err := s.repo.DeleteMenu(id, orgID); err != nil {
		return errors.New("menu not found")
	}
	return nil
}

// ── Menu Items ────────────────────────────────────────────────────────────────

func (s *MenuService) CreateMenuItem(menuID uuid.UUID, orgID uuid.UUID, req *models.CreateMenuItemRequest) (*models.MenuItem, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Price < 0 {
		return nil, errors.New("price must be >= 0")
	}
	// Verify menu belongs to org
	if _, err := s.repo.GetMenuByID(menuID, orgID); err != nil {
		return nil, errors.New("menu not found")
	}
	item := &models.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: true,
	}
	if err := s.repo.CreateMenuItem(item, menuID, orgID); err != nil {
		return nil, err
	}
	return s.repo.GetMenuItemByID(item.ID, orgID)
}

func (s *MenuService) GetMenuItemByID(id uuid.UUID, orgID uuid.UUID) (*models.MenuItem, error) {
	item, err := s.repo.GetMenuItemByID(id, orgID)
	if err != nil {
		return nil, errors.New("menu item not found")
	}
	return item, nil
}

func (s *MenuService) UpdateMenuItem(id uuid.UUID, orgID uuid.UUID, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	if req.Price != nil && *req.Price < 0 {
		return nil, errors.New("price must be >= 0")
	}
	item, err := s.repo.UpdateMenuItem(id, orgID, req)
	if err != nil {
		return nil, errors.New("menu item not found")
	}
	return item, nil
}

func (s *MenuService) DeleteMenuItem(id uuid.UUID, orgID uuid.UUID) error {
	if err := s.repo.DeleteMenuItem(id, orgID); err != nil {
		return errors.New("menu item not found")
	}
	return nil
}

// ── Guest (public) ────────────────────────────────────────────────────────────

func (s *MenuService) GuestListMenus(orgID *uuid.UUID, page, pageSize int) ([]models.Menu, int, error) {
	return s.repo.GuestListMenus(orgID, page, pageSize)
}
