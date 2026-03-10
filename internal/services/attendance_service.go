package services

import (
	"errors"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type AttendanceService struct {
	repo        *repository.AttendanceRepository
	holidayRepo *repository.HolidayRepository
	empRepo     *repository.EmployeeRepository
}

func NewAttendanceService(
	repo *repository.AttendanceRepository,
	holidayRepo *repository.HolidayRepository,
	empRepo *repository.EmployeeRepository,
) *AttendanceService {
	return &AttendanceService{repo: repo, holidayRepo: holidayRepo, empRepo: empRepo}
}

func (s *AttendanceService) ClockIn(employeeID uuid.UUID, notes string) (*models.Attendance, error) {
	if _, err := s.empRepo.GetByID(employeeID); err != nil {
		return nil, errors.New("employee not found")
	}
	now := time.Now()
	date := now.Truncate(24 * time.Hour)

	// Determine status
	status := models.AttendanceStatusPresent
	if IsWeekend(date) {
		status = models.AttendanceStatusWeekend
	} else {
		isHol, err := s.holidayRepo.IsHoliday(date.Format("2006-01-02"), "")
		if err == nil && isHol {
			status = models.AttendanceStatusHoliday
		}
	}

	a, err := s.repo.SetClockIn(employeeID, date, now)
	if err != nil {
		return nil, err
	}
	a.Status = status
	a.Notes = notes
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AttendanceService) ClockOut(employeeID uuid.UUID, notes string) (*models.Attendance, error) {
	now := time.Now()
	date := now.Truncate(24 * time.Hour)
	a, err := s.repo.SetClockOut(employeeID, date, now)
	if err != nil {
		return nil, errors.New("no clock-in found for today")
	}
	if notes != "" {
		a.Notes = notes
	}
	return a, s.repo.Update(a)
}

func (s *AttendanceService) GetByEmployeeAndDate(employeeID uuid.UUID, date time.Time) (*models.Attendance, error) {
	return s.repo.GetByEmployeeAndDate(employeeID, date)
}

func (s *AttendanceService) ListByEmployee(employeeID uuid.UUID, from, to time.Time) ([]models.Attendance, error) {
	return s.repo.ListByEmployee(employeeID, from, to)
}

func (s *AttendanceService) ListByDepartmentAndDate(departmentID uuid.UUID, date time.Time) ([]models.Attendance, error) {
	return s.repo.ListByDepartmentAndDate(departmentID, date)
}

func (s *AttendanceService) CreateManual(a *models.Attendance) error {
	if a.EmployeeID == uuid.Nil {
		return errors.New("employee_id is required")
	}
	if a.Date.IsZero() {
		return errors.New("date is required")
	}
	a.Source = models.AttendanceSourceManual
	return s.repo.Create(a)
}

func (s *AttendanceService) Update(a *models.Attendance) error {
	existing, err := s.repo.GetByID(a.ID)
	if err != nil {
		return errors.New("attendance record not found")
	}
	a.EmployeeID = existing.EmployeeID
	a.Date = existing.Date
	a.Source = models.AttendanceSourceManual
	return s.repo.Update(a)
}

func (s *AttendanceService) GetMonthlySummary(employeeID uuid.UUID, month, year int) (map[string]interface{}, error) {
	return s.repo.GetMonthlySummary(employeeID, month, year)
}
