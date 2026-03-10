package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterLeaveRoutes(
	ltH *handlers.LeaveTypeHandler,
	lbH *handlers.LeaveBalanceHandler,
	lrH *handlers.LeaveRequestHandler,
) {
	// Leave Types
	http.HandleFunc("GET /api/v1/hr/leave-types",
		withAuth(ltH.List))

	http.HandleFunc("POST /api/v1/hr/leave-types",
		withAuthAndRole(ltH.Create, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/leave-types/{id}",
		withAuth(ltH.GetByID))

	http.HandleFunc("PUT /api/v1/hr/leave-types/{id}",
		withAuthAndRole(ltH.Update, models.RoleSuperAdmin, models.RoleHRManager))

	// Leave Balances
	http.HandleFunc("GET /api/v1/hr/leave-balances/me",
		withAuth(lbH.GetMyBalances))

	http.HandleFunc("POST /api/v1/hr/leave-balances/initialize/{year}",
		withAuthAndRole(lbH.Initialize, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/leave-balances/employee/{id}",
		withAuthAndRole(lbH.GetByEmployee, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("POST /api/v1/hr/leave-balances/adjust/{id}",
		withAuthAndRole(lbH.Adjust, models.RoleSuperAdmin, models.RoleHRManager))

	// Leave Requests
	http.HandleFunc("GET /api/v1/hr/leave-requests/me",
		withAuth(lrH.GetMyRequests))

	http.HandleFunc("GET /api/v1/hr/leave-requests",
		withAuthAndRole(lrH.List, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("POST /api/v1/hr/leave-requests",
		withAuth(lrH.Create))

	http.HandleFunc("GET /api/v1/hr/leave-requests/{id}",
		withAuth(lrH.GetByID))

	http.HandleFunc("POST /api/v1/hr/leave-requests/{id}/cancel",
		withAuth(lrH.Cancel))

	http.HandleFunc("POST /api/v1/hr/leave-requests/{id}/approve",
		withAuthAndRole(lrH.Approve, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("POST /api/v1/hr/leave-requests/{id}/reject",
		withAuthAndRole(lrH.Reject, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))
}

func RegisterAttendanceRoutes(h *handlers.AttendanceHandler) {
	http.HandleFunc("POST /api/v1/hr/attendance/clock-in",
		withAuth(h.ClockIn))

	http.HandleFunc("POST /api/v1/hr/attendance/clock-out",
		withAuth(h.ClockOut))

	http.HandleFunc("GET /api/v1/hr/attendance/me",
		withAuth(h.GetMyAttendance))

	http.HandleFunc("GET /api/v1/hr/attendance/summary",
		withAuthAndRole(h.GetSummary, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("POST /api/v1/hr/attendance/manual",
		withAuthAndRole(h.CreateManual, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/attendance/employee/{id}",
		withAuthAndRole(h.GetByEmployee, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("GET /api/v1/hr/attendance/department/{id}",
		withAuthAndRole(h.GetByDepartment, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	http.HandleFunc("PUT /api/v1/hr/attendance/{id}",
		withAuthAndRole(h.Update, models.RoleSuperAdmin, models.RoleHRManager))
}

func RegisterHolidayRoutes(h *handlers.HolidayHandler) {
	http.HandleFunc("GET /api/v1/hr/holidays",
		withAuth(h.List))

	http.HandleFunc("POST /api/v1/hr/holidays",
		withAuthAndRole(h.Create, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("PUT /api/v1/hr/holidays/{id}",
		withAuthAndRole(h.Update, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("DELETE /api/v1/hr/holidays/{id}",
		withAuthAndRole(h.Delete, models.RoleSuperAdmin, models.RoleHRManager))
}
