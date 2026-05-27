package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterMenuRoutes(h *handlers.MenuHandler) {
	// Single menu per org — no ID in path
	http.HandleFunc("GET /api/v1/menu",
		withAuth(h.GetMenu))
	http.HandleFunc("PUT /api/v1/menu",
		withAuthAndRole(h.UpsertMenu, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	// Menu items
	http.HandleFunc("POST /api/v1/menu/items",
		withAuthAndRole(h.CreateMenuItem, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/menu/items/{item_id}",
		withAuthAndRole(h.UpdateMenuItem, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("DELETE /api/v1/menu/items/{item_id}",
		withAuthAndRole(h.DeleteMenuItem, models.RoleBranchAdmin, models.RoleManager))
}
