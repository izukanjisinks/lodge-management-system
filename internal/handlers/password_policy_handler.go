package handlers

import (
	"net/http"

	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/internal/utils/password"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type PasswordPolicyHandler struct {
	policyService *services.PasswordPolicyService
	userService   *services.UserService
}

func NewPasswordPolicyHandler(policyService *services.PasswordPolicyService, userService *services.UserService) *PasswordPolicyHandler {
	return &PasswordPolicyHandler{
		policyService: policyService,
		userService:   userService,
	}
}

// GetPasswordPolicy retrieves the password policy (global default for now)
func (h *PasswordPolicyHandler) GetPasswordPolicy(w http.ResponseWriter, r *http.Request) {
	policy := h.policyService.GetPolicy()
	if policy == nil {
		utils.RespondError(w, http.StatusNotFound, "No password policy configured")
		return
	}

	utils.RespondJSON(w, http.StatusOK, policy)
}

// UpdatePasswordPolicy updates or creates the global password policy
func (h *PasswordPolicyHandler) UpdatePasswordPolicy(w http.ResponseWriter, r *http.Request) {
	// Only super_admin or hr_manager can update password policy
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if user has permission
	if user.Role == nil || (user.Role.Name != models.RoleSuperAdmin && user.Role.Name != models.RoleHRManager) {
		utils.RespondError(w, http.StatusForbidden, "Insufficient permissions to update password policy")
		return
	}

	var req models.CreatePasswordPolicyRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update password policy
	policy, err := h.policyService.UpsertPolicy(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, policy)
}

// ChangePassword allows authenticated user to change their password
func (h *PasswordPolicyHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		utils.RespondError(w, http.StatusBadRequest, "Old password and new password are required")
		return
	}

	if err := h.userService.ChangePassword(user.UserID, req.OldPassword, req.NewPassword); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

// GeneratePassword generates a password based on the current password policy
func (h *PasswordPolicyHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	// Get current password policy
	policy := h.policyService.GetPolicy()
	if policy == nil {
		utils.RespondError(w, http.StatusNotFound, "No password policy configured")
		return
	}

	// Generate password using the policy settings
	generatedPassword, err := password.GeneratePassword(
		policy.MinLength,
		policy.RequireUppercase,
		policy.RequireLowercase,
		policy.RequireNumbers,
		policy.RequireSpecialChars,
	)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate password")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"password": generatedPassword,
	})
}

// ResetUserPassword allows admin to reset a user's password
func (h *PasswordPolicyHandler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	// Only super_admin or hr_manager can reset passwords
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if user.Role == nil || (user.Role.Name != models.RoleSuperAdmin && user.Role.Name != models.RoleHRManager) {
		utils.RespondError(w, http.StatusForbidden, "Insufficient permissions to reset passwords")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.ResetPassword(userID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Password reset successfully. A temporary password has been sent to the user's email. User must change password on next login.",
	})
}
