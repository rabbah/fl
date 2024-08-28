package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type apiGenerateCommandInput struct {
	Input struct {
		Prompt   string `json:"prompt"`
		Language string `json:"language"`
		FLID     string `json:"flid"`
	} `json:"Input"`
}

type apiGenerateCommandOutput struct {
	Output GeneratedCommandResult `json:"Output"`
}

type GeneratedCommandResult struct {
	Valid bool   `json:"valid"`
	Quota int    `json:"quota"`
	Cmd   string `json:"cmd"`
}

func GenerateCommand(prompt string, language string, flid string) (*GeneratedCommandResult, error) {
	body := apiGenerateCommandInput{}
	body.Input.Prompt = prompt
	body.Input.Language = language
	body.Input.FLID = flid

	statusCode, response, err := utils.PostJSON(GenerateCmd, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to generate command: %s", string(response))
		return nil, err
	}

	res := apiGenerateCommandOutput{}
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return nil, err
	}

	return &res.Output, nil
}
