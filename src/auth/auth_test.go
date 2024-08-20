package auth

import (
	"testing"
)

func TestRegisterIp(t *testing.T) {
	ip := "1.1.1.1"
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjEuMS4xLjMiLCJmbGlkIjoiZGZjM2YzZjMtYjMyNy00NGZiLWJlYWMtZTY4ZTY0MmIyZThhIiwidmVyc2lvbiI6InYxIn0=.P5JWgTAJU91QPfLPGItu715fKOc1ImwWxJWL+FLb24g="

	result, err := RegisterIp(ip)
	jwt := result.Output.Jwt
	if jwt != expected || err != nil {
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
	expectedResult := true
	expectedFlid := "ada20f6a-b0a6-4aa9-ac13-d8a4ccc25167"
	expectedVersion := "v1"

	result, err := VerifyJwt(jwt)
	valid := result.Output.Valid
	flid := result.Output.Flid.Flid
	version := result.Output.Flid.Version
	if valid != expectedResult || flid != expectedFlid || version != expectedVersion || err != nil {
		t.Fatalf(`VerifyJwt(%s) = (%v, %s, %s, %v). Expected (%v, %s, %s, %v)`, jwt, result, flid, version, err, expectedResult, expectedFlid, expectedVersion, nil)
	}
}

func TestVerifyJwtFail(t *testing.T) {
	jwt := "false jwt"
	expectedResult := false
	expectedFlid := ""
	expectedVersion := ""

	result, err := VerifyJwt(jwt)
	valid := result.Output.Valid
	flid := result.Output.Flid.Flid
	version := result.Output.Flid.Version
	if valid != expectedResult || flid != expectedFlid || version != expectedVersion || err != nil {
		t.Fatalf(`VerifyJwt(%s) = (%v, %s, %s, %v). Expected (%v, %s, %s, %v)`, jwt, result, flid, version, err, expectedResult, expectedFlid, expectedVersion, nil)
	}
}

func TestGetIp(t *testing.T) {
	t.Skip("This requires knowing your ext ip. Only execute test if ExternalIP is suspected to fail.")
	expected := ""

	ip, err := GetExternalIP()

	if ip != expected || err != nil {
		t.Fatalf(`ExternalIP) = %s, %v. Expected %s, %v`, ip, err, expected, nil)
	}
}
