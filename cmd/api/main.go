package main

import (
	"fmt"
	"log"
	"net/http"

	"lodge-system/internal/config"
	"lodge-system/internal/database"
	"lodge-system/internal/handlers"
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

	workflowService := services.NewWorkflowService(workflowRepo, instanceRepo, taskRepo, historyRepo, userRepo, emailService)

	// Seed predefined roles
	if err := roleService.InitializePredefinedRoles(); err != nil {
		log.Fatalf("Failed to initialize roles: %v", err)
	}
	log.Println("Roles initialized")

	// Seed default admin
	if err := userService.SeedSuperAdmin("admin@lodge.dev", "Admin@123"); err != nil {
		log.Printf("Warning: failed to seed admin: %v", err)
	} else {
		log.Println("Admin ready (admin@lodge.dev)")
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(services.NewRoomService(roomRepo))
	clientHandler := handlers.NewClientHandler(services.NewClientService(clientRepo))
	workflowHandler := handlers.NewWorkflowHandler(workflowService)
	workflowAdminHandler := handlers.NewWorkflowAdminHandler(workflowRepo)
	passwordPolicyHandler := handlers.NewPasswordPolicyHandler(passwordPolicyService, userService)

	// Register routes
	routes.RegisterRoutes(authHandler, userHandler, roomHandler, clientHandler, workflowHandler, workflowAdminHandler)
	routes.RegisterPasswordPolicyRoutes(passwordPolicyHandler)

	// Apply CORS middleware globally
	handler := middleware.CORS(http.DefaultServeMux)

	addr := ":" + cfg.ServerPort
	log.Printf("Lodge Management System running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
