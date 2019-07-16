package user

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/contextgg/go-sdk/autha"
	"github.com/contextgg/go-sdk/httpbuilder"
)

const headerName = "_ctx_type_name_"

type provider struct {
	functionName string
	namespace    uuid.UUID
	username     string
	password     string
}

// connection string, id *autha.Identity, token autha.Token
func (p *provider) Login(m *autha.UserLogin) (*autha.IdentityID, error) {
	ns := uuid.NewSHA1(p.namespace, []byte(m.Connection))
	id := uuid.NewSHA1(ns, []byte(m.Identity.ID))

	// Inject an aggregate id.
	raw := struct {
		*autha.UserLogin
		AggregateID string
	}{
		m,
		id.String(),
	}

	status, err := httpbuilder.NewFaaS().
		SetAuthBasic(p.username, p.password).
		SetFunction(p.functionName).
		SetMethod(http.MethodPost).
		SetBody(&raw).
		AddHeader(headerName, "Login").
		Do()

	if err != nil {
		return nil, err
	}
	if status < 200 && status > 400 {
		return nil, errors.New("Invalid http status")
	}
	return &autha.IdentityID{ID: id.String()}, nil
}

// NewProvider creates a new user provider
func NewProvider(functionName, authDNS, username, password string) autha.UserProvider {
	base := uuid.NewSHA1(uuid.NameSpaceURL, []byte(authDNS))

	return &provider{
		functionName: functionName,
		namespace:    base,
		username:     username,
		password:     password,
	}
}
