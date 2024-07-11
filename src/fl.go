package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

/**********************
 * globals/helpers
 *********************/

/*
 * To add a flag:
 * 1. Update its entry in Usage()
 * 2. Add it to the string list called 'flags' below and update numFlags
 * 3. Create a handler function and add it to the switch-case in argParse()
 */
// useage definition functions to explain command and its args
var Usage = func () {
	fmt.Println("Usage of fl:")
	fmt.Println("\t-h,--help: show command usage")
}

// list of allowed flags
var (
	flags = []string { 
		"-h", "--help" ,
		"-y",
	}
)

// number of distinct flags
const (
	numFlags = 2	
) 

/**********************
 * arg parsing/ handlers
 *********************/

// handler for when --help or -h are provided
func flagHandleHelp() {
	Usage()
	os.Exit(1)
}

// handler for when -y is provided
func flagHandleExecuteCmd() {

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

func main() {
	// parse arguments and recieve prompt
	prompt := argParse()
	// Concatenate all arguments to form the payload
	apiUrl := "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"

	// Make the API call
	response, err := http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Parse the response body as JSON
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to parse Flows API response: %s\n", err)
		os.Exit(1)
	}

	// Check if the "output" field exists
	result, ok := data["output"].(string)
	if !ok {
		fmt.Println("Error: Expected output field not found in Flows API response")
		os.Exit(1)
	}

	// Emit the result
	fmt.Println(result)
}
