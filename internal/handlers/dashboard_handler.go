package handlers

import (
	"net/http"
	"time"

	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/pkg/utils"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(svc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: svc}
}

func (h *DashboardHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time

	if v := r.URL.Query().Get("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid 'from' date format, use YYYY-MM-DD")
			return
		}
		from = &t
	}
	if v := r.URL.Query().Get("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid 'to' date format, use YYYY-MM-DD")
			return
		}
		to = &t
	}

	stats, err := h.service.GetAdminDashboard(from, to)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to load admin dashboard")
		return
	}
	utils.RespondJSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) GetMyDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stats, err := h.service.GetEmployeeDashboard(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to load dashboard")
		return
	}
	if stats == nil {
		utils.RespondError(w, http.StatusNotFound, "No employee record linked to your account")
		return
	}

	utils.RespondJSON(w, http.StatusOK, stats)
}
