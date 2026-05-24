package services

import (
	"errors"
	"fmt"

	"lodge-system/internal/database"
	"lodge-system/internal/models"
	"lodge-system/internal/repository"
	"lodge-system/internal/utils/email"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type BackofficeOrganizationService struct {
	orgRepo      *repository.OrganizationRepository
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	emailService *email.EmailService
}

func NewBackofficeOrganizationService(
	orgRepo *repository.OrganizationRepository,
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
) *BackofficeOrganizationService {
	return &BackofficeOrganizationService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *BackofficeOrganizationService) SetEmailService(svc *email.EmailService) {
	s.emailService = svc
}

func (s *BackofficeOrganizationService) List() ([]models.Organization, error) {
	return s.orgRepo.List()
}

func (s *BackofficeOrganizationService) GetByID(id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.GetByID(id)
}

func (s *BackofficeOrganizationService) Update(id uuid.UUID, req models.OrgDetails) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("organization not found")
	}
	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Email != "" {
		org.Email = req.Email
	}
	if req.Phone != "" {
		org.Phone = req.Phone
	}
	if req.StreetAddress != "" {
		org.StreetAddress = req.StreetAddress
	}
	if req.City != "" {
		org.City = req.City
	}
	if req.Country != "" {
		org.Country = req.Country
	}
	if req.LogoURL != "" {
		org.LogoURL = req.LogoURL
	}
	if err := s.orgRepo.Update(org); err != nil {
		return nil, err
	}
	return s.orgRepo.GetByID(id)
}

func (s *BackofficeOrganizationService) Delete(id uuid.UUID) error {
	return s.orgRepo.Delete(id)
}

// Provision creates an organization and its first admin user in a single transaction.
// The admin password is randomly generated and emailed to the admin.
func (s *BackofficeOrganizationService) Provision(req models.ProvisionOrgRequest) (*models.Organization, *models.User, error) {
	if req.Organization.Name == "" {
		return nil, nil, errors.New("organization name is required")
	}
	if req.Admin.FullName == "" || req.Admin.Email == "" {
		return nil, nil, errors.New("admin full_name and email are required")
	}

	adminRole, err := s.roleRepo.GetRoleByName(models.RoleAdmin)
	if err != nil {
		return nil, nil, fmt.Errorf("admin role not found: %w", err)
	}

	password, err := utils.GenerateRandomPassword()
	if err != nil {
		return nil, nil, err
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	org := &models.Organization{
		Name:          req.Organization.Name,
		Email:         req.Organization.Email,
		Phone:         req.Organization.Phone,
		StreetAddress: req.Organization.StreetAddress,
		City:          req.Organization.City,
		Country:       req.Organization.Country,
		LogoURL:       req.Organization.LogoURL,
	}
	if err = s.orgRepo.CreateTx(tx, org); err != nil {
		return nil, nil, fmt.Errorf("failed to create organization: %w", err)
	}

	admin := &models.User{
		FullName:       req.Admin.FullName,
		Email:          req.Admin.Email,
		Password:       hashed,
		RoleID:         &adminRole.RoleID,
		Role:           adminRole,
		IsActive:       true,
		ChangePassword: true,
		OrgID:          &org.ID,
	}
	if err = s.userRepo.CreateTx(tx, admin); err != nil {
		return nil, nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}

	if s.emailService != nil {
		go func() {
			body := email.WelcomeUserTemplate(req.Admin.FullName, req.Admin.Email, password)
			if sendErr := s.emailService.SendEmail(
				[]string{req.Admin.Email},
				fmt.Sprintf("Your admin account for %s", req.Organization.Name),
				body,
			); sendErr != nil {
				fmt.Printf("warning: failed to send admin welcome email to %s: %v\n", req.Admin.Email, sendErr)
			}
		}()
	}

	return org, admin, nil
}
