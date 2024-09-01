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
		return fmt.Errorf("please login first")
	}

	status, err := api.StartSubscription(flags.FLID)
	if err != nil {
		return err
	}

	if status.Subscription {
		fmt.Println("You are already subscribed to the service.")
		return nil
	}

	url := status.SubscriptionURL + "?client_reference_id=" + flags.FLID
	err = utils.OpenURL(url)
	if err != nil {
		fmt.Println("Could not open the browser automatically, so please navigate to the following URL to subscribe:\n\t", url)
	}

	return nil
}

func cancelSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return fmt.Errorf("please login first")
	}

	status, err := api.CancelSubscription(flags.FLID)
	if err != nil {
		return err
	}

	canceledAtTime, _ := status.Canceled_At.Int64()
	canceledAt := time.Unix(canceledAtTime, 0)

	cancelAtTime, _ := status.Cancel_At.Int64()
	cancelAt := time.Unix(cancelAtTime, 0)

	fmt.Println("Your subscription was canceled on", canceledAt.Format("2006-01-02"))
	fmt.Println("Your may continue using the tool until the end of your billing cycle on", cancelAt.Format("2006-01-02"))

	return nil
}

func statusSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return fmt.Errorf("please login first")
	}

	status, err := api.StatusOfSubscription(flags.FLID)
	if err != nil {
		return err
	}

	//TODO/sub status
	canceledAtTime, _ := status.Canceled_At.Int64()
	canceledAt := time.Unix(canceledAtTime, 0)

	cancelAtTime, _ := status.Cancel_At.Int64()
	cancelAt := time.Unix(cancelAtTime, 0)

	fmt.Println("Your subscription was canceled on", canceledAt.Format("2006-01-02"))
	fmt.Println("Your may continue using the tool until the end of your billing cycle on", cancelAt.Format("2006-01-02"))

	return nil
}
