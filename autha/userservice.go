package autha

import (
	"context"
)

// BaseUser struct
type BaseUser struct {
	Provider   string   `json:"provider"`
	Connection string   `json:"connection"`
	Token      Token    `json:"token"`
	Profile    *Profile `json:"profile"`
}

// PersistUser struct
type PersistUser struct {
	*BaseUser

	PrimaryUserID      *UserID `json:"primary_user_id,omitempty"`
	IsConnectedProfile bool    `json:"is_connected_profile"`
}

// ConnectUser struct
type ConnectUser struct {
	*BaseUser
}

// NewPersistUser return a new persist user struct
func NewPersistUser(provider, connection string, token Token, profile *Profile, primaryUserID *UserID, isConnectedProfile bool) *PersistUser {
	return &PersistUser{
		BaseUser: &BaseUser{
			Provider:   provider,
			Connection: connection,
			Token:      token,
			Profile:    profile,
		},
		PrimaryUserID:      primaryUserID,
		IsConnectedProfile: isConnectedProfile,
	}
}

// NewConnectUser return a new persist user struct
func NewConnectUser(provider, connection string, token Token, profile *Profile) *ConnectUser {
	return &ConnectUser{
		BaseUser: &BaseUser{
			Provider:   provider,
			Connection: connection,
			Token:      token,
			Profile:    profile,
		},
	}
}

// UserService is the common interface for users
type UserService interface {
	Persist(context.Context, *PersistUser) (*UserID, error)
	Connect(context.Context, *UserID, *ConnectUser) error
}
