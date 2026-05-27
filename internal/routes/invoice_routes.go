package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterInvoiceRoutes(h *handlers.InvoiceHandler) {
	// List / get — admin and manager
	http.HandleFunc("GET /api/v1/invoices",
		withAuthAndRole(h.List, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/invoices/{id}",
		withAuthAndRole(h.GetByID, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Lookup by booking — all authenticated staff
	http.HandleFunc("GET /api/v1/invoices/booking/{booking_id}",
		withAuth(h.GetByBookingID))

	// Status update — admin and manager only
	http.HandleFunc("PATCH /api/v1/invoices/{id}/status",
		withAuthAndRole(h.UpdateStatus, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}
