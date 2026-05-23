package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterBranchRoutes(h *handlers.BranchHandler) {
	http.HandleFunc("GET /api/v1/branches",
		withAuth(h.List))

	http.HandleFunc("GET /api/v1/branches/{id}",
		withAuth(h.GetByID))

	http.HandleFunc("POST /api/v1/branches",
		withAuthAndRole(h.Create, models.RoleAdmin))

	http.HandleFunc("PUT /api/v1/branches/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin))

	http.HandleFunc("DELETE /api/v1/branches/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin))
}
