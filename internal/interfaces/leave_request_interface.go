package interfaces

import (
	"time"

	"hr-system/internal/models"

	"github.com/google/uuid"
)

type LeaveRequestFilter struct {
	EmployeeID   *uuid.UUID
	Status       string
	StartDateGTE *time.Time
	EndDateLTE   *time.Time
	DepartmentID *uuid.UUID
}

type LeaveRequestInterface interface {
	Create(req *models.LeaveRequest) error
	GetByID(id uuid.UUID) (*models.LeaveRequest, error)
	List(filter LeaveRequestFilter, page, pageSize int) ([]models.LeaveRequest, int, error)
	Cancel(id, employeeID uuid.UUID) error
	Approve(id, reviewerID uuid.UUID, comment string) error
	Reject(id, reviewerID uuid.UUID, comment string) error
	HasOverlap(employeeID uuid.UUID, start, end time.Time, excludeID *uuid.UUID) (bool, error)
}
