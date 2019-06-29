package user

import "github.com/contextgg/go-sdk/autha"

type provider struct {
}

func (p *provider) Login(connection string, id *autha.Identity, token autha.Token) (*autha.User, error) {
	// fake the login for now!
	u := &autha.User{
		Connection: connection,
		Provider:   "none",
		ID:         id.ID,
		State:      "OK",
	}
	return u, nil
}

// NewProvider creates a new user provider
func NewProvider(name string) autha.UserProvider {
	return &provider{}
}
