package models

import (
	"time"

	"github.com/google/uuid"
)

type PayrollStatus string

const (
	PayrollStatusOpen       PayrollStatus = "OPEN"
	PayrollStatusProcessing PayrollStatus = "PROCESSING"
	PayrollStatusCompleted  PayrollStatus = "COMPLETED"
	PayrollStatusCancelled  PayrollStatus = "CANCELLED"
)

type Payroll struct {
	ID              uuid.UUID     `json:"id"`
	StartDate       time.Time     `json:"start_date"`
	EndDate         time.Time     `json:"end_date"`
	Status          PayrollStatus `json:"status"`
	ProcessedBy     *uuid.UUID    `json:"processed_by"`
	ProcessedAt     *time.Time    `json:"processed_at,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Relations (populated on demand)
	ProcessedByName string    `json:"processed_by_name,omitempty"`
	Payslips        []Payslip `json:"payslips,omitempty"`
	TotalNetSalary  float64   `json:"total_net_salary,omitempty"`
	EmployeeCount   int       `json:"employee_count,omitempty"`
}
