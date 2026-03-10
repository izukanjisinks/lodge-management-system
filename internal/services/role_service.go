package services

import (
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) InitializePredefinedRoles() error {
	for _, r := range models.GetPredefinedRoles() {
		_, err := s.repo.GetByName(r.Name)
		if err != nil {
			// Role doesn't exist — create it
			if err := s.repo.Create(&r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *RoleService) CreateRole(role *models.Role) error {
	return s.repo.Create(role)
}

func (s *RoleService) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	return s.repo.GetByID(id)
}

func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	return s.repo.GetByName(name)
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) UpdateRole(role *models.Role) error {
	return s.repo.Update(role)
}

func (s *RoleService) DeleteRole(id uuid.UUID) error {
	return s.repo.Delete(id)
}
