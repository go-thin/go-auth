# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Current State

The library was briefly taken in an over-engineered direction (component-based architecture, multiple framework integrations, SQL backends). That was reverted. The project is now back to a focused, understandable core:

- `AuthService` with four methods: `Register`, `Login`, `ValidateAccessToken`, `ValidateRefreshToken`
- Argon2id password hashing (OWASP-recommended parameters)
- JWT access + refresh token pair (HS256/HS384/HS512/RS256)
- `go-events` bus integration — publishes `UserRegistered` and `UserLoggedIn` events
- Pluggable `storage.Storage` interface — bring your own backend
- Structured `AuthError` types with HTTP status mapping
- In-memory storage implementation for testing (internal package)
