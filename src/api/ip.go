package api

import (
	"encoding/json"
	"fl/utils"
)

type apiInputRegisterIP struct {
	Input struct {
		IP string `json:"ip"`
	} `json:"Input"`
}

type apiInputRegisterGitHub struct {
	Input struct {
		Profile utils.GitHubProfile `json:"profile"`
	} `json:"Input"`
}

type apiOutput struct {
	Output interface{} `json:"Output"`
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

	input := apiInputRegisterIP{}
	input.Input.IP = ip

	_, response, err := utils.PostJSON(RegisterIP, input)
	if err != nil {
		return
	}

	output := apiOutput{}
	err = json.Unmarshal([]byte(response), &output)
	if err != nil {
		return
	}

	flid = output.Output.(string)
	return
}

/**
 * Register a user by their GitHub profile and return their FLID.
 * If the registration fails, an empty string or error is returned.
 */
func RegisterUserByGitHubProfile(profile utils.GitHubProfile) (flid string, err error) {
	input := apiInputRegisterGitHub{}
	input.Input.Profile = profile

	_, response, err := utils.PostJSON(RegisterGitHub, input)
	if err != nil {
		return
	}

	output := apiOutput{}
	err = json.Unmarshal([]byte(response), &output)
	if err != nil {
		return
	}

	flid = output.Output.(string)
	return
}
