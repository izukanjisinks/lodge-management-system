package handlers

import (
	"net/http"

	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve roles")
		return
	}
	utils.RespondJSON(w, http.StatusOK, roles)
}

func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.service.GetRoleByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Role not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, role)
}
