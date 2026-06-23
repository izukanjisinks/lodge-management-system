package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

func RegisterBackofficeRoutes(
	authHandler *handlers.BackofficeAuthHandler,
	userHandler *handlers.BackofficeUserHandler,
	orgHandler *handlers.BackofficeOrganizationHandler,
) {
	// Auth — public
	http.HandleFunc("POST /api/v1/backoffice/auth/login", withPublic(authHandler.Login))

	// Auth — authenticated
	http.HandleFunc("POST /api/v1/backoffice/auth/change-password", withBackofficeAuth(authHandler.ChangePassword))
	http.HandleFunc("GET /api/v1/backoffice/auth/me", withBackofficeAuth(authHandler.Me))

	// Backoffice user management
	http.HandleFunc("GET /api/v1/backoffice/users", withBackofficeAuth(userHandler.List))
	http.HandleFunc("POST /api/v1/backoffice/users", withBackofficeAuth(userHandler.Create))
	http.HandleFunc("GET /api/v1/backoffice/users/{id}", withBackofficeAuth(userHandler.GetByID))
	http.HandleFunc("PUT /api/v1/backoffice/users/{id}", withBackofficeAuth(userHandler.Update))
	http.HandleFunc("DELETE /api/v1/backoffice/users/{id}", withBackofficeAuth(userHandler.Delete))
	http.HandleFunc("POST /api/v1/backoffice/users/{id}/reset-password", withBackofficeAuth(userHandler.ResetPassword))
	http.HandleFunc("POST /api/v1/backoffice/users/{id}/lock", withBackofficeAuth(userHandler.Lock))
	http.HandleFunc("POST /api/v1/backoffice/users/{id}/unlock", withBackofficeAuth(userHandler.Unlock))

	// Organization management
	http.HandleFunc("GET /api/v1/backoffice/organizations", withBackofficeAuth(orgHandler.List))
	http.HandleFunc("POST /api/v1/backoffice/organizations/provision", withBackofficeAuth(orgHandler.Provision))
	http.HandleFunc("GET /api/v1/backoffice/organizations/{id}", withBackofficeAuth(orgHandler.GetByID))
	http.HandleFunc("PUT /api/v1/backoffice/organizations/{id}", withBackofficeAuth(orgHandler.Update))
	http.HandleFunc("PATCH /api/v1/backoffice/organizations/{id}/status", withBackofficeAuth(orgHandler.ToggleStatus))
	http.HandleFunc("DELETE /api/v1/backoffice/organizations/{id}", withBackofficeAuth(orgHandler.Delete))
}
