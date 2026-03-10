package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Code                  string    `json:"code"`
	Description           string    `json:"description"`
	DefaultDaysPerYear    int       `json:"default_days_per_year"`
	IsPaid                bool      `json:"is_paid"`
	IsCarryForwardAllowed bool      `json:"is_carry_forward_allowed"`
	MaxCarryForwardDays   int       `json:"max_carry_forward_days"`
	RequiresApproval      bool      `json:"requires_approval"`
	RequiresDocument      bool      `json:"requires_document"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// DefaultLeaveTypes returns the seeded leave types per the spec.
func DefaultLeaveTypes() []LeaveType {
	return []LeaveType{
		{Code: "AL", Name: "Annual Leave", DefaultDaysPerYear: 21, IsPaid: true, IsCarryForwardAllowed: true, MaxCarryForwardDays: 5, RequiresApproval: true},
		{Code: "SL", Name: "Sick Leave", DefaultDaysPerYear: 15, IsPaid: true, IsCarryForwardAllowed: false, RequiresApproval: true, RequiresDocument: true},
		{Code: "PL", Name: "Parental Leave", DefaultDaysPerYear: 90, IsPaid: true, IsCarryForwardAllowed: false, RequiresApproval: true},
		{Code: "UL", Name: "Unpaid Leave", DefaultDaysPerYear: 0, IsPaid: false, IsCarryForwardAllowed: false, RequiresApproval: true},
		{Code: "CL", Name: "Compassionate Leave", DefaultDaysPerYear: 5, IsPaid: true, IsCarryForwardAllowed: false, RequiresApproval: true},
		{Code: "ML", Name: "Marriage Leave", DefaultDaysPerYear: 3, IsPaid: true, IsCarryForwardAllowed: false, RequiresApproval: true},
	}
}
