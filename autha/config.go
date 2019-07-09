package autha

import (
	"fmt"
	"log"
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
		log.Print(fmt.Errorf("Error Session Load: %s", err.Error()))
		return
	}

	url, err := c.authProvider.BeginAuth(session)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("auth"), http.StatusFound)
		log.Print(fmt.Errorf("Error Begin Auth: %s", err.Error()))
		return
	}

	// save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Save: %s", err.Error()))
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// Callback for the provider
func (c *Config) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Load: %s", err.Error()))
		return
	}

	token, err := c.authProvider.Authorize(session, r.URL.Query())
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
		log.Print(fmt.Errorf("Error Authorize: %s", err.Error()))
		return
	}

	identity, err := c.authProvider.LoadIdentity(token, session)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("identity"), http.StatusFound)
		log.Print(fmt.Errorf("Error Load Identity: %s", err.Error()))
		return
	}

	// TODO what about if we are linking?

	// if we have an id store it!
	id, err := c.userProvider.Login(NewUserLogin(c.connection, identity, token))
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("user"), http.StatusFound)
		log.Print(fmt.Errorf("Error Login: %s", err.Error()))
		return
	}

	if err := c.userStore.Save(id, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
		log.Print(fmt.Errorf("Error IdentityID Save: %s", err.Error()))
		return
	}

	// save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Save: %s", err.Error()))
		return
	}

	// what's the next step?
	http.Redirect(w, r, c.loginURL, http.StatusFound)
}
