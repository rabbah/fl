package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	registerUrl = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/register"
	verifyUrl   = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/verify"
	extIpUrl    = "https://api.ipify.org"
)

type Input interface{}
type Output interface {
	parse(res *http.Response)
}
type Request struct {
	input Input
}

func (req Request) send(apiRequest string) (res *http.Response, err error) {
	reqJSON, _ := json.Marshal(req.input)

	res, _, err = reqFlows(apiRequest, reqJSON)
	if err != nil {
		return res, err
	}

	return res, nil
}

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

type VerifyInput struct {
	Input string `json:"Input"`
}

type VerifyOutput struct {
	Output struct {
		Valid bool `json:"valid"`
		Quota int  `json:"quota"`
		Flid  struct {
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

func ExternalIP() (string, error) {
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

func RegisterIp(ip string) (result string, err error) {
	input := RegisterInput{
		Input: struct {
			Ip string `json:"ip"`
		}{Ip: ip},
	}
	output := RegisterOutput{}
	req := Request{input}

	res, err := req.send(registerUrl)
	if err != nil {
		return "", err
	}

	output, err = output.parse(res)
	if err != nil {
		return "", err
	}

	if output.Output.Error != "" {
		return "", errors.New(output.Output.Error)
	}

	result = output.Output.Jwt

	return result, nil
}

func VerifyJwt(jwt string) (result bool, flid string, version string, err error) {
	input := VerifyInput{
		jwt,
	}
	output := VerifyOutput{}
	req := Request{input}

	res, err := req.send(verifyUrl)
	if err != nil {
		return false, "", "", err
	}

	output, err = output.parse(res)
	if err != nil {
		return false, "", "", err
	}

	result = output.Output.Valid
	flid = output.Output.Flid.Flid
	version = output.Output.Flid.Version

	return result, flid, version, nil
}

func reqFlows(apiUrl string, reqJSON []byte) (res *http.Response, msg string, err error) {
	res, err = http.Post(apiUrl, "application/json", bytes.NewReader(reqJSON))

	if err != nil {
		return nil, "Failed to send response", err
	}

	return res, "", err
}
