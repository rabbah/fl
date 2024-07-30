package helpers

import (
	"errors"
	"fmt"
	"strings"
)

/*
 * To add a new flag:
 * 1. Add it to FlagStruct and ConstructFlags (increment the constructor's Len as well)
 * 2. Create and add the flag handler to the switch-case in argParse()
 * 3. Update the flag's entry in Usage()
 */

// global flag structure
type FlagStruct struct {
	Verbose, Help, Noexec, Tui, Output bool
	Outfile, Prompt                    string
	Len                                int
}

func ConstructFlags() (Flags FlagStruct) {
	return FlagStruct{
		Verbose: false,
		Help:    false,
		Noexec:  false,
		Tui:     false,
		Output:  false,
		Outfile: "",
		Prompt:  "",
		Len:     5,
	}
}

// useage definition functions to explain command and its args
var Usage = func() {
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

func flagsHandlerNoexec(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Noexec = true
}

func flagsHandlerOutput(Flags *FlagStruct, startPromptIndex *int, outfile string) {
	*startPromptIndex += 2
	Flags.Output = true
	Flags.Outfile = outfile // pass name of outfile
}

/**********************
 * ArgParse
 *********************/

// parse the user input for potential prompts
func ArgParse(args []string, Flags *FlagStruct) (err error) {
	// Check if command line arguments are provided
	if len(args) < 2 {
		// expecting at least 2 arguments
		return errors.New("expecting at least 2 args")
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
			return nil // exit early if help flag
		case "-v": // handle program verbosity
			fallthrough
		case "--verbose":
			flagsHandlerVerbose(Flags, &startPromptIndex)
		case "-n":
			flagsHandlerNoexec(Flags, &startPromptIndex)
		case "-o":
			flagsHandlerOutput(Flags, &startPromptIndex, args[i+1])
			i++ // skip next arg (it should be filename)
		default:
			// skip searching for switches if invalid arg is found (assume it is prompt)
			validArg = false
		}
	}

	// if -o raised but empty filename passed, use default filename
	// (this implies no prompt was passed either, but still safety check)
	if Flags.Output && Flags.Outfile == "" {
		return errors.New("outfile cannot be empty")
	}

	Flags.Prompt = strings.Join(args[startPromptIndex:], " ")
	if Flags.Prompt == "" {
		return errors.New("prompt cannot be empty")
	}

	return nil
}
