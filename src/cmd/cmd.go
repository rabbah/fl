package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	configCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf"},
		Short:   "Tool configuration",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			reset, _ := cmd.Flags().GetBool("reset")
			if reset {
				flags.FLID = ""
				flags.AutoExecuteConf = false
				flags.LangtoolConf = ""
				return WriteConfig(filepath, *flags)
			}
			return cmd.Help()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	configGetSubCmd := &cobra.Command{
		Use:   "get",
		Short: "Get properties",
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
		Use:           "set",
		Short:         "Set properties",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.AutoExecuteConf, _ = cmd.Flags().GetBool("run")
			flags.LangtoolConf, _ = cmd.Flags().GetString("langtool")
			return WriteConfig(filepath, *flags)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
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
	rootCmd.AddCommand(subscribeCmd)
	subscribeCmd.AddCommand(subStartCmd)
	subscribeCmd.AddCommand(subCancelCmd)
	subscribeCmd.AddCommand(subRestoreCmd)

	// config commands
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().Bool("reset", false, "Reset configuration")

	configCmd.AddCommand(configGetSubCmd)
	configGetSubCmd.PersistentFlags().BoolP("run", "r", false, "Get auto-execute setting")
	configGetSubCmd.PersistentFlags().BoolP("langtool", "l", false, "Get shell or tool setting")
	configGetSubCmd.PersistentFlags().BoolP("flid", "f", false, "Get login info")

	configCmd.AddCommand(configSetSubCmd)
	configSetSubCmd.PersistentFlags().BoolP("run", "r", flags.AutoExecuteConf, "Set auto-execute")
	configSetSubCmd.PersistentFlags().StringP("langtool", "l", flags.LangtoolConf, "Set default shell or a tool or use")

	exitAfterHelp(rootCmd, 0)
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

func exitAfterHelp(c *cobra.Command, exitCode int) {
	helpFunc := c.HelpFunc()
	c.SetHelpFunc(func(c *cobra.Command, s []string) {
		helpFunc(c, s)
		os.Exit(exitCode)
	})
}
