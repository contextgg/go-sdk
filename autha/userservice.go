package autha

import (
	"context"
)

// PersistUser struct
type PersistUser struct {
	Provider   string   `json:"provider"`
	Connection string   `json:"connection"`
	Token      Token    `json:"token"`
	Profile    *Profile `json:"profile"`

	PrimaryUserID      *UserID `json:"primary_user_id,omitempty"`
	IsConnectedProfile bool    `json:"is_connected_profile"`
}

// NewPersistUser return a new persist user struct
func NewPersistUser(provider, connection string, token Token, profile *Profile, primaryUserID *UserID, isConnectedProfile bool) *PersistUser {
	return &PersistUser{
		Provider:           provider,
		Connection:         connection,
		Token:              token,
		Profile:            profile,
		PrimaryUserID:      primaryUserID,
		IsConnectedProfile: isConnectedProfile,
	}
}

// UserService is the common interface for users
type UserService interface {
	Persist(context.Context, *PersistUser) (*UserID, error)
	Connect(context.Context, *UserID, *PersistUser) error
}
