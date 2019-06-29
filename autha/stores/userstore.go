package stores

import (
	"time"

	"github.com/gorilla/sessions"

	"github.com/contextgg/go-sdk/autha"
)

type userStore struct {
	cookieStore *sessions.CookieStore
}

// NewUserStore creates a new session store
func NewUserStore(keypairs ...[]byte) (autha.UserStore, error) {
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
