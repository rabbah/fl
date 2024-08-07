package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	explainUrl           = "https://flow.pstmn-beta.io/api/30292fd2914e417a8b2d61e76b73edeb"
	commandGenerationUrl = "https://flow.pstmn-beta.io/api/38a029541f794a65afb284a7f4e7d3b3"
)

func ExplainCommand(command string, language string) (result string, err error) {
	var data map[string]interface{}

	// convert to JSON for proper escaping of strings which may be in the command
	req := map[string]string{
		"command":  command,
		"language": language,
	}
	reqJSON, _ := json.Marshal(req)
	response, err := http.Post(explainUrl, "application/json", bytes.NewBuffer(reqJSON))

	if err != nil {
		return "Failed to send response", err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "Failed to parse Flows API response: %s\n", err
	}

	result, ok := data["output"].(string)
	if !ok {
		return "", errors.New("expected output field not found in Flows API response")
	}

	return result, nil
}

func GenerateCommand(prompt string, language string) (result string, err error) {
	var data map[string]interface{}

	// convert to JSON for proper escaping of strings which may be in the command
	req := map[string]string{
		"prompt":   prompt,
		"language": language,
	}
	reqJSON, _ := json.Marshal(req)
	response, err := http.Post(commandGenerationUrl, "application/json", bytes.NewBuffer(reqJSON))

	if err != nil {
		return "Failed to send response", err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "Failed to parse Flows API response", err
	}

	result, ok := data["output"].(string)
	if !ok {
		return "", errors.New("expected output field not found in Flows API response")
	}

	return result, nil
}

func PromptAI(apiUrl string, prompt string, language string) (res *http.Response, err error) {
	// convert to JSON for proper escaping of strings which may be in the command
	req := map[string]string{
		"prompt":   prompt,
		"language": language,
	}
	reqJSON, _ := json.Marshal(req)
	res, err = http.Post(apiUrl, "application/json", bytes.NewBuffer(reqJSON))
	return res, err
}

func ParseResponse(res *http.Response) (result string, err error) {
	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return "Failed to parse Flows API response: %s\n", err
	}

	result, ok := data["output"].(string)

	if !ok {
		return "", errors.New("expected output field not found in Flows API response")
	}

	return result, nil
}
