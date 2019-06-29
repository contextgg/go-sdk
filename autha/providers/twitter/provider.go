package twitter

import (
	"github.com/mrjones/oauth"

	"github.com/contextgg/go-sdk/autha"
)

var (
	requestURL      = "https://api.twitter.com/oauth/request_token"
	authorizeURL    = "https://api.twitter.com/oauth/authorize"
	authenticateURL = "https://api.twitter.com/oauth/authenticate"
	tokenURL        = "https://api.twitter.com/oauth/access_token"
	endpointProfile = "https://api.twitter.com/1.1/account/verify_credentials.json"
)

type provider struct {
	consumer *oauth.Consumer
}

func (p *provider) Name() string {
	return "twitter"
}

func (p *provider) BeginAuth(session autha.Session) string {
	return ""
}

func (p *provider) Authorize(session autha.Session, params autha.Params) (autha.Token, error) {
	return nil, nil
}

func (p *provider) LoadIdentity(token autha.Token, session autha.Session) (*autha.Identity, error) {
	return nil, nil
}

// NewProvider creates a new Provider
func NewProvider(clientID, clientSecret, callbackURL string) autha.AuthProvider {
	return &provider{
		consumer: newConsumer(clientID, clientSecret, callbackURL, true),
	}
}

func newConsumer(clientID, clientSecret, authURL string, debug bool) *oauth.Consumer {
	c := oauth.NewConsumer(
		clientID,
		clientSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   requestURL,
			AuthorizeTokenUrl: authURL,
			AccessTokenUrl:    tokenURL,
		})

	c.Debug(debug)
	return c
}
