package routes

import (
	"lodge-system/internal/handlers"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	workflowHandler *handlers.WorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
) {
	RegisterPublicRoutes(authHandler)
	RegisterAuthRoutes(authHandler)
	RegisterUserRoutes(userHandler)
	RegisterWorkflowRoutes(workflowHandler)
	RegisterWorkflowAdminRoutes(workflowAdminHandler)
}
