package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FlagConfig struct {
	Verbose, Explain       bool
	PromptRun, AutoExecute bool
	Outfile, Langtool      string
	Login                  bool
	Config                 bool
	Prompt                 string
	FLID                   string
}

func ParseCommandLine(args []string, flags *FlagConfig) {
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

	configCmd := &cobra.Command{
		Use:   "conf",
		Short: "Configure settings",
		Run: func(cmd *cobra.Command, args []string) {
			flags.Config = true
		},
	}

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show command usage")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.PersistentFlags().BoolVarP(&flags.PromptRun, "prompt", "p", false, "Prompt to run generated commands")
	configCmd.PersistentFlags().BoolVarP(&flags.AutoExecute, "run", "r", false, "Configure fl to always run the generated command without prompting")
	rootCmd.PersistentFlags().BoolVarP(&flags.Explain, "explain", "e", false, "Explain the generated command")

	rootCmd.PersistentFlags().StringVarP(&flags.Outfile, "outfile", "o", "", "Write generated command to file")
	rootCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", "", "Generate command for specific shell or a tool")

	configCmd.PersistentFlags().BoolVarP(&flags.AutoExecute, "run", "r", false, "Configure fl to always run the generated command without prompting")
	configCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", "", "Set default shell or a tool or use")

	rootCmd.SetArgs(args)
	rootCmd.Execute()
}

func ReadConfig(filepath string, flags *FlagConfig) error {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	flags.AutoExecute = viper.GetBool("run")
	flags.Langtool = viper.GetString("langtool")
	flags.FLID = viper.GetString("flid")
	return nil
}

func WriteConfig(filepath string, flags FlagConfig) error {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")

	viper.Set("run", flags.AutoExecute)
	viper.Set("langtool", flags.Langtool)
	viper.Set("flid", flags.FLID)

	return viper.WriteConfig()
}
