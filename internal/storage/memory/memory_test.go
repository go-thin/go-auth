package memory

import (
	"testing"

	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"
)

var _ storage.Storage = (*InMemoryStorage)(nil)

func TestDeleteUser_RemovesUser(t *testing.T) {
	store := NewInMemoryStorage()

	user := models.User{
		ID:           "u-delete",
		Username:     "delete-me",
		Email:        "delete-me@example.com",
		PasswordHash: "h",
	}
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
