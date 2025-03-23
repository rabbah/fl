package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type GenerateCommandInput struct {
	Prompt   string `json:"prompt"`
	Language string `json:"language"`
	FLID     string `json:"flid"`
}
type GeneratedCommandResult struct {
	Valid bool   `json:"valid"`
	Quota bool   `json:"quota"`
	Cmd   string `json:"cmd"`
}

func GenerateCommand(prompt string, language string, flid string) (*GeneratedCommandResult, error) {
	body := GenerateCommandInput{}
	body.Prompt = prompt
	body.Language = language
	body.FLID = flid

	statusCode, response, err := utils.PostJSON(GenerateCmdAPI, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to generate command: %s", string(response))
		return nil, err
	}

	res := GeneratedCommandResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
