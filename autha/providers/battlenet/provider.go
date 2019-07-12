package battlenet

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/contextgg/go-sdk/gen"
	"github.com/contextgg/go-sdk/httpbuilder"

	"golang.org/x/oauth2"

	"github.com/contextgg/go-sdk/autha"
)

const (
	authURL      string = "https://us.battle.net/oauth/authorize"
	tokenURL     string = "https://us.battle.net/oauth/token"
	endpointUser string = "https://us.battle.net/oauth/userinfo"
)

// CurrentUser the object representing the current discord user
type CurrentUser struct {
	ID        int64  `json:"id"`
	Battletag string `json:"battletag"`
}

type provider struct {
	config *oauth2.Config
}

func (p *provider) Name() string {
	return "battlenet"
}

func (p *provider) BeginAuth(session autha.Session) (string, error) {
	// state for the oauth grant!
	state := gen.RandomString(64)

	// set the state
	session.Set("state", state)

	// generate the url
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

func (p *provider) Authorize(session autha.Session, params autha.Params) (autha.Token, error) {
	state := params.Get("state")
	if len(state) == 0 {
		return nil, errors.New("No state value in params")
	}

	if !autha.SessionHasValue(session, "state", state) {
		return nil, errors.New("Invalid state")
	}

	code := params.Get("code")
	if len(code) == 0 {
		return nil, errors.New("No code value in params")
	}

	token, err := p.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("Invalid token received from provider")
	}

	// TODO what to do with the token.

	return token, nil
}

func (p *provider) LoadIdentity(token autha.Token, session autha.Session) (*autha.Identity, error) {
	t, ok := token.(*oauth2.Token)
	if !ok {
		return nil, errors.New("Wrong token type")
	}

	authType := t.TokenType
	accessToken := t.AccessToken

	// todo get the user!
	var user CurrentUser
	status, err := httpbuilder.New().
		SetURL(endpointUser).
		SetAuthToken(authType, accessToken).
		SetOut(&user).
		Do()
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("Invalid Status Code %d", status)
	}

	id := &autha.Identity{
		Provider: p.Name(),
		ID:       fmt.Sprintf("%d", user.ID),
		Profile:  user,
	}
	return id, nil
}

// NewProvider creates a new Provider
func NewProvider(clientID, clientSecret, callbackURL string, scopes ...string) autha.AuthProvider {
	return &provider{
		config: newConfig(clientID, clientSecret, callbackURL, scopes),
	}
}

func newConfig(clientID, clientSecret, callbackURL string, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	if len(scopes) > 0 {
		for _, scope := range scopes {
			c.Scopes = append(c.Scopes, scope)
		}
	} else {
		c.Scopes = []string{ScopeIdentify}
	}

	return c
}
