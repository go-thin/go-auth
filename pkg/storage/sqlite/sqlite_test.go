package sqlite_test

import (
	"testing"

	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"
	"github.com/go-thin/go-auth/pkg/storage/sqlite"
)

func newTestStore(t *testing.T) *sqlite.Store {
	t.Helper()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("sqlite.New: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestCreateUser_and_GetUserByUsername(t *testing.T) {
	store := newTestStore(t)

	user := models.User{
		ID:           "u1",
		Username:     "alice",
		Email:        "alice@example.com",
		PasswordHash: "hashvalue",
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	got, err := store.GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("ID = %q, want %q", got.ID, user.ID)
	}
	if got.Email != user.Email {
		t.Errorf("Email = %q, want %q", got.Email, user.Email)
	}
	if got.PasswordHash != user.PasswordHash {
		t.Errorf("PasswordHash = %q, want %q", got.PasswordHash, user.PasswordHash)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	store := newTestStore(t)

	user := models.User{ID: "u1", Username: "bob", Email: "bob@example.com", PasswordHash: "h"}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}

	dup := models.User{ID: "u2", Username: "bob", Email: "other@example.com", PasswordHash: "h"}
	if err := store.CreateUser(dup); err == nil {
		t.Error("expected error for duplicate username, got nil")
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	store := newTestStore(t)

	_, err := store.GetUserByUsername("nobody")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != storage.ErrUserNotFound {
		t.Errorf("error = %v, want storage.ErrUserNotFound", err)
	}
}

func TestDeleteUser_RemovesUser(t *testing.T) {
	store := newTestStore(t)

	user := models.User{ID: "u-delete", Username: "delete-me", Email: "delete-me@example.com", PasswordHash: "h"}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	if err := store.DeleteUser("u-delete"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	_, err := store.GetUserByUsername("delete-me")
	if err == nil {
		t.Fatal("expected user to be deleted")
	}
	if err != storage.ErrUserNotFound {
		t.Errorf("error = %v, want storage.ErrUserNotFound", err)
	}
}
