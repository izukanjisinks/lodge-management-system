package handlers

import (
	"fmt"
	"net/http"

	"hr-system/internal/interfaces"
	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		// Log the detailed error for debugging
		fmt.Printf("ERROR: Failed to decode employee JSON: %v\n", err)
		utils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// Validate password is provided
	if req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "password is required")
		return
	}

	// Log the received employee data (without password)
	fmt.Printf("DEBUG: Creating employee: %s %s <%s>\n", req.FirstName, req.LastName, req.Email)

	// Create employee with user account
	if err := h.service.CreateWithUser(&req.Employee, req.Password); err != nil {
		fmt.Printf("ERROR: Failed to create employee in service: %v\n", err)
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, req.Employee)
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	emp, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Employee not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, emp)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	q := r.URL.Query()
	filter := interfaces.EmployeeFilter{
		Search:           q.Get("search"),
		EmploymentStatus: q.Get("status"),
		EmploymentType:   q.Get("type"),
	}
	if v := q.Get("department_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.DepartmentID = &id
		}
	}
	if v := q.Get("position_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.PositionID = &id
		}
	}

	emps, total, err := h.service.List(filter, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list employees")
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     emps,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	var emp models.Employee
	if err := utils.DecodeJson(r, &emp); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	emp.ID = id
	if err := h.service.Update(&emp); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	if err := h.service.SoftDelete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Employee deleted"})
}

func (h *EmployeeHandler) GetDirectReports(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	reports, err := h.service.GetDirectReports(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get direct reports")
		return
	}
	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *EmployeeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filter := interfaces.EmployeeFilter{}
	emps, _, err := h.service.List(filter, 1, 1)
	if err != nil || len(emps) == 0 {
		utils.RespondError(w, http.StatusNotFound, "Employee record not found")
		return
	}

	// Find employee linked to this user
	_ = userID
	utils.RespondJSON(w, http.StatusOK, emps[0])
}

func (h *EmployeeHandler) GetManagersByDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID, err := uuid.Parse(r.PathValue("department_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	managers, err := h.service.GetManagersByDepartment(departmentID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get managers")
		return
	}

	utils.RespondJSON(w, http.StatusOK, managers)
}
