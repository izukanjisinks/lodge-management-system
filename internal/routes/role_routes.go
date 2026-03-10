package routes

import (
	"net/http"

	"hr-system/internal/handlers"
)

func RegisterRoleRoutes(h *handlers.RoleHandler) {
	http.HandleFunc("GET /api/v1/roles", withAuth(h.GetAll))
	http.HandleFunc("GET /api/v1/roles/{id}", withAuth(h.GetByID))
}
