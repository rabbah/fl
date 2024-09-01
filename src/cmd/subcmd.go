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
			_, err := login(flags, filepath, api.GitHubClientID)
			return err
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

	subscribeCmd.AddCommand(subLoginCmd)
	subscribeCmd.AddCommand(subStartCmd)
	subscribeCmd.AddCommand(subCancelCmd)
	subscribeCmd.AddCommand(subStatusCmd)

	rootCmd.AddCommand(subscribeCmd)
}
