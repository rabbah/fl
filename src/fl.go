package main

import (
	"fl/exec"
	"fl/helpers"
	"fl/io"
	"fmt"
	"os"
)

/**********************
 * globals
 *********************/

const (
	// url of the Flow endpoint
	apiUrl = "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"
)

/**********************
 * main
 *********************/

func main() {

	// initialize flags struct
	Flags := helpers.ConstructFlags()

	// parse arguments and recieve prompt
	prompt, err := helpers.ArgParse(os.Args, &Flags)

	// exit if -h/--help flags found
	if Flags.Help {
		// handler for when --help or -h are provided
		helpers.Usage()
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		helpers.Usage()
		os.Exit(1)
	}

	helpers.Print(Flags.Verbose, "Prompt extracted:", prompt)

	// Make the API call
	helpers.Print(Flags.Verbose, "Sending prompt...")
	res, err := helpers.PromptAI(apiUrl, prompt)
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// Parse the response body as JSON
	helpers.Print(Flags.Verbose, "Parsing response...")
	result, err := helpers.ParseResponse(res)
	if err != nil {
		fmt.Printf(result, err)
		os.Exit(1)
	}

	// Emit the result
	helpers.Print(Flags.Verbose, "Output: \n")
	fmt.Println(result)
	fmt.Println()

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
