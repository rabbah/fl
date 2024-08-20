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

type RequestRegister struct {
	Input map[string]string
}

type ResponseRegister struct {
	Output map[string]string
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
	req := RequestRegister{
		map[string]string{"ip": ip},
	}
	reqJSON, _ := json.Marshal(req)

	response, msg, err := reqFlows(registerUrl, reqJSON)
	if err != nil {
		return msg, err
	}
	defer response.Body.Close()

	result, msg, err = parseIpRegister(response)
	if err != nil {
		return msg, err
	}

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

func parseIpRegister(res *http.Response) (result string, msg string, err error) {
	var data ResponseRegister

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", "Failed to parse Flows API response", err
	}

	errStr := data.Output["error"]
	if errStr != "" {
		return "", "", errors.New(errStr)
	}

	return data.Output["jwt"], "", err
}
