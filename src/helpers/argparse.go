package helpers

import (
	"errors"
	"fl/io"
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
	Verbose, Help, PromptExec, Output, Explain bool
	Outfile, Prompt, Language                  string
	Len                                        int
}

func ConstructFlags(Config io.Config) (Flags FlagStruct) {
	return FlagStruct{
		Verbose:    false,
		Help:       false,
		PromptExec: false,
		Output:     false,
		Explain:    false,
		Outfile:    "",
		Prompt:     "",
		Language:   Config.Language,
		Len:        7,
	}
}

/**********************
 * Other args and parsing helpers
 *********************/

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// useage definition functions to explain command and its args
var Usage = func() {
	fmt.Print(`
fl by itself will open the graphical interface. Otherwise, prompt is required.

Usage: fl [-hnvt] [-o filename] [-l language] prompt...

 -h,--help              show command usage
 -p                     prompt for running generated command
 -v,--verbose           display updates of the command progress
 -e                     for extra verification, ask AI what the generated command does
 -o outfile             output generated command to the passed textfile
 -l language            target language to generate code for, default is bash/unix commands
                        ┃ target language can also be a desired command name!

Config: fl conf <config param>

 flid reset             delete your locally saved flid
 flid restore           restore your local flid
 --autoexecute=BOOL     enable autoexecution
 --language=STRING      change the DEFAULT language/environment to generate a command for

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

func flagsHandlerExplain(Flags *FlagStruct, startPromptIndex *int) {
	*startPromptIndex++
	Flags.Explain = true
}

func flagsHandlerOutput(Flags *FlagStruct, startPromptIndex *int, outfile string) {
	*startPromptIndex += 2
	Flags.Output = true
	Flags.Outfile = outfile // pass name of outfile
}

func flagsHandlerLanguage(Flags *FlagStruct, startPromptIndex *int, language string) {
	*startPromptIndex += 2
	Flags.Language = language
}

/**********************
 * ArgParse
 *********************/

// parse the user input for potential prompts
func ArgParse(args []string, Flags *FlagStruct) (err error) {
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
		case "-e":
			flagsHandlerExplain(Flags, &startPromptIndex)
		case "-o":
			flagsHandlerOutput(Flags, &startPromptIndex, args[i+1])
			i++ // skip next arg (it should be filename)
		case "-l":
			flagsHandlerLanguage(Flags, &startPromptIndex, args[i+1])
			i++
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

	// fl by itself (no prompt) is an error
	Flags.Prompt = strings.Join(args[startPromptIndex:], " ")
	if Flags.Prompt == "" {
		return errors.New("prompt cannot be empty")
	}

	return nil
}
