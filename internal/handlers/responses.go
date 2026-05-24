package handlers

import "lodge-system/internal/models"

// userResponse converts a User model to the SystemUser shape expected by the frontend.
func userResponse(u *models.User) map[string]interface{} {
	roleName := ""
	if u.Role != nil {
		roleName = u.Role.Name
	} else if u.RoleName != "" {
		roleName = u.RoleName
	}
	status := "active"
	if !u.IsActive {
		status = "inactive"
	}
	res := map[string]interface{}{
		"id":          u.UserID,
		"full_name":   u.FullName,
		"email":       u.Email,
		"role":        roleName,
		"status":      status,
		"branch_id":   u.BranchID,
		"branch_name": u.BranchName,
		"created_at":  u.CreatedAt,
	}
	if u.LastLoginAt != nil {
		res["last_login"] = u.LastLoginAt
	}
	return res
}
