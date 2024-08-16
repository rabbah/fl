package auth

import (
	"testing"
)

func TestRegisterIp(t *testing.T) {
	result, err := VerifyIp("1.1.1.1")
	expected := `{
    "Output": {
        "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjEuMS4xLjMiLCJmbGlkIjoiZGZjM2YzZjMtYjMyNy00NGZiLWJlYWMtZTY4ZTY0MmIyZThhIiwidmVyc2lvbiI6InYxIn0=.P5JWgTAJU91QPfLPGItu715fKOc1ImwWxJWL+FLb24g="
    }
}`

	if result != expected || err != nil {
		t.Fatalf(`parseResponse(...) = ("%s",%v). Expected ('%s',%v)`, result, err, expected, nil)
	}
}

func TestRegisterIpFail(t *testing.T) {
	result, err := VerifyIp("abcd")
	expected := `{
    "Output": {
        "error": "invalid registration"
    }
}`

	if result != expected || err != nil {
		t.Fatalf(`parseResponse(...) = ("%s",%v). Expected ('%s',%v)`, result, err, expected, nil)
	}
}
