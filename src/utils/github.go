package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type GitHubAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GitHubProfile struct {
	Id    json.Number `json:"id"`
	Login string      `json:"login"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
}

func GetGitHubAccessToken(clientID string) (token GitHubAccessToken, err error) {
	githubDeviceLogin := fmt.Sprintf("https://github.com/login/device/code?client_id=%s&scope=user", clientID)

	_, body, err := PostJSON(githubDeviceLogin, nil)
	if err != nil {
		return
	}

	// Parse the response body to extract device_code, user_code, verification_uri, and interval
	var response struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		Interval        int    `json:"interval"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	deviceCode := response.DeviceCode
	interval := response.Interval

	Clip(response.UserCode)
	fmt.Println("Press 'return' key to open your browser automatically and login to Github.")
	fmt.Println("You will paste the following GitHub device code to login (already copied to your clipboard):", response.UserCode)
	// Wait until any key is pressed to continue so the user has time to read the message
	fmt.Scanln()
	err = OpenURL(response.VerificationURI + "?code=" + response.UserCode)
	if err != nil {
		fmt.Println("Could not open the browser automatically, so please navigate to the following URL to login to GitHub:\n\t", response.VerificationURI)
	}

	// Poll until the device and user codes expire or the user has successfully authorized the app with a valid user code
	fmt.Print("Polling for access token..")
	for i := 0; i < 12; i++ {
		fmt.Print(".")
		// Make a POST request to the GitHub API to exchange the device code for an access token
		accessTokenURL := "https://github.com/login/oauth/access_token?client_id=" + clientID + "&device_code=" + deviceCode + "&grant_type=urn:ietf:params:oauth:grant-type:device_code"
		_, body, err = PostJSON(accessTokenURL, nil)
		if err != nil {
			fmt.Println()
			return
		}

		// Parse the response body to extract the access token
		err = json.Unmarshal(body, &token)
		if err != nil {
			fmt.Println()
			return
		}

		// Check if the access token is empty
		if token.AccessToken != "" {
			fmt.Println()
			return
		}

		// Polling interval
		time.Sleep(time.Duration(interval) * time.Second)
	}

	if token.AccessToken == "" {
		fmt.Println(" timed out!")
		err = fmt.Errorf("timed out waiting for GitHub access token")
	}

	return
}

func ExchangeTokenForGitHubUserProfile(token string) (profile string, err error) {
	_, profile, err = GetJSON("https://api.github.com/user", nil, token)
	return
}

func GetGitHubUserProfile(clientID string) (profile GitHubProfile, err error) {
	token, err := GetGitHubAccessToken(clientID)
	if err != nil {
		return
	}

	if token.AccessToken == "" {
		err = fmt.Errorf("failed to get GitHub access token")
		return
	}

	response, err := ExchangeTokenForGitHubUserProfile(token.AccessToken)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(response), &profile)
	return
}
