package handlers

import (
	"net/http"

	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Parse pagination
	pag := utils.ParsePagination(r)
	search := r.URL.Query().Get("search")

	// Parse role_id filter
	var roleID *uuid.UUID
	if roleIDStr := r.URL.Query().Get("role_id"); roleIDStr != "" {
		if id, err := uuid.Parse(roleIDStr); err == nil {
			roleID = &id
		}
	}

	// Parse is_active filter
	var isActive *bool
	if activeStr := r.URL.Query().Get("is_active"); activeStr != "" {
		if activeStr == "true" {
			val := true
			isActive = &val
		} else if activeStr == "false" {
			val := false
			isActive = &val
		}
	}

	users, total, err := h.service.ListUsers(search, roleID, isActive, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Clear sensitive password data
	for i := range users {
		users[i].Password = ""
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     users,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Clear sensitive password data
	user.Password = ""

	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	profile, err := h.service.GetProfile(user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RoleID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Role ID is required")
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	user, err := h.service.ChangeUserRole(id, roleID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Password = ""
	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) Lock(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.LockUser(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "User account locked successfully",
	})
}

func (h *UserHandler) Unlock(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.UnlockUser(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "User account unlocked successfully",
	})
}
