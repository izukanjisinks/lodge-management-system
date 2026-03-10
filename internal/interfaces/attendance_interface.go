package interfaces

import (
	"time"

	"hr-system/internal/models"

	"github.com/google/uuid"
)

type AttendanceSummary struct {
	EmployeeID    uuid.UUID `json:"employee_id"`
	Month         int       `json:"month"`
	Year          int       `json:"year"`
	DaysPresent   int       `json:"days_present"`
	DaysAbsent    int       `json:"days_absent"`
	DaysOnLeave   int       `json:"days_on_leave"`
	DaysLate      int       `json:"days_late"`
	TotalHours    float64   `json:"total_hours"`
	OvertimeHours float64   `json:"overtime_hours"`
}

type AttendanceInterface interface {
	ClockIn(employeeID uuid.UUID, notes string) (*models.Attendance, error)
	ClockOut(employeeID uuid.UUID, notes string) (*models.Attendance, error)
	GetByEmployeeAndDate(employeeID uuid.UUID, date time.Time) (*models.Attendance, error)
	ListByEmployee(employeeID uuid.UUID, from, to time.Time) ([]models.Attendance, error)
	ListByDepartmentAndDate(departmentID uuid.UUID, date time.Time) ([]models.Attendance, error)
	CreateManual(a *models.Attendance) error
	Update(a *models.Attendance) error
	GetSummary(employeeID uuid.UUID, month, year int) (*AttendanceSummary, error)
}
