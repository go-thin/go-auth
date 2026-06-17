package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"
	_ "modernc.org/sqlite"
)

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL
)`

// Store is a SQLite-backed implementation of storage.Storage.
type Store struct {
	db *sql.DB
}

// New opens a SQLite database at dsn and ensures the users table exists.
// Use ":memory:" for an ephemeral in-process database.
func New(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: ping: %w", err)
	}
	if _, err := db.Exec(createUsersTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: create table: %w", err)
	}
	return &Store{db: db}, nil
}

// Close releases the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying *sql.DB for advanced use cases (e.g., migrations, transactions).
func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) CreateUser(user models.User) error {
	_, err := s.db.Exec(
		`INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("sqlite: create user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password_hash FROM users WHERE username = ?`,
		username,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("sqlite: get user: %w", err)
	}
	return &u, nil
}

func (s *Store) DeleteUser(id string) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("sqlite: delete user: %w", err)
	}
	return nil
}
