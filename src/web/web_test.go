package web

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseWebResponse(t *testing.T) {
	expected := "Test response body from server"
	expectedJSON := fmt.Sprintf(`{"output":"%s"}`, expected)
	res := &http.Response{
		Body: io.NopCloser(strings.NewReader(expectedJSON)),
	}

	result, _, err := parseResponse(res)
	if result != expected || err != nil {
		t.Fatalf(`parseResponse(...) = ("%s",%v). Expected ('%s',%v)`, result, err, expected, nil)
	}
}
