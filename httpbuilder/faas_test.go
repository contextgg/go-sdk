package httpbuilder

import (
	"net/http"
	"testing"
)

func TestInvokeFunction(t *testing.T) {
	result := ""

	builder := NewFaaS().
		SetFunction("echo").
		SetMethod(http.MethodPost).
		SetBody("hello?").
		SetOut(&result).
		SetLogger(t.Logf)

	status, err := builder.Do()
	if err != nil {
		t.Error(err)
	}

	if status != http.StatusOK {
		t.Errorf("Wrong status: %d", status)
	}

	if result != "hello?" {
		t.Errorf("Result: %s", result)
	}
}
