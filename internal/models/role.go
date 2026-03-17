package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin        = "admin"
	RoleManager      = "manager"
	RoleReceptionist = "receptionist"
	RoleCleaner      = "cleaner"
)

type Role struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetPredefinedRoles() []Role {
	return []Role{
		{Name: RoleAdmin, Description: "Full system access — manages users, rooms, bookings, and configuration"},
		{Name: RoleManager, Description: "Oversees operations — approves bookings, views reports, manages rooms"},
		{Name: RoleReceptionist, Description: "Front-desk staff — handles bookings, clients, and invoices"},
		{Name: RoleCleaner, Description: "Housekeeping staff — views assigned rooms and cleaning schedule"},
	}
}
