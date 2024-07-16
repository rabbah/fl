package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"fl/helpers"
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

	Flags := helpers.ConstructFlags()

	// parse arguments and recieve prompt
	prompt, err := helpers.ArgParse(os.Args, &Flags)

	// exit if -h/--help flags found
	if Flags.Help == true {
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

	response, err := http.Post(apiUrl, "application/json", strings.NewReader(fmt.Sprintf(`{"prompt": "%s"}`, prompt)))
	if err != nil {
		fmt.Printf("Failed to call Flows API: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Parse the response body as JSON
	helpers.Print(Flags.Verbose, "Parsing response...")

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Failed to parse Flows API response: %s\n", err)
		os.Exit(1)
	}

	// Check if the "output" field exists
	helpers.Print(Flags.Verbose, "Checking AI output field...")
	
	result, ok := data["output"].(string)
	if !ok {
		fmt.Println("Error: Expected output field not found in Flows API response")
		os.Exit(1)
	}
	
	// Emit the result
	helpers.Print(Flags.Verbose, "Output: \n")

	fmt.Println(result)
	fmt.Println()

	// if not skipping prompt, ask user if they would like to execute
	userExecute := false
	if !Flags.Noexec && !Flags.Autoexecute {
		var userInput string
		fmt.Print("Would you like to execute the command? (y/n): ")
		fmt.Scanln(&userInput)
		userInput = strings.ToLower(userInput)
		if userInput == "y" || userInput == "yes" {
			userExecute = true
		}
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if Flags.Autoexecute || userExecute {
		// convert to arr of values (exec requires a specific format)
		fullCmd := strings.Split(result, " ")
		cmd := fullCmd[0]
		args := []string{}

		if len(fullCmd) > 1 {
			args=fullCmd[1:]
		}

		helpers.Print(Flags.Verbose, "Executing the result...")

		out, err := exec.Command(cmd, args...).Output()
	
		if err != nil {
			panic(err)
		}
	
		// Print the output
		fmt.Println(string(out))
	}
}
