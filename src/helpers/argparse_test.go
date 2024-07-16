package helpers

import(
    "testing"
    "strings"
)

// test ArgParse extracts prompt with no flags
func TestArgParseNoFlags(t *testing.T) {
    Flags := ConstructFlags()

    prompt := "This is an example prompt"
    cli_input := "fl" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if parsed != prompt || err != nil {
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// one flag with no prompt
func TestArgParseOneFlagNoPrompt(t *testing.T) {
    Flags := ConstructFlags()

    prompt := ""
    cli_input := "fl -v" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if parsed != prompt || err == nil { // should cause an error when no prompt given
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test help activates despite invalid prompt
func TestArgParseHelpNoPrompt(t *testing.T) {
    Flags := ConstructFlags()

    prompt := ""
    cli_input := "fl -v -h" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if !Flags.Help {
        t.Fatalf(`help = 'false'. Expected 'true'`)
    }
    if parsed != prompt || err != nil {
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test help activates and skips prompt parsing
func TestArgParseHelp(t *testing.T) {
    Flags := ConstructFlags()

    prompt := "This is an example prompt"
    expectedPrompt := "" // skip prompt when -h is found
    cli_input := "fl -v -h" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if !Flags.Help {
        t.Fatalf(`help = 'false'. Expected 'true'`)
    }
    if parsed != expectedPrompt || err != nil {
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, expectedPrompt)
    }
}

// test ArgParse raises err when prompt is empty
func TestArgParseEmpty(t *testing.T) {
    Flags := ConstructFlags()

    prompt := ""
    cli_input := "fl" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if err == nil {
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}

// test ArgParse with all flags + prompt
func TestArgParseAllFlags(t *testing.T) {
    Flags := ConstructFlags()

    prompt := "This is an example prompt"
    cli_input := "fl -y -n -v" + " " + prompt
    parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
    if Flags.Autoexecute || !Flags.Noexec || !Flags.Verbose {
        t.Fatalf(`ArgParse("%s") yields flags y: %t, n: %t, v: %t. Expected y: false, n: true, v: true`, cli_input, Flags.Autoexecute, Flags.Noexec, Flags.Verbose)
    }
    if parsed != prompt || err != nil {
        t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
    }
}