package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMiddlewareTestService(t *testing.T) *Service {
	t.Helper()

	svc, err := New(Config{
		Storage: newSimpleStorage(),
		JWT: JWTConfig{
			AccessSecret:    []byte("access-secret"),
			RefreshSecret:   []byte("refresh-secret"),
			AccessTokenTTL:  time.Minute,
			RefreshTokenTTL: time.Hour,
			SigningMethod:   HS256,
		},
	})
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}

	return svc
}

func newMiddlewareTestToken(t *testing.T, svc *Service) string {
	t.Helper()

	if _, err := svc.Register(RegisterPayload{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}

	tokens, err := svc.Login("alice", "password123", map[string]interface{}{
		"role": "admin",
	})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	return tokens.AccessToken
}

func TestMiddleware_StoresClaimsInRequestContext(t *testing.T) {
	svc := newMiddlewareTestService(t)
	token := newMiddlewareTestToken(t, svc)

	handler := Middleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("claims missing from request context")
		}
		if claims["username"] != "alice" {
			t.Errorf("username claim = %v, want %q", claims["username"], "alice")
		}
		if claims["role"] != "admin" {
			t.Errorf("role claim = %v, want %q", claims["role"], "admin")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestMiddleware_RejectsMissingAuthorizationHeader(t *testing.T) {
	svc := newMiddlewareTestService(t)
	called := false
	handler := Middleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler was called")
	}
}

func TestMiddleware_RejectsInvalidAuthorizationHeader(t *testing.T) {
	svc := newMiddlewareTestService(t)
	called := false
	handler := Middleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token not-a-bearer-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler was called")
	}
}

func TestMiddleware_RejectsInvalidBearerToken(t *testing.T) {
	svc := newMiddlewareTestService(t)
	called := false
	handler := Middleware(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler was called")
	}
}
