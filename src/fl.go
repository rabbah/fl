package main

import (
	"fl/api" // Add this line to import the auth package
	"fl/cmd"
	"fl/exec"
	"fl/utils"
	"fmt"
	"os"
	"path/filepath"
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
		flid, err := api.RegisterUserByIP()
		if err != nil {
			fmt.Printf("Failed to register you. You can try again or use 'fl login': %v\n", err)
			os.Exit(1)
		}

		flags.FLID = flid
		cmd.WriteConfig(filepath, flags)
	}

	utils.Log(flags.Verbose, "Flags: %+v\n", flags)

	runFL(flags)
}

func runFL(flags cmd.FlagConfig) {
	res, err := api.GenerateCommand(flags.Prompt, flags.Langtool, flags.FLID)
	if err != nil {
		fmt.Printf("Error generating a command: %v\n", err)
		os.Exit(1)
	}

	// invalid token, no command
	if !res.Valid {
		fmt.Println(`
Warning: Your access code is invalid.
Use 'fl subscription start' to start a new subscription or
'fl subscription restore' to restore an existing subscription.
To get a new guest access code, use 'fl config --reset'.`)
		return
	}

	fmt.Println(res.Cmd)

	if res.Quota {
		fmt.Println(`
Warning: You have exhausted your allowed quota.
Features will be limited and your access may get cut off entirely.
Use 'fl subscription start' to subscribe and continue using the tool.`)
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
