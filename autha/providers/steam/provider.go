package twitter

import (
	"fmt"
	"net/url"

	"github.com/contextgg/go-sdk/autha"
)

const (
	// Steam API Endpoints
	apiLoginEndpoint       = "https://steamcommunity.com/openid/login"
	apiUserSummaryEndpoint = "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s"

	// OpenID settings
	openIDMode       = "checkid_setup"
	openIDNs         = "http://specs.openid.net/auth/2.0"
	openIDIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
)

type provider struct {
	apiKey      string
	callbackURL string
}

func (p *provider) Name() string {
	return "steam"
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
func NewProvider(apiKey string, callbackURL string) autha.AuthProvider {
	return &provider{
		apiKey:      apiKey,
		callbackURL: callbackURL,
	}
}

// getAuthURL is an internal function to build the correct
// authentication url to redirect the user to Steam.
func (p *provider) getAuthURL() (*url.URL, error) {
	callbackURL, err := url.Parse(p.callbackURL)
	if err != nil {
		return nil, err
	}

	urlValues := map[string]string{
		"openid.claimed_id": openIDIdentifier,
		"openid.identity":   openIDIdentifier,
		"openid.mode":       openIDMode,
		"openid.ns":         openIDNs,
		"openid.realm":      fmt.Sprintf("%s://%s", callbackURL.Scheme, callbackURL.Host),
		"openid.return_to":  callbackURL.String(),
	}

	u, err := url.Parse(apiLoginEndpoint)
	if err != nil {
		return nil, err
	}

	v := u.Query()
	for key, value := range urlValues {
		v.Set(key, value)
	}
	u.RawQuery = v.Encode()

	return u, nil
}
