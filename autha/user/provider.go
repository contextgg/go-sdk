package user

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/contextgg/go-sdk/autha"
	"github.com/contextgg/go-sdk/httpbuilder"
)

type provider struct {
	functionName string
	namespace    uuid.UUID
	username     string
	password     string
}

// connection string, id *autha.Identity, token autha.Token
func (p *provider) Login(m *autha.UserLogin) (*autha.IdentityID, error) {
	var result autha.IdentityID

	aggregateID := uuid.NewSHA1(
			uuid.NewSHA1(p.namespace, []byte(m.Connection)), []byte(m.Identity.ID)
		).
		String()

	// Inject an aggregate id.
	raw := struct {
		*autha.UserLogin
		AggregateID string
	}{
		m,
		aggregateID,
	}

	status, err := httpbuilder.NewFaaS().
		SetAuthBasic(p.username, p.password).
		SetFunction(p.functionName).
		SetMethod(http.MethodPost).
		SetBody(&raw).
		SetOut(&result).
		Do()

	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.New("Invalid http status")
	}
	return &result, nil
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
