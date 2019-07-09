package user

import (
	"errors"
	"net/http"

	"github.com/contextgg/go-sdk/autha"
	"github.com/contextgg/go-sdk/httpbuilder"
)

type provider struct {
	functionName string
	username     string
	password     string
}

// connection string, id *autha.Identity, token autha.Token
func (p *provider) Login(m *autha.UserLogin) (*autha.IdentityID, error) {
	var result autha.IdentityID

	status, err := httpbuilder.NewFaaS().
		SetAuthBasic(p.username, p.password).
		SetFunction(p.functionName).
		SetMethod(http.MethodPost).
		SetBody(m).
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
func NewProvider(functionName, username, password string) autha.UserProvider {
	return &provider{
		functionName: functionName,
		username:     username,
		password:     password,
	}
}
