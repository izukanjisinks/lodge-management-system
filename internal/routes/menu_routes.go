package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterMenuRoutes(h *handlers.MenuHandler) {
	// List and get — all authenticated staff
	http.HandleFunc("GET /api/v1/menus",
		withAuth(h.ListMenus))
	http.HandleFunc("GET /api/v1/menus/{id}",
		withAuth(h.GetMenu))

	// Create / update / delete menus — admin and manager
	http.HandleFunc("POST /api/v1/menus",
		withAuthAndRole(h.CreateMenu, models.RoleAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/menus/{id}",
		withAuthAndRole(h.UpdateMenu, models.RoleAdmin, models.RoleManager))
	http.HandleFunc("DELETE /api/v1/menus/{id}",
		withAuthAndRole(h.DeleteMenu, models.RoleAdmin, models.RoleManager))

	// Menu items — admin and manager
	http.HandleFunc("POST /api/v1/menus/{id}/items",
		withAuthAndRole(h.CreateMenuItem, models.RoleAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/menus/{id}/items/{item_id}",
		withAuthAndRole(h.UpdateMenuItem, models.RoleAdmin, models.RoleManager))
	http.HandleFunc("DELETE /api/v1/menus/{id}/items/{item_id}",
		withAuthAndRole(h.DeleteMenuItem, models.RoleAdmin, models.RoleManager))
}
