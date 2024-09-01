package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type FlagConfig struct {
	Verbose                bool   // verbose output while running
	Explain                bool   // explainer command
	PromptRun, AutoExecute bool   // prompt to run generated commands or auto run
	Outfile                string // write generated command to file
	Langtool               string // generate command for specific shell or a tool
	Prompt                 string // command prompt

	// these are properties from config file
	AutoExecuteConf bool
	LangtoolConf    string
	FLID            string
}

func ParseCommandLine(args []string, filepath string, flags *FlagConfig) error {
	rootCmd := &cobra.Command{
		Use:   "fl <prompt>",
		Short: "A command-line tool for generating command line scripts using AI",
		Long:  "A command-line tool for generating command line scripts using AI.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flags.Prompt = strings.Join(args[0:], " ")
		},
	}

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show command usage")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.PersistentFlags().BoolVarP(&flags.PromptRun, "prompt", "p", false, "Prompt to run generated commands")
	rootCmd.PersistentFlags().BoolVarP(&flags.AutoExecute, "run", "r", flags.AutoExecuteConf, "Automatically execute generated commands")
	//TODO//rootCmd.PersistentFlags().BoolVarP(&flags.Explain, "explain", "e", false, "Explain the generated command")

	rootCmd.PersistentFlags().StringVarP(&flags.Outfile, "outfile", "o", "", "Write generated command to file")
	//TODO//rootCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", flags.LangtoolConf, "Generate command for specific shell or a tool")

	// subscribe commands
	addSubscribeCommand(rootCmd, filepath, flags)

	// config commands
	addConfCommand(rootCmd, filepath, flags)

	exitAfterHelp(rootCmd, 0)
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func exitAfterHelp(c *cobra.Command, exitCode int) {
	helpFunc := c.HelpFunc()
	c.SetHelpFunc(func(c *cobra.Command, s []string) {
		helpFunc(c, s)
		os.Exit(exitCode)
	})
}
