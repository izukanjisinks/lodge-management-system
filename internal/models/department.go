package models

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Code               string     `json:"code"`
	Description        string     `json:"description"`
	ParentDepartmentID *uuid.UUID `json:"parent_department_id,omitempty"`
	ManagerID          *uuid.UUID `json:"manager_id,omitempty"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`

	// Resolved names (populated by List queries)
	ParentDepartmentName string `json:"parent_department_name,omitempty"`
	ManagerName          string `json:"manager_name,omitempty"`

	// Relations (populated on demand)
	Children []*Department `json:"children,omitempty"`
}
