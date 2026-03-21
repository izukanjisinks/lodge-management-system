package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterMealPlanRoutes(h *handlers.MealPlanHandler) {
	// Read — public (website guests select meal plans during reservation)
	http.HandleFunc("GET /api/v1/meal-plans",
		withPublic(h.List))

	http.HandleFunc("GET /api/v1/meal-plans/{id}",
		withPublic(h.GetByID))

	// Write — admin and manager only
	http.HandleFunc("POST /api/v1/meal-plans",
		withAuthAndRole(h.Create, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/meal-plans/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("DELETE /api/v1/meal-plans/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin, models.RoleManager))
}
