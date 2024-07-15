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
    prompt := "This is an example prompt"
    cli_input := "fl " + prompt
    parsed, err := argParse(strings.Split(cli_input, " "))
    if parsed != prompt {
        t.Fatalf(`argParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}