package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveBalance struct {
	ID              uuid.UUID `json:"id"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	LeaveTypeID     uuid.UUID `json:"leave_type_id"`
	Year            int       `json:"year"`
	TotalEntitled   int       `json:"total_entitled"` //total entitled for the year, excluding carried forward and earned leave days
	Used            int       `json:"used"`
	Pending         int       `json:"pending"`           //pending leave requests that are not yet approved
	CarriedForward  int       `json:"carried_forward"`   //unused days from previous year that are carried forward
	EarnedLeaveDays int       `json:"earned_leave_days"` //days earned through tenure, not counted in total_entitled
	// Balance is computed: total_entitled + carried_forward + earned_leave_days - used - pending
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (populated on demand)
	LeaveType *LeaveType `json:"leave_type,omitempty"`
}
