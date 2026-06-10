package auth

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	goevents "github.com/go-thin/go-events"
	"github.com/go-thin/go-auth/internal/jwtutils"
	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"

	"github.com/google/uuid"
)

// AuthService provides high-level authentication operations.
// It orchestrates user storage, password hashing, and token generation.
type AuthService struct {
	storage    storage.Storage
	jwtManager jwtutils.TokenManager
	eventBus   goevents.Bus
}

// NewAuthServiceLegacy creates a new AuthService from a configuration.
// Prefer auth.New, auth.NewSQLite, or auth.NewPostgres for new code.
func NewAuthServiceLegacy(cfg Config) (*AuthService, error) {
	// Internal validation of the config
	if cfg.Storage == nil {
		return nil, errors.New("storage implementation cannot be nil")
	}
	if cfg.JWT.AccessSecret == nil || cfg.JWT.RefreshSecret == nil {
		return nil, errors.New("JWT access and refresh secrets must be provided")
	}

	// Create the internal JWT manager from the public config.
	// The user of the library no longer needs to know about jwtutils.
	jwtManager := jwtutils.NewJWTManager(jwtutils.JWTConfig{
		AccessSecret:    cfg.JWT.AccessSecret,
		RefreshSecret:   cfg.JWT.RefreshSecret,
		Issuer:          cfg.JWT.Issuer,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		SigningMethod:   cfg.JWT.SigningMethod,
	})

	return &AuthService{
		storage:    cfg.Storage,
		jwtManager: jwtManager,
		eventBus:   cfg.EventBus,
	}, nil
}

// RegisterPayload defines the data required to register a new user.
// Using a struct makes the API cleaner and more extensible.
type RegisterPayload struct {
	Username string
	Email    string
	Password string
}

// Register creates a new user, hashes their password, and saves them to storage.
// It returns the newly created user.
func (s *AuthService) Register(payload RegisterPayload) (*models.User, error) {
	// Basic validation (in a real app, you'd add more, e.g., password strength)
	if payload.Username == "" || payload.Password == "" {
		return nil, errors.New("username and password cannot be empty")
	}

	// Check if user already exists
	if _, err := s.storage.GetUserByUsername(payload.Username); err == nil {
		return nil, fmt.Errorf("user with username '%s' already exists", payload.Username)
	}

	passwordHash, err := HashPassword(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := models.User{
		ID:           uuid.New().String(), // Generate a new unique ID for the user.
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: passwordHash,
	}

	if err := s.storage.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if s.eventBus != nil {
		s.eventBus.Publish(UserRegistered{
			UserID:   newUser.ID,
			Username: newUser.Username,
			Email:    newUser.Email,
		})
	}

	return &newUser, nil
}

// LoginResponse is the structure returned upon a successful login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login authenticates a user and returns an access and refresh token pair.
// It accepts customClaims to be embedded in the access token for authorization purposes.
func (s *AuthService) Login(username, password string, customClaims map[string]interface{}) (*LoginResponse, error) {
	user, err := s.storage.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password") // Generic error for security
	}

	match, err := CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		log.Printf("error checking password hash for user %s: %v", username, err)
		return nil, errors.New("invalid username or password")
	}
	if !match {
		return nil, errors.New("invalid username or password")
	}

	// Add standard claims like username to the token by default.
	claims := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}
	// Merge with any custom claims provided by the caller.
	for k, v := range customClaims {
		claims[k] = v
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if s.eventBus != nil {
		s.eventBus.Publish(UserLoggedIn{
			UserID:   user.ID,
			Username: user.Username,
		})
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateAccessToken validates an access token string.
// It returns the claims if the token is valid, otherwise an error.
func (s *AuthService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return s.jwtManager.ValidateAccessToken(tokenString)
}

// ValidateRefreshToken validates a refresh token string.
// It returns the claims if the token is valid, otherwise an error.
func (s *AuthService) ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return s.jwtManager.ValidateRefreshToken(tokenString)
}
