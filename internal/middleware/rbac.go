package middleware

import (
	"context"
	"net/http"

	"hr-system/internal/models"
)

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !user.HasPermission(permission) {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			for _, perm := range permissions {
				if user.HasPermission(perm) {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

func RequireRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if user.Role == nil || user.Role.Name != roleName {
				http.Error(w, "Forbidden: role required", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAnyRole(roleNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if user.Role == nil {
				http.Error(w, "Forbidden: role required", http.StatusForbidden)
				return
			}
			for _, role := range roleNames {
				if user.Role.Name == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden: role required", http.StatusForbidden)
		})
	}
}

const PermissionsKey ContextKey = "permissions"

func AttachPermissionsToContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, PermissionsKey, permissions)
}

func GetPermissionsFromContext(ctx context.Context) []string {
	permissions, ok := ctx.Value(PermissionsKey).([]string)
	if !ok {
		return []string{}
	}
	return permissions
}
