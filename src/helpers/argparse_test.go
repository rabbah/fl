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
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// one flag with no prompt
func TestArgParseOneFlagNoPrompt(t *testing.T) {
	Flags := ConstructFlags()

	prompt := ""
	cli_input := "fl -v" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err == nil { // should cause an error when no prompt given
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test ArgParse sets TUI if no args passed
func TestArgParseEmpty(t *testing.T) {
	Flags := ConstructFlags()

	prompt := ""
	cli_input := "fl" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help || Flags.Output || !Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "o", Flags)
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
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || !Flags.Help || Flags.Output || Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "h", Flags)
	}
	if Flags.Prompt != expectedPrompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, expectedPrompt)
	}
}

// test tui
func TestArgParseTui(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -t" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help || Flags.Output || !Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "t", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test verbose
func TestArgParseVerbose(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -v" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || !Flags.Verbose || Flags.Help || Flags.Output || Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "v", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test autoexec
func TestArgParseAutoexec(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -y" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help || Flags.Output || Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "y", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test noexec
func TestArgParseNoexec(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -n" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || !Flags.Noexec || Flags.Verbose || Flags.Help || Flags.Output || Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "n", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test output
func TestArgParseOutput(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	outfile := "outfile"
	cli_input := "fl -o " + outfile + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || Flags.Noexec || Flags.Verbose || Flags.Help || !Flags.Output || Flags.Tui {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "o", Flags)
	}
	if Flags.Outfile != outfile {
		t.Fatalf(`ArgParse("%s") yields outfile = '%s'. Expected '%s'`, cli_input, Flags.Outfile, outfile)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
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
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !Flags.Help {
		t.Fatalf(`help = 'false'. Expected 'true'`)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test help with multiple flags, verify skip prompt parsing
func TestArgParseHelpMultipleFlags(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	expectedPrompt := "" // skip prompt when -h is found
	cli_input := "fl -v -h" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !Flags.Help {
		t.Fatalf(`help = 'false'. Expected 'true'`)
	}
	if Flags.Prompt != expectedPrompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, expectedPrompt)
	}
}

// test ArgParse with all flags + prompt
func TestArgParseAllFlags(t *testing.T) {
	Flags := ConstructFlags()

	prompt := "This is an example prompt"
	cli_input := "fl -y -n -v" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Autoexecute || !Flags.Noexec || !Flags.Verbose {
		t.Fatalf(`ArgParse("%s") yields flags %+v. Expected y: false, n: true, v: true`, cli_input, Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}
