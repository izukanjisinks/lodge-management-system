package handlers

import (
	"net/http"

	"lodge-system/internal/services"
	"lodge-system/pkg/utils"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) StaffStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStaffStats()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to load dashboard stats")
		return
	}
	utils.RespondJSON(w, http.StatusOK, stats)
}
