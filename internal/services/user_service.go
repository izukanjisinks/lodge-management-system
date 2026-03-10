package services

import (
	"errors"
	"fmt"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/internal/utils/email"
	"hr-system/internal/utils/password"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo                  *repository.UserRepository
	roleRepo              *repository.RoleRepository
	empRepo               *repository.EmployeeRepository
	passwordPolicyService *PasswordPolicyService
	emailService          *email.EmailService
}

func NewUserService(repo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		repo:     repo,
		roleRepo: roleRepo,
		empRepo:  repository.NewEmployeeRepository(),
	}
}

// SetEmailService sets the email service (called after initialization)
func (s *UserService) SetEmailService(emailService *email.EmailService) {
	s.emailService = emailService
}

// SetPasswordPolicyService sets the password policy service (called after initialization)
func (s *UserService) SetPasswordPolicyService(policyService *PasswordPolicyService) {
	s.passwordPolicyService = policyService
}

func (s *UserService) Register(user *models.User) error {
	exists, err := s.repo.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	// Validate password against policy if available
	if s.passwordPolicyService != nil {
		if err := s.passwordPolicyService.ValidateNewPassword(uuid.Nil, user.Password, ""); err != nil {
			return err
		}
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	// Default to employee role if none specified
	if user.RoleID == nil {
		role, err := s.roleRepo.GetByName(models.RoleEmployee)
		if err == nil {
			user.RoleID = &role.RoleID
		}
	}

	user.IsActive = true
	now := time.Now()
	user.PasswordChangedAt = &now

	// Set password expiry if policy requires it
	if s.passwordPolicyService != nil {
		user.PasswordExpiresAt = s.passwordPolicyService.CalculatePasswordExpiry()
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	// Record password in history
	if s.passwordPolicyService != nil {
		_ = s.passwordPolicyService.RecordPasswordChange(user.UserID, hashed)
	}

	return nil
}

func (s *UserService) Login(email, password string) (map[string]interface{}, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Check account lockout status
	if s.passwordPolicyService != nil {
		locked, reason := s.passwordPolicyService.CheckAccountLockout(user)
		if locked {
			return nil, errors.New(reason)
		}
	}

	// Verify password
	if err := utils.ComparePasswords(user.Password, password); err != nil {
		// Failed login attempt - increment counter
		if s.passwordPolicyService != nil {
			user.FailedLoginAttempts++

			// Check if account should be locked
			if s.passwordPolicyService.ShouldLockAccount(user.FailedLoginAttempts) {
				lockoutTime := s.passwordPolicyService.CalculateLockoutTime()
				user.LockedUntil = &lockoutTime
			}

			// Update user with failed attempt info
			_ = s.repo.Update(user)
		}
		return nil, errors.New("invalid email or password")
	}

	// Successful login - reset failed attempts and unlock if temporarily locked
	if s.passwordPolicyService != nil {
		user.FailedLoginAttempts = 0
		user.LockedUntil = nil
		_ = s.repo.Update(user)
	}

	// Check password expiry
	var passwordExpired bool
	var passwordExpiringSoon bool
	var daysUntilExpiry int
	if s.passwordPolicyService != nil {
		passwordExpired, passwordExpiringSoon, daysUntilExpiry = s.passwordPolicyService.CheckPasswordExpiry(user)
		if passwordExpired {
			return nil, errors.New("password has expired, please change your password")
		}
	}

	// Generate token with appropriate expiry
	var tokenExpiry time.Duration
	if s.passwordPolicyService != nil {
		tokenExpiry = time.Duration(s.passwordPolicyService.GetSessionTimeout()) * time.Second
	} else {
		tokenExpiry = 24 * time.Hour // Default 24 hours
	}

	token, err := utils.GenerateToken(user.Email, user.UserID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"token":      token,
		"expires_at": time.Now().Add(tokenExpiry),
		"user": map[string]interface{}{
			"user_id":         user.UserID,
			"email":           user.Email,
			"role":            user.Role,
			"is_active":       user.IsActive,
			"change_password": user.ChangePassword,
		},
	}

	// Add password expiry warnings if applicable
	if passwordExpiringSoon {
		response["password_warning"] = map[string]interface{}{
			"message":           "Your password is expiring soon",
			"days_until_expiry": daysUntilExpiry,
		}
	}

	return response, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

// ListUsers returns paginated users with optional filtering
func (s *UserService) ListUsers(search string, roleID *uuid.UUID, isActive *bool, page, pageSize int) ([]models.User, int, error) {
	return s.repo.List(search, roleID, isActive, page, pageSize)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(updates *models.User) (*models.User, error) {
	existing, err := s.repo.GetUserByID(updates.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if updates.Email != "" && updates.Email != existing.Email {
		exists, err := s.repo.EmailExists(updates.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already in use")
		}
		existing.Email = updates.Email
	}

	if updates.RoleID != nil {
		existing.RoleID = updates.RoleID
	}
	existing.IsActive = updates.IsActive

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(existing.UserID)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

// ChangeUserRole changes a user's role
func (s *UserService) ChangeUserRole(userID uuid.UUID, roleID uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify role exists
	_, err = s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.RoleID = &roleID
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.repo.GetUserByID(userID)
}

func (s *UserService) DeactivateUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	user.IsActive = false
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	_, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}

// LockUser locks a user account (admin action)
func (s *UserService) LockUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	user.IsLocked = true
	user.LockedUntil = nil // permanent lock until admin unlocks
	return s.repo.Update(user)
}

// UnlockUser unlocks a user account and resets failed login attempts
func (s *UserService) UnlockUser(id uuid.UUID) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	user.IsLocked = false
	user.LockedUntil = nil
	user.FailedLoginAttempts = 0
	return s.repo.Update(user)
}

// ChangePassword allows a user to change their password
func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := utils.ComparePasswords(user.Password, oldPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password against policy
	if s.passwordPolicyService != nil {
		if err := s.passwordPolicyService.ValidateNewPassword(userID, newPassword, user.Password); err != nil {
			return err
		}
	}

	// Hash new password
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user
	user.Password = hashed
	now := time.Now()
	user.PasswordChangedAt = &now
	user.ChangePassword = false // Clear force change flag

	// Set new expiry
	if s.passwordPolicyService != nil {
		user.PasswordExpiresAt = s.passwordPolicyService.CalculatePasswordExpiry()
	}

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Record in password history
	if s.passwordPolicyService != nil {
		_ = s.passwordPolicyService.RecordPasswordChange(userID, hashed)
	}

	return nil
}

// ResetPassword allows admin to reset a user's password with an auto-generated password and emails it to the user
func (s *UserService) ResetPassword(userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate a random password
	newPassword, err := password.GenerateTemporaryPassword()
	if err != nil {
		return fmt.Errorf("failed to generate temporary password: %w", err)
	}

	// Hash new password
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user - force password change on next login
	user.Password = hashed
	now := time.Now()
	user.PasswordChangedAt = &now
	user.ChangePassword = true // Force change on next login
	user.FailedLoginAttempts = 0
	user.IsLocked = false
	user.LockedUntil = nil

	// Set new expiry
	if s.passwordPolicyService != nil {
		user.PasswordExpiresAt = s.passwordPolicyService.CalculatePasswordExpiry()
	}

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Record in password history
	if s.passwordPolicyService != nil {
		_ = s.passwordPolicyService.RecordPasswordChange(userID, hashed)
	}

	// Send the new password to the user's email in a goroutine
	if s.emailService != nil {
		userEmail := user.Email
		go func() {
			htmlBody := email.PasswordResetTemplate(newPassword)
			if err := s.emailService.SendEmail([]string{userEmail}, "Your Password Has Been Reset", htmlBody); err != nil {
				fmt.Printf("WARNING: Failed to send password reset email to %s: %v\n", userEmail, err)
			}
		}()
	}

	return nil
}

func (s *UserService) SeedSuperAdmin(email, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	existing, err := s.repo.GetUserByEmail(email)
	if err == nil {
		// User exists — update password and ensure active
		existing.Password = hashed
		existing.IsActive = true
		now := time.Now()
		existing.PasswordChangedAt = &now
		if err := s.repo.Update(existing); err != nil {
			return err
		}
		if s.passwordPolicyService != nil {
			_ = s.passwordPolicyService.RecordPasswordChange(existing.UserID, hashed)
		}
		return nil
	}

	role, err := s.roleRepo.GetByName(models.RoleSuperAdmin)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &models.User{
		Email:             email,
		Password:          hashed,
		RoleID:            &role.RoleID,
		IsActive:          true,
		PasswordChangedAt: &now,
	}

	if s.passwordPolicyService != nil {
		user.PasswordExpiresAt = s.passwordPolicyService.CalculatePasswordExpiry()
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	if s.passwordPolicyService != nil {
		_ = s.passwordPolicyService.RecordPasswordChange(user.UserID, hashed)
	}

	return nil
}

// GetProfile retrieves the complete profile (user + employee data) for a user
func (s *UserService) GetProfile(userID uuid.UUID) (*models.ProfileResponse, error) {
	// Get user data
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Clear sensitive data
	user.Password = ""

	profile := &models.ProfileResponse{
		User: user,
	}

	// Try to get employee data (not all users have employee records)
	employee, err := s.empRepo.GetByUserID(userID)
	if err == nil {
		profile.Employee = employee
	}

	return profile, nil
}
