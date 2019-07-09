package autha

import "net/http"

// UserStore to save and load user
type UserStore interface {
	// Save the user to the request
	Save(*IdentityID, http.ResponseWriter, *http.Request) error
}
