package main

import (
	"fmt"
	"log"
	"net/http"

	"lodge-system/internal/config"
	"lodge-system/internal/database"
	"lodge-system/internal/handlers"
	"lodge-system/internal/jobs"
	"lodge-system/internal/middleware"
	"lodge-system/internal/repositories"
	"lodge-system/internal/repository"
	"lodge-system/internal/routes"
	"lodge-system/internal/services"
	"lodge-system/internal/utils/email"
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

	// Repositories
	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	roomRepo := repository.NewRoomRepository()
	clientRepo := repository.NewClientRepository()
	bookingRepo := repository.NewBookingRepository()
	mealPlanRepo := repository.NewMealPlanRepository()
	invoiceRepo := repository.NewInvoiceRepository()
	dashboardRepo := repository.NewDashboardRepository()
	auditLogRepo := repository.NewAuditLogRepository()
	auditLogHandler := handlers.NewAuditLogHandler(services.NewAuditLogService(auditLogRepo))
	orgSettingsRepo := repository.NewOrganizationSettingsRepository()
	orgSettingsHandler := handlers.NewOrganizationSettingsHandler(services.NewOrganizationSettingsService(orgSettingsRepo))

	workflowRepo := repository.NewWorkflowRepository()
	instanceRepo := repository.NewWorkflowInstanceRepository()
	taskRepo := repository.NewAssignedTaskRepository()
	historyRepo := repository.NewWorkflowHistoryRepository()

	passwordPolicyRepo := repositories.NewPasswordPolicyRepository()
	passwordHistoryRepo := repositories.NewPasswordHistoryRepository()

	// Services
	roleService := services.NewRoleService(userRepo, roleRepo)
	userService := services.NewUserService(userRepo, roleRepo)

	passwordPolicyService := services.NewPasswordPolicyService(passwordPolicyRepo, passwordHistoryRepo)
	log.Println("Password policy service initialized")

	userService.SetPasswordPolicyService(passwordPolicyService)

	emailService := email.NewEmailService(&cfg.Email)
	log.Println("Email service initialized")

	userService.SetEmailService(emailService)

	guestRepo := repository.NewGuestRepository()

	workflowService := services.NewWorkflowService(workflowRepo, instanceRepo, taskRepo, historyRepo, userRepo, clientRepo, guestRepo, emailService)

	// Seed predefined roles
	if err := roleService.InitializePredefinedRoles(); err != nil {
		log.Fatalf("Failed to initialize roles: %v", err)
	}
	log.Println("Roles initialized")

	// Handlers
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(services.NewRoomService(roomRepo))
	clientHandler := handlers.NewClientHandler(services.NewClientService(clientRepo))
	bookingSvc := services.NewBookingService(bookingRepo, roomRepo, clientRepo)
	invoiceSvc := services.NewInvoiceService(invoiceRepo, bookingRepo, roomRepo)
	bookingSvc.SetInvoiceService(invoiceSvc)

	bookingHandler := handlers.NewBookingHandler(bookingSvc)
	mealPlanHandler := handlers.NewMealPlanHandler(services.NewMealPlanService(mealPlanRepo))
	invoiceHandler := handlers.NewInvoiceHandler(invoiceSvc)
	dashboardHandler := handlers.NewDashboardHandler(services.NewDashboardService(dashboardRepo))
	workflowHandler := handlers.NewWorkflowHandler(workflowService)
	workflowAdminHandler := handlers.NewWorkflowAdminHandler(workflowRepo)
	passwordPolicyHandler := handlers.NewPasswordPolicyHandler(passwordPolicyService, userService)

	orgRepo := repository.NewOrganizationRepository()

	guestAuthSvc := services.NewGuestAuthService(guestRepo)
	guestAuthSvc.SetEmailService(emailService)
	guestBookingSvc := services.NewGuestBookingService(bookingRepo, roomRepo, guestAuthSvc)
	guestBookingSvc.SetWorkflowService(workflowService)
	guestAuthHandler := handlers.NewGuestAuthHandler(guestAuthSvc, orgRepo)
	guestBookingHandler := handlers.NewGuestBookingHandler(guestBookingSvc)

	backofficeUserRepo := repository.NewBackofficeUserRepository()

	backofficeAuthSvc := services.NewBackofficeAuthService(backofficeUserRepo)
	backofficeAuthSvc.SetEmailService(emailService)

	backofficeUserSvc := services.NewBackofficeUserService(backofficeUserRepo)
	backofficeUserSvc.SetEmailService(emailService)

	backofficeOrgSvc := services.NewBackofficeOrganizationService(orgRepo, userRepo, roleRepo)
	backofficeOrgSvc.SetEmailService(emailService)

	backofficeAuthHandler := handlers.NewBackofficeAuthHandler(backofficeAuthSvc)
	backofficeUserHandler := handlers.NewBackofficeUserHandler(backofficeUserSvc)
	backofficeOrgHandler := handlers.NewBackofficeOrganizationHandler(backofficeOrgSvc)

	menuRepo := repository.NewMenuRepository()
	orderRepo := repository.NewOrderRepository()
	menuHandler := handlers.NewMenuHandler(services.NewMenuService(menuRepo))
	orderSvc := services.NewOrderService(orderRepo, invoiceRepo, bookingRepo, auditLogRepo)
	orderHandler := handlers.NewOrderHandler(orderSvc)

	branchRepo := repository.NewBranchRepository()
	branchHandler := handlers.NewBranchHandler(services.NewBranchService(branchRepo))

	reviewRepo := repository.NewReviewRepository()
	reviewHandler := handlers.NewReviewHandler(services.NewReviewService(reviewRepo, bookingRepo, guestAuthSvc))

	// Background jobs
	jobs.NewOverdueCheckoutJob(bookingRepo, invoiceRepo, auditLogRepo, orgSettingsRepo).Start()
	log.Println("Overdue checkout job scheduled")
	jobs.NewCloseOrdersJob(orderSvc, orgSettingsRepo).Start()
	log.Println("Close orders job scheduled")

	// Register routes
	routes.RegisterRoutes(authHandler,
		userHandler,
		roomHandler,
		clientHandler,
		bookingHandler,
		mealPlanHandler,
		invoiceHandler,
		dashboardHandler,
		workflowHandler,
		workflowAdminHandler,
		menuHandler,
		orderHandler,
		guestAuthHandler,
		guestBookingHandler,
		reviewHandler,
		backofficeAuthHandler,
		backofficeUserHandler,
		backofficeOrgHandler,
		auditLogHandler,
		orgSettingsHandler,
		branchHandler)
	routes.RegisterPasswordPolicyRoutes(passwordPolicyHandler)

	// Apply CORS middleware globally
	handler := middleware.CORS(http.DefaultServeMux)

	addr := ":" + cfg.ServerPort
	log.Printf("Lodge Management System running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
