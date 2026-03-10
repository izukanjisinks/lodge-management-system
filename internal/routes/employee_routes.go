package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterEmployeeRoutes(
	empH *handlers.EmployeeHandler,
	docH *handlers.EmployeeDocumentHandler,
	ecH *handlers.EmergencyContactHandler,
) {
	http.HandleFunc("GET /api/v1/hr/employees",
		withAuthAndRole(empH.List, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("POST /api/v1/hr/employees",
		withAuthAndRole(empH.Create, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/employees/{id}",
		withAuth(empH.GetByID))

	http.HandleFunc("PUT /api/v1/hr/employees/{id}",
		withAuthAndRole(empH.Update, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("DELETE /api/v1/hr/employees/{id}",
		withAuthAndRole(empH.Delete, models.RoleSuperAdmin))

	http.HandleFunc("GET /api/v1/hr/employees/{id}/direct-reports",
		withAuth(empH.GetDirectReports))

	http.HandleFunc("GET /api/v1/hr/departments/{department_id}/managers",
		withAuth(empH.GetManagersByDepartment))

	http.HandleFunc("GET /api/v1/hr/employees/{id}/documents",
		withAuth(docH.ListByEmployee))

	http.HandleFunc("POST /api/v1/hr/employees/{id}/documents",
		withAuth(docH.Create))

	http.HandleFunc("POST /api/v1/hr/employees/{id}/documents/{did}/verify",
		withAuthAndRole(docH.Verify, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/employees/{id}/emergency-contacts",
		withAuth(ecH.ListByEmployee))

	http.HandleFunc("POST /api/v1/hr/employees/{id}/emergency-contacts",
		withAuth(ecH.Create))
}
