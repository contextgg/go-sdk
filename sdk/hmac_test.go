package sdk

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func Test_GenerateInvalid_GivesError(t *testing.T) {

	input := []byte("test")
	signature := "ab"
	secretKey := "key"
	err := Validate(input, signature, secretKey)
	if err == nil {
		t.Errorf("expected error when signature didn't have at least 5 characters in length")
		t.Fail()
		return
	}

	wantErr := "invalid encodedHash, should have at least 5 characters"
	if err.Error() != wantErr {
		t.Errorf("want: %s, got: %s", wantErr, err.Error())
		t.Fail()
	}
}

func Test_ValidateWithoutSha1PrefixFails(t *testing.T) {
	digest := "sign this message"
	key := "my key"

	encodedHash := "6791a762f7568f945c2e1e396cea243e944100a6"

	valid := Validate([]byte(digest), encodedHash, key)

	if valid == nil {
		t.Errorf("Expected error due to missing prefix")
		t.Fail()
	}
}

func Test_ValidateWithSha1Prefix(t *testing.T) {
	digest := "sign this message"
	key := "my key"

	encodedHash := "sha1=" + "6791a762f7568f945c2e1e396cea243e944100a6"

	valid := Validate([]byte(digest), encodedHash, key)

	if valid != nil {
		t.Errorf("Expected no error, but got: %s", valid.Error())
		t.Fail()
	}
}

func Test_SignWithKey(t *testing.T) {
	digest := "sign this message"
	key := []byte("my key")

	wantHash := "6791a762f7568f945c2e1e396cea243e944100a6"

	hash := Sign([]byte(digest), key)
	encodedHash := hex.EncodeToString(hash)

	if encodedHash != wantHash {
		t.Errorf("Sign want hash: %s, got: %s", wantHash, encodedHash)
		t.Fail()
	}
}

func Test_validHMACWithSecretKey_validSecret(t *testing.T) {

	data := []byte("Store this string")
	key := []byte("key-goes-here")
	signed := hmac.Sign(data, key)
	digest := fmt.Sprintf("sha1=%s", hex.EncodeToString(signed))
	err := validHMACWithSecretKey(&data, string(key), digest)

	if err != nil {
		t.Errorf("with %s, found error: %s", digest, err)
		t.Fail()
	}
}

func Test_validHMACWithSecretKey_invalidSecret(t *testing.T) {

	data := []byte("Store this string")
	key := []byte("key-goes-here")
	signed := hmac.Sign(data, key)
	digest := fmt.Sprintf("sha1=%s", hex.EncodeToString(signed))
	err := validHMACWithSecretKey(&data, string(key[:4]), digest)

	if err == nil {
		t.Errorf("with %s, expected to find error", digest)
		t.Fail()
	}
}

func Test_HmacEnabled(t *testing.T) {
	tests := []struct {
		title        string
		value        string
		expectedBool bool
	}{
		{
			title:        "environmental variable `validate_hmac` is unset",
			value:        "",
			expectedBool: true,
		},
		{
			title:        "environmental variable `validate_hmac` is set with random value",
			value:        "random",
			expectedBool: true,
		},
		{
			title:        "environmental variable `validate_hmac` is set with explicit `0`",
			value:        "0",
			expectedBool: false,
		},
		{
			title:        "environmental variable `validate_hmac` is set with explicit `false`",
			value:        "false",
			expectedBool: false,
		},
	}
	hmacEnvVar := "validate_hmac"
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			os.Setenv(hmacEnvVar, test.value)
			value := HmacEnabled()
			if value != test.expectedBool {
				t.Errorf("Expected value: %v got: %v", test.expectedBool, value)
			}
		})
	}
}
