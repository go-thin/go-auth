package auth

// UserRegistered is published when a new user successfully registers.
type UserRegistered struct {
	UserID   string
	Username string
	Email    string
}

func (e UserRegistered) Topic() string { return "auth.user.registered" }

// UserLoggedIn is published when a user successfully authenticates.
type UserLoggedIn struct {
	UserID   string
	Username string
}

func (e UserLoggedIn) Topic() string { return "auth.user.logged_in" }
