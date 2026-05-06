package handlers

import (
	"net/http"
	"time"

	"lodge-system/internal/middleware"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"
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

	var from, to *time.Time
	if v := r.URL.Query().Get("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid from date, expected YYYY-MM-DD")
			return
		}
		from = &t
	}
	if v := r.URL.Query().Get("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid to date, expected YYYY-MM-DD")
			return
		}
		end := t.Add(24*time.Hour - time.Second)
		to = &end
	}

	logs, total, err := h.service.List(orgID, entityType, entityID, action, from, to, p.Page, p.PageSize)
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
