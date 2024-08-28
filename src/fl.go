package main

import (
	"fl/api"
	"fl/auth"
	"fl/cmd"
	"fl/exec"
	"fl/helpers"
	"fl/io"
	"fl/utils"
	"fl/web"
	"fmt"
	"os"
	"path/filepath"
)

func runFL(flags cmd.FlagConfig) {
	helpers.Print(flags.Verbose, "Prompt extracted:", flags.Prompt)

	// Validate executing user
	helpers.Print(flags.Verbose, "Authorizing and sending prompt...")
	result, msg, err := auth.ValidateUserGetCmd(flags.Prompt, flags.Langtool, flags.FLID)
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}

	// Emit the result
	helpers.Print(flags.Verbose, "Output: \n")
	fmt.Println(result)
	fmt.Println()

	// copy to clipboard
	utils.Clip(result)

	// check if explain flag, then call explain
	if flags.Explain {
		helpers.Print(flags.Verbose, "Sending command for explanation...")

		explanation, err := web.ExplainCommand(result, flags.Langtool)
		if err != nil {
			fmt.Printf("Failed to call Flows API - %s: %v\n", explanation, err)
			os.Exit(1)
		}

		fmt.Println(explanation)
	}

	// if not skipping prompt, ask user if they would like to execute
	userExecute := false
	if flags.PromptRun {
		userExecute = exec.PromptExec()
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if (flags.AutoExecute && !flags.PromptRun) || userExecute {
		helpers.Print(flags.Verbose, "Executing the result...")

		Cmd := exec.Command(result)
		out, err := Cmd.Exec()

		if err != nil {
			panic(err)
		}

		// Print the output
		fmt.Println(out)
	}

	if flags.Outfile != "" {
		err := io.Output(flags.Outfile, result)

		if err != nil {
			fmt.Printf("Failed save output to file: %s\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	home, _ := os.UserHomeDir()
	filepath := filepath.Join(home, ".flconf")
	flags := cmd.FlagConfig{}

	err := cmd.ReadConfig(filepath, &flags)
	if err != nil {
		panic(err)
	}

	cmd.ParseCommandLine(os.Args[1:], &flags)

	if flags.Config {
		cmd.WriteConfig(filepath, flags)
		os.Exit(0)
	}

	if flags.Login {
		os.Exit(0)
	}

	if flags.FLID == "" {
		flid, err := api.RegisterUserByIP()
		if err != nil {
			fmt.Printf("Failed to register you. You can try again or use 'fl login': %v\n", err)
			os.Exit(1)
		}

		flags.FLID = flid
		newFlags := cmd.FlagConfig{}
		err = cmd.ReadConfig(filepath, &newFlags)
		if err != nil {
			panic(err)
		}
		cmd.WriteConfig(filepath, newFlags)
	}

	runFL(flags)
}
