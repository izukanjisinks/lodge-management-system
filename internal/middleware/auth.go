package middleware

import (
	"context"
	"net/http"
	"strings"

	"hr-system/internal/repository"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
	UserEmail ContextKey = "userEmail"
	UserKey   ContextKey = "user"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format. Expected: Bearer <token>")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetUserByID(claims.UserID)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "User not found")
			return
		}

		if !user.IsActive {
			utils.RespondError(w, http.StatusUnauthorized, "Account is inactive")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmail, claims.Email)
		ctx = context.WithValue(ctx, UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmail).(string)
	return email, ok
}
