package main

import (
	"fmt"
	"log"
	"net/http"

	"hr-system/internal/config"
	"hr-system/internal/database"
	"hr-system/internal/handlers"
	"hr-system/internal/jobs"
	"hr-system/internal/middleware"
	"hr-system/internal/repositories"
	"hr-system/internal/repository"
	"hr-system/internal/routes"
	"hr-system/internal/services"
	"hr-system/internal/utils/email"
)

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	if err := database.Connect(connStr); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Database connected")

	// Repositories — Phase 1
	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	deptRepo := repository.NewDepartmentRepository()
	posRepo := repository.NewPositionRepository()
	empRepo := repository.NewEmployeeRepository()
	docRepo := repository.NewEmployeeDocumentRepository()
	ecRepo := repository.NewEmergencyContactRepository()

	// Repositories — Phase 2
	ltRepo := repository.NewLeaveTypeRepository()
	lbRepo := repository.NewLeaveBalanceRepository()
	lrRepo := repository.NewLeaveRequestRepository()
	holidayRepo := repository.NewHolidayRepository()
	attRepo := repository.NewAttendanceRepository()

	// Workflow Repositories
	workflowRepo := repository.NewWorkflowRepository()
	instanceRepo := repository.NewWorkflowInstanceRepository()
	taskRepo := repository.NewAssignedTaskRepository()
	historyRepo := repository.NewWorkflowHistoryRepository()

	// Password Policy Repositories (from repositories package)
	passwordPolicyRepo := repositories.NewPasswordPolicyRepository()
	passwordHistoryRepo := repositories.NewPasswordHistoryRepository()

	// Services — Phase 1
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	deptService := services.NewDepartmentService(deptRepo)
	posService := services.NewPositionService(posRepo, deptRepo)

	// Services — Phase 2 (create some services early for workflow dependencies)
	ltService := services.NewLeaveTypeService(ltRepo)
	lbService := services.NewLeaveBalanceService(lbRepo, ltRepo, empRepo)

	// Password Policy Service
	passwordPolicyService := services.NewPasswordPolicyService(passwordPolicyRepo, passwordHistoryRepo)
	log.Println("Password policy service initialized")

	// Set password policy service on user service
	userService.SetPasswordPolicyService(passwordPolicyService)

	// Email Service
	emailService := email.NewEmailService(&cfg.Email)
	log.Println("Email service initialized")

	// Set email service on user service (for password reset emails)
	userService.SetEmailService(emailService)

	// Employee Service (created after email service for dependency)
	empService := services.NewEmployeeService(empRepo, deptRepo, posRepo, userService, emailService)
	docService := services.NewEmployeeDocumentService(docRepo, empRepo)
	ecService := services.NewEmergencyContactService(ecRepo, empRepo)

	// Workflow Service (create after leave balance service for dependency injection)
	// Note: LeaveRequestRepo, LeaveBalanceService, and EmailService are passed to enable full workflow functionality
	workflowService := services.NewWorkflowService(workflowRepo, instanceRepo, taskRepo, historyRepo, userRepo, empRepo, lrRepo, lbService, emailService)

	// Services — Phase 2 (continued)
	lrService := services.NewLeaveRequestService(lrRepo, lbService, ltRepo, holidayRepo, empRepo, workflowService)
	holidayService := services.NewHolidayService(holidayRepo)
	attService := services.NewAttendanceService(attRepo, holidayRepo, empRepo)

	// Seed predefined roles
	if err := roleService.InitializePredefinedRoles(); err != nil {
		log.Fatalf("Failed to initialize roles: %v", err)
	}
	log.Println("Roles initialized")

	// Seed default super admin
	if err := userService.SeedSuperAdmin("admin@hr-system.com", "Admin@123"); err != nil {
		log.Printf("Warning: failed to seed super admin: %v", err)
	} else {
		log.Println("Super admin ready (admin@hr-system.com)")
	}

	// Seed default leave types
	if err := ltService.SeedDefaults(); err != nil {
		log.Printf("Warning: failed to seed leave types: %v", err)
	} else {
		log.Println("Leave types initialized")
	}

	// Handlers — Phase 1
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	deptHandler := handlers.NewDepartmentHandler(deptService)
	posHandler := handlers.NewPositionHandler(posService)
	empHandler := handlers.NewEmployeeHandler(empService)
	docHandler := handlers.NewEmployeeDocumentHandler(docService)
	ecHandler := handlers.NewEmergencyContactHandler(ecService)

	// Handlers — Phase 2
	ltHandler := handlers.NewLeaveTypeHandler(ltService)
	lbHandler := handlers.NewLeaveBalanceHandler(lbService, empService)
	lrHandler := handlers.NewLeaveRequestHandler(lrService, empService)
	holidayHandler := handlers.NewHolidayHandler(holidayService)
	attHandler := handlers.NewAttendanceHandler(attService, empService)

	// Payslip
	payslipRepo := repository.NewPayslipRepository()
	payslipService := services.NewPayslipService(payslipRepo, empRepo, posRepo, lbRepo)
	payslipHandler := handlers.NewPayslipHandler(payslipService)

	// Payroll
	payrollRepo := repository.NewPayrollRepository()
	payrollService := services.NewPayrollService(payrollRepo, payslipService, empRepo, emailService)
	payrollHandler := handlers.NewPayrollHandler(payrollService)

	// Dashboard
	adminDashRepo := repository.NewAdminDashboardRepository()
	dashboardService := services.NewDashboardService(empRepo, posRepo, deptRepo, lbRepo, lrRepo, adminDashRepo)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Workflow Handler
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	// Workflow Admin Handler
	workflowAdminHandler := handlers.NewWorkflowAdminHandler(workflowRepo)

	// Password Policy Handler
	passwordPolicyHandler := handlers.NewPasswordPolicyHandler(passwordPolicyService, userService)

	// Background jobs
	jobs.NewMonthlyLeaveAccrualJob(empRepo, lbRepo, ltRepo).Start()
	log.Println("Monthly leave accrual job scheduled")
	jobs.NewYearEndCarryForwardJob(lbRepo, ltRepo).Start()
	log.Println("Year-end carry-forward job scheduled")

	// Register routes
	routes.RegisterRoutes(
		authHandler, userHandler, roleHandler, deptHandler, posHandler, empHandler, docHandler, ecHandler,
		ltHandler, lbHandler, lrHandler, attHandler, holidayHandler, dashboardHandler,
		workflowHandler, workflowAdminHandler,
	)
	routes.RegisterPasswordPolicyRoutes(passwordPolicyHandler)
	routes.RegisterPayslipRoutes(payslipHandler)
	routes.RegisterPayrollRoutes(payrollHandler)

	// Apply CORS middleware globally to the default mux
	handler := middleware.CORS(http.DefaultServeMux)

	addr := ":" + cfg.ServerPort
	log.Printf("HR System running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
