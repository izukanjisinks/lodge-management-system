package handlers

import (
	"net/http"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type DepartmentHandler struct {
	service *services.DepartmentService
}

func NewDepartmentHandler(service *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dept models.Department
	if err := utils.DecodeJson(r, &dept); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.Create(&dept); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, dept)
}

func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}
	dept, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Department not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, dept)
}

func (h *DepartmentHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	filter := interfaces.DepartmentFilter{
		Search: r.URL.Query().Get("search"),
	}
	if v := r.URL.Query().Get("is_active"); v == "true" {
		t := true
		filter.IsActive = &t
	} else if v == "false" {
		f := false
		filter.IsActive = &f
	}

	depts, total, err := h.service.List(filter, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list departments")
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     depts,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}
	var dept models.Department
	if err := utils.DecodeJson(r, &dept); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	dept.ID = id
	if err := h.service.Update(&dept); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}
	if err := h.service.SoftDelete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Department deleted"})
}

func (h *DepartmentHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.service.GetTree()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to build department tree")
		return
	}
	utils.RespondJSON(w, http.StatusOK, tree)
}

// uuidFromPath extracts the last path segment as a UUID.
// Kept for handlers that are not yet on named-param routes.
func uuidFromPath(path string) (uuid.UUID, error) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return uuid.Parse(path[i+1:])
		}
	}
	return uuid.Parse(path)
}
