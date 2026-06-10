# go-auth

JWT authentication library for Go. Part of the [go-thin](https://github.com/go-thin) ecosystem.

## Install

```bash
go get github.com/pragneshbagary/go-auth
```

## Quick start

```go
import (
    "github.com/pragneshbagary/go-auth/pkg/auth"
    "github.com/pragneshbagary/go-auth/internal/storage/memory"
)

svc, err := auth.NewAuthService(auth.Config{
    Storage: memory.NewInMemoryStorage(),
    JWT: auth.JWTConfig{
        AccessSecret:  []byte("access-secret"),
        RefreshSecret: []byte("refresh-secret"),
        SigningMethod:  auth.HS256,
    },
})

// Register
user, err := svc.Register(auth.RegisterPayload{
    Username: "alice",
    Email:    "alice@example.com",
    Password: "hunter2",
})

// Login
resp, err := svc.Login("alice", "hunter2", nil)
// resp.AccessToken  — JWT access token
// resp.RefreshToken — JWT refresh token

// Validate
claims, err := svc.ValidateAccessToken(resp.AccessToken)
```

## Storage

`AuthService` accepts any backend through a two-method interface:

```go
type Storage interface {
    CreateUser(user models.User) error
    GetUserByUsername(username string) (*models.User, error)
}
```

`internal/storage/memory` ships an in-memory implementation for tests and local development.

## Event bus

Pass an optional `goevents.Bus` to publish auth events to the rest of your application:

```go
import goevents "github.com/go-thin/go-events"

bus := goevents.New()

// go-audit: record every event
bus.SubscribeAll(func(env goevents.Envelope) error {
    log.Printf("[%s] %s", env.PublishedAt.Format(time.RFC3339), env.Event.Topic())
    return nil
})

// go-notify: send welcome email on registration
bus.Subscribe("auth.user.registered", func(env goevents.Envelope) error {
    e := env.Event.(auth.UserRegistered)
    return mailer.SendWelcome(e.Email, e.Username)
})

svc, err := auth.NewAuthService(auth.Config{
    Storage:  store,
    JWT:      jwtCfg,
    EventBus: bus, // nil = no events published
})
```

### Published events

**`auth.user.registered`** — fired after a successful `Register` call.
Fields: `UserID`, `Username`, `Email`.

**`auth.user.logged_in`** — fired after a successful `Login` call.
Fields: `UserID`, `Username`.

## JWT config fields

- `AccessSecret` — HMAC secret for access tokens (required)
- `RefreshSecret` — HMAC secret for refresh tokens (required)
- `SigningMethod` — `HS256`, `HS384`, or `HS512` (required)
- `Issuer` — JWT `iss` claim (optional)
- `AccessTokenTTL` — access token lifetime, zero means no expiry (optional)
- `RefreshTokenTTL` — refresh token lifetime, zero means no expiry (optional)

## Testing

```bash
go test ./...
```

## License

MIT
