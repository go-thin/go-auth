package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL
)`

// Store is a PostgreSQL-backed implementation of storage.Storage.
type Store struct {
	db *sql.DB
}

// New opens a PostgreSQL connection using connString and ensures the users table exists.
// Example: postgres.New("postgres://user:pass@localhost/mydb")
func New(connString string) (*Store, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	if _, err := db.Exec(createUsersTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: create table: %w", err)
	}
	return &Store{db: db}, nil
}

// Close releases the underlying database connection pool.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying *sql.DB for advanced use cases (e.g., migrations, transactions).
func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) CreateUser(user models.User) error {
	_, err := s.db.Exec(
		`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("postgres: create user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password_hash FROM users WHERE username = $1`,
		username,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("postgres: get user: %w", err)
	}
	return &u, nil
}
