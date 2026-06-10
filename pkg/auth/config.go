package auth

import (
	"time"

	goevents "github.com/go-thin/go-events"
	"github.com/go-thin/go-auth/pkg/storage"
)

// SigningMethod defines the type for JWT signing methods.
const (
	HS256 = "HS256"
	HS384 = "HS384"
	HS512 = "HS512"
	RS256 = "RS256"
)

// JWTConfig holds the configuration for JWT generation and validation.
// These settings are used to create the internal JWTManager.
type JWTConfig struct {
	AccessSecret    []byte
	RefreshSecret   []byte
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningMethod   string
}

// Config is the main configuration struct for the AuthService.
// It holds all the necessary components and settings.
type Config struct {
	Storage  storage.Storage
	JWT      JWTConfig
	EventBus goevents.Bus // optional; nil disables event publishing
}
