package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	registerUrl = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/register"
	verifyUrl   = "https://add9d90f-2d32-483d-835f-3dd2cb814764.mock.pstmn.io/verify"
)

type RequestRegister struct {
	Input map[string]string
}

type ResponseRegister struct {
	Output map[string]string
}

type RequestVerify struct {
	Input string
}

type ResponseVerify struct {
	Output struct {
		Valid bool `json:"valid"`
		Flid  struct {
			Flid    string `json:"flid"`
			Version string `json:"version"`
		} `json:"flid"`
	} `json:"Output"`
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
	req := RequestVerify{
		jwt,
	}
	reqJSON, _ := json.Marshal(req)

	fmt.Println("JWT: ", string(reqJSON))

	response, _, err := reqFlows(verifyUrl, reqJSON)
	if err != nil {
		return false, "", "", err
	}
	defer response.Body.Close()

	var parsedResponse ResponseVerify
	parsedResponse, _, err = parseVerifyJwt(response)
	if err != nil {
		return false, "", "", err
	}

	result = parsedResponse.Output.Valid
	flid = parsedResponse.Output.Flid.Flid
	version = parsedResponse.Output.Flid.Version

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

func parseVerifyJwt(res *http.Response) (data ResponseVerify, msg string, err error) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return data, "Failed to parse Flows API response", err
	}

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return data, "Failed to parse Flows API response", err
	}

	return data, "", err
}
