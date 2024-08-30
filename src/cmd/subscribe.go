package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
	"os"
)

const (
	clientID = "Ov23liak5XRTpeHgGDtx"
)

func Subscribe(flags *FlagConfig, filepath string) {
	token, err := utils.GetGitHubAccessToken(clientID)
	if err != nil {
		panic(err)
	}

	if token.AccessToken == "" {
		err = fmt.Errorf("failed to get GitHub access token")
		panic(err)
	}

	flags.FLID, err = api.LoginCommand(flags.FLID, token.AccessToken)
	if err != nil {
		panic(err)
	}

	WriteConfig(filepath, *flags)
	os.Exit(0)
}
