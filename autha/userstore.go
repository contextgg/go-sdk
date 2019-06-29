package autha

import "net/http"

// UserStore to save and load user
type UserStore interface {
	// Save the user to the request
	Save(*User, http.ResponseWriter, *http.Request) error
}
