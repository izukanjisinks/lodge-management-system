package services

import (
	"errors"
	"fmt"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"
	"lodge-system/internal/utils/email"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type GuestAuthService struct {
	guestRepo  *repository.GuestRepository
	emailService *email.EmailService
}

func NewGuestAuthService(guestRepo *repository.GuestRepository) *GuestAuthService {
	return &GuestAuthService{guestRepo: guestRepo}
}

func (s *GuestAuthService) SetEmailService(svc *email.EmailService) {
	s.emailService = svc
}

func (s *GuestAuthService) Register(req *models.GuestRegisterRequest) (*models.Guest, error) {
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("full_name, email, and password are required")
	}

	exists, err := s.guestRepo.EmailExists(req.Email)
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

	guest := &models.Guest{
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  hashed,
		Phone:     req.Phone,
		IsActive:  true,
	}

	if err := s.guestRepo.Create(guest); err != nil {
		return nil, fmt.Errorf("failed to create guest account: %w", err)
	}

	if err := s.guestRepo.CreateIndividualProfile(guest.ID, guest); err != nil {
		return nil, fmt.Errorf("failed to create guest profile: %w", err)
	}

	if s.emailService != nil {
		go func() {
			body := email.GuestWelcomeTemplate(req.FullName)
			if sendErr := s.emailService.SendEmail([]string{req.Email}, "Welcome to The Sanctuary", body); sendErr != nil {
				fmt.Printf("warning: failed to send guest welcome email to %s: %v\n", req.Email, sendErr)
			}
		}()
	}

	return guest, nil
}

func (s *GuestAuthService) Login(emailAddr, password string) (*models.Guest, string, error) {
	guest, err := s.guestRepo.GetByEmail(emailAddr)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !guest.IsActive {
		return nil, "", errors.New("account is inactive")
	}

	if err := utils.ComparePasswords(guest.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateGuestToken(guest.Email, guest.ID)
	if err != nil {
		return nil, "", err
	}

	return guest, token, nil
}

func (s *GuestAuthService) GetByID(id uuid.UUID) (*models.Guest, error) {
	return s.guestRepo.GetByID(id)
}

func (s *GuestAuthService) GetProfileByGuestID(guestID uuid.UUID) (*models.IndividualClient, error) {
	return s.guestRepo.GetIndividualProfileByGuestID(guestID)
}

func (s *GuestAuthService) UpdateProfileIDPassport(profileID uuid.UUID, idPassport string) error {
	return s.guestRepo.UpdateIndividualProfileIDPassport(profileID, idPassport)
}

func (s *GuestAuthService) UpdateProfile(id uuid.UUID, req *models.GuestUpdateRequest) (*models.Guest, error) {
	guest, err := s.guestRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("guest not found")
	}

	if req.FullName != "" {
		guest.FullName = req.FullName
	}
	if req.Phone != "" {
		guest.Phone = req.Phone
	}

	if err := s.guestRepo.Update(guest); err != nil {
		return nil, err
	}
	return s.guestRepo.GetByID(id)
}
