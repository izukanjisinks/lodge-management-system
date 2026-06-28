package middleware

import (
	"context"
	"net/http"

	"lodge-system/internal/models"
	"lodge-system/pkg/utils"
)

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if !user.HasPermission(permission) {
				utils.RespondError(w, http.StatusForbidden, "You do not have permission to perform this action")
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
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			for _, perm := range permissions {
				if user.HasPermission(perm) {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.RespondError(w, http.StatusForbidden, "You do not have permission to perform this action")
		})
	}
}

func RequireRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*models.User)
			if !ok || user == nil {
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if user.Role == nil || user.Role.Name != roleName {
				utils.RespondError(w, http.StatusForbidden, "This action requires a branch-level role")
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
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if user.Role == nil {
				utils.RespondError(w, http.StatusForbidden, "This action requires a branch-level role")
				return
			}
			for _, role := range roleNames {
				if user.Role.Name == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.RespondError(w, http.StatusForbidden, "This action requires a branch-level role")
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
