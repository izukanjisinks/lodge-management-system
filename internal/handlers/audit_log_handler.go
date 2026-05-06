package handlers

import (
	"lodge-system/internal/middleware"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"
	"net/http"
)

type AuditLogHandler struct {
	service *services.AuditLogService
}

func NewAuditLogHandler(service *services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

func (h *AuditLogHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)

	entityType := r.URL.Query().Get("entity_type")
	entityID := r.URL.Query().Get("entity_id")
	action := r.URL.Query().Get("action")

	logs, total, err := h.service.List(orgID, entityType, entityID, action, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     logs,
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    total,
	})
}
