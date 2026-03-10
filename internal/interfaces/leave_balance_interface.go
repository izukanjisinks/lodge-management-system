package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type AdjustBalanceInput struct {
	LeaveBalanceID uuid.UUID
	Delta          int
	Reason         string
}

type LeaveBalanceInterface interface {
	GetByEmployeeAndYear(employeeID uuid.UUID, year int) ([]models.LeaveBalance, error)
	GetByEmployeeTypeYear(employeeID, leaveTypeID uuid.UUID, year int) (*models.LeaveBalance, error)
	InitializeForEmployee(employeeID uuid.UUID, year int) error
	Adjust(input AdjustBalanceInput) error
	IncrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error
	DecrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error
	ApproveLeave(employeeID, leaveTypeID uuid.UUID, year, days int) error
}
