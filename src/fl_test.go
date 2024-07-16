package main

import(
    "testing"
    "strings"
)

/**********************
 * argparse
 *********************/

// test argparse extracts prompt with no flags
func TestArgParseNoFlags(t *testing.T) {
    Flags := new_flags()

    prompt := "This is an example prompt"
    cli_input := "fl" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if parsed != prompt || err != nil {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// one flag with no prompt
func TestArgParseOneFlagNoPrompt(t *testing.T) {
    Flags := new_flags()

    prompt := ""
    cli_input := "fl -v" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if parsed != prompt || err == nil { // should cause an error when no prompt given
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test help activates despite invalid prompt
func TestArgParseHelpNoPrompt(t *testing.T) {
    Flags := new_flags()

    prompt := ""
    cli_input := "fl -v -h" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if !Flags.help {
        t.Fatalf(`help = 'false'. Expected 'true'`)
    }
    if parsed != prompt || err != nil {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test help activates and skips prompt parsing
func TestArgParseHelp(t *testing.T) {
    Flags := new_flags()

    prompt := "This is an example prompt"
    expectedPrompt := "" // skip prompt when -h is found
    cli_input := "fl -v -h" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if !Flags.help {
        t.Fatalf(`help = 'false'. Expected 'true'`)
    }
    if parsed != expectedPrompt || err != nil {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, expectedPrompt)
    }
}

// test argparse raises err when prompt is empty
func TestArgParseEmpty(t *testing.T) {
    Flags := new_flags()

    prompt := ""
    cli_input := "fl" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if err == nil {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test argparse with all flags + prompt
func TestArgParseAllFlags(t *testing.T) {
    Flags := new_flags()

    prompt := "This is an example prompt"
    cli_input := "fl -y -n -v" + " " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "), &Flags)
    if Flags.autoexecute || !Flags.noexec || !Flags.verbose {
        t.Fatalf(`argParse("%s") yields flags y: %t, n: %t, v: %t. Expected y: false, n: true, v: true`, cli_input, Flags.autoexecute, Flags.noexec, Flags.verbose)
    }
    if parsed != prompt || err != nil {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}