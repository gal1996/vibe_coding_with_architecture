package interactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo    repository.UserRepository
	authService port.AuthService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository, authService port.AuthService) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		authService: authService,
	}
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Username string
	Password string
	IsAdmin  bool
}

// Register registers a new user
func (uc *UserUseCase) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	// Check if username already exists
	existingUser, _ := uc.userRepo.FindByUsername(ctx, input.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Generate user ID
	userID := generateUserID()

	// Create user entity
	user, err := entity.NewUser(userID, input.Username, input.Password, input.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	// Save to repository
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// LoginInput represents the input for user login
type LoginInput struct {
	Username string
	Password string
}

// LoginOutput represents the output of user login
type LoginOutput struct {
	User  *entity.User
	Token string
}

// Login authenticates a user and returns a token
func (uc *UserUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Find user by username
	user, err := uc.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check password
	err = user.CheckPassword(input.Password)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate token
	token, err := uc.authService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginOutput{
		User:  user,
		Token: token,
	}, nil
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// generateUserID generates a unique user ID
func generateUserID() string {
	// In a real implementation, this would use a proper ID generator
	return fmt.Sprintf("USR-%d", timeUser.Now().Unix())
}

var timeUser = struct {
	Now func() struct {
		Unix func() int64
	}
}{
	Now: func() struct {
		Unix func() int64
	} {
		return struct {
			Unix func() int64
		}{
			Unix: func() int64 {
				return 1234567891 // Placeholder
			},
		}
	},
}