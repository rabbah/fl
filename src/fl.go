package main

import (
	"fl/exec"
	"fl/helpers"
	"fl/io"
	"fl/ui"
	"fl/web"
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

// REMOVE THIS LATER
const (
	// url of the Flow endpoint
	apiUrl = "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Failed to initialize clipboard.")
		panic(err)
	}
}

/**********************
 * fl in-line execution
 *********************/

func noTui(Flags helpers.FlagStruct) {

	helpers.Print(Flags.Verbose, "Prompt extracted:", Flags.Prompt)

	// Make the API call
	helpers.Print(Flags.Verbose, "Sending prompt...")
	res, err := web.PromptAI(apiUrl, Flags.Prompt)
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// Parse the response body as JSON
	helpers.Print(Flags.Verbose, "Parsing response...")
	result, err := web.ParseResponse(res)
	if err != nil {
		fmt.Printf(result, err)
		os.Exit(1)
	}

	// Emit the result
	helpers.Print(Flags.Verbose, "Output: \n")
	fmt.Println(result)
	fmt.Println()

	clipboard.Write(clipboard.FmtText, []byte(result))

	// if not skipping prompt, ask user if they would like to execute
	userExecute := false
	if !Flags.Noexec && !Flags.Autoexecute {
		userExecute = exec.PromptExec()
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if Flags.Autoexecute || userExecute {
		helpers.Print(Flags.Verbose, "Executing the result...")

		out, err := exec.Exec(result)

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

	// initialize flags struct
	Flags := helpers.ConstructFlags()

	// parse arguments and recieve prompt
	err := helpers.ArgParse(os.Args, &Flags)

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
		err = ui.RunProgram(&Flags)
		if err != nil {
			fmt.Printf("Error running TUI: %v", err)
			os.Exit(1)
		}
	} else {
		// execute in-line if TUI flag not set
		noTui(Flags)
	}
}
