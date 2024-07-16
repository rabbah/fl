package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"errors"
)

/**********************
 * globals/helpers
 *********************/

// global flag structure
type FlagStruct struct {
	verbose, help, autoexecute, noexec bool
	len int
}

const (
	// url of the Flow endpoint
	apiUrl = "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"
) 

func new_flags()(Flags FlagStruct) {
	return FlagStruct {
		verbose: false,
		help: false,
		autoexecute: false,
		noexec: false,
		len: 4,
	}
}

/*
 * To add a new flag:
 * 1. Update its entry in Usage()
 * 2. Add it to the list of flags above and update numFlags
 * 3. Create a handler function and add it to the switch-case in argParse()
 */
// useage definition functions to explain command and its args
var Usage = func () {
	fmt.Println("Usage: fl [-hynv] prompt...")
	// for formatting - please start with a space and ensure descruption alignment with tabs
	fmt.Println(" -h,--help\t\tshow command usage")
	fmt.Println(" -y\t\t\tautoexecute the generated command")
	fmt.Println(" -n\t\t\tdo not prompt for or run generated command (takes priority over -y)")
	fmt.Println(" -v,--verbose\t\tdisplay updates of the command progress")
}

// print passed prompt if verbose check set
func print(verbose bool, str ...interface{}) {
	if verbose {
		fmt.Println(str...)
	}
}

/**********************
 * arg parsing
 *********************/

// parse the user input for potential prompts
var argParse = func (args []string, Flags *FlagStruct) (prompt string, err error) {
	// Check if command line arguments are provided
	if len(args) < 2 {
		// expecting at least 2 arguments
		return "", errors.New("expecting at least 2 args")
	}

	// Start of the user prompt (after args have been parsed)
	startPromptIndex := 1
	// flag to exit for loop if non-flag detected
	validArg := true

	/* check for flags (add 1 bc first index is command path)
	 * validarg: if non-flag found, check for prompt
	 * also exit if args runs out
	 */
	for i := 1; i < Flags.len+1 && i < len(args) && validArg; i++ {
		switch args[i] {
		case "-h": // help commands (just display useage)
			fallthrough
		case "--help":
			startPromptIndex++
			Flags.help = true
			return "", nil // exit early if help flag
		case "-v": // handle program verbosity
			fallthrough
		case "--verbose":
			startPromptIndex++
			Flags.verbose = true
		case "-y": // execute command automatically
			startPromptIndex++
			Flags.autoexecute = true
		case "-n":
			startPromptIndex++
			Flags.noexec = true
		default:
			// skip searching for switches if invalid arg is found (assume it is prompt)
			validArg = false
			break
		}
	}

	// noexec takes priority over autoexecute, turn off autoexec
	// guarentees mutual exclusivity
	if Flags.noexec && Flags.autoexecute {
		Flags.autoexecute = false
	}

	prompt = strings.Join(args[startPromptIndex:], " ")
	if prompt == "" {
		return "", errors.New("Prompt cannot be empty") 
	}

	return prompt, nil
}

/**********************
 * main
 *********************/

func main() {

	Flags := new_flags()

	// parse arguments and recieve prompt
	prompt, err := argParse(os.Args, &Flags)

	// exit if -h/--help flags found
	if Flags.help == true {
		// handler for when --help or -h are provided
		Usage()
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		Usage()
		os.Exit(1)
	}

	print(Flags.verbose, "Prompt extracted:", prompt)

	// Make the API call
	print(Flags.verbose, "Sending prompt...")

	response, err := http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Parse the response body as JSON
	print(Flags.verbose, "Parsing response...")

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to parse Flows API response: %s\n", err)
		os.Exit(1)
	}

	// Check if the "output" field exists
	print(Flags.verbose, "Checking AI output field...")
	
	result, ok := data["output"].(string)
	if !ok {
		fmt.Println("Error: Expected output field not found in Flows API response")
		os.Exit(1)
	}
	
	// Emit the result
	print(Flags.verbose, "Output: \n")

	fmt.Println(result)
	fmt.Println()

	// if not skipping prompt, ask user if they would like to execute
	userExecute := false
	if !Flags.noexec && !Flags.autoexecute {
		var userInput string
		fmt.Print("Would you like to execute the command? (y/n): ")
		fmt.Scanln(&userInput)
		userInput = strings.ToLower(userInput)
		if userInput == "y" || userInput == "yes" {
			userExecute = true
		}
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if Flags.autoexecute || userExecute {
		// convert to arr of values (exec requires a specific format)
		fullCmd := strings.Split(result, " ")
		cmd := fullCmd[0]
		args := []string{}

		if len(fullCmd) > 1 {
			args=fullCmd[1:]
		}

		print(Flags.verbose, "Executing the result...")

		out, err := exec.Command(cmd, args...).Output()
	
		if err != nil {
			panic(err)
		}
	
		// Print the output
		fmt.Println(string(out))
	}
}
