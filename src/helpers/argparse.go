package helpers

import (
	"errors"
	"fl/io"
	"fmt"
	"regexp"
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
	Verbose, Help, PromptExec, Tui, Output bool
	Outfile, Prompt                        string
	Len                                    int
}

func ConstructFlags() (Flags FlagStruct) {
	return FlagStruct{
		Verbose:    false,
		Help:       false,
		PromptExec: false,
		Tui:        false,
		Output:     false,
		Outfile:    "",
		Prompt:     "",
		Len:        5,
	}
}

// useage definition functions to explain command and its args
var Usage = func() {
	fmt.Println(`
fl by itself will open the graphical interface. Otherwise, prompt is required.

Usage: fl [-hnvt] [-o filename] prompt...

 -h,--help              show command usage
 -p                     prompt for running generated command
 -v,--verbose           display updates of the command progress
 -o                     output generated command to the passed textfile
 -t                     enter the graphical interface (TUI)

Config: fl config <config param>

 --autoexecute=BOOL     enable autoexecution
`)
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

func flagsHandlerPrompt(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.PromptExec = true
}

func flagsHandlerOutput(Flags *FlagStruct, startPromptIndex *int, outfile string) {
	*startPromptIndex += 2
	Flags.Output = true
	Flags.Outfile = outfile // pass name of outfile
}

func confHandlerAutoexec(Config *io.Config, arg string) {
	// rewrite autoexec with right side of '='
	value := strings.Split(arg, "=")[1]
	if strings.ToLower(value) == "true" {
		Config.Autoexec = true
	} else {
		Config.Autoexec = false
	}
}

/**********************
 * Other parse helpers
 *********************/

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

/**********************
 * ConfParse
 *********************/

var (
	regex_autoexecute = regexp.MustCompile("--autoexecute=")
)

// check if this is a config command - follow format 'fl config <CONFIGCMD>'
func ConfParse(args []string, Config *io.Config) (confCmd bool, err error) {
	if args[1] == "config" {
		if regex_autoexecute.MatchString(args[2]) {
			confHandlerAutoexec(Config, args[2])
		} else {
			return true, errors.New("config param not found")
		}
		return true, nil
	} else {
		return false, nil
	}
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
		case "-p":
			flagsHandlerPrompt(Flags, &startPromptIndex)
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
