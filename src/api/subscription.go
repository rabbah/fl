package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type apiSubscriptionInput struct {
	Input struct {
		FLID string `json:"flid"`
	} `json:"Input"`
}

type SubscriptionResult struct {
	Subscription    bool   `json:"subscription"`
	SubscriptionURL string `json:"subscriptionURL"`
	Error           string `json:"error"`
}

type apiSubscriptionOutput struct {
	Output SubscriptionResult `json:"Output"`
}

func StartSubscription(flid string) (*SubscriptionResult, error) {
	body := apiSubscriptionInput{}
	body.Input.FLID = flid

	statusCode, response, err := utils.PostJSON(StartSubscriptionAPI, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to start a new subscription: %s", string(response))
		return nil, err
	}

	res := apiSubscriptionOutput{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	if res.Output.Error != "" {
		err = fmt.Errorf("failed to get GitHub access token")
		return nil, err
	}

	return &res.Output, nil
}
