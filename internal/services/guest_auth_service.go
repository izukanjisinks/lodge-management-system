package services

import (
	"errors"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"
	"lodge-system/internal/repository"
	"lodge-system/internal/utils/email"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type GuestRegisterRequest struct {
	FullName         string `json:"full_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Phone            string `json:"phone"`
	IDPassportNumber string `json:"id_passport_number"`
	Nationality      string `json:"nationality"`
}

type GuestUpdateProfileRequest struct {
	FullName         string `json:"full_name"`
	Phone            string `json:"phone"`
	IDPassportNumber string `json:"id_passport_number"`
	Nationality      string `json:"nationality"`
}

type GuestAuthService struct {
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	clientRepo   *repository.ClientRepository
	emailService *email.EmailService
}

func NewGuestAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	clientRepo *repository.ClientRepository,
) *GuestAuthService {
	return &GuestAuthService{userRepo: userRepo, roleRepo: roleRepo, clientRepo: clientRepo}
}

func (s *GuestAuthService) SetEmailService(emailService *email.EmailService) {
	s.emailService = emailService
}

// Register creates a users row (role=guest) and a linked individual_profiles row atomically.
func (s *GuestAuthService) Register(req *GuestRegisterRequest) (*models.User, error) {
	if req.FullName == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		return nil, errors.New("full_name, email, password, and phone are required")
	}

	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetRoleByName(models.RoleGuest)
	if err != nil {
		return nil, fmt.Errorf("guest role not found: %w", err)
	}

	now := time.Now()
	user := &models.User{
		FullName:          req.FullName,
		Email:             req.Email,
		Password:          hashed,
		RoleID:            &role.RoleID,
		RoleName:          models.RoleGuest,
		Role:              role,
		IsActive:          true,
		PasswordChangedAt: &now,
	}

	profile := &models.IndividualClient{
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		IDPassportNumber: req.IDPassportNumber,
		Nationality:      req.Nationality,
		Status:           models.ClientStatusActive,
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.userRepo.CreateTx(tx, user); err != nil {
		return nil, fmt.Errorf("failed to create user account: %w", err)
	}

	if err = s.clientRepo.CreateIndividualTx(tx, profile, user.UserID); err != nil {
		return nil, fmt.Errorf("failed to create guest profile: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if s.emailService != nil {
		go func() {
			body := email.GuestWelcomeTemplate(req.FullName)
			if sendErr := s.emailService.SendEmail([]string{req.Email}, "Welcome to The Sanctuary", body); sendErr != nil {
				fmt.Printf("warning: failed to send guest welcome email to %s: %v\n", req.Email, sendErr)
			}
		}()
	}

	return user, nil
}

// GetProfileByUserID resolves the individual_profiles record for a logged-in guest.
func (s *GuestAuthService) GetProfileByUserID(userID uuid.UUID) (*models.IndividualClient, error) {
	return s.clientRepo.GetIndividualByUserID(userID)
}

// UpdateProfile updates the editable fields on a guest's individual profile.
func (s *GuestAuthService) UpdateProfile(userID uuid.UUID, req *GuestUpdateProfileRequest) (*models.IndividualClient, error) {
	profile, err := s.clientRepo.GetIndividualByUserID(userID)
	if err != nil {
		return nil, errors.New("guest profile not found")
	}

	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.IDPassportNumber != "" {
		profile.IDPassportNumber = req.IDPassportNumber
	}
	if req.Nationality != "" {
		profile.Nationality = req.Nationality
	}

	if err := s.clientRepo.UpdateIndividualByUserID(userID, profile); err != nil {
		return nil, err
	}
	return s.clientRepo.GetIndividualByUserID(userID)
}
