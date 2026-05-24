package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
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
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	branchID, err := middleware.ResolveBranchID(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	stats, err := h.service.GetStaffStats(orgID, branchID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to load dashboard stats")
		return
	}
	utils.RespondJSON(w, http.StatusOK, stats)
}
