package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin        = "admin"
	RoleBranchAdmin  = "branch_admin"
	RoleManager      = "manager"
	RoleReceptionist = "receptionist"
	RoleCleaner      = "cleaner"
	RoleGuest        = "guest"
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
		{Name: RoleAdmin, Description: "Admin — full access including branch management"},
		{Name: RoleBranchAdmin, Description: "Branch admin — full access scoped to their assigned branch"},
		{Name: RoleManager, Description: "Oversees operations — approves bookings, views reports, manages rooms"},
		{Name: RoleReceptionist, Description: "Front-desk staff — handles bookings, clients, and invoices"},
		{Name: RoleCleaner, Description: "Housekeeping staff — views assigned rooms and cleaning schedule"},
		{Name: RoleGuest, Description: "Guest user — limited access to view available rooms and make bookings via website"},
	}
}
