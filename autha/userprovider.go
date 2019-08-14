package autha

import (
	"context"
	"errors"
)

// ErrUserNotFound when a user doesn't exist
var ErrUserNotFound = errors.New("User not found")

// User of the system
type User struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	Connection string `json:"connection"`
	Provider   string `json:"provider"`
}

// UserLogin for a user authing
type UserLogin struct {
	Connection string    `json:"connection"`
	Identity   *Identity `json:"identity"`
	Token      Token     `json:"token"`
}

// IdentityID so we can lookup the User by ID
type IdentityID struct {
	ID string `json:"id"`
}

// NewUserLogin create a new login model
func NewUserLogin(connection string, identity *Identity, token Token) *UserLogin {
	return &UserLogin{
		Connection: connection,
		Identity:   identity,
		Token:      token,
	}
}

// UserProvider is the common interface for users
type UserProvider interface {
	Login(context.Context, *UserLogin) (*IdentityID, error)
}
