package auth

import (
	"time"

	"github.com/go-thin/go-auth/pkg/storage"
	goevents "github.com/go-thin/go-events"
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
	Storage      storage.Storage
	JWT          JWTConfig
	Argon2Params *Argon2Params // optional; nil uses DefaultParams
	EventBus     goevents.Bus  // optional; nil disables event publishing
}
