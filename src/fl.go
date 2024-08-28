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
		panic(err)
	}

	err = cmd.ParseCommandLine(os.Args[1:], &flags)
	if err != nil {
		panic(err)
	}

	if flags.Config {
		cmd.WriteConfig(filepath, &flags)
		os.Exit(0)
	}

	if flags.Login {
		os.Exit(0)
	}

	if flags.FLID == "" {
		flid, err := api.RegisterUserByIP()
		if err != nil {
			fmt.Printf("Failed to register you. You can try again or use 'fl login': %v\n", err)
			os.Exit(1)
		}

		flags.FLID = flid
		cmd.WriteConfig(filepath, &flags)
	}

	utils.Log(flags.Verbose, "Flags: %+v\n", flags)

	runFL(flags)
}

func runFL(flags cmd.FlagConfig) {
	res, err := api.GenerateCommand(flags.Prompt, flags.Langtool, flags.FLID)
	if err != nil {
		fmt.Printf("Failed to generate a command: %v\n", err)
		os.Exit(1)
	}

	utils.Clip(res.Cmd)
	fmt.Println(res.Cmd)
	fmt.Println()

	if flags.Outfile != "" {
		err = os.WriteFile(flags.Outfile, []byte(res.Cmd), 0755)
		if err != nil {
			fmt.Printf("Failed save output to file: %s\n", err)
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
			panic(err)
		}

		// Print the output
		fmt.Println(out)
	}
}
