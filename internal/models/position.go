package models

import (
	"time"

	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type Position struct {
	ID                 uuid.UUID  `json:"id"`
	Title              string     `json:"title"`
	Code               string     `json:"code"`
	DepartmentID       uuid.UUID  `json:"department_id"`
	RoleID             *uuid.UUID `json:"role_id,omitempty"`
	GradeLevel         string     `json:"grade_level"`
	BaseSalary         float64    `json:"base_salary"`
	HousingAllowance   float64    `json:"housing_allowance"`
	TransportAllowance float64    `json:"transport_allowance"`
	MedicalAllowance   float64    `json:"medical_allowance"`
	IncomeTax          float64    `json:"income_tax"`
	Description        string     `json:"description"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`

	// Resolved names (populated by List queries)
	DepartmentName string `json:"department_name,omitempty"`

	// Relations (populated on demand)
	Role *Role `json:"role,omitempty"`
}

// CalculateSalaryComponents computes allowances and income tax from the base salary
func (p *Position) CalculateSalaryComponents() {
	b := utils.CalculateSalaryBreakdown(p.BaseSalary)
	p.HousingAllowance = b.HousingAllowance
	p.TransportAllowance = b.TransportAllowance
	p.MedicalAllowance = b.MedicalAllowance
	p.IncomeTax = b.IncomeTax
}

// GrossSalary returns base salary + all allowances
func (p *Position) GrossSalary() float64 {
	return p.BaseSalary + p.HousingAllowance + p.TransportAllowance + p.MedicalAllowance
}

// NetSalary returns gross salary minus income tax
func (p *Position) NetSalary() float64 {
	return p.GrossSalary() - p.IncomeTax
}
