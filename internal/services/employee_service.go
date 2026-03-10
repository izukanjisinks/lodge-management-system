package services

import (
	"errors"
	"fmt"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/internal/utils/email"

	"github.com/google/uuid"
)

type EmployeeService struct {
	repo         *repository.EmployeeRepository
	deptRepo     *repository.DepartmentRepository
	posRepo      *repository.PositionRepository
	userService  *UserService
	emailService *email.EmailService
}

func NewEmployeeService(
	repo *repository.EmployeeRepository,
	deptRepo *repository.DepartmentRepository,
	posRepo *repository.PositionRepository,
	userService *UserService,
	emailService *email.EmailService,
) *EmployeeService {
	return &EmployeeService{
		repo:         repo,
		deptRepo:     deptRepo,
		posRepo:      posRepo,
		userService:  userService,
		emailService: emailService,
	}
}

func (s *EmployeeService) Create(emp *models.Employee) error {
	if err := s.validateEmployee(emp, nil); err != nil {
		return err
	}

	// Generate employee number: EMP-YYYYMMDD-XXXX
	datePrefix := time.Now().Format("20060102")
	count, err := s.repo.CountByDatePrefix(datePrefix)
	if err != nil {
		return err
	}
	emp.EmployeeNumber = fmt.Sprintf("EMP-%s-%04d", datePrefix, count+1)
	emp.EmploymentStatus = models.EmploymentStatusActive

	return s.repo.Create(emp)
}

// CreateWithUser creates an employee and automatically creates a linked user account
func (s *EmployeeService) CreateWithUser(emp *models.Employee, password string) error {
	if password == "" {
		return errors.New("password is required when creating employee with user account")
	}

	// Validate employee first
	if err := s.validateEmployee(emp, nil); err != nil {
		return err
	}

	// Generate employee number: EMP-YYYYMMDD-XXXX
	datePrefix := time.Now().Format("20060102")
	count, err := s.repo.CountByDatePrefix(datePrefix)
	if err != nil {
		return err
	}
	emp.EmployeeNumber = fmt.Sprintf("EMP-%s-%04d", datePrefix, count+1)
	emp.EmploymentStatus = models.EmploymentStatusActive

	// Create employee record first
	if err := s.repo.Create(emp); err != nil {
		return err
	}

	// Get position's role_id to assign to user
	position, err := s.posRepo.GetByID(emp.PositionID)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}

	// Create user account
	user := &models.User{
		Email:          emp.Email,
		Password:       password,
		IsActive:       true,
		ChangePassword: true, // Force password change on first login
	}

	// Assign role from position, or default to employee role if position has no role
	if position.RoleID != nil {
		user.RoleID = position.RoleID
	} else {
		// Fallback to employee role if position doesn't have a role assigned
		roleRepo := repository.NewRoleRepository()
		role, err := roleRepo.GetByName(models.RoleEmployee)
		if err == nil {
			user.RoleID = &role.RoleID
		}
	}

	// Register user (this validates password against policy and hashes it)
	if err := s.userService.Register(user); err != nil {
		return fmt.Errorf("employee created but user account failed: %w", err)
	}

	// Send welcome email in background (non-blocking)
	go func() {
		welcomeHTML := email.WelcomeEmployeeTemplate(emp.FirstName, emp.LastName, emp.Email, password)

		if err := s.emailService.SendEmail([]string{emp.Email}, "Welcome to HR System", welcomeHTML); err != nil {
			fmt.Printf("Failed to send welcome email to %s: %v\n", emp.Email, err)
		}
	}()

	// Link user to employee
	emp.UserID = &user.UserID
	if err := s.repo.Update(emp); err != nil {
		return fmt.Errorf("employee and user created but linking failed: %w", err)
	}

	return nil
}

func (s *EmployeeService) GetByID(id uuid.UUID) (*models.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) GetByEmployeeNumber(number string) (*models.Employee, error) {
	return s.repo.GetByEmployeeNumber(number)
}

func (s *EmployeeService) List(filter interfaces.EmployeeFilter, page, pageSize int) ([]models.Employee, int, error) {
	return s.repo.List(filter, page, pageSize)
}

func (s *EmployeeService) Update(emp *models.Employee) error {
	existing, err := s.repo.GetByID(emp.ID)
	if err != nil {
		return errors.New("employee not found")
	}

	// Validate termination fields
	if emp.EmploymentStatus == models.EmploymentStatusTerminated ||
		emp.EmploymentStatus == models.EmploymentStatusResigned {
		if emp.TerminationDate == nil {
			return errors.New("termination_date is required when status is terminated or resigned")
		}
		if emp.TerminationReason == "" {
			return errors.New("termination_reason is required when status is terminated or resigned")
		}
	}

	// Don't allow changing employee_number
	emp.EmployeeNumber = existing.EmployeeNumber

	if err := s.validateEmployee(emp, &emp.ID); err != nil {
		return err
	}

	return s.repo.Update(emp)
}

func (s *EmployeeService) SoftDelete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("employee not found")
	}
	return s.repo.SoftDelete(id)
}

func (s *EmployeeService) GetDirectReports(managerID uuid.UUID) ([]models.Employee, error) {
	return s.repo.GetDirectReports(managerID)
}

func (s *EmployeeService) GetManagersByDepartment(departmentID uuid.UUID) ([]models.Employee, error) {
	return s.repo.GetManagersByDepartment(departmentID)
}

func (s *EmployeeService) GetOrgSubtree(rootID uuid.UUID) ([]*models.Employee, error) {
	// BFS to collect org subtree
	visited := map[uuid.UUID]bool{}
	queue := []uuid.UUID{rootID}
	var result []*models.Employee

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		reports, err := s.repo.GetDirectReports(current)
		if err != nil {
			return nil, err
		}
		for i := range reports {
			r := &reports[i]
			result = append(result, r)
			queue = append(queue, r.ID)
		}
	}
	return result, nil
}

func (s *EmployeeService) validateEmployee(emp *models.Employee, excludeID *uuid.UUID) error {
	if emp.FirstName == "" || emp.LastName == "" {
		return errors.New("first name and last name are required")
	}
	if emp.Email == "" {
		return errors.New("email is required")
	}
	if emp.HireDate.IsZero() {
		return errors.New("hire_date is required")
	}

	// Check email uniqueness
	exists, err := s.repo.EmailActiveExists(emp.Email, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already in use by another active employee")
	}

	// Verify department
	dept, err := s.deptRepo.GetByID(emp.DepartmentID)
	if err != nil || !dept.IsActive {
		return errors.New("department not found or inactive")
	}

	// Verify position
	pos, err := s.posRepo.GetByID(emp.PositionID)
	if err != nil || !pos.IsActive {
		return errors.New("position not found or inactive")
	}

	// Validate manager (no self-reference)
	if emp.ManagerID != nil && *emp.ManagerID == emp.ID {
		return errors.New("employee cannot be their own manager")
	}

	return nil
}
