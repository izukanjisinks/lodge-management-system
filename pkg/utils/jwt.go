package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeStaff     TokenType = "staff"
	TokenTypeGuest     TokenType = "guest"
	TokenTypeBackoffice TokenType = "backoffice"
)

type CustomClaims struct {
	Email     string     `json:"email"`
	UserID    uuid.UUID  `json:"userId"`
	OrgID     uuid.UUID  `json:"orgId"`
	BranchID  *uuid.UUID `json:"branchId,omitempty"`
	Role      string     `json:"role"`
	TokenType TokenType  `json:"tokenType"`
	jwt.RegisteredClaims
}

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "hr-system-secret-key"
	}
	return []byte(secret)
}

func GenerateStaffToken(email string, userID, orgID uuid.UUID, branchID *uuid.UUID, role string) (string, error) {
	claims := CustomClaims{
		Email:     email,
		UserID:    userID,
		OrgID:     orgID,
		BranchID:  branchID,
		Role:      role,
		TokenType: TokenTypeStaff,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func GenerateGuestToken(email string, guestID uuid.UUID) (string, error) {
	claims := CustomClaims{
		Email:     email,
		UserID:    guestID,
		TokenType: TokenTypeGuest,
		Role:      "guest",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func GenerateBackofficeToken(email string, userID uuid.UUID) (string, error) {
	claims := CustomClaims{
		Email:     email,
		UserID:    userID,
		TokenType: TokenTypeBackoffice,
		Role:      "backoffice",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func ExtractEmailFromToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// GenerateToken is kept for backward compatibility during migration.
// Prefer GenerateStaffToken, GenerateGuestToken, or GenerateBackofficeToken.
func GenerateToken(email string, userId uuid.UUID) (string, error) {
	return GenerateStaffToken(email, userId, uuid.Nil, nil, "")
}
