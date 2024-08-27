package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

type FlagConfig struct {
	Verbose, Run, Explain bool
	Outfile, Langtool     string
	Login                 bool
	Prompt                string
}

func ParseCommandLine(args []string) FlagConfig {
	flags := FlagConfig{}

	rootCmd := &cobra.Command{
		Use:   "fl <prompt>",
		Short: "A command-line tool for generating command line scripts using AI",
		Long:  "A command-line tool for generating command line scripts using AI.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flags.Prompt = strings.Join(args[0:], " ")
		},
	}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the fl service",
		Run: func(cmd *cobra.Command, args []string) {
			flags.Login = true
		},
	}

	rootCmd.AddCommand(loginCmd)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show command usage")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.PersistentFlags().BoolVarP(&flags.Run, "run", "r", false, "Prompt to run generated commands")
	rootCmd.PersistentFlags().BoolVarP(&flags.Explain, "explain", "e", false, "Include generated command explanation")

	rootCmd.PersistentFlags().StringVarP(&flags.Outfile, "outfile", "o", "", "Write generated command to file")
	rootCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", "", "Generate command for specific shell or a tool")

	rootCmd.SetArgs(args)
	rootCmd.Execute()
	return flags
}
