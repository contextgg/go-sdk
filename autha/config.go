package autha

import (
	"net/http"
	"net/url"
)

var (
	// StatusLogin when user has already been created
	StatusLogin = "login"
)

// Config for our common authentication pattern
type Config struct {
	connection   string
	loginURL     string
	errorURL     string
	authProvider AuthProvider
	sessionStore SessionStore
	userProvider UserProvider
	userStore    UserStore
}

// NewConfig will return a new config
func NewConfig(
	connection string,
	loginURL string,
	errorURL string,
	authProvider AuthProvider,
	sessionStore SessionStore,
	userProvider UserProvider,
	userStore UserStore) *Config {
	return &Config{
		connection:   connection,
		loginURL:     loginURL,
		errorURL:     errorURL,
		authProvider: authProvider,
		sessionStore: sessionStore,
		userProvider: userProvider,
		userStore:    userStore,
	}
}

func (c *Config) fullErrorURL(errorType string) string {
	str := c.errorURL
	if len(str) == 0 {
		str = c.loginURL
	}
	u, _ := url.Parse(str)

	q := u.Query()
	q.Set("error.type", errorType)

	u.RawQuery = q.Encode()

	return u.String()
}

// Begin the auth method
func (c *Config) Begin(w http.ResponseWriter, r *http.Request) {
	// get the current session!
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		return
	}

	url := c.authProvider.BeginAuth(session)

	// save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// Callback for the provider
func (c *Config) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		return
	}

	token, err := c.authProvider.Authorize(session, r.URL.Query())
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
		return
	}

	id, err := c.authProvider.LoadIdentity(token, session)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
		return
	}

	// TODO what about if we are linking?

	// if we have an id store it!
	user, err := c.userProvider.Login(c.connection, id, token)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("user"), http.StatusFound)
		return
	}

	if err := c.userStore.Save(user, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("user"), http.StatusFound)
		return
	}

	// // save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		return
	}

	// what's the next step?
	http.Redirect(w, r, c.loginURL, http.StatusFound)
}
