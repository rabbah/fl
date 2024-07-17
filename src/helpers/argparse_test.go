package helpers

import (
	"strings"
	"testing"
)

/**********************
 * Basic pos/neg tests
 *********************/

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

/**********************
 * Individual flag interactions
 *********************/

// test help
func TestArgParseHelp(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	expectedPrompt := "" // skip prompt when -h is found
	cli_input := "fl -h" + " " + prompt
	parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || !Flags.Help {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual - y: %t, n: %t, v: %t, h: %t`, cli_input, "h", Flags.Autoexecute, Flags.Noexec, Flags.Verbose, Flags.Help)
	}
	if parsed != expectedPrompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, expectedPrompt)
	}
}

// test verbose
func TestArgParseVerbose(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -v" + " " + prompt
	parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || !Flags.Verbose || Flags.Help {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual - y: %t, n: %t, v: %t, h: %t`, cli_input, "v", Flags.Autoexecute, Flags.Noexec, Flags.Verbose, Flags.Help)
	}
	if parsed != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
	}
}

// test autoexec
func TestArgParseAutoexec(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -y" + " " + prompt
	parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual - y: %t, n: %t, v: %t, h: %t`, cli_input, "y", Flags.Autoexecute, Flags.Noexec, Flags.Verbose, Flags.Help)
	}
	if parsed != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
	}
}

// test noexec
func TestArgParseNoexec(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -n" + " " + prompt
	parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || !Flags.Noexec || Flags.Verbose || Flags.Help {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual - y: %t, n: %t, v: %t, h: %t`, cli_input, "n", Flags.Autoexecute, Flags.Noexec, Flags.Verbose, Flags.Help)
	}
	if parsed != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
	}
}

// test output
func TestArgParseOutput(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	outfile := "outfile"
	cli_input := "fl -o " + outfile + " " + prompt
	parsed, err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help || !Flags.Output {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual - y: %t, n: %t, v: %t, h: %t`, cli_input, "n", Flags.Autoexecute, Flags.Noexec, Flags.Verbose, Flags.Help)
	}
	if Flags.Outfile != outfile {
		t.Fatalf(`ArgParse("%s") yields outfile = '%s'. Expected '%s'`, cli_input, Flags.Outfile, outfile)
	}
	if parsed != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, parsed, err, prompt)
	}
}

/**********************
 * Validate multiple flag interactions
 *********************/

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

// test help with multiple flags, verify skip prompt parsing
func TestArgParseHelpMultipleFlags(t *testing.T) {
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
