package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type LoginInput struct {
	Token string `json:"token"`
}
type LoginResult struct {
	FLID string `json:"flid"`
}

func LoginCommand(token string) (string, error) {
	body := LoginInput{}
	body.Token = token

	statusCode, response, err := utils.PostJSON(LoginGitHubAPI, body)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to login: %s", string(response))
		return "", err
	}

	res := LoginResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}

	flid := res.FLID

	if flid == "" {
		err = fmt.Errorf("failed to login: %s", string(response))
	}

	return flid, err
}
