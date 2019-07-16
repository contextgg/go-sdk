package discord

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/contextgg/go-sdk/gen"
	"github.com/contextgg/go-sdk/httpbuilder"

	"golang.org/x/oauth2"

	"github.com/contextgg/go-sdk/autha"
)

const (
	authURL      = "https://discordapp.com/api/oauth2/authorize"
	tokenURL     = "https://discordapp.com/api/oauth2/token"
	userEndpoint = "https://discordapp.com/api/users/@me"
)

const (
	// ScopeIdentify allows /users/@me without email
	ScopeIdentify string = "identify"
	// ScopeEmail enables /users/@me to return an email
	ScopeEmail string = "email"
	// ScopeConnections allows /users/@me/connections to return linked Twitch and YouTube accounts
	ScopeConnections string = "connections"
	// ScopeGuilds allows /users/@me/guilds to return basic information about all of a user's guilds
	ScopeGuilds string = "guilds"
	// ScopeJoinGuild allows /invites/{invite.id} to be used for joining a user's guild
	ScopeJoinGuild string = "guilds.join"
	// ScopeGroupDMjoin allows your app to join users to a group dm
	ScopeGroupDMjoin string = "gdm.join"
	// ScopeBot for oauth2 bots, this puts the bot in the user's selected guild by default
	ScopeBot string = "bot"
	// ScopeWebhook this generates a webhook that is returned in the oauth token response for authorization code grants
	ScopeWebhook string = "webhook.incoming"
)

// CurrentUser the object representing the current discord user
type CurrentUser struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	Avatar        *string `json:"avatar"`
	Bot           bool    `json:"bot"`
	MFAEnabled    bool    `json:"mfa_enabled"`
	Locale        string  `json:"locale"`
	Verified      bool    `json:"verified"`
	Email         string  `json:"email"`
	Flags         int     `json:"flags"`
	PremiumType   int     `json:"premium_type"`
}

type provider struct {
	config *oauth2.Config
}

func (p *provider) Name() string {
	return "discord"
}

func (p *provider) BeginAuth(ctx context.Context, session autha.Session) (string, error) {
	// state for the oauth grant!
	state := gen.RandomString(64)

	// set the state
	session.Set("state", state)

	// generate the url
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

func (p *provider) Authorize(ctx context.Context, session autha.Session, params autha.Params) (autha.Token, error) {
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

func (p *provider) LoadIdentity(ctx context.Context, token autha.Token, session autha.Session) (*autha.Identity, error) {
	t, ok := token.(*oauth2.Token)
	if !ok {
		return nil, errors.New("Wrong token type")
	}

	authType := t.TokenType
	accessToken := t.AccessToken

	// todo get the user!
	var user CurrentUser
	status, err := httpbuilder.New().
		SetURL(userEndpoint).
		SetAuthToken(authType, accessToken).
		SetOut(&user).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("Invalid Status Code %d", status)
	}

	id := &autha.Identity{
		Provider: p.Name(),
		ID:       user.ID,
		Username: fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
		Email:    user.Email,
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
