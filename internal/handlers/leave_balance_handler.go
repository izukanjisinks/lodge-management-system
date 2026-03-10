package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type LeaveBalanceHandler struct {
	service    *services.LeaveBalanceService
	empService *services.EmployeeService
}

func NewLeaveBalanceHandler(svc *services.LeaveBalanceService, empSvc *services.EmployeeService) *LeaveBalanceHandler {
	return &LeaveBalanceHandler{service: svc, empService: empSvc}
}

func (h *LeaveBalanceHandler) GetMyBalances(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	year := time.Now().Year()
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}
	filter := interfaces.EmployeeFilter{}
	emps, _, err := h.empService.List(filter, 1, 100)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to find employee")
		return
	}
	var employeeID uuid.UUID
	for _, e := range emps {
		if e.UserID != nil && *e.UserID == user.UserID {
			employeeID = e.ID
			break
		}
	}
	if employeeID == uuid.Nil {
		utils.RespondError(w, http.StatusNotFound, "No employee record linked to your account")
		return
	}

	balances, err := h.service.GetByEmployeeAndYear(employeeID, year)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get balances")
		return
	}
	utils.RespondJSON(w, http.StatusOK, balances)
}

func (h *LeaveBalanceHandler) GetByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	year := time.Now().Year()
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}

	balances, err := h.service.GetByEmployeeAndYear(employeeID, year)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get balances")
		return
	}
	utils.RespondJSON(w, http.StatusOK, balances)
}

func (h *LeaveBalanceHandler) Initialize(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid year")
		return
	}

	var req struct {
		EmployeeID uuid.UUID `json:"employee_id"`
	}
	if err := utils.DecodeJson(r, &req); err != nil || req.EmployeeID == uuid.Nil {
		utils.RespondError(w, http.StatusBadRequest, "employee_id is required")
		return
	}

	if err := h.service.InitializeForEmployee(req.EmployeeID, year); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Balances initialized"})
}

func (h *LeaveBalanceHandler) Adjust(w http.ResponseWriter, r *http.Request) {
	balanceID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid balance ID")
		return
	}

	var req struct {
		Delta  int    `json:"delta"`
		Reason string `json:"reason"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Adjust(interfaces.AdjustBalanceInput{
		LeaveBalanceID: balanceID,
		Delta:          req.Delta,
		Reason:         req.Reason,
	}); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Balance adjusted"})
}
