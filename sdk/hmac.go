package sdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

// CheckMAC verifies hash checksum
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}

// Sign a message with the key and return bytes.
// Note: for human readable output see encoding/hex and
// encode string functions.
func Sign(message, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	signed := mac.Sum(nil)
	return signed
}

// Validate validate an encodedHash taken
// from GitHub via X-Hub-Signature HTTP Header.
// Note: if using another source, just add a 5 letter prefix such as "sha1="
func Validate(bytesIn []byte, encodedHash string, secretKey string) error {
	var validated error

	if len(encodedHash) > 5 {

		hashingMethod := encodedHash[:5]
		if hashingMethod != "sha1=" {
			return fmt.Errorf("unexpected hashing method: %s", hashingMethod)
		}

		messageMAC := encodedHash[5:] // first few chars are: sha1=
		messageMACBuf, _ := hex.DecodeString(messageMAC)

		res := CheckMAC(bytesIn, []byte(messageMACBuf), []byte(secretKey))
		if res == false {
			validated = fmt.Errorf("invalid message digest or secret")
		}
	} else {
		return fmt.Errorf("invalid encodedHash, should have at least 5 characters")
	}

	return validated
}

// HmacEnabled uses validate_hmac env-var to verify if the
// feature is disabled
func HmacEnabled() bool {
	if val, exists := os.LookupEnv("validate_hmac"); exists {
		return val != "false" && val != "0"
	}
	return true
}

// ValidHMAC returns an error if HMAC could not be validated or if
// the signature could not be loaded.
func ValidHMAC(payload *[]byte, secretKey string, digest string) error {
	key, err := ReadSecret(secretKey)
	if err != nil {
		return fmt.Errorf("unable to load HMAC symmetric key, %s", err.Error())
	}

	return validHMACWithSecretKey(payload, key, digest)
}

func validHMACWithSecretKey(payload *[]byte, secretText string, digest string) error {
	validated := Validate(*payload, digest, secretText)

	if validated != nil {
		return fmt.Errorf("unable to validate HMAC")
	}
	return nil
}

func readBool(key string) bool {
	if val, exists := os.LookupEnv(key); exists {
		return val != "false" && val != "0"
	}
	return true
}

func init() {

}
