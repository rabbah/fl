package helpers

import (
	"strings"
	"errors"
	"fmt"
)

/*
 * To add a new flag:
 * 1. Add it to FlagStruct and ConstructFlags (increment the constructor's Len as well)
 * 2. Create and add the flag handler to the switch-case in argParse()
 * 3. Update the flag's entry in Usage()
 */

// global flag structure
type FlagStruct struct {
	Verbose, Help, Autoexecute, Noexec bool
	Len int
}

func ConstructFlags()(Flags FlagStruct) {
	return FlagStruct {
		Verbose: false,
		Help: false,
		Autoexecute: false,
		Noexec: false,
		Len: 4,
	}
}

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
func Print(verbose bool, str ...interface{}) {
	if verbose {
		fmt.Println(str...)
	}
}

/**********************
 * Private Flag Handlers
 *********************/

func flagsHandlerHelp(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Help = true
}

func flagsHandlerVerbose(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Verbose = true
}

func flagsHandlerAutoexecute(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Autoexecute = true
}

func flagsHandlerNoexec(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Noexec = true
}

/**********************
 * ArgParse
 *********************/

// parse the user input for potential prompts
func ArgParse(args []string, Flags *FlagStruct) (prompt string, err error) {
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
	for i := 1; i < Flags.Len+1 && i < len(args) && validArg; i++ {
		switch args[i] {
		case "-h": // help commands (just display useage)
			fallthrough
		case "--help":
			flagsHandlerHelp(Flags, &startPromptIndex)
			return "", nil // exit early if help flag
		case "-v": // handle program verbosity
			fallthrough
		case "--verbose":
			flagsHandlerVerbose(Flags, &startPromptIndex)
		case "-y": // execute command automatically
			flagsHandlerAutoexecute(Flags, &startPromptIndex)
		case "-n":
			flagsHandlerNoexec(Flags, &startPromptIndex)
		default:
			// skip searching for switches if invalid arg is found (assume it is prompt)
			validArg = false
			break
		}
	}

	// noexec takes priority over autoexecute, turn off autoexec
	// guarentees mutual exclusivity
	if Flags.Noexec && Flags.Autoexecute {
		Flags.Autoexecute = false
	}

	prompt = strings.Join(args[startPromptIndex:], " ")
	if prompt == "" {
		return "", errors.New("Prompt cannot be empty") 
	}

	return prompt, nil
}