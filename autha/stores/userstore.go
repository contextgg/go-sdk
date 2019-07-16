package stores

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/contextgg/go-sdk/autha"
)

const userStoreKey = "_ctx_user"

type userStore struct {
	cookieStore *sessions.CookieStore
}

func (s *userStore) Save(identityID *autha.IdentityID, w http.ResponseWriter, r *http.Request) error {
	if identityID == nil {
		return errors.New("No identityID provided")
	}

	// load up the session
	sess, err := s.cookieStore.Get(r, userStoreKey)
	if err != nil {
		return err
	}

	sess.Values["id"] = identityID.ID

	return sess.Save(r, w)
}

func (s *userStore) Load(r *http.Request) (*autha.IdentityID, bool, error) {
	// load up the session
	sess, err := s.cookieStore.Get(r, userStoreKey)
	if err != nil {
		return nil, false, err
	}

	if sess.IsNew {
		return nil, false, nil
	}

	// try convert it!
	id, ok := sess.Values["id"].(string)
	if !ok {
		return nil, false, errors.New("No ID found in session")
	}

	return &autha.IdentityID{
		ID: id,
	}, true, nil
}

// NewUserStore creates a new session store
func NewUserStore(secure bool, keypairs ...[]byte) (autha.UserStore, error) {
	if len(keypairs) == 0 {
		return nil, ErrNeedKeys
	}

	// create a new session store!
	cookieStore := sessions.NewCookieStore(keypairs...)
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.Secure = secure

	// TODO what about remembering people?
	cookieStore.MaxAge(int((time.Hour * 12) / time.Second))

	return &userStore{
		cookieStore: cookieStore,
	}, nil
}
