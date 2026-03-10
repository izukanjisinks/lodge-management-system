package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleHRManager  = "hr_manager"
	RoleManager    = "manager"
	RoleEmployee   = "employee"
)

type Role struct {
	RoleID      uuid.UUID  `json:"role_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func GetPredefinedRoles() []Role {
	return []Role{
		{Name: RoleSuperAdmin, Description: "Full system access"},
		{Name: RoleHRManager, Description: "Manage employees, payroll, recruitment, leave"},
		{Name: RoleManager, Description: "Approve leave, view team, review team"},
		{Name: RoleEmployee, Description: "Self-service access"},
	}
}
