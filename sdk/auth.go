package sdk

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

var (
	usernameSecret string
	passwordSecret string
)

func init() {
	// optional, add error handling, or read on each request / use sync.Once()

	usernameRaw := MustReadSecret("fn-basic-auth-username", "")
	usernameSecret = strings.TrimSpace(string(usernameRaw))

	passwordRaw := MustReadSecret("fn-basic-auth-password", "")
	passwordSecret = strings.TrimSpace(string(passwordRaw))
}

// AuthEnabled uses validate_hmac env-var to verify if the
// feature is disabled
func AuthEnabled() bool {
	if val, exists := os.LookupEnv("validate_auth"); exists {
		return val != "false" && val != "0"
	}
	return true
}

// EnsureAuth a function invocation
func EnsureAuth(w http.ResponseWriter, req *http.Request) error {
	if !isAuthorized(req) {
		message := "You must authorize."
		w.Header().Set("", `Basic realm="Restricted"`)
		http.Error(w, message, http.StatusUnauthorized)
		return errors.New(message)
	}
	return nil
}

func isAuthorized(req *http.Request) bool {
	if username, password, ok := req.BasicAuth(); ok &&
		username == usernameSecret &&
		password == passwordSecret {
		return true
	}
	return false
}
