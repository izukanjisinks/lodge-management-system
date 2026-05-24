package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterPasswordPolicyRoutes(h *handlers.PasswordPolicyHandler) {
	http.HandleFunc("GET /api/v1/password-policy",
		withAuthAndRole(h.GetPasswordPolicy, models.RoleAdmin, models.RoleBranchAdmin))

	http.HandleFunc("PUT /api/v1/password-policy",
		withAuthAndRole(h.UpdatePasswordPolicy, models.RoleAdmin, models.RoleBranchAdmin))

	http.HandleFunc("POST /api/v1/auth/change-password",
		withAuth(h.ChangePassword))

	http.HandleFunc("GET /api/v1/auth/generate-password",
		withAuth(h.GeneratePassword))

	http.HandleFunc("POST /api/v1/admin/users/{id}/reset-password",
		withAuthAndRole(h.ResetUserPassword, models.RoleAdmin, models.RoleBranchAdmin))
}
