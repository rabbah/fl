package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
)

func loginGitHub(verbose bool, githubClientId string) (string, error) {
	token, err := utils.GetGitHubAccessToken(githubClientId)
	if err != nil {
		return "", err
	}

	if token.AccessToken == "" {
		err = fmt.Errorf("failed to get GitHub access token")
		return "", err
	}

	if verbose {
		fmt.Println("GitHub access token:", token.AccessToken)
	}

	flid, err := api.LoginCommand(token.AccessToken)
	if err != nil {
		return "", err
	}

	return flid, nil
}

func loginGuest() (string, error) {
	flid, err := api.LoginGuestUserByIP()
	if err != nil {
		err = fmt.Errorf("failed to get a Guest access token: %v", err)
		return "", err
	}

	return flid, nil
}

func LoginMessage(guest bool) error {
	if guest {
		fmt.Println("Please login in first. Use the following command: fl subscription login --guest")
	} else {
		fmt.Println("Please login in first. Use the following command: fl subscription login")
	}
	return nil
}
