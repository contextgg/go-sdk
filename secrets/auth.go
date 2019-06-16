package secrets

// BasicAuthCredentials for credentials
type BasicAuthCredentials struct {
	Username string
	Password string
}

// LoadBasicAuth will load the secrets from disk
func LoadBasicAuth(prefix string) *BasicAuthCredentials {
	usernameKey := prefix + "-basic-auth-username"
	username := MustReadSecret(usernameKey, "")
	if username == "" {
		return nil
	}

	passwordKey := prefix + "-basic-auth-password"
	password, err := ReadSecret(passwordKey)
	if err != nil {
		return nil
	}

	return &BasicAuthCredentials{
		Username: username,
		Password: password,
	}
}
