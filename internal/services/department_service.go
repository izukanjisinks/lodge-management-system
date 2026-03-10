package services

import (
	"errors"
	"strings"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type DepartmentService struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) Create(dept *models.Department) error {
	dept.Code = strings.ToUpper(strings.TrimSpace(dept.Code))
	if dept.Code == "" {
		return errors.New("department code is required")
	}
	if dept.Name == "" {
		return errors.New("department name is required")
	}

	exists, err := s.repo.CodeExists(dept.Code, nil)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("department code already in use")
	}

	dept.IsActive = true
	return s.repo.Create(dept)
}

func (s *DepartmentService) GetByID(id uuid.UUID) (*models.Department, error) {
	return s.repo.GetByID(id)
}

func (s *DepartmentService) List(filter interfaces.DepartmentFilter, page, pageSize int) ([]models.Department, int, error) {
	return s.repo.List(filter, page, pageSize)
}

func (s *DepartmentService) Update(dept *models.Department) error {
	dept.Code = strings.ToUpper(strings.TrimSpace(dept.Code))

	existing, err := s.repo.GetByID(dept.ID)
	if err != nil {
		return errors.New("department not found")
	}

	if dept.Code != existing.Code {
		exists, err := s.repo.CodeExists(dept.Code, &dept.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("department code already in use")
		}
	}

	return s.repo.Update(dept)
}

func (s *DepartmentService) SoftDelete(id uuid.UUID) error {
	count, err := s.repo.GetEmployeeCount(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete department with active employees")
	}
	return s.repo.SoftDelete(id)
}

func (s *DepartmentService) GetTree() ([]*models.Department, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Build tree from flat list
	byID := make(map[uuid.UUID]*models.Department)
	for i := range all {
		d := &all[i]
		d.Children = []*models.Department{}
		byID[d.ID] = d
	}

	var roots []*models.Department
	for i := range all {
		d := &all[i]
		if d.ParentDepartmentID == nil {
			roots = append(roots, d)
		} else if parent, ok := byID[*d.ParentDepartmentID]; ok {
			parent.Children = append(parent.Children, d)
		}
	}
	return roots, nil
}

func (s *DepartmentService) GetEmployeeCount(id uuid.UUID) (int, error) {
	return s.repo.GetEmployeeCount(id)
}
