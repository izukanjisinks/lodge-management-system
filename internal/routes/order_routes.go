package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterOrderRoutes(h *handlers.OrderHandler) {
	// Read — all authenticated staff
	http.HandleFunc("GET /api/v1/orders",
		withAuth(h.List))
	http.HandleFunc("GET /api/v1/orders/{id}",
		withAuth(h.GetByID))

	// Place orders — admin, manager, receptionist
	http.HandleFunc("POST /api/v1/orders",
		withAuthAndRole(h.PlaceOrder, models.RoleAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("POST /api/v1/orders/walk-in",
		withAuthAndRole(h.PlaceWalkInOrder, models.RoleAdmin, models.RoleManager, models.RoleReceptionist))

	// Add items to an existing order
	http.HandleFunc("POST /api/v1/orders/{id}/items",
		withAuthAndRole(h.AddItems, models.RoleAdmin, models.RoleManager, models.RoleReceptionist))
}
