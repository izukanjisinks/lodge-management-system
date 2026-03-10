package services

import (
	"errors"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type LeaveBalanceService struct {
	repo         *repository.LeaveBalanceRepository
	leaveTypeRepo *repository.LeaveTypeRepository
	empRepo      *repository.EmployeeRepository
}

func NewLeaveBalanceService(
	repo *repository.LeaveBalanceRepository,
	ltRepo *repository.LeaveTypeRepository,
	empRepo *repository.EmployeeRepository,
) *LeaveBalanceService {
	return &LeaveBalanceService{repo: repo, leaveTypeRepo: ltRepo, empRepo: empRepo}
}

func (s *LeaveBalanceService) GetByEmployeeAndYear(employeeID uuid.UUID, year int) ([]models.LeaveBalance, error) {
	return s.repo.GetByEmployeeAndYear(employeeID, year)
}

func (s *LeaveBalanceService) GetByEmployeeTypeYear(employeeID, leaveTypeID uuid.UUID, year int) (*models.LeaveBalance, error) {
	return s.repo.GetByEmployeeTypeYear(employeeID, leaveTypeID, year)
}

// InitializeForEmployee creates leave balance records for all active leave types for a given year.
// It prorates entitlement for new employees hired mid-year.
func (s *LeaveBalanceService) InitializeForEmployee(employeeID uuid.UUID, year int) error {
	emp, err := s.empRepo.GetByID(employeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	leaveTypes, err := s.leaveTypeRepo.List(true)
	if err != nil {
		return err
	}

	for _, lt := range leaveTypes {
		entitled := lt.DefaultDaysPerYear
		// Prorate if hired in the current year
		if emp.HireDate.Year() == year {
			entitled = ProrateEntitlement(lt.DefaultDaysPerYear, emp.HireDate, year)
		}

		lb := &models.LeaveBalance{
			EmployeeID:    employeeID,
			LeaveTypeID:   lt.ID,
			Year:          year,
			TotalEntitled: entitled,
		}
		if err := s.repo.Upsert(lb); err != nil {
			return err
		}
	}
	return nil
}

func (s *LeaveBalanceService) Adjust(input interfaces.AdjustBalanceInput) error {
	if input.Reason == "" {
		return errors.New("adjustment reason is required")
	}
	return s.repo.Adjust(input.LeaveBalanceID, input.Delta)
}

func (s *LeaveBalanceService) IncrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	return s.repo.IncrementPending(employeeID, leaveTypeID, year, days)
}

func (s *LeaveBalanceService) DecrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	return s.repo.DecrementPending(employeeID, leaveTypeID, year, days)
}

func (s *LeaveBalanceService) ApproveLeave(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	return s.repo.ApproveLeave(employeeID, leaveTypeID, year, days)
}

// HasSufficientBalance checks if the employee has enough balance for the given days.
func (s *LeaveBalanceService) HasSufficientBalance(employeeID, leaveTypeID uuid.UUID, year, days int) (bool, error) {
	lb, err := s.repo.GetByEmployeeTypeYear(employeeID, leaveTypeID, year)
	if err != nil {
		// Auto-initialize if missing
		if initErr := s.InitializeForEmployee(employeeID, time.Now().Year()); initErr != nil {
			return false, initErr
		}
		lb, err = s.repo.GetByEmployeeTypeYear(employeeID, leaveTypeID, year)
		if err != nil {
			return false, err
		}
	}
	return lb.Balance >= days, nil
}
