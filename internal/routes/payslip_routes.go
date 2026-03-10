package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterPayslipRoutes(h *handlers.PayslipHandler) {
	// Get my payslips - any authenticated employee
	http.HandleFunc("GET /api/v1/hr/payslips/me",
		withAuth(h.GetMyPayslips))

	// List all payslips - requires SuperAdmin or HRManager
	http.HandleFunc("GET /api/v1/hr/payslips",
		withAuthAndRole(h.List, models.RoleSuperAdmin, models.RoleHRManager))

	// Generate payslip - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/hr/payslips",
		withAuthAndRole(h.Generate, models.RoleSuperAdmin, models.RoleHRManager))

	// Get payslip by ID - any authenticated user
	http.HandleFunc("GET /api/v1/hr/payslips/{id}",
		withAuth(h.GetByID))

	// Delete payslip - requires SuperAdmin
	http.HandleFunc("DELETE /api/v1/hr/payslips/{id}",
		withAuthAndRole(h.Delete, models.RoleSuperAdmin))
}
