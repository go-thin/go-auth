# go-auth v2.0.0: Modern Authentication for Go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue.svg)](https://golang.org/)
[![Version](https://img.shields.io/github/v/tag/pragneshbagary/go-auth?label=version&color=green)](https://github.com/pragneshbagary/go-auth/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pragneshbagary/go-auth)](https://goreportcard.com/report/github.com/pragneshbagary/go-auth)
[![Documentation](https://pkg.go.dev/badge/github.com/pragneshbagary/go-auth.svg)](https://pkg.go.dev/github.com/pragneshbagary/go-auth)
[![Tests](https://github.com/pragneshbagary/go-auth/workflows/Tests/badge.svg)](https://github.com/pragneshbagary/go-auth/actions)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](https://github.com/pragneshbagary/go-auth)

**A production-ready, feature-rich authentication library for Go applications with zero breaking changes from v1.**

## Why Choose go-auth v2?

- **Production Ready**: Used in production with comprehensive testing
- **Security First**: Built with security best practices and regular audits
- **Scalable**: Designed for high-performance applications
- **Developer Friendly**: Intuitive API with excellent documentation
- **Future Proof**: Regular updates and active maintenance
- **Community Driven**: Open source with active community support
- **Zero Downtime Migration**: Upgrade from v1 without breaking changes

---

## What's New in v2.0.0

- **Simplified API** with intuitive constructors
- **Component-based architecture** (Users, Tokens, Middleware)
- **Enhanced configuration** with environment variable support
- **Framework-specific middleware** (Gin, Echo, Fiber)
- **Built-in monitoring, metrics, and health checks**
- **Advanced security features** and best practices
- **Automatic database migration** system
- **Comprehensive documentation** and examples
- **Full backward compatibility** with v1
- **Automated migration tools**

**Zero Breaking Changes!** Your existing v1 code continues to work without modification.

### v1 vs v2 Comparison

| Feature | v1 | v2 |
|---------|----|----|
| **API Complexity** | Complex configuration | Simple constructors |
| **Architecture** | Monolithic | Component-based |
| **Framework Support** | Manual middleware | Built-in adapters |
| **Configuration** | Code-only | Environment variables |
| **Monitoring** | None | Built-in metrics |
| **Migration Tools** | None | Automated tools |
| **Security** | Basic | Advanced features |
| **Performance** | Good | Optimized |
| **Documentation** | Basic | Comprehensive |
| **Backward Compatibility** | N/A | 100% compatible |

---

## Installation

```bash
go get github.com/pragneshbagary/go-auth@v2.0.0
```

### Requirements

- Go 1.19 or higher
- SQLite (included) or PostgreSQL (optional)
- Redis (optional, for caching)

---

## Quick Start

### Simple Setup (New v2 API)

```go
package main

import (
    "log"
    "github.com/pragneshbagary/go-auth/pkg/auth"
)

func main() {
    // Initialize with SQLite database
    authService, err := auth.New("auth.db", "your-jwt-secret")
    if err != nil {
        log.Fatal(err)
    }
    
    // Register a user
    user, err := authService.Register(auth.RegisterRequest{
        Username: "alice",
        Email:    "alice@example.com",
        Password: "secure_password123",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Login with custom claims
    loginResult, err := authService.Login("alice", "secure_password123", map[string]interface{}{
        "role": "admin",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate access token
    claims, err := authService.ValidateAccessToken(loginResult.AccessToken)
    if err != nil {
        log.Fatal(err)
    }
    
    // Refresh tokens
    refreshResult, err := authService.RefreshToken(loginResult.RefreshToken)
    if err != nil {
        log.Fatal(err)
    }
    
    // Advanced user management
    users := authService.Users()
    profile, err := users.Get(user.ID)
    if err != nil {
        log.Fatal(err)
    }
    
    // Advanced token management
    tokens := authService.Tokens()
    err = tokens.Revoke(refreshResult.AccessToken)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Even Simpler Setup

```go
// For quick prototyping
auth, err := auth.Quick("your-jwt-secret")

// Load configuration from environment
authService, err := auth.NewFromEnv()
```

---

## Component-Based Architecture

### Users Component

```go
users := authService.Users()

// Get user profile (safe, no sensitive data)
profile, err := users.Get(userID)
profile, err := users.GetByEmail("user@example.com")
profile, err := users.GetByUsername("username")

// Update user profile
err = users.Update(userID, auth.UserUpdate{
    Email: &newEmail,
    Metadata: map[string]interface{}{
        "role": "admin",
        "department": "engineering",
    },
})

// Change password
err = users.ChangePassword(userID, oldPassword, newPassword)

// Password reset workflow
resetToken, err := users.CreateResetToken("user@example.com")
err = users.ResetPassword(resetToken.Token, newPassword)

// List users
userList, err := users.List(10, 0)
err = users.Delete(userID)
```

### Tokens Component

```go
tokens := authService.Tokens()

// Advanced token validation
user, err := tokens.Validate(accessToken)
isValid := tokens.IsValid(accessToken)

// Batch validation for performance
results := tokens.ValidateBatch([]string{token1, token2, token3})

// Token refresh with automatic rotation
refreshResult, err := tokens.Refresh(refreshToken)

// Token revocation
err = tokens.Revoke(accessToken)           // Single token
err = tokens.RevokeAll(userID)             // All user tokens

// Session management
sessions, err := tokens.ListActiveSessions(userID)
sessionInfo, err := tokens.GetSessionInfo(accessToken)

// Cleanup
err = tokens.CleanupExpired()
```

---

## Framework Integration

### Gin

```go
r := gin.Default()
middleware := authService.Middleware()

// Protected routes
protected := r.Group("/api/protected")
protected.Use(middleware.Gin())
{
    protected.GET("/profile", func(c *gin.Context) {
        // User automatically injected into context
        user, ok := auth.GetUserFromGin(c)
        if !ok {
            c.JSON(500, gin.H{"error": "User not found"})
            return
        }
        c.JSON(200, gin.H{"user": user})
    })
}

// Optional authentication
r.GET("/api/optional", middleware.GinOptional(), handler)
```

### Echo

```go
e := echo.New()
middleware := authService.Middleware()

protected := e.Group("/api/protected")
protected.Use(middleware.Echo())

protected.GET("/profile", func(c echo.Context) error {
    user, ok := auth.GetUserFromEcho(c)
    if !ok {
        return c.JSON(500, map[string]string{"error": "User not found"})
    }
    return c.JSON(200, map[string]interface{}{"user": user})
})
```

### Fiber

```go
app := fiber.New()
middleware := authService.Middleware()

protected := app.Group("/api/protected")
protected.Use(middleware.Fiber())

protected.Get("/profile", func(c *fiber.Ctx) error {
    user, ok := auth.GetUserFromFiber(c)
    if !ok {
        return c.Status(500).JSON(fiber.Map{"error": "User not found"})
    }
    return c.JSON(fiber.Map{"user": user})
})
```

### Standard HTTP

```go
mux := http.NewServeMux()

// Protected endpoint
mux.Handle("/protected", authService.Protect(protectedHandler))
mux.Handle("/optional", authService.Optional(optionalHandler))
```

---

## Configuration Options

### Environment Variables

Create a `.env` file or set environment variables:

```bash
AUTH_JWT_ACCESS_SECRET="your-access-secret"
AUTH_JWT_REFRESH_SECRET="your-refresh-secret"
AUTH_DB_TYPE="sqlite"
AUTH_DB_URL="auth.db"
AUTH_JWT_ISSUER="my-app"
AUTH_ACCESS_TOKEN_TTL="15m"
AUTH_REFRESH_TOKEN_TTL="168h"
AUTH_APP_NAME="My Application"
AUTH_LOG_LEVEL="info"
```

Then use:

```go
// Automatically loads from environment
authService, err := auth.NewFromEnv()
```

### Advanced Configuration

```go
config := &auth.AuthConfig{
    DatabasePath:     "auth.db",
    JWTSecret:        "your-jwt-secret",
    JWTRefreshSecret: "your-refresh-secret",
    JWTIssuer:        "my-app",
    AccessTokenTTL:   15 * time.Minute,
    RefreshTokenTTL:  7 * 24 * time.Hour,
    AppName:          "My Application",
    Version:          "1.0.0",
    LogLevel:         "info",
}

authService, err := auth.NewWithConfig(config)
```

---

## Monitoring & Health Checks

```go
// Health checks
err := authService.Health()
info := authService.GetSystemInfo()

// Metrics
metrics := authService.GetMetrics()
collector := authService.MetricsCollector()
loginRate := collector.GetLoginSuccessRate()

// Structured logging
logger := authService.Logger()
logger.Info("Custom operation", map[string]interface{}{
    "user_id": userID,
    "operation": "login",
})

// HTTP endpoints for monitoring
mux := http.NewServeMux()
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    if err := authService.Health(); err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
})
```

---

## Database Support

### SQLite (Default)

```go
authService, err := auth.NewSQLite("auth.db", "jwt-secret")
```

### PostgreSQL

```go
authService, err := auth.NewPostgres("postgres://user:pass@localhost/db", "jwt-secret")
```

### In-Memory (Testing)

```go
authService, err := auth.NewInMemory("jwt-secret")
```

---

## Migration from v1

### Automatic Compatibility

Your existing v1 code continues to work:

```go
// v1 code still works!
cfg := auth.Config{
    Storage: storage,
    JWT: auth.JWTConfig{
        AccessSecret: []byte("secret"),
        // ... v1 config
    },
}

authService, err := auth.NewAuthService(cfg)
```

### Migration Tools

Use the automated migration tool:

```bash
go install github.com/pragneshbagary/go-auth/cmd/migrate@v2.0.0
migrate -path . -output migration-report.txt
```

### Migration Guide

See the comprehensive [Migration Guide](MIGRATION.md) for detailed instructions.

---

## Examples

Explore the comprehensive examples:

- [Basic Usage](examples/basic_usage_example.go)
- [Advanced Features](examples/advanced_usage_example.go)
- [Framework Integration](examples/)
- [Middleware Usage](examples/middleware_example.go)
- [Error Handling](examples/error_handling_example.go)
- [Token Management](examples/token_management_example.go)
- [User Management](examples/users_management_example.go)
- [Migration Examples](examples/migration_example.go)

---

## Security Features

- **Argon2id Password Hashing** - Industry-standard secure hashing
- **JWT with Refresh Tokens** - Secure token-based authentication
- **SQL Injection Protection** - Parameterized queries and input validation
- **Timing Attack Resistance** - Constant-time comparisons
- **Rate Limiting Support** - Built-in protection against brute force
- **Secure Token Storage** - Automatic token revocation and cleanup
- **Input Validation** - Comprehensive validation for all inputs
- **Security Headers** - Automatic security header management

---

## Performance

- **Connection Pooling** - Automatic database connection management
- **Batch Operations** - Efficient batch token validation
- **Caching Support** - Optional Redis caching for improved performance
- **Optimized Queries** - Efficient database operations
- **Memory Management** - Optimized memory usage and cleanup

---

## Testing

Run the comprehensive test suite:

```bash
# Unit tests
go test ./...

# With coverage
go test ./... -cover

# Race condition testing
go test ./... -race

# Benchmarks
go test ./... -bench=.

# Security tests
go test ./pkg/auth -run TestSecurity

# Performance tests
go test ./pkg/auth -run TestPerformance
```

### Test Coverage

- **Unit Tests**: 95%+ coverage
- **Integration Tests**: Database operations, middleware, components
- **Security Tests**: SQL injection, timing attacks, data exposure
- **Performance Tests**: Load testing, benchmarks, memory usage
- **Compatibility Tests**: v1 API backward compatibility

---

## Documentation

- [API Documentation](https://pkg.go.dev/github.com/pragneshbagary/go-auth)
- [Migration Guide](MIGRATION.md)
- [Examples](examples/)
- [Release Notes](https://github.com/pragneshbagary/go-auth/releases)

---

## Contributing

Contributions are welcomed! Please see [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Changelog

### v2.0.0 (Latest)
- **New**: Component-based architecture (Users, Tokens, Middleware)
- **New**: Framework-specific middleware (Gin, Echo, Fiber)
- **New**: Environment-based configuration
- **New**: Built-in monitoring and metrics
- **New**: Automated migration tools
- **New**: Advanced security features
- **New**: Performance optimizations
- **Backward Compatible**: All v1 APIs continue to work
- **Documentation**: Comprehensive examples and guides

### v1.x
- Basic authentication functionality
- JWT token management
- SQLite storage support

[View Full Changelog](https://github.com/pragneshbagary/go-auth/releases)

---

## Support

- [Documentation](https://pkg.go.dev/github.com/pragneshbagary/go-auth)
- [Issue Tracker](https://github.com/pragneshbagary/go-auth/issues)
- [Discussions](https://github.com/pragneshbagary/go-auth/discussions)
- [Email Support](mailto:pragneshbagary1699@gmail.com)
- [Migration Guide](MIGRATION.md) - Upgrade from v1 to v2

---

**If you find go-auth useful, please consider giving it a star on GitHub!**
