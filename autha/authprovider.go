package autha

import "context"

// Identity represents the identity of a discord user
type Identity struct {
	Provider string      `json:"provider"`
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Profile  interface{} `json:"profile,omitempty"`
}

// AuthProvider is the common interface for doing auth
type AuthProvider interface {
	// Name of the provider IE Discord, Twitter, Twitch
	Name() string

	// BeginAuth the start of a token exchange
	BeginAuth(context.Context, Session) (string, error)

	// Authorize confirm everything is ok
	Authorize(context.Context, Session, Params) (Token, error)

	// LoadIdentity will try to load the current users identity
	LoadIdentity(context.Context, Token, Session) (*Identity, error)
}
