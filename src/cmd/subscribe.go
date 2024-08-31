package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
)

const (
	clientID = "Ov23liak5XRTpeHgGDtx"
	stripURL = "https://buy.stripe.com/test_00g9D094f5n1bAI288"
)

func Subscribe(flags *FlagConfig, filepath string) error {
	token, err := utils.GetGitHubAccessToken(clientID)
	if err != nil {
		return err
	}

	if token.AccessToken == "" {
		err = fmt.Errorf("failed to get GitHub access token")
		return err
	}

	if flags.Verbose {
		fmt.Println("GitHub access token:", token.AccessToken)
	}

	flags.FLID, err = api.LoginCommand(token.AccessToken)
	if err != nil {
		return err
	}

	err = WriteConfig(filepath, *flags)
	if err != nil {
		return err

	}
	url := stripURL + "?client_reference_id=" + flags.FLID
	err = utils.OpenURL(url)
	if err != nil {
		fmt.Println("Could not open the browser automatically, so please navigate to the following URL to subscribe:\n\t", url)
	}

	return nil
}
