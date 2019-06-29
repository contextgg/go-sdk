package stores

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/contextgg/go-sdk/autha"
)

type sessionStore struct {
	cookieStore *sessions.CookieStore
}

func (s *sessionStore) Load(string, *http.Request) (autha.Session, error) {
	return nil, nil
}

func (s *sessionStore) Save(string, http.ResponseWriter, *http.Request) error {
	return nil
}

// NewSessionStore creates a new session store
func NewSessionStore(keypairs ...[]byte) (autha.SessionStore, error) {
	if len(keypairs) == 0 {
		return nil, ErrNeedKeys
	}

	// create a new session store!
	cookieStore := sessions.NewCookieStore(keypairs...)

	// 12 hours, set this to something because if we don't then sessions
	// may never expire as long as the browser remains opened.
	cookieStore.MaxAge(int((time.Hour * 12) / time.Second))
	cookieStore.Options.HttpOnly = true
	// cookieStore.Options.Secure = true

	return &sessionStore{
		cookieStore: cookieStore,
	}, nil
}
