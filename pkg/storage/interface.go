package storage

import (
	"errors"

	"github.com/go-thin/go-auth/pkg/models"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// Storage is the interface that storage backends must implement.
type Storage interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (*models.User, error)
}
