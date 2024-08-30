package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type apiLoginInput struct {
	Input struct {
		JWT   string `json:"jwt"`
		Token string `json:"token"`
	} `json:"Input"`
}

type apiLoginOutput struct {
	Output LoginResult `json:"Output"`
}

type LoginResult struct {
	FLID string `json:"flid"`
}

func LoginCommand(jwt string, token string) (string, error) {
	body := apiLoginInput{}
	body.Input.JWT = jwt
	body.Input.Token = token

	statusCode, response, err := utils.PostJSON(LoginIP, body)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to login: %s", string(response))
		return "", err
	}

	res := apiLoginOutput{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}

	return res.Output.FLID, nil
}
