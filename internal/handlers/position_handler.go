package handlers

import (
	"net/http"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type PositionHandler struct {
	service *services.PositionService
}

func NewPositionHandler(service *services.PositionService) *PositionHandler {
	return &PositionHandler{service: service}
}

func (h *PositionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pos models.Position
	if err := utils.DecodeJson(r, &pos); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.Create(&pos); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, pos)
}

func (h *PositionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}
	pos, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Position not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, pos)
}

func (h *PositionHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	filter := interfaces.PositionFilter{}
	if v := r.URL.Query().Get("department_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.DepartmentID = &id
		}
	}
	if v := r.URL.Query().Get("grade_level"); v != "" {
		filter.GradeLevel = v
	}
	if v := r.URL.Query().Get("is_active"); v == "true" {
		t := true
		filter.IsActive = &t
	}

	positions, total, err := h.service.List(filter, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list positions")
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     positions,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *PositionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}
	var pos models.Position
	if err := utils.DecodeJson(r, &pos); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	pos.ID = id
	if err := h.service.Update(&pos); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *PositionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}
	if err := h.service.SoftDelete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Position deleted"})
}
