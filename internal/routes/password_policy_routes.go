package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterPasswordPolicyRoutes(handler *handlers.PasswordPolicyHandler) {
	// Password policy management routes (authenticated)
	http.HandleFunc("GET /api/v1/password-policy",
		withAuthAndRole(handler.GetPasswordPolicy, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("PUT /api/v1/password-policy",
		withAuthAndRole(handler.UpdatePasswordPolicy, models.RoleSuperAdmin, models.RoleHRManager))

	// Password change routes (authenticated)
	http.HandleFunc("POST /api/v1/auth/change-password",
		withAuth(handler.ChangePassword))

	http.HandleFunc("POST /api/v1/auth/reset-password",
		withAuthAndRole(handler.ResetUserPassword, models.RoleSuperAdmin, models.RoleHRManager))

	// Password generation route (accessible to authenticated users)
	http.HandleFunc("GET /api/v1/password-policy/generate",
		withAuth(handler.GeneratePassword))
}
