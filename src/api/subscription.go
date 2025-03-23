package api

import (
	"encoding/json"
	"fl/utils"
	"fmt"
)

type SubscriptionInput struct {
	FLID string `json:"flid"`
}

type SubscriptionResult struct {
	Status          string      `json:"status"` // one of 'guest', 'paid' or 'cancelling'
	Created         json.Number `json:"created"`
	Canceled_At     json.Number `json:"canceled_at"`
	Cancel_At       json.Number `json:"cancel_at"`
	SubscriptionURL string      `json:"subscriptionURL"`
	Error           string      `json:"error"`
}

func StartSubscription(flid string) (*SubscriptionResult, error) {
	body := SubscriptionInput{}
	body.FLID = flid

	statusCode, response, err := utils.PostJSON(StartSubscriptionAPI, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to start a new subscription: %s", string(response))
		return nil, err
	}

	res := SubscriptionResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		err = fmt.Errorf(res.Error)
		return nil, err
	}

	return &res, nil
}

func CancelSubscription(flid string) (*SubscriptionResult, error) {
	body := SubscriptionInput{}
	body.FLID = flid

	statusCode, response, err := utils.PostJSON(CancelSubscriptionAPI, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to cancel subscription: %s", string(response))
		return nil, err
	}

	res := SubscriptionResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		err = fmt.Errorf(res.Error)
		return nil, err
	}

	return &res, nil
}

func StatusOfSubscription(flid string) (*SubscriptionResult, error) {
	body := SubscriptionInput{}
	body.FLID = flid

	statusCode, response, err := utils.PostJSON(StatusOfSubscriptionAPI, body)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		err = fmt.Errorf("failed to check status of subscription: %s", string(response))
		return nil, err
	}

	res := SubscriptionResult{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		err = fmt.Errorf(res.Error)
		return nil, err
	}

	return &res, nil
}
