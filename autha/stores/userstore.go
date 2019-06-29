package stores

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/contextgg/go-sdk/autha"
)

const userStoreKey = "_ctx_user"

type userStore struct {
	cookieStore *sessions.CookieStore
}

func (s *userStore) Save(user *autha.User, w http.ResponseWriter, r *http.Request) error {
	// load up the session
	sess, err := s.cookieStore.Get(r, userStoreKey)
	if err != nil {
		return err
	}

	sess.Values["id"] = user.ID
	sess.Values["state"] = user.State
	sess.Values["connection"] = user.Connection
	sess.Values["provider"] = user.Provider

	return sess.Save(r, w)
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

	return &userStore{
		cookieStore: cookieStore,
	}, nil
}
