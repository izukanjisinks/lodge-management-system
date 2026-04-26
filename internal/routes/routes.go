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
	invoiceHandler *handlers.InvoiceHandler,
	dashboardHandler *handlers.DashboardHandler,
	workflowHandler *handlers.WorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
	guestAuthHandler *handlers.GuestAuthHandler,
	guestBookingHandler *handlers.GuestBookingHandler,
	reviewHandler *handlers.ReviewHandler,
	backofficeAuthHandler *handlers.BackofficeAuthHandler,
	backofficeUserHandler *handlers.BackofficeUserHandler,
	backofficeOrgHandler *handlers.BackofficeOrganizationHandler,
) {
	RegisterPublicRoutes(authHandler)
	RegisterAuthRoutes(authHandler)
	RegisterUserRoutes(userHandler)
	RegisterRoomRoutes(roomHandler)
	RegisterClientRoutes(clientHandler)
	RegisterBookingRoutes(bookingHandler)
	RegisterMealPlanRoutes(mealPlanHandler)
	RegisterInvoiceRoutes(invoiceHandler)
	RegisterDashboardRoutes(dashboardHandler)
	RegisterWorkflowRoutes(workflowHandler)
	RegisterWorkflowAdminRoutes(workflowAdminHandler)
	RegisterGuestRoutes(guestAuthHandler, guestBookingHandler, roomHandler, mealPlanHandler)
	RegisterReviewRoutes(reviewHandler)
	RegisterBackofficeRoutes(backofficeAuthHandler, backofficeUserHandler, backofficeOrgHandler)
}
