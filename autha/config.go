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
	sessionStore SessionStore
	userStore    UserStore
	authProvider AuthProvider
	userService  UserService
}

// NewConfig will return a new config
func NewConfig(
	connection string,
	loginURL string,
	errorURL string,
	sessionStore SessionStore,
	userStore UserStore,
	authProvider AuthProvider,
	userService UserService,
) *Config {
	return &Config{
		connection:   connection,
		loginURL:     loginURL,
		errorURL:     errorURL,
		sessionStore: sessionStore,
		userStore:    userStore,
		authProvider: authProvider,
		userService:  userService,
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
	ctx := r.Context()

	// get the current session!
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Load: %w", err))
		return
	}

	url, err := c.authProvider.BeginAuth(ctx, session, r.URL.Query())
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("auth"), http.StatusFound)
		log.Print(fmt.Errorf("Error Begin Auth: %w", err))
		return
	}

	// save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Save: %w", err))
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// Callback for the provider
func (c *Config) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Load: %w", err))
		return
	}

	token, err := c.authProvider.Authorize(ctx, session, r.URL.Query())
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
		log.Print(fmt.Errorf("Error Authorize: %w", err))
		return
	}

	profile, err := c.authProvider.LoadProfile(ctx, token, session)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("identity"), http.StatusFound)
		log.Print(fmt.Errorf("Error Load Identity: %w", err))
		return
	}

	currentUserID, isConnected, err := c.userStore.Load(r)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("identity"), http.StatusFound)
		log.Print(fmt.Errorf("Error loading current user id: %w", err))
		return
	}

	// if we are linking we need to tell the user this is a secondary account!
	pu := NewPersistUser(
		c.authProvider.Name(),
		c.connection,
		token,
		profile,
		currentUserID,
		isConnected,
	)
	// if we have an id store it!
	id, err := c.userService.Persist(r.Context(), pu)
	if err != nil {
		http.Redirect(w, r, c.fullErrorURL("user"), http.StatusFound)
		log.Print(fmt.Errorf("Error Login: %w", err))
		return
	}

	if isConnected && *currentUserID != *id {
		cu := NewConnectUser(
			c.authProvider.Name(),
			c.connection,
			token,
			profile,
		)
		// we need to connect the accounts.
		if err := c.userService.Connect(r.Context(), currentUserID, cu); err != nil {
			http.Redirect(w, r, c.fullErrorURL("user"), http.StatusFound)
			log.Print(fmt.Errorf("Error connecting profiles: %w", err))
			return
		}
	} else {
		// this is to auth the user!
		if err := c.userStore.Save(id, w, r); err != nil {
			http.Redirect(w, r, c.fullErrorURL("id"), http.StatusFound)
			log.Print(fmt.Errorf("Error Profile Store Save: %w", err))
			return
		}
	}

	// save the session
	if err := c.sessionStore.Save(session, w, r); err != nil {
		http.Redirect(w, r, c.fullErrorURL("session"), http.StatusFound)
		log.Print(fmt.Errorf("Error Session Save: %w", err))
		return
	}

	// what's the next step?
	http.Redirect(w, r, c.loginURL, http.StatusFound)
}
