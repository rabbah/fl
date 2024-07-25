package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

/**********************
 * globals
 *********************/

const (
	// url of the Flow endpoint
	commandGenerationUrl = "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"
)

/* This is a combination of PromptAI and ParseResponse.
 * The other two functions should be removed appropriately when possible
 */
func GenerateCommand(prompt string) (result string, err error) {
	var data map[string]interface{}

	response, err := http.Post(commandGenerationUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
	if err != nil {
		return "", err
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

func PromptAI(apiUrl string, prompt string) (res *http.Response, err error) {
	res, err = http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
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
