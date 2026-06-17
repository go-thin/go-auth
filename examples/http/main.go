package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-thin/go-auth/internal/storage/memory"
	"github.com/go-thin/go-auth/pkg/auth"
)

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	svc, err := auth.New(auth.Config{
		Storage: memory.NewInMemoryStorage(),
		JWT: auth.JWTConfig{
			AccessSecret:    []byte("dev-access-secret-change-me"),
			RefreshSecret:   []byte("dev-refresh-secret-change-me"),
			Issuer:          "go-auth-http-example",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
			SigningMethod:   auth.HS256,
		},
	})
	if err != nil {
		log.Fatalf("create auth service: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler(svc))
	mux.HandleFunc("/login", loginHandler(svc))
	mux.HandleFunc("/protected", protectedHandler(svc))

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func registerHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "use POST")
			return
		}

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		user, err := svc.Register(auth.RegisterPayload{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, user)
	}
}

func loginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "use POST")
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		tokens, err := svc.Login(req.Username, req.Password, nil)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, tokens)
	}
}

func protectedHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "use GET")
			return
		}

		token, ok := bearerToken(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "missing Bearer token")
			return
		}

		claims, err := svc.ValidateAccessToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"message": "welcome to the protected route",
			"claims":  claims,
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	return token, token != ""
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
