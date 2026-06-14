package routes

import (
	"net/http"

	"lodge-system/internal/middleware"
)

func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply JWT auth
		middleware.JWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

func withAuthAndRole(handler http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply JWT auth and role check
		middleware.JWTAuth(
			middleware.RequireAnyRole(roles...)(http.HandlerFunc(handler)),
		).ServeHTTP(w, r)
	}
}

func withPublic(handler http.HandlerFunc) http.HandlerFunc {
	return handler
}

func withGuestAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.GuestJWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

func withBackofficeAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.BackofficeJWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

func withWebUserAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WebUserJWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}
