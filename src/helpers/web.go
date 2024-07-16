package helpers

import (
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"errors"
)

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
		return "", errors.New("Error: Expected output field not found in Flows API response")
	}

	return result, nil
}