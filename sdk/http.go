package sdk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path"
	"time"
)

// ContextSignatureHeader the key for the HMAC Signature
const ContextSignatureHeader = "X-Context-Signature"

// DefaultHTTPClient for making things easier
var DefaultHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// Post for http
func Post(functionName string, data interface{}, result interface{}) error {
	gatewayURL := os.Getenv("gateway_url")
	if len(gatewayURL) == 0 {
		return errors.New("No gateway defined")
	}

	prefix := "function" // What about async?
	url := path.Join(gatewayURL, prefix, functionName)

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// build the request!
	bodyReader := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if HmacEnabled() {
		key, err := ReadSecret("payload-secret")
		if err != nil {
			return err
		}

		digest := Sign(body, []byte(key))
		req.Header.Add(ContextSignatureHeader, "sha1="+hex.EncodeToString(digest))
	}

	if AuthEnabled() {
		req.SetBasicAuth(usernameSecret, passwordSecret)
	}

	res, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		// parse the response
		return json.NewDecoder(res.Body).Decode(&result)
	}

	// here !
	return errors.New("Invalid response code")
}
