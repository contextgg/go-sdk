package autha

import "errors"

// ErrUserNotFound when a user doesn't exist
var ErrUserNotFound = errors.New("User not found")

// User of the system
type User struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	Connection string `json:"connection"`
	Provider   string `json:"provider"`
}

// UserProvider is the common interface for users
type UserProvider interface {
	Login(string, *Identity, Token) (*User, error)
}
