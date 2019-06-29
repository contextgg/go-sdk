package user

import "github.com/contextgg/go-sdk/autha"

type provider struct {
}

func (p *provider) Login(connection string, id *autha.Identity, token autha.Token) (*autha.User, error) {
	return nil, nil
}

// NewProvider creates a new user provider
func NewProvider(name string) autha.UserProvider {
	return &provider{}
}
