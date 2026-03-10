package handlers

import (
	"net/http"
	"strconv"

	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type PayslipHandler struct {
	service *services.PayslipService
}

func NewPayslipHandler(service *services.PayslipService) *PayslipHandler {
	return &PayslipHandler{service: service}
}

// Generate creates a payslip for an employee for the given month/year
func (h *PayslipHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID string `json:"employee_id"`
		Month      int    `json:"month"`
		Year       int    `json:"year"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	payslip, err := h.service.Generate(empID, req.Month, req.Year)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, payslip)
}

// GetByID returns a specific payslip
func (h *PayslipHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payslip ID")
		return
	}

	payslip, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Payslip not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, payslip)
}

// List returns paginated payslips with optional filters
func (h *PayslipHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)

	var employeeID *uuid.UUID
	if v := r.URL.Query().Get("employee_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			employeeID = &id
		}
	}

	var month *int
	if v := r.URL.Query().Get("month"); v != "" {
		if m, err := strconv.Atoi(v); err == nil {
			month = &m
		}
	}

	var year *int
	if v := r.URL.Query().Get("year"); v != "" {
		if y, err := strconv.Atoi(v); err == nil {
			year = &y
		}
	}

	payslips, total, err := h.service.List(employeeID, month, year, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list payslips")
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     payslips,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// GetMyPayslips returns the authenticated employee's own payslips
func (h *PayslipHandler) GetMyPayslips(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	emp, err := h.service.GetEmployeeByUserID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Employee record not found")
		return
	}

	pag := utils.ParsePagination(r)
	empID := emp.ID

	var year *int
	if v := r.URL.Query().Get("year"); v != "" {
		if y, err := strconv.Atoi(v); err == nil {
			year = &y
		}
	}

	payslips, total, err := h.service.List(&empID, nil, year, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve payslips")
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     payslips,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// Delete removes a payslip
func (h *PayslipHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payslip ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Payslip deleted successfully",
	})
}
