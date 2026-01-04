package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// contextKey is a type for context keys
type contextKey string

const (
	// ContextKeyUser is the context key for the current user
	ContextKeyUser contextKey = "user"
	// SecretKey is the secret key for JWT signing (in production, this should be from environment)
	SecretKey = "your-secret-key-change-in-production"
)

// JWTAuthService implements AuthService using JWT
type JWTAuthService struct {
	userRepo repository.UserRepository
}

// NewJWTAuthService creates a new JWT auth service
func NewJWTAuthService(userRepo repository.UserRepository) port.AuthService {
	return &JWTAuthService{
		userRepo: userRepo,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func (s *JWTAuthService) GenerateToken(user *entity.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *JWTAuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token claims")
}

// GetCurrentUser gets the current user from the context
func (s *JWTAuthService) GetCurrentUser(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(ContextKeyUser).(*entity.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

// SetUserInContext sets the user in the context
func SetUserInContext(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, user)
}