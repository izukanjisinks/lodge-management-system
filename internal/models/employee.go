package models

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentType string
type EmploymentStatus string

const (
	EmploymentTypeFullTime  EmploymentType = "full_time"
	EmploymentTypePartTime  EmploymentType = "part_time"
	EmploymentTypeContract  EmploymentType = "contract"
	EmploymentTypeIntern    EmploymentType = "intern"

	EmploymentStatusActive     EmploymentStatus = "active"
	EmploymentStatusOnLeave    EmploymentStatus = "on_leave"
	EmploymentStatusSuspended  EmploymentStatus = "suspended"
	EmploymentStatusTerminated EmploymentStatus = "terminated"
	EmploymentStatusResigned   EmploymentStatus = "resigned"
)

type Employee struct {
	ID                 uuid.UUID        `json:"id"`
	UserID             *uuid.UUID       `json:"user_id,omitempty"`
	EmployeeNumber     string           `json:"employee_number"`
	FirstName          string           `json:"first_name"`
	LastName           string           `json:"last_name"`
	Email              string           `json:"email"`
	PersonalEmail      string           `json:"personal_email"`
	Phone              string           `json:"phone"`
	DateOfBirth        *time.Time       `json:"date_of_birth,omitempty"`
	Gender             string           `json:"gender"`
	NationalID         string           `json:"national_id"`
	MaritalStatus      string           `json:"marital_status"`
	Address            string           `json:"address"`
	City               string           `json:"city"`
	State              string           `json:"state"`
	Country            string           `json:"country"`
	DepartmentID       uuid.UUID        `json:"department_id"`
	PositionID         uuid.UUID        `json:"position_id"`
	ManagerID          *uuid.UUID       `json:"manager_id,omitempty"`
	HireDate           time.Time        `json:"hire_date"`
	ProbationEndDate   *time.Time       `json:"probation_end_date,omitempty"`
	EmploymentType     EmploymentType   `json:"employment_type"`
	EmploymentStatus   EmploymentStatus `json:"employment_status"`
	TerminationDate    *time.Time       `json:"termination_date,omitempty"`
	TerminationReason  string           `json:"termination_reason"`
	ProfilePhotoURL    string           `json:"profile_photo_url"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          *time.Time       `json:"deleted_at,omitempty"`

	// Resolved names (populated by List queries)
	DepartmentName string `json:"department_name,omitempty"`
	PositionName   string `json:"position_name,omitempty"`
	ManagerName    string `json:"manager_name,omitempty"`

	// Relations (populated on demand)
	Department *Department `json:"department,omitempty"`
	Position   *Position   `json:"position,omitempty"`
}

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// CreateEmployeeRequest is the request payload for creating a new employee
type CreateEmployeeRequest struct {
	Employee
	Password string `json:"password"` // Initial password for the user account
}
