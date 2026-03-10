package handlers

import (
	"net/http"
	"time"

	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type PayrollHandler struct {
	service *services.PayrollService
}

func NewPayrollHandler(service *services.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// Create opens a new payroll period
func (h *PayrollHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
		return
	}

	payroll, err := h.service.Create(startDate, endDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, payroll)
}

// GetByID returns a payroll with its payslips
func (h *PayrollHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	payroll, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Payroll not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, payroll)
}

// List returns paginated payrolls
func (h *PayrollHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	status := r.URL.Query().Get("status")

	payrolls, total, err := h.service.List(status, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list payrolls")
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     payrolls,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// Process generates payslips for all active employees and completes the payroll
func (h *PayrollHandler) Process(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	payroll, err := h.service.Process(id, userID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, payroll)
}

// Cancel marks a payroll as cancelled
func (h *PayrollHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	if err := h.service.Cancel(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Payroll cancelled successfully",
	})
}

// Delete removes a payroll (only if OPEN)
func (h *PayrollHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Payroll deleted successfully",
	})
}
