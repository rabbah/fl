package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func AddSubscribeCommand(rootCmd *cobra.Command, filepath string, flags *FlagConfig) {
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

	subStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start subscription",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Subscribe(flags, filepath)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subCancelCmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel subscription",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			//TODO//cancel subscription
			panic("Cancelling subscription... not implemented yet")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subRestoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore a previous subscription",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO//restore subscription
			panic("Restoring subscription... not implemented yet")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	subscribeCmd.AddCommand(subStartCmd)
	subscribeCmd.AddCommand(subCancelCmd)
	subscribeCmd.AddCommand(subRestoreCmd)

	rootCmd.AddCommand(subscribeCmd)
}
