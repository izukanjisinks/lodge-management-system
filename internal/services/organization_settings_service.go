package services

import (
	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type OrganizationSettingsService struct {
	repo *repository.OrganizationSettingsRepository
}

func NewOrganizationSettingsService(repo *repository.OrganizationSettingsRepository) *OrganizationSettingsService {
	return &OrganizationSettingsService{repo: repo}
}

func (s *OrganizationSettingsService) Get(orgID uuid.UUID) (*models.OrganizationSettings, error) {
	return s.repo.GetForOrg(orgID)
}

func (s *OrganizationSettingsService) Upsert(orgID uuid.UUID, req *models.UpdateOrganizationSettingsRequest) (*models.OrganizationSettings, error) {
	return s.repo.Upsert(orgID, req)
}
