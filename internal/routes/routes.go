package routes

import (
	"hr-system/internal/handlers"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roleHandler *handlers.RoleHandler,
	deptHandler *handlers.DepartmentHandler,
	posHandler *handlers.PositionHandler,
	empHandler *handlers.EmployeeHandler,
	docHandler *handlers.EmployeeDocumentHandler,
	ecHandler *handlers.EmergencyContactHandler,
	ltHandler *handlers.LeaveTypeHandler,
	lbHandler *handlers.LeaveBalanceHandler,
	lrHandler *handlers.LeaveRequestHandler,
	attHandler *handlers.AttendanceHandler,
	holidayHandler *handlers.HolidayHandler,
	dashboardHandler *handlers.DashboardHandler,
	workflowHandler *handlers.WorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
) {
	RegisterPublicRoutes(authHandler)
	RegisterAuthRoutes(authHandler)
	RegisterUserRoutes(userHandler)
	RegisterRoleRoutes(roleHandler)
	RegisterDepartmentRoutes(deptHandler)
	RegisterPositionRoutes(posHandler)
	RegisterEmployeeRoutes(empHandler, docHandler, ecHandler)
	RegisterLeaveRoutes(ltHandler, lbHandler, lrHandler)
	RegisterAttendanceRoutes(attHandler)
	RegisterHolidayRoutes(holidayHandler)
	RegisterDashboardRoutes(dashboardHandler)
	RegisterWorkflowRoutes(workflowHandler)
	RegisterWorkflowAdminRoutes(workflowAdminHandler)
}
