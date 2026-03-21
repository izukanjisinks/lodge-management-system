package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterReviewRoutes(h *handlers.ReviewHandler) {
	// Public — anyone can see the lodge rating summary
	http.HandleFunc("GET /api/v1/reviews/summary", withPublic(h.GetSummary))

	// Authenticated guest — submit a review after check-out
	http.HandleFunc("POST /api/v1/guest/reviews",
		withAuthAndRole(h.Submit, models.RoleGuest))
}
