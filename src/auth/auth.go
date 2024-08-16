package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	verifyUrl = "https://5709c812-4c6b-4bab-af35-a1e5a1413eda.mock.pstmn.io/register"
)

type InputVerify struct {
	Input map[string]string
}

type OutputVerify struct {
	Output map[string]string
}

func VerifyIp(ip string) (result string, err error) {
	req := InputVerify{
		map[string]string{"ip": ip},
	}
	reqJSON, _ := json.Marshal(req)

	response, msg, err := reqFlows(verifyUrl, reqJSON)
	if err != nil {
		return msg, err
	}
	defer response.Body.Close()

	result, msg, err = parseResponse(response)
	if err != nil {
		return msg, err
	}

	return result, nil
}

func reqFlows(apiUrl string, reqJSON []byte) (res *http.Response, msg string, err error) {
	res, err = http.Post(apiUrl, "application/json", bytes.NewReader(reqJSON))

	if err != nil {
		return nil, "Failed to send response", err
	}

	return res, "", err
}

func parseResponse(res *http.Response) (result string, msg string, err error) {
	var data OutputVerify

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
