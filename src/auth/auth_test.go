package auth

import (
	"testing"
)

func TestRegisterIp(t *testing.T) {
	ip := "1.1.1.1"
	result, err := VerifyIp(ip)
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjEuMS4xLjMiLCJmbGlkIjoiZGZjM2YzZjMtYjMyNy00NGZiLWJlYWMtZTY4ZTY0MmIyZThhIiwidmVyc2lvbiI6InYxIn0=.P5JWgTAJU91QPfLPGItu715fKOc1ImwWxJWL+FLb24g="

	if result != expected || err != nil {
		t.Fatalf(`VerifyIp(%s) = ("%s",%v). Expected ('%s',%v)`, ip, result, err, expected, nil)
	}
}

func TestRegisterIpFail(t *testing.T) {
	ip := "abcd"
	_, err := VerifyIp(ip)

	if err == nil {
		t.Fatalf(`VerifyIp(%s) = should raise an error`, ip)
	}
}
