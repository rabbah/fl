package main

import (
	"fl/exec"
	"fl/helpers"
	"fl/io"
	"fl/web"
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Failed to initialize clipboard.")
		panic(err)
	}
}

/**********************
 * TUI-based logic
 *********************/

func startTui(Flags helpers.FlagStruct, Config io.Config) {

	/*
	 *@DISABLED until TUI available for production
	 */
	// err := ui.RunProgram(&Flags)
	// if err != nil {
	// 	fmt.Printf("Error running TUI: %v", err)
	// 	os.Exit(1)
	// }

	fmt.Println("TUI is disabled for now - see '-h' for CLI usage!")
}

func noTui(Flags helpers.FlagStruct, Config io.Config) {

	helpers.Print(Flags.Verbose, "Prompt extracted:", Flags.Prompt)

	// Make the API call
	helpers.Print(Flags.Verbose, "Sending prompt...")
	result, err := web.GenerateCommand(Flags.Prompt, Flags.Language)
	if err != nil {
		fmt.Printf("Failed to call Flows API - %s: %v\n", result, err)
		os.Exit(1)
	}

	// Emit the result
	helpers.Print(Flags.Verbose, "Output: \n")
	fmt.Println(result)
	fmt.Println()

	// copy to clipboard
	clipboard.Write(clipboard.FmtText, []byte(result))

	// check if explain flag, then look
	if Flags.Explain {
		helpers.Print(Flags.Verbose, "Sending command for explanation...")

		explanation, err := web.ExplainCommand(result, Flags.Language)
		if err != nil {
			fmt.Printf("Failed to call Flows API - %s: %v\n", explanation, err)
			os.Exit(1)
		}

		fmt.Println(explanation)
	}

	// if not skipping prompt, ask user if they would like to execute
	userExecute := false
	if Flags.PromptExec {
		userExecute = exec.PromptExec()
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if (Config.Autoexec && !Flags.PromptExec) || userExecute {
		helpers.Print(Flags.Verbose, "Executing the result...")

		Cmd := exec.Command(result)
		out, err := Cmd.Exec()

		if err != nil {
			panic(err)
		}

		// Print the output
		fmt.Println(out)
	}

	if Flags.Output {
		err := io.Output(Flags.Outfile, result)

		if err != nil {
			fmt.Printf("Failed save output to file: %s\n", err)
			os.Exit(1)
		}
	}
}

/**********************
 * main
 *********************/

func main() {

	// get config data
	Config, err := io.ReadConf()
	if err != nil {
		fmt.Printf("Config read error: %s\n", err)
		helpers.Usage()
		os.Exit(1)
	}

	// check if the entered command was a conf parse command
	wasConfCmd, err := helpers.ConfParse(os.Args, &Config)
	if err != nil {
		fmt.Println(err)
		helpers.Usage()
		os.Exit(1)
	}
	if wasConfCmd {
		// save conf and exit
		err = Config.SaveConf()
		if err != nil {
			fmt.Printf("Config write error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// initialize flags struct
	Flags := helpers.ConstructFlags(Config)
	// parse arguments and recieve prompt
	err = helpers.ArgParse(os.Args, &Flags)

	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		helpers.Usage()
		os.Exit(1)
	}

	// exit if -h/--help flags found
	if Flags.Help {
		// handler for when --help or -h are provided
		helpers.Usage()
		os.Exit(0)
	}

	// Otherwise check for TUI flag
	if Flags.Tui {
		startTui(Flags, Config)
	} else {
		// execute in-line if TUI flag not set
		noTui(Flags, Config)
	}
}
