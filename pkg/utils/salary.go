package utils

import "math"

// Allowance percentages (proportion of base salary)
const (
	HousingAllowancePct   = 0.20  // 20% of base salary
	TransportAllowancePct = 0.10  // 10% of base salary
	MedicalAllowancePct   = 0.08  // 8% of base salary
	LeaveDayRate          = 100.0 // fixed amount per unused leave day
)

// SalaryBreakdown holds the calculated salary components
type SalaryBreakdown struct {
	BaseSalary         float64
	HousingAllowance   float64
	TransportAllowance float64
	MedicalAllowance   float64
	GrossSalary        float64
	IncomeTax          float64
	NetSalary          float64
}

// CalculateSalaryBreakdown computes all salary components from a base salary
func CalculateSalaryBreakdown(baseSalary float64) SalaryBreakdown {
	housing := round(baseSalary * HousingAllowancePct)
	transport := round(baseSalary * TransportAllowancePct)
	medical := round(baseSalary * MedicalAllowancePct)
	gross := baseSalary + housing + transport + medical
	tax := CalculatePAYE(gross)

	return SalaryBreakdown{
		BaseSalary:         baseSalary,
		HousingAllowance:   housing,
		TransportAllowance: transport,
		MedicalAllowance:   medical,
		GrossSalary:        gross,
		IncomeTax:          tax,
		NetSalary:          gross - tax,
	}
}

// CalculatePAYE computes income tax using progressive PAYE brackets
//
//	Band 1: 0 – 4,800       → 0%
//	Band 2: 4,801 – 6,800   → 20%
//	Band 3: 6,801+           → 30%
func CalculatePAYE(grossIncome float64) float64 {
	var tax float64
	switch {
	case grossIncome <= 4800:
		tax = 0
	case grossIncome <= 6800:
		tax = (grossIncome - 4800) * 0.20
	default:
		tax = (2000 * 0.20) + (grossIncome-6800)*0.30
	}
	return round(tax)
}

// round2 rounds to 2 decimal places
func round(v float64) float64 {
	return math.Round(v*100) / 100
}
