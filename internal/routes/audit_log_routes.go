package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterAuditLogRoutes(h *handlers.AuditLogHandler) {
	http.HandleFunc("GET /api/v1/audit-logs",
		withAuthAndRole(h.List, models.RoleAdmin, models.RoleManager))
}
