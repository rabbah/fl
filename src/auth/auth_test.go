package auth

import (
	"testing"
)

func TestRegisterIp(t *testing.T) {
	ip := "1.1.1.1"
	result, err := RegisterIp(ip)
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjEuMS4xLjMiLCJmbGlkIjoiZGZjM2YzZjMtYjMyNy00NGZiLWJlYWMtZTY4ZTY0MmIyZThhIiwidmVyc2lvbiI6InYxIn0=.P5JWgTAJU91QPfLPGItu715fKOc1ImwWxJWL+FLb24g="

	if result != expected || err != nil {
		t.Fatalf(`RegisterIp(%s) = ("%s",%v). Expected ('%s',%v)`, ip, result, err, expected, nil)
	}
}

func TestRegisterIpFail(t *testing.T) {
	ip := "abcd"
	_, err := RegisterIp(ip)

	if err == nil {
		t.Fatalf(`RegisterIp(%s) = should raise an error`, ip)
	}
}

func TestVerifyJwt(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjEuMS4xLjMiLCJmbGlkIjoiZGZjM2YzZjMtYjMyNy00NGZiLWJlYWMtZTY4ZTY0MmIyZThhIiwidmVyc2lvbiI6InYxIn0=.P5JWgTAJU91QPfLPGItu715fKOc1ImwWxJWL+FLb24g="
	result, flid, version, err := VerifyJwt(jwt)
	expectedResult := true
	expectedFlid := "ada20f6a-b0a6-4aa9-ac13-d8a4ccc25167"
	expectedVersion := "v1"

	if result != expectedResult || flid != expectedFlid || version != expectedVersion || err != nil {
		t.Fatalf(`VerifyJwt(%s) = (%v, %s, %s, %v). Expected (%v, %s, %s, %v)`, jwt, result, flid, version, err, expectedResult, expectedFlid, expectedVersion, nil)
	}
}

func TestVerifyJwtFail(t *testing.T) {
	jwt := "false jwt"
	result, flid, version, err := VerifyJwt(jwt)
	expectedResult := false
	expectedFlid := ""
	expectedVersion := ""

	if result != expectedResult || flid != expectedFlid || version != expectedVersion || err != nil {
		t.Fatalf(`VerifyJwt(%s) = (%v, %s, %s, %v). Expected (%v, %s, %s, %v)`, jwt, result, flid, version, err, expectedResult, expectedFlid, expectedVersion, nil)
	}
}
