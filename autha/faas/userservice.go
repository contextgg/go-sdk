package faas

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/contextgg/go-sdk/autha"
	"github.com/contextgg/go-sdk/httpbuilder"
)

const queryName = "type_name"

// PersistCommand struct
type PersistCommand struct {
	*autha.PersistUser

	AggregateID string `json:"aggregate_id"`
}

// ConnectCommand struct
type ConnectCommand struct {
	*autha.ConnectUser

	AggregateID string `json:"aggregate_id"`
}

type provider struct {
	functionName string
	namespace    uuid.UUID
	username     string
	password     string
}

// connection string, id *autha.Identity, token autha.Token
func (p *provider) Persist(ctx context.Context, m *autha.PersistUser) (*autha.UserID, error) {
	ns := uuid.NewSHA1(p.namespace, []byte(m.Connection))
	uid := uuid.NewSHA1(ns, []byte(m.Profile.ID))
	id := uid.String()

	// HACK fix the connected profile if id's are the same!
	if m.IsConnectedProfile && string(*m.PrimaryUserID) == id {
		m.IsConnectedProfile = false
		m.PrimaryUserID = nil
	}

	raw := PersistCommand{
		m,
		id,
	}

	status, err := httpbuilder.NewFaaS().
		SetAuthBasic(p.username, p.password).
		SetFunction(p.functionName).
		SetMethod(http.MethodPost).
		SetBody(&raw).
		AddQuery(queryName, "PersistUserCommand").
		Do(ctx)

	if err != nil {
		return nil, err
	}
	if status < 200 || status > 400 {
		return nil, errors.New("Invalid http status")
	}

	userID := autha.UserID(id)
	return &userID, nil
}

// connection string, id *autha.Identity, token autha.Token
func (p *provider) Connect(ctx context.Context, userID *autha.UserID, m *autha.ConnectUser) error {
	if userID == nil {
		return fmt.Errorf("No userid supplied")
	}

	raw := ConnectCommand{
		m,
		string(*userID),
	}

	status, err := httpbuilder.NewFaaS().
		SetAuthBasic(p.username, p.password).
		SetFunction(p.functionName).
		SetMethod(http.MethodPost).
		SetBody(&raw).
		AddQuery(queryName, "ConnectUserCommand").
		Do(ctx)

	if err != nil {
		return fmt.Errorf("Calling user service failed %w", err)
	}
	if status < 200 || status > 400 {
		return errors.New("Invalid http status")
	}

	return nil
}

// NewService creates a new user provider
func NewService(functionName, authDNS, username, password string) autha.UserService {
	base := uuid.NewSHA1(uuid.NameSpaceURL, []byte(authDNS))

	return &provider{
		functionName: functionName,
		namespace:    base,
		username:     username,
		password:     password,
	}
}
