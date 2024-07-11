package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

/**********************
 * globals/helpers
 *********************/

/*
 * To add a new flag:
 * 1. Update its entry in Usage()
 * 2. Add it to the list of flags below and update numFlags
 * 3. Create a handler function and add it to the switch-case in argParse()
 */
// useage definition functions to explain command and its args
var Usage = func () {
	fmt.Println("Usage of fl:")
	fmt.Println("\t-h,--help: show command usage")
}

// list of global flags
var (
	verbose = false
	autoexecute = false
)

// number of distinct flags
const (
	numFlags = 3	
) 

/**********************
 * arg parsing/ handlers
 *********************/

// handler for when --help or -h are provided
func flagHandleHelp() {
	Usage()
	os.Exit(0)
}

// verbose handler
func flagHandleVerbose() {
	verbose = true
}

// handler for when -y is provided
func flagHandleExecuteCmd() {
	autoexecute = true
}

// parse the user input for potential prompts
var argParse = func () (prompt string) {
	// Check if command line arguments are provided
	if len(os.Args) < 2 {
		// expecting at least 2 arguments
		Usage()
		os.Exit(1)
	}

	// Start of the user prompt (after args have been parsed)
	startPromptIndex := 1
	// flag to exit for loop if non-flag detected
	validArg := true

	// check for flags (add 1 bc first index is command path)
	for i := 1; i < numFlags+1 && validArg; i++ {
		switch os.Args[i] {
		case "-h": // help commands (just display useage)
			fallthrough
		case "--help":
			startPromptIndex++
			flagHandleHelp()
		case "-v": // handle program verbosity
			fallthrough
		case "--verbose":
			startPromptIndex++
			flagHandleVerbose()
		case "-y": // execute command automatically
			startPromptIndex++
			flagHandleExecuteCmd()
		default:
			// skip searching for switches if invalid arg is found (assume it is prompt)
			validArg = false
			break
		}
	}

	prompt = strings.Join(os.Args[startPromptIndex:], " ")
	if prompt == "" {
		fmt.Println("Prompt cannot be empty\n")
		Usage()
		os.Exit(1)
	}

	return (prompt)
}

/**********************
 * main
 *********************/

func main() {
	// parse arguments and recieve prompt
	prompt := argParse()
	// Concatenate all arguments to form the payload
	apiUrl := "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"

	// @verbose
	if verbose { fmt.Println("Prompt extracted:", prompt) }

	// Make the API call
	// @verbose
	if verbose { fmt.Println("Sending prompt...") }
	response, err := http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Parse the response body as JSON
	// @verbose
	if verbose { fmt.Println("Parsing response...") }
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to parse Flows API response: %s\n", err)
		os.Exit(1)
	}

	// Check if the "output" field exists
	// @verbose
	if verbose { fmt.Println("Checking AI output field...") }
	result, ok := data["output"].(string)
	if !ok {
		fmt.Println("Error: Expected output field not found in Flows API response")
		os.Exit(1)
	}
	
	// Emit the result
	// @verbose
	if verbose { fmt.Println("Output: \n") }
	fmt.Println(result)
	fmt.Println()

	// perform the command if autoexecute enabled
	if autoexecute {
		// convert to arr of values (exec requires a specific format)
		fullCmd := strings.Split(result, " ")
		cmd := fullCmd[0]
		args := []string{}

		if len(fullCmd) > 1 {
			args=fullCmd[1:]
		}

		// @verbose
		if verbose { fmt.Println("Executing the result...") }

		out, err := exec.Command(cmd, args...).Output()
	
		if err != nil {
			panic(err)
		}
	
		// Print the output
		fmt.Println(string(out))
	}
}
