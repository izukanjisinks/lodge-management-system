package models

import (
	"time"

	"github.com/google/uuid"
)

type Payslip struct {
	ID                 uuid.UUID `json:"id"`
	EmployeeID         uuid.UUID `json:"employee_id"`
	Month              int       `json:"month"`
	Year               int       `json:"year"`
	BaseSalary         float64   `json:"base_salary"`
	HousingAllowance   float64   `json:"housing_allowance"`
	TransportAllowance float64   `json:"transport_allowance"`
	MedicalAllowance   float64   `json:"medical_allowance"`
	GrossSalary        float64   `json:"gross_salary"`
	IncomeTax          float64   `json:"income_tax"`
	LeaveDays          float64   `json:"leave_days"`
	NetSalary          float64   `json:"net_salary"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relations (populated on demand)
	EmployeeName string `json:"employee_name,omitempty"`
	PositionName string `json:"position_name,omitempty"`
}
