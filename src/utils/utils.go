package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	myIpAPI = "https://api.ipify.org"
)

func GetJSON(urlStr string, queryParams map[string]string, bearer string) (int, string, error) {
	// Step 1: Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", fmt.Errorf("error parsing URL: %w", err)
	}

	// Step 2: Add query parameters to the URL
	if queryParams != nil {
		query := parsedURL.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		parsedURL.RawQuery = query.Encode()
	}

	// Step 3: Create a new HTTP GET request
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return 0, "", fmt.Errorf("error creating request: %w", err)
	}

	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	// Step 4: Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Step 5: Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("error reading response body: %w", err)
	}

	// Step 7: Return the status code, the response map, and no error
	return resp.StatusCode, string(body), nil
}

func PostJSON(url string, payload interface{}) (int, []byte, error) {
	// Step 1: Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Step 2: Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("error creating request: %w", err)
	}

	// Step 3: Set the Accept and Content-Type header to application/json
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Step 4: Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Step 5: Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Step 7: Return the status code, the response map, and no error
	return resp.StatusCode, body, nil
}

/**
 * Get external IP address
 * @return IP as string, error
 */
func GetExternalIP() (string, error) {
	resp, err := http.Get(myIpAPI)
	if err != nil {
		return "", fmt.Errorf("error determining IP address: %w", err)
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading IP address from response: %w", err)
	}

	return string(ip), nil
}

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// print passed prompt if verbose enabled
func Log(verbose bool, msg string, rest ...interface{}) {
	if verbose {
		fmt.Printf(msg, rest...)
	}
}
