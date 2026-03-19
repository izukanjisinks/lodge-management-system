package routes

import (
	"lodge-system/internal/handlers"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roomHandler *handlers.RoomHandler,
	clientHandler *handlers.ClientHandler,
	bookingHandler *handlers.BookingHandler,
	mealPlanHandler *handlers.MealPlanHandler,
	workflowHandler *handlers.WorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
) {
	RegisterPublicRoutes(authHandler)
	RegisterAuthRoutes(authHandler)
	RegisterUserRoutes(userHandler)
	RegisterRoomRoutes(roomHandler)
	RegisterClientRoutes(clientHandler)
	RegisterBookingRoutes(bookingHandler)
	RegisterMealPlanRoutes(mealPlanHandler)
	RegisterWorkflowRoutes(workflowHandler)
	RegisterWorkflowAdminRoutes(workflowAdminHandler)
}
