package services

import (
	"errors"
	"strings"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type BranchService struct {
	repo *repository.BranchRepository
}

func NewBranchService(repo *repository.BranchRepository) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) Create(orgID uuid.UUID, req *models.CreateBranchRequest) (*models.Branch, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("branch name is required")
	}
	if strings.TrimSpace(req.BranchCode) == "" {
		return nil, errors.New("branch_code is required")
	}
	b := &models.Branch{
		OrgID:      orgID,
		Name:       strings.TrimSpace(req.Name),
		BranchCode: strings.ToUpper(strings.TrimSpace(req.BranchCode)),
		Address:    strings.TrimSpace(req.Address),
		Location:   strings.TrimSpace(req.Location),
	}
	if err := s.repo.Create(b); err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, errors.New("branch_code already exists for this organization")
		}
		return nil, err
	}
	return b, nil
}

func (s *BranchService) GetByID(id, orgID uuid.UUID) (*models.Branch, error) {
	b, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	return b, nil
}

func (s *BranchService) List(orgID uuid.UUID) ([]models.Branch, error) {
	return s.repo.List(orgID)
}

func (s *BranchService) Update(id, orgID uuid.UUID, req *models.UpdateBranchRequest) (*models.Branch, error) {
	b, err := s.repo.GetByID(id, orgID)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, errors.New("branch name cannot be empty")
		}
		b.Name = strings.TrimSpace(*req.Name)
	}
	if req.Address != nil {
		b.Address = strings.TrimSpace(*req.Address)
	}
	if req.Location != nil {
		b.Location = strings.TrimSpace(*req.Location)
	}
	if err := s.repo.Update(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BranchService) Delete(id, orgID uuid.UUID) error {
	return s.repo.Delete(id, orgID)
}
