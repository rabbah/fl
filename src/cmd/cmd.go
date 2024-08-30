package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FlagConfig struct {
	Verbose, Explain       bool
	PromptRun, AutoExecute bool
	Outfile, Langtool      string
	Prompt                 string

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

	subscribeCmd := &cobra.Command{
		Use:   "subscribe",
		Short: "Manage your fl subscription",
		Run: func(cmd *cobra.Command, args []string) {
			Subscribe(flags, filepath)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Tool configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	configGetSubCmd := &cobra.Command{
		Use:   "get",
		Short: "Get setting(s)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			run, _ := cmd.Flags().GetBool("run")
			langtool, _ := cmd.Flags().GetBool("langtool")
			flid, _ := cmd.Flags().GetBool("flid")
			all := !run && !langtool && !flid

			if all || flid {
				fmt.Println("flid:", flags.FLID)
			}

			if all || run {
				fmt.Println("run:", flags.AutoExecuteConf)
			}

			if all || langtool {
				fmt.Println("langtool:", flags.LangtoolConf)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	configSetSubCmd := &cobra.Command{
		Use:   "set",
		Short: "Set setting(s)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			flags.AutoExecuteConf, _ = cmd.Flags().GetBool("run")
			flags.LangtoolConf, _ = cmd.Flags().GetString("langtool")

			err := WriteConfig(filepath, *flags)
			if err != nil {
				panic(err)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show command usage")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.PersistentFlags().BoolVarP(&flags.PromptRun, "prompt", "p", false, "Prompt to run generated commands")
	rootCmd.PersistentFlags().BoolVarP(&flags.AutoExecute, "run", "r", flags.AutoExecuteConf, "Automatically execute generated commands")
	//TODO
	//rootCmd.PersistentFlags().BoolVarP(&flags.Explain, "explain", "e", false, "Explain the generated command")

	rootCmd.PersistentFlags().StringVarP(&flags.Outfile, "outfile", "o", "", "Write generated command to file")
	//TODO
	//rootCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", flags.LangtoolConf, "Generate command for specific shell or a tool")

	rootCmd.AddCommand(subscribeCmd)

	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetSubCmd)
	configGetSubCmd.PersistentFlags().BoolP("run", "r", false, "Get auto-execute setting")
	configGetSubCmd.PersistentFlags().BoolP("langtool", "l", false, "Get shell or tool setting")
	configGetSubCmd.PersistentFlags().BoolP("flid", "f", false, "Get login info")

	configCmd.AddCommand(configSetSubCmd)
	configSetSubCmd.PersistentFlags().BoolP("run", "r", flags.AutoExecuteConf, "Set auto-execute")
	configSetSubCmd.PersistentFlags().StringP("langtool", "l", flags.LangtoolConf, "Set default shell or a tool or use")

	applyExitOnHelp(rootCmd, 0)
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func ReadConfig(filepath string, flags *FlagConfig) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	flags.AutoExecuteConf = viper.GetBool("run")
	flags.LangtoolConf = viper.GetString("langtool")
	flags.FLID = viper.GetString("flid")

	return nil
}

func WriteConfig(filepath string, flags FlagConfig) error {
	viper.Set("run", flags.AutoExecuteConf)
	viper.Set("langtool", flags.LangtoolConf)
	viper.Set("flid", flags.FLID)

	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")
	return viper.WriteConfig()
}

func applyExitOnHelp(c *cobra.Command, exitCode int) {
	helpFunc := c.HelpFunc()
	c.SetHelpFunc(func(c *cobra.Command, s []string) {
		helpFunc(c, s)
		os.Exit(exitCode)
	})
}
