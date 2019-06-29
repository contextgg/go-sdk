package user

import "github.com/contextgg/go-sdk/autha"

type provider struct {
}

// NewProvider creates a new user provider
func NewProvider(name string) autha.UserProvider {
	return &provider{}
}
