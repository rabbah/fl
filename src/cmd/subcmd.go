package cmd

import (
	"fl/api"
	"os"

	"github.com/spf13/cobra"
)

func addSubscribeCommand(rootCmd *cobra.Command, filepath string, flags *FlagConfig) {
	subscribeCmd := &cobra.Command{
		Use:           "subscription",
		Aliases:       []string{"sub"},
		Short:         "Manage your subscription",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subLoginCmd := &cobra.Command{
		Use:           "login",
		Short:         "Create a login",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flid := ""
			err := error(nil)

			if guest, _ := cmd.Flags().GetBool("guest"); guest {
				flid, err = loginGuest()
			} else {
				flid, err = loginGitHub(flags.Verbose, api.GitHubClientID)
			}

			if err != nil {
				return err
			}

			// save the flid to the configuration file for future use
			flags.FLID = flid

			err = writeConfig(filepath, *flags)
			if err != nil {
				return err
			}

			if subscribe, _ := cmd.Flags().GetBool("subscribe"); subscribe {
				return startSubscription(flags)
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subStartCmd := &cobra.Command{
		Use:           "start",
		Short:         "Start subscription",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startSubscription(flags)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subCancelCmd := &cobra.Command{
		Use:           "cancel",
		Short:         "Cancel subscription",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cancelSubscription(flags)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subStatusCmd := &cobra.Command{
		Use:           "status",
		Short:         "Check status of subscription",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusSubscription(flags)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subLoginCmd.PersistentFlags().BoolP("guest", "g", false, "Guest login")
	subLoginCmd.PersistentFlags().BoolP("subscribe", "s", false, "Start subscription")
	subLoginCmd.MarkFlagsMutuallyExclusive("guest", "subscribe")

	subscribeCmd.AddCommand(subLoginCmd)
	subscribeCmd.AddCommand(subStartCmd)
	subscribeCmd.AddCommand(subCancelCmd)
	subscribeCmd.AddCommand(subStatusCmd)

	rootCmd.AddCommand(subscribeCmd)
}
