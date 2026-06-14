package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

func RegisterWebUserRoutes(h *handlers.WebUserAuthHandler) {
	// Public — no auth
	http.HandleFunc("POST /api/v1/web/auth/register", withPublic(h.Register))
	http.HandleFunc("POST /api/v1/web/auth/login", withPublic(h.Login))

	// Authenticated web user
	http.HandleFunc("GET /api/v1/web/profile", withWebUserAuth(h.GetProfile))
	http.HandleFunc("PUT /api/v1/web/profile", withWebUserAuth(h.UpdateProfile))
	http.HandleFunc("PUT /api/v1/web/auth/change-password", withWebUserAuth(h.ChangePassword))
}
