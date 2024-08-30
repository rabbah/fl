package cmd

import (
	"fl/api"
	"fl/utils"
	"fmt"
	"os"
	"path/filepath"
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

func ParseCommandLine(args []string, flags *FlagConfig) error {
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
			token, err := utils.GetGitHubAccessToken(utils.ClientID)
			if err != nil {
				panic(err)
			}

			if token.AccessToken == "" {
				err = fmt.Errorf("failed to get GitHub access token")
				panic(err)
			}

			flags.FLID, err = api.LoginCommand(flags.FLID, token.AccessToken)
			if err != nil {
				panic(err)
			}

			filepath := configDefaultPath()
			flags := &FlagConfig{}

			err = ReadConfig(configDefaultPath(), flags)
			if err != nil {
				err = fmt.Errorf("cannot save login information: %v", err)
				panic(err)
			}

			WriteConfig(filepath, flags)
			os.Exit(0)
		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Tool configuration",
		Args:  cobra.NoArgs,
	}

	configGetSubCmd := &cobra.Command{
		Use:   "get",
		Short: "Get setting(s)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			flags := &FlagConfig{}

			err := ReadConfig(configDefaultPath(), flags)
			if err != nil {
				panic(err)
			}

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

			os.Exit(0)
		},
	}

	configSetSubCmd := &cobra.Command{
		Use:   "set",
		Short: "Set setting(s)",
		Run: func(cmd *cobra.Command, args []string) {
			filepath := configDefaultPath()
			flags := &FlagConfig{}

			err := ReadConfig(filepath, flags)
			if err != nil {
				panic(err)
			}

			flags.AutoExecuteConf, _ = cmd.Flags().GetBool("run")
			flags.LangtoolConf, _ = cmd.Flags().GetString("langtool")

			err = WriteConfig(filepath, flags)
			if err != nil {
				panic(err)
			}

			os.Exit(0)
		},
	}

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show command usage")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	rootCmd.PersistentFlags().BoolVarP(&flags.PromptRun, "prompt", "p", false, "Prompt to run generated commands")
	rootCmd.PersistentFlags().BoolVarP(&flags.AutoExecute, "run", "r", flags.AutoExecuteConf, "Configure fl to always run the generated command without prompting")
	//TODO
	//rootCmd.PersistentFlags().BoolVarP(&flags.Explain, "explain", "e", false, "Explain the generated command")

	rootCmd.PersistentFlags().StringVarP(&flags.Outfile, "outfile", "o", "", "Write generated command to file")
	//TODO
	//rootCmd.PersistentFlags().StringVarP(&flags.Langtool, "langtool", "l", flags.LangtoolConf, "Generate command for specific shell or a tool")

	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetSubCmd)
	configGetSubCmd.PersistentFlags().BoolP("run", "r", false, "Get prompt to run setting")
	configGetSubCmd.PersistentFlags().BoolP("langtool", "l", false, "Get shell or tool setting")
	configGetSubCmd.PersistentFlags().BoolP("flid", "f", false, "Get login info")

	configCmd.AddCommand(configSetSubCmd)
	configSetSubCmd.PersistentFlags().BoolP("run", "r", flags.AutoExecuteConf, "Set prompt to run setting")
	configSetSubCmd.PersistentFlags().StringP("langtool", "l", flags.LangtoolConf, "Set default shell or a tool or use")

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func configDefaultPath() string {
	home, _ := os.UserHomeDir()
	filepath := filepath.Join(home, ".flconf")
	return filepath
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

func WriteConfig(filepath string, flags *FlagConfig) error {
	viper.Set("run", flags.AutoExecuteConf)
	viper.Set("langtool", flags.LangtoolConf)
	viper.Set("flid", flags.FLID)

	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")
	return viper.WriteConfig()
}
