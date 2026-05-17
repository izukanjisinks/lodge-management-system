package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

func RegisterAuditLogRoutes(h *handlers.AuditLogHandler) {
	http.HandleFunc("GET /api/v1/audit-logs",
		withAuth(h.List))
}
