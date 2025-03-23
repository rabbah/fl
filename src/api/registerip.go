package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type RegisterInput struct {
	IP string `json:"ip"`
}

/**
 * Register a user by their IP and return their FLID.
 * If the registration fails, an empty string or error is returned.
 */
func LoginGuestUserByIP() (flid string, err error) {
	ip, err := utils.GetExternalIP()
	if ip == "" || err != nil {
		return
	}

	input := RegisterInput{}
	input.IP = ip

	_, response, err := utils.PostJSON(LoginGuestAPI, input)
	if err != nil {
		return
	}

	res := LoginResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	flid = res.FLID

	if flid == "" {
		err = fmt.Errorf("failed to register as a guest: %s", string(response))
	}

	return
}
