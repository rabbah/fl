package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Check if command line arguments are provided
	if len(os.Args) < 2 {
		// expecting at least 2 arguments
		fmt.Println("Usage: fl <prompt>")
		os.Exit(1)
	}

	// Concatenate all arguments to form the payload
	apiUrl := "https://flow.pstmn-beta.io/api/4e5b4cfcdec54831a31d9f38aaf1a938"
	prompt := strings.Join(os.Args[1:], " ")

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
