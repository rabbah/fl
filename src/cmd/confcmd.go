package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addConfCommand(rootCmd *cobra.Command, filepath string, flags *FlagConfig) {
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
				return writeConfig(filepath, *flags)
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
			return writeConfig(filepath, *flags)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}

	configCmd.PersistentFlags().Bool("reset", false, "Reset configuration")

	configCmd.AddCommand(configGetSubCmd)
	configGetSubCmd.PersistentFlags().BoolP("run", "r", false, "Get auto-execute setting")
	configGetSubCmd.PersistentFlags().BoolP("langtool", "l", false, "Get shell or tool setting")
	configGetSubCmd.PersistentFlags().BoolP("flid", "f", false, "Get login info")

	configCmd.AddCommand(configSetSubCmd)
	configSetSubCmd.PersistentFlags().BoolP("run", "r", flags.AutoExecuteConf, "Set auto-execute")
	configSetSubCmd.PersistentFlags().StringP("langtool", "l", flags.LangtoolConf, "Set default shell or a tool or use")

	rootCmd.AddCommand(configCmd)
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

func writeConfig(filepath string, flags FlagConfig) error {
	viper.Set("run", flags.AutoExecuteConf)
	viper.Set("langtool", flags.LangtoolConf)
	viper.Set("flid", flags.FLID)

	viper.SetConfigFile(filepath)
	viper.SetConfigType("json")
	return viper.WriteConfig()
}
