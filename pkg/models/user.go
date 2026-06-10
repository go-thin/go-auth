package models

// User defines the structure for a user in the system.
// PasswordHash must always store a hashed password, never plaintext.
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
