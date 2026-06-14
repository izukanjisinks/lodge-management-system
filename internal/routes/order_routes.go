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
	http.HandleFunc("GET /api/v1/orders/in-house-guests",
		withAuth(h.ListInHouseGuests))
	http.HandleFunc("GET /api/v1/orders/{id}",
		withAuth(h.GetByID))

	// Place orders — admin, branch admin, manager, receptionist
	http.HandleFunc("POST /api/v1/orders",
		withAuthAndRole(h.PlaceOrder, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("POST /api/v1/orders/walk-in",
		withAuthAndRole(h.PlaceWalkInOrder, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Add / remove items on an existing order
	http.HandleFunc("POST /api/v1/orders/{id}/items",
		withAuthAndRole(h.AddItems, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("DELETE /api/v1/orders/{id}/items/{item_id}",
		withAuthAndRole(h.RemoveItem, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Manually close all open orders for the org — admin, branch admin, manager
	http.HandleFunc("PATCH /api/v1/orders/close-all",
		withAuthAndRole(h.CloseAllOrders, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
}
