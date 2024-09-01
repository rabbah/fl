package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
	"time"
)

func login(flags *FlagConfig, filepath string, githubClientId string) (string, error) {
	token, err := utils.GetGitHubAccessToken(githubClientId)
	if err != nil {
		return "", err
	}

	if token.AccessToken == "" {
		err = fmt.Errorf("failed to get GitHub access token")
		return "", err
	}

	if flags.Verbose {
		fmt.Println("GitHub access token:", token.AccessToken)
	}

	flags.FLID, err = api.LoginCommand(token.AccessToken)
	if err != nil {
		return "", err
	}

	err = WriteConfig(filepath, *flags)
	if err != nil {
		return "", err
	}

	return flags.FLID, nil
}

func startSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return loginMessage()
	}

	status, err := api.StartSubscription(flags.FLID)
	if err != nil {
		return err
	}

	if status.Status == "guest" {
		url := status.SubscriptionURL + "?client_reference_id=" + flags.FLID
		err = utils.OpenURL(url)
		if err != nil {
			fmt.Println("Could not open the browser automatically, so please navigate to the following URL to subscribe:\n\t", url)
		}
	} else {
		printStatus(status)
	}

	return nil
}

func cancelSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return loginMessage()
	}

	status, err := api.CancelSubscription(flags.FLID)
	if err != nil {
		return err
	}

	printStatus(status)
	return nil
}

func statusSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return loginMessage()
	}

	status, err := api.StatusOfSubscription(flags.FLID)
	if err != nil {
		return err
	}

	printStatus(status)
	return nil
}

func loginMessage() error {
	fmt.Println("Please login in first. Use the following command: fl subscription login")
	return nil
}

func printStatus(status *api.SubscriptionResult) {
	switch status.Status {
	case "guest":
		fmt.Println("You do not have an active subscription.")

	case "paid":
		createdTime, _ := status.Created.Int64()
		created := time.Unix(createdTime, 0)

		fmt.Printf("Your subscription is active. It was created on %s.\n", created.Format("2006-01-02"))

	case "canceling":
		canceledAtTime, _ := status.Canceled_At.Int64()
		cancelAtTime, _ := status.Cancel_At.Int64()

		canceledAt := time.Unix(canceledAtTime, 0)
		cancelAt := time.Unix(cancelAtTime, 0)

		fmt.Printf("Your subscription was canceled on %s.\n", canceledAt.Format("2006-01-02"))
		fmt.Printf("Your may continue using the tool until the end of your billing cycle on %s.\n", cancelAt.Format("2006-01-02"))
	}
}
