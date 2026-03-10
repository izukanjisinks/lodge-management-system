package services

import (
	"errors"
	"fmt"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type PayslipService struct {
	repo    *repository.PayslipRepository
	empRepo *repository.EmployeeRepository
	posRepo *repository.PositionRepository
	lbRepo  *repository.LeaveBalanceRepository
}

func NewPayslipService(
	repo *repository.PayslipRepository,
	empRepo *repository.EmployeeRepository,
	posRepo *repository.PositionRepository,
	lbRepo *repository.LeaveBalanceRepository,
) *PayslipService {
	return &PayslipService{
		repo:    repo,
		empRepo: empRepo,
		posRepo: posRepo,
		lbRepo:  lbRepo,
	}
}

// Generate creates a payslip for an employee for the given month/year.
// It pulls salary data from the employee's position and unused leave days from leave balances.
func (s *PayslipService) Generate(employeeID uuid.UUID, month, year int) (*models.Payslip, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	// Check if payslip already exists for this period
	existing, err := s.repo.GetByEmployeeAndPeriod(employeeID, month, year)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("payslip already exists for %s %d", time.Month(month), year)
	}

	// Get employee
	emp, err := s.empRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	if emp.EmploymentStatus != models.EmploymentStatusActive {
		return nil, errors.New("employee is not active")
	}

	// Get position for salary data
	pos, err := s.posRepo.GetByID(emp.PositionID)
	if err != nil {
		return nil, errors.New("employee position not found")
	}

	// Calculate salary breakdown from position's base salary
	breakdown := utils.CalculateSalaryBreakdown(pos.BaseSalary)

	// Calculate leave days compensation (unused leave days × fixed rate)
	leaveDaysAmount := s.calculateLeaveDaysCompensation(employeeID, year)

	// Build payslip
	netSalary := breakdown.NetSalary + leaveDaysAmount

	payslip := &models.Payslip{
		EmployeeID:         employeeID,
		Month:              month,
		Year:               year,
		BaseSalary:         breakdown.BaseSalary,
		HousingAllowance:   breakdown.HousingAllowance,
		TransportAllowance: breakdown.TransportAllowance,
		MedicalAllowance:   breakdown.MedicalAllowance,
		GrossSalary:        breakdown.GrossSalary,
		IncomeTax:          breakdown.IncomeTax,
		LeaveDays:          leaveDaysAmount,
		NetSalary:          netSalary,
	}

	if err := s.repo.Create(payslip); err != nil {
		return nil, fmt.Errorf("failed to create payslip: %w", err)
	}

	// Re-fetch to populate relations
	return s.repo.GetByID(payslip.ID)
}

// calculateLeaveDaysCompensation computes compensation for unused leave days
func (s *PayslipService) calculateLeaveDaysCompensation(employeeID uuid.UUID, year int) float64 {
	balances, err := s.lbRepo.GetByEmployeeAndYear(employeeID, year)
	if err != nil {
		return 0
	}

	totalUnused := 0
	for _, b := range balances {
		if b.Balance > 0 {
			totalUnused += b.EarnedLeaveDays
		}
	}

	return float64(totalUnused) * utils.LeaveDayRate
}

func (s *PayslipService) GetByID(id uuid.UUID) (*models.Payslip, error) {
	return s.repo.GetByID(id)
}

func (s *PayslipService) GetByEmployeeAndPeriod(employeeID uuid.UUID, month, year int) (*models.Payslip, error) {
	return s.repo.GetByEmployeeAndPeriod(employeeID, month, year)
}

func (s *PayslipService) List(employeeID *uuid.UUID, month *int, year *int, page, pageSize int) ([]models.Payslip, int, error) {
	return s.repo.List(employeeID, month, year, page, pageSize)
}

func (s *PayslipService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// GetEmployeeByUserID returns the employee record linked to the given user ID
func (s *PayslipService) GetEmployeeByUserID(userID uuid.UUID) (*models.Employee, error) {
	return s.empRepo.GetByUserID(userID)
}
