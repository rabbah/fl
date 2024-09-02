package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
	"time"
)

func startSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return LoginMessage(false)
	}

	status, err := api.StartSubscription(flags.FLID)
	if err != nil {
		return err
	}

	if status.Status == "guest" {
		url := status.SubscriptionURL + "?client_reference_id=" + flags.FLID
		fmt.Println(`Continue in your browser. If the link does not open automatically, please navigate to the following URL to subscribe:`)
		fmt.Println(url)
		utils.OpenURL(url)
	} else {
		printStatus(status)
	}

	return nil
}

func cancelSubscription(flags *FlagConfig) error {
	if flags.FLID == "" {
		return LoginMessage(false)
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
		return LoginMessage(false)
	}

	status, err := api.StatusOfSubscription(flags.FLID)
	if err != nil {
		return err
	}

	printStatus(status)
	return nil
}

func printStatus(status *api.SubscriptionResult) {
	switch status.Status {
	case "guest", "default":
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
		fmt.Println("Please wait until the end of your billing cycle to start a new subscription.")
	}
}
