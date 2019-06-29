package autha

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	// StatusLogin when user has already been created
	StatusLogin = "login"
)

// Config for our common authentication pattern
type Config struct {
	connection   string
	loginURL     string
	authProvider AuthProvider
	sessionStore SessionStore
	userProvider UserProvider
	userStore    UserStore
}

// NewConfig will return a new config
func NewConfig(
	connection string,
	loginURL string,
	authProvider AuthProvider,
	sessionStore SessionStore,
	userProvider UserProvider,
	userStore UserStore) *Config {
	return &Config{
		connection:   connection,
		loginURL:     loginURL,
		authProvider: authProvider,
		sessionStore: sessionStore,
		userProvider: userProvider,
		userStore:    userStore,
	}
}

// Begin the auth method
func (c *Config) Begin(w http.ResponseWriter, r *http.Request) {
	// get the current session!
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		// TODO go back to the error page.
		return
	}

	url := c.authProvider.BeginAuth(session)

	// save the session
	if err := c.sessionStore.Save(c.connection, w, r); err != nil {
		// TODO go back to the error page.
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// Callback for the provider
func (c *Config) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionStore.Load(c.connection, r)
	if err != nil {
		// TODO go back to the error page.
		return
	}

	token, err := c.authProvider.Authorize(session, r.URL.Query())
	if err != nil {
		// TODO go back to the error page.
		return
	}

	id, err := c.authProvider.LoadIdentity(token, session)
	if err != nil {
		// TODO go back to the error page.
		return
	}

	data, _ := json.Marshal(id)
	log.Printf("ID: %s", string(data))

	// userSession, err := c.userStore.Load(r)
	// if err != nil {
	// 	// TODO go back to the error page.
	// 	return
	// }

	// // TODO what about if we are linking?

	// // if we have an id store it!
	// user, err := c.userProvider.Login(c.connection, id, token)
	// if err != nil {
	// 	// TODO go back to the error page.
	// 	return
	// }

	// if err != c.userStore.Save(w, r, id, user.Status); err != nil {
	// 	// TODO go back to the error page.
	// 	return
	// }

	// // save the session
	// if err := c.sessionStore.Save(c.connection, w, r); err != nil {
	// 	// TODO go back to the error page.
	// 	return
	// }

	// what's the next step?

}
