package main

import (
	"fl/api" // Add this line to import the auth package
	"fl/cmd"
	"fl/exec"
	"fl/utils"
	"fmt"
	"os"
	"path/filepath"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func main() {
	home, _ := os.UserHomeDir()
	filepath := filepath.Join(home, ".flconf")
	flags := cmd.FlagConfig{}

	err := cmd.ReadConfig(filepath, &flags)
	if err != nil {
		fmt.Printf("Error reading saved configuration: %s\n", err)
		os.Exit(1)
	}

	err = cmd.ParseCommandLine(os.Args[1:], filepath, &flags)
	if err != nil {
		fmt.Printf("Error handling command line arguments: %s\n", err)
		os.Exit(1)
	}

	if flags.FLID == "" {
		cmd.LoginMessage(true)
		os.Exit(0)
	}

	utils.Log(flags.Verbose, "Flags: %+v\n", flags)

	if flags.Prompt == "" {
		intro()
	} else {
		runFL(flags)
	}
}

func runFL(flags cmd.FlagConfig) {
	res, err := api.GenerateCommand(flags.Prompt, flags.Langtool, flags.FLID)
	if err != nil {
		fmt.Printf("Error generating a command: %v\n", err)
		os.Exit(1)
	}

	// invalid token, no command
	if !res.Valid {
		fmt.Println("Your access code is invalid.")
		cmd.LoginMessage(true)
		return
	}

	fmt.Println(res.Cmd)

	if res.Quota {
		fmt.Println(`
Warning: You have exhausted your allowed quota.
Features will be limited and your access may get cut off entirely.
Use 'fl subscription login --subscribe' to subscribe and continue using the tool.`)
		return
	}

	// no quota -> no clipboard, prompt or auto-run
	utils.Clip(res.Cmd)

	if flags.Outfile != "" {
		err = os.WriteFile(flags.Outfile, []byte(res.Cmd), 0755)
		if err != nil {
			fmt.Printf("Error saving output to file: %s\n", err)
			os.Exit(1)
		}
	}

	//TODO
	/*
		if flags.Explain {
			explanation, err := web.ExplainCommand(result, flags.Langtool)
			if err != nil {
				fmt.Printf("Failed to call Flows API - %s: %v\n", explanation, err)
				os.Exit(1)
			}

			fmt.Println(explanation)
		}
	*/

	runIt := false
	if flags.PromptRun && !flags.AutoExecute {
		runIt = exec.PromptExec()
	}

	// perform the command if autoexecute enabled or user prompted to exec
	if flags.AutoExecute || runIt {
		utils.Log(flags.Verbose, "Executing the generated command...")

		Cmd := exec.Command(res.Cmd)
		out, err := Cmd.Exec()

		if err != nil {
			fmt.Printf("Error while executing command: %s\n", err)
			os.Exit(1)
		}

		// Print the output
		fmt.Println(out)
	}
}

func intro() {
	source := `
Here are some sample calls for using 'fl':

**Description**: Remove a directory and all its contents.

	fl remove a directory and all its contents

	Sample output:
	rm -r directory_name

**Description**: Search for files containing a specific keyword in the src directory.

    fl search for files containing keyword in src directory

	Sample output:
	grep -r "keyword" src

**Description**: Process a CSV file to extract a column and count unique occurences of a value.

    fl count the number of unique values that appear in the second column of a csv file, make sure the count is case insensitive, report the total count only

	Sample output:
	awk -F, '{print tolower($2)}' file.csv | sort -u | wc -l

**Description**: Call an authenticated API and pass in some JSON data.

    fl call an api that returns JSON and sends some data {"foo":"bar"} as json where the api uses basic auth and the secret is an environment variable called API_KEY

	Sample output:
	curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic $API_KEY' -d '{"foo":"bar"}' https://api.example.com/endpoint`

	result := markdown.Render(source, 80, 0)
	fmt.Println(string(result))
}
