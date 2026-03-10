package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID               uuid.UUID  `json:"user_id"`
	Email                string     `json:"email"`
	Password             string     `json:"-"`
	RoleID               *uuid.UUID `json:"role_id,omitempty"`
	Role                 *Role      `json:"role,omitempty"`
	IsActive             bool       `json:"is_active"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	ChangePassword       bool       `json:"change_password"`
	PasswordChangedAt    *time.Time `json:"password_changed_at,omitempty"`
	PasswordExpiresAt    *time.Time `json:"password_expires_at,omitempty"`
	FailedLoginAttempts  int        `json:"-"` // Never expose in JSON
	IsLocked             bool       `json:"is_locked"`
	LockedUntil          *time.Time `json:"locked_until,omitempty"`
}

// HasPermission checks role-based access using HR role names.
// The permission string follows the format "resource:action" or just a role name.
func (u *User) HasPermission(permission string) bool {
	if u.Role == nil {
		return false
	}
	switch u.Role.Name {
	case RoleSuperAdmin:
		return true
	case RoleHRManager:
		// HR managers can do everything except super_admin-only operations
		restricted := map[string]bool{
			"roles:delete": true,
			"users:delete": true,
		}
		return !restricted[permission]
	case RoleManager:
		allowed := map[string]bool{
			"employees:read":       true,
			"leave_requests:read":  true,
			"leave_requests:approve": true,
			"attendance:read":      true,
			"performance:read":     true,
			"performance:write":    true,
			"goals:read":           true,
			"goals:write":          true,
		}
		return allowed[permission]
	case RoleEmployee:
		selfService := map[string]bool{
			"own_profile:read":     true,
			"own_profile:write":    true,
			"leave_requests:write": true,
			"leave_balances:read":  true,
			"attendance:write":     true,
			"payslips:read":        true,
			"performance:self":     true,
			"goals:write":          true,
		}
		return selfService[permission]
	}
	return false
}
