package models

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string
type AttendanceSource string

const (
	AttendanceStatusPresent  AttendanceStatus = "present"
	AttendanceStatusAbsent   AttendanceStatus = "absent"
	AttendanceStatusHalfDay  AttendanceStatus = "half_day"
	AttendanceStatusOnLeave  AttendanceStatus = "on_leave"
	AttendanceStatusHoliday  AttendanceStatus = "holiday"
	AttendanceStatusWeekend  AttendanceStatus = "weekend"

	AttendanceSourceManual    AttendanceSource = "manual"
	AttendanceSourceSystem    AttendanceSource = "system"
	AttendanceSourceBiometric AttendanceSource = "biometric"
)

type Attendance struct {
	ID            uuid.UUID        `json:"id"`
	EmployeeID    uuid.UUID        `json:"employee_id"`
	Date          time.Time        `json:"date"`
	ClockIn       *time.Time       `json:"clock_in,omitempty"`
	ClockOut      *time.Time       `json:"clock_out,omitempty"`
	TotalHours    float64          `json:"total_hours"`
	Status        AttendanceStatus `json:"status"`
	OvertimeHours float64          `json:"overtime_hours"`
	Notes         string           `json:"notes"`
	Source        AttendanceSource `json:"source"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
