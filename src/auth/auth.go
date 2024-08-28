package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fl/io"
	"fmt"
	"net/http"
)

const (
	registerUrl = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/register"
	verifyUrl   = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/verify"
	extIpUrl    = "https://api.ipify.org"
	stripeUrl   = "Not yet implemented"
)

func reqFlows(apiUrl string, reqJSON []byte) (res *http.Response, msg string, err error) {
	res, err = http.Post(apiUrl, "application/json", bytes.NewReader(reqJSON))

	if err != nil {
		return nil, "Failed to send response", err
	}

	return res, "", err
}

// helper structs
type Input interface{}
type Output interface {
	parse(res *http.Response)
}
type Request struct {
	input Input
}

// private Request.send, get http response
func (req Request) send(apiRequest string) (res *http.Response, err error) {
	reqJSON, _ := json.Marshal(req.input)

	res, _, err = reqFlows(apiRequest, reqJSON)
	if err != nil {
		return res, err
	}

	return res, nil
}

/**
 * Register IP request/response structures/functions
 */
type RegisterInput struct {
	Input struct {
		Ip string `json:"ip"`
	} `json:"Input"`
}

type RegisterOutput struct {
	Output struct {
		Error string `json:"error"`
		Jwt   string `json:"jwt"`
	} `json:"Output"`
}

func (RegisterOutput) parse(res *http.Response) (RegisterOutput, error) {
	var tmp RegisterOutput

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return tmp, err
	}

	err = json.Unmarshal(bodyBytes, &tmp)
	if err != nil {
		return tmp, err
	}

	data := tmp
	res.Body.Close()
	return data, nil
}

/**
 * Verify JWT request/response structures/functions
 */
type VerifyInput struct {
	Input struct {
		Jwt      string `json:"jwt"`
		Flid     string `json:"flid"`
		Prompt   string `json:"prompt"`
		Language string `json:"language"`
	} `json:"Input"`
}

type VerifyOutput struct {
	Output struct {
		Valid      bool   `json:"valid"`
		AboveQuota bool   `json:"abovequota"`
		Cmd        string `json:"cmd"`
		Flid       struct {
			Flid    string `json:"flid"`
			Version string `json:"version"`
		} `json:"flid"`
	} `json:"Output"`
}

func (VerifyOutput) parse(res *http.Response) (VerifyOutput, error) {
	var tmp VerifyOutput

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return tmp, err
	}

	err = json.Unmarshal(bodyBytes, &tmp)
	if err != nil {
		return tmp, err
	}

	data := tmp
	res.Body.Close()
	return data, nil
}

/**
 * Private helpers functions
 */
func getExternalIP() (string, error) {
	resp, err := http.Get(extIpUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func registerIp(ip string) (output RegisterOutput, err error) {
	input := RegisterInput{
		Input: struct {
			Ip string `json:"ip"`
		}{Ip: ip},
	}
	req := Request{input}

	res, err := req.send(registerUrl)
	if err != nil {
		return output, err
	}

	output, err = output.parse(res)
	if err != nil {
		return output, err
	}

	if output.Output.Error != "" {
		return output, errors.New(output.Output.Error)
	}

	return output, nil
}

func verifyAndGenCmd(prompt string, language string, jwt string) (output VerifyOutput, err error) {
	input := VerifyInput{
		Input: struct {
			Jwt      string `json:"jwt"`
			Flid     string `json:"flid"`
			Prompt   string `json:"prompt"`
			Language string `json:"language"`
		}{
			Jwt:      jwt,
			Flid:     "",
			Prompt:   prompt,
			Language: language,
		},
	}
	req := Request{input}

	res, err := req.send(verifyUrl)
	if err != nil {
		return output, err
	}

	output, err = output.parse(res)
	if err != nil {
		return output, err
	}

	return output, nil
}

func genCmdWithFlid(prompt string, language string, flid string) (output VerifyOutput, err error) {
	input := VerifyInput{
		Input: struct {
			Jwt      string `json:"jwt"`
			Flid     string `json:"flid"`
			Prompt   string `json:"prompt"`
			Language string `json:"language"`
		}{
			Jwt:      "",
			Flid:     flid,
			Prompt:   prompt,
			Language: language,
		},
	}
	req := Request{input}

	res, err := req.send(verifyUrl)
	if err != nil {
		return output, err
	}

	output, err = output.parse(res)
	if err != nil {
		return output, err
	}

	return output, nil
}

/**
 * flid not found in conf, do an ip lookup and register ip.
 * return verifyoutput
 */
func flidNotFound(prompt string, language string) (VerifyOutput VerifyOutput, msg string, err error) {
	// Grab this user's public IP
	ip, err := getExternalIP()
	if err != nil {
		return VerifyOutput, "Failed to retrieve ip", err
	}

	// Use ip to register/check registration
	RegisterOutput, err := registerIp(ip)
	if err != nil {
		return VerifyOutput, "Failed to register user", err
	}

	// User returned jwt to check validation
	jwt := RegisterOutput.Output.Jwt
	VerifyOutput, err = verifyAndGenCmd(prompt, language, jwt)
	if err != nil {
		return VerifyOutput, "Failed to verify user credentials", err
	}

	// Exit if invalid jwt given
	if !VerifyOutput.Output.Valid {
		return VerifyOutput, "Failed to verify user credentials", errors.New("failed to validate user")
	}

	return VerifyOutput, "", nil
}

/**
 * Check for FLID. if FLID not found, get public IP and register it into a flid. Also gen command.
 * Otherwise, use flid and call backend to gen command.
 * Assume success iff err = nil.
 */
func ValidateUserGetCmd(prompt string, language string, flid string) (cmd string, msg string, err error) {
	var VerifyOutput VerifyOutput

	// no flid found
	if flid == "" {
		// get new flid and also send the cmd information
		VerifyOutput, msg, err = flidNotFound(prompt, language)
		if err != nil {
			return "", msg, err
		}
		// Save FLID to config if not found
		//Config.FLID = VerifyOutput.Output.Flid.Flid
		//Config.SaveConf()
	} else {
		// use flid to generate command
		VerifyOutput, err = genCmdWithFlid(prompt, language, flid)
		if err != nil {
			return "", "", err
		}
	}

	// Logic to check the quota
	// by specification, do not cancel upon meeting quota
	if VerifyOutput.Output.AboveQuota {
		fmt.Printf("You have exceeded the quota. Please register for payments here: %s\n\n", stripeUrl)
	}

	cmd = VerifyOutput.Output.Cmd
	return cmd, "", nil
}

/**
 * Restore FLID by calling external webhook
 */
func RestoreFLID() (err error) {
	fmt.Println("RestoreFLID() Not Yet Implemented.")
	return err
}
