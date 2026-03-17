package services

import (
	"lodge-system/internal/models"
	"lodge-system/internal/repository"
)

type RoleService struct {
	userRepo *repository.UserRepository
}

func NewRoleService(userRepo *repository.UserRepository) *RoleService {
	return &RoleService{userRepo: userRepo}
}

func (s *RoleService) InitializePredefinedRoles() error {
	for _, r := range models.GetPredefinedRoles() {
		_, err := s.userRepo.GetRoleByName(r.Name)
		if err != nil {
			if err := s.userRepo.CreateRole(&r); err != nil {
				return err
			}
		}
	}
	return nil
}
