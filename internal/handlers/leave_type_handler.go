package handlers

import (
	"net/http"

	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type LeaveTypeHandler struct {
	service *services.LeaveTypeService
}

func NewLeaveTypeHandler(service *services.LeaveTypeService) *LeaveTypeHandler {
	return &LeaveTypeHandler{service: service}
}

func (h *LeaveTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var lt models.LeaveType
	if err := utils.DecodeJson(r, &lt); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.Create(&lt); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, lt)
}

func (h *LeaveTypeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid leave type ID")
		return
	}
	lt, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Leave type not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, lt)
}

func (h *LeaveTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") != "false"
	lts, err := h.service.List(activeOnly)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list leave types")
		return
	}
	utils.RespondJSON(w, http.StatusOK, lts)
}

func (h *LeaveTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid leave type ID")
		return
	}
	var lt models.LeaveType
	if err := utils.DecodeJson(r, &lt); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	lt.ID = id
	if err := h.service.Update(&lt); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}
