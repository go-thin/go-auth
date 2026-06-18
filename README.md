# go-auth

A focused JWT authentication library for Go. Handles user registration, login, and token validation with Argon2id password hashing and an event bus for auth lifecycle hooks.

## Installation

```bash
go get github.com/go-thin/go-auth
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/go-thin/go-auth/pkg/auth"
)

func main() {
    svc, err := auth.New(auth.Config{
        Storage: myStorage, // your storage.Storage implementation
        JWT: auth.JWTConfig{
            AccessSecret:    []byte("your-access-secret"),
            RefreshSecret:   []byte("your-refresh-secret"),
            Issuer:          "my-app",
            AccessTokenTTL:  15 * time.Minute,
            RefreshTokenTTL: 7 * 24 * time.Hour,
            SigningMethod:   auth.HS256,
        },
    })
    if err != nil {
        panic(err)
    }

    // Register
    user, err := svc.Register(auth.RegisterPayload{
        Username: "alice",
        Email:    "alice@example.com",
        Password: "s3cr3t",
    })
    _ = user

    // Login
    tokens, err := svc.Login("alice", "s3cr3t", nil)
    fmt.Println(tokens.AccessToken, tokens.RefreshToken)

    // Validate
    claims, err := svc.ValidateAccessToken(tokens.AccessToken)
    fmt.Println(claims["username"])
}
```

## Configuration

```go
auth.Config{
    Storage:  myStorage,     // required — implements storage.Storage
    JWT:      auth.JWTConfig{...},
    EventBus: bus,           // optional — nil disables event publishing
}
```

### JWTConfig fields

| Field             | Type            | Description                                              |
|-------------------|-----------------|----------------------------------------------------------|
| `AccessSecret`    | `[]byte`        | HMAC secret (or RSA private key bytes for RS256)         |
| `RefreshSecret`   | `[]byte`        | Separate secret for refresh tokens                       |
| `Issuer`          | `string`        | JWT `iss` claim                                          |
| `AccessTokenTTL`  | `time.Duration` | Access token lifetime                                    |
| `RefreshTokenTTL` | `time.Duration` | Refresh token lifetime                                   |
| `SigningMethod`   | `string`        | `auth.HS256`, `auth.HS384`, `auth.HS512`, `auth.RS256`   |

## API

### Register

```go
user, err := svc.Register(auth.RegisterPayload{
    Username: "alice",
    Email:    "alice@example.com",
    Password: "s3cr3t",
})
```

Returns `*models.User`. Password is hashed with Argon2id before storage — the plaintext is never persisted.

### Login

```go
tokens, err := svc.Login(username, password string, customClaims map[string]interface{})
// tokens.AccessToken  — short-lived JWT
// tokens.RefreshToken — long-lived JWT
```

`customClaims` is merged into the access token payload alongside the default `username` and `email` claims. Pass `nil` if no extra claims are needed.

### ValidateAccessToken / ValidateRefreshToken

```go
claims, err := svc.ValidateAccessToken(tokenString)
claims, err := svc.ValidateRefreshToken(tokenString)
// returns jwt.MapClaims on success
```

Use `ValidateRefreshToken` to confirm a refresh token is still valid; re-authenticate via `Login` to issue a fresh access token.

## Storage backends

Three backends ship with the library:

```go
// SQLite — file or in-memory
import "github.com/go-thin/go-auth/pkg/storage/sqlite"
store, err := sqlite.New("./auth.db")   // file
store, err := sqlite.New(":memory:")    // ephemeral

// PostgreSQL
import "github.com/go-thin/go-auth/pkg/storage/postgres"
store, err := postgres.New("postgres://user:pass@localhost/mydb")

// Pass to auth.New
svc, err := auth.New(auth.Config{Storage: store, JWT: auth.JWTConfig{...}})
```

Both backends auto-create the `users` table on first open. `store.DB()` exposes the underlying `*sql.DB` if you need raw access for migrations or transactions.

## Implementing a custom storage backend

Implement the three-method `storage.Storage` interface to plug in any other backend:

```go
import (
    "github.com/go-thin/go-auth/pkg/models"
    "github.com/go-thin/go-auth/pkg/storage"
)

type MyDB struct{}

func (db *MyDB) CreateUser(user models.User) error                        { /* insert */ return nil }
func (db *MyDB) GetUserByUsername(username string) (*models.User, error)  { /* select */ return nil, nil }
func (db *MyDB) DeleteUser(id string) error                               { /* delete */ return nil }
// Return storage.ErrUserNotFound when the user does not exist.
```

Pass your implementation as `Config.Storage`.

## Event Bus

When `Config.EventBus` is set (using [go-events](https://github.com/go-thin/go-events)), the service publishes events on auth actions:

| Event            | Topic                  | Fields                        |
|------------------|------------------------|-------------------------------|
| `UserRegistered` | `auth.user.registered` | `UserID`, `Username`, `Email` |
| `UserLoggedIn`   | `auth.user.logged_in`  | `UserID`, `Username`          |

```go
import goevents "github.com/go-thin/go-events"

bus := goevents.New()
bus.Subscribe("auth.user.registered", func(e auth.UserRegistered) {
    fmt.Println("new user:", e.Username)
})

svc, _ := auth.NewAuthServiceLegacy(auth.Config{
    // ...
    EventBus: bus,
})
```

## Error Handling

The `auth` package exports structured `AuthError` types with machine-readable codes for use in HTTP handlers:

```go
// Construct structured errors
err := auth.ErrInvalidCredentials()
err := auth.ErrUserExists("username")
err := auth.ErrInvalidToken()

// Write JSON error responses
auth.WriteJSONError(w, err)            // auto-maps to HTTP status
auth.WriteErrorResponse(w, err, 401)   // explicit status
```

## Testing

```bash
go test ./...        # all tests; postgres tests skip without TEST_POSTGRES_DSN
go test ./... -race
go test ./... -cover

# Run postgres integration tests against a real database:
TEST_POSTGRES_DSN="postgres://user:pass@localhost/testdb" go test ./pkg/storage/postgres/...
```

## License

MIT
