package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterPayrollRoutes(h *handlers.PayrollHandler) {
	// List payrolls - requires SuperAdmin or HRManager
	http.HandleFunc("GET /api/v1/hr/payrolls",
		withAuthAndRole(h.List, models.RoleSuperAdmin, models.RoleHRManager))

	// Create payroll period - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/hr/payrolls",
		withAuthAndRole(h.Create, models.RoleSuperAdmin, models.RoleHRManager))

	// Get payroll by ID (includes payslips) - requires SuperAdmin or HRManager
	http.HandleFunc("GET /api/v1/hr/payrolls/{id}",
		withAuthAndRole(h.GetByID, models.RoleSuperAdmin, models.RoleHRManager))

	// Process payroll (generate all payslips) - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/hr/payrolls/{id}/process",
		withAuthAndRole(h.Process, models.RoleSuperAdmin, models.RoleHRManager))

	// Cancel payroll - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/hr/payrolls/{id}/cancel",
		withAuthAndRole(h.Cancel, models.RoleSuperAdmin, models.RoleHRManager))

	// Delete payroll (only OPEN) - requires SuperAdmin
	http.HandleFunc("DELETE /api/v1/hr/payrolls/{id}",
		withAuthAndRole(h.Delete, models.RoleSuperAdmin))
}
