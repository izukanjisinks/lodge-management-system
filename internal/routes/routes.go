package routes

import (
	"lodge-system/internal/handlers"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roomHandler *handlers.RoomHandler,
	clientHandler *handlers.ClientHandler,
	workflowHandler *handlers.WorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
) {
	RegisterPublicRoutes(authHandler)
	RegisterAuthRoutes(authHandler)
	RegisterUserRoutes(userHandler)
	RegisterRoomRoutes(roomHandler)
	RegisterClientRoutes(clientHandler)
	RegisterWorkflowRoutes(workflowHandler)
	RegisterWorkflowAdminRoutes(workflowAdminHandler)
}
