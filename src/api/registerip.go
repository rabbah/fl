package api

import (
	"encoding/json"
	"fl/utils"
)

type apiRegisterInput struct {
	Input struct {
		IP string `json:"ip"`
	} `json:"Input"`
}

type apiRegisterOutput struct {
	Output LoginResult `json:"Output"`
}

/**
 * Register a user by their IP and return their FLID.
 * If the registration fails, an empty string or error is returned.
 */
func RegisterUserByIP() (flid string, err error) {
	ip, err := utils.GetExternalIP()
	if ip == "" || err != nil {
		return
	}

	input := apiRegisterInput{}
	input.Input.IP = ip

	_, response, err := utils.PostJSON(RegisterIP, input)
	if err != nil {
		return
	}

	output := apiRegisterOutput{}
	err = json.Unmarshal(response, &output)
	if err != nil {
		return
	}

	flid = output.Output.FLID
	return
}
