package helpers

import (
	"fl/io"
	"strings"
	"testing"
)

/**********************
 * Private helpers for checking equality of flags
 *********************/

// check if array contains value
func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func testFlag(flagName string, Flags FlagStruct, t *testing.T) bool {
	switch flagName {
	case "Verbose":
		return Flags.Verbose
	case "Help":
		return Flags.Help
	case "PromptExec":
		return Flags.PromptExec
	case "Tui":
		return Flags.Tui
	case "Output":
		return Flags.Output
	case "Explain":
		return Flags.Explain
	}

	t.Fatalf(`Flag %s is not implemented`, flagName)
	return false
}

// helper for validating all flags per test
func testFlagsExpected(expected []string, Flags FlagStruct, t *testing.T) bool {
	allFlags := []string{"Verbose", "Help", "PromptExec", "Tui", "Output", "Explain"}
	for _, flagName := range allFlags {
		v := testFlag(flagName, Flags, t)
		if contains(expected, flagName) && !v {
			// flag is in expected but actually false
			return false
		} else if !contains(expected, flagName) && v {
			// flag is not expected but is found to be true
			return false
		}
	}

	return true
}

/**********************
 * Basic pos/neg tests
 *********************/

// test ArgParse extracts prompt with no flags
func TestArgParseNoFlags(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	prompt := "This is an example prompt"
	cli_input := "fl" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// one flag with no prompt
func TestArgParseOneFlagNoPrompt(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	prompt := ""
	cli_input := "fl -v" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err == nil { // should cause an error when no prompt given
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test ArgParse sets tui flag if no args passed
func TestArgParseEmpty(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Tui"}

	prompt := ""
	cli_input := "fl" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "t", Flags)
	}
}

/**********************
 * Individual flag interactions
 *********************/

// test help
func TestArgParseHelp(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Help"}

	prompt := "This is an example prompt"
	expectedPrompt := "" // skip prompt when -h is found
	cli_input := "fl -h" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "h", Flags)
	}
	if Flags.Prompt != expectedPrompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, expectedPrompt)
	}
}

// test tui
func TestArgParseTui(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Tui"}

	prompt := "This is an example prompt"
	cli_input := "fl -t" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "t", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test verbose
func TestArgParseVerbose(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Verbose"}

	prompt := "This is an example prompt"
	cli_input := "fl -v" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "v", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test noexec
func TestArgParsePromptExec(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"PromptExec"}

	prompt := "This is an example prompt"
	cli_input := "fl -p" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "p", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

func TestArgParseExplain(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Explain"}

	prompt := "This is an example prompt"
	cli_input := "fl -e" + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "n", Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

// test output
func TestArgParseOutput(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"Output"}

	prompt := "This is an example prompt"
	outfile := "outfile"
	cli_input := "fl -o " + outfile + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects only the %s flag. Actual: %+v`, cli_input, "o", Flags)
	}
	if Flags.Outfile != outfile {
		t.Fatalf(`ArgParse("%s") yields outfile = '%s'. Expected '%s'`, cli_input, Flags.Outfile, outfile)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

func TestArgParseLanguage(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{}

	prompt := "This is an example prompt"
	language := "powershell"
	cli_input := "fl -l " + language + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") should not set these flags. Actual: %+v`, cli_input, Flags)
	}
	if Flags.Language != language {
		t.Fatalf(`ArgParse("%s") yields language = '%s'. Expected '%s'`, cli_input, Flags.Language, language)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
}

/*
 * config opts
 */
// test autoexec
func TestConfParseAutoexec(t *testing.T) {
	Config := io.NewConf()

	prompt := "This is an example prompt"
	cli_input := "fl conf --autoexecute=true" + " " + prompt
	wasConfCmd, err := ConfParse(strings.Split(cli_input, " "), &Config)
	if !wasConfCmd || err != nil {
		t.Fatalf(`ConfParse("%s") = (%v, %v). Expected (true, nil)`, cli_input, wasConfCmd, err)
	}
	if !Config.Autoexec {
		t.Fatalf(`ConfParse("%s") should set only Autoexec: %v`, cli_input, Config)
	}
}

func TestConfParseLanguage(t *testing.T) {
	Config := io.NewConf()

	prompt := "This is an example prompt"
	language := "powershell"
	cli_input := "fl conf --language=" + language + " " + prompt
	wasConfCmd, err := ConfParse(strings.Split(cli_input, " "), &Config)
	if !wasConfCmd || err != nil {
		t.Fatalf(`ConfParse("%s") = (%v, %v). Expected (true, nil)`, cli_input, wasConfCmd, err)
	}
	if Config.Language != language {
		t.Fatalf(`ConfParse("%s") = %s. Expected '%s'`, cli_input, Config.Language, language)
	}
}

/**********************
* Validate multiple flag interactions
*********************/

// test help activates despite invalid prompt
func TestArgParseHelpNoPrompt(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

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
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

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

// test ArgParse with all flags + prompt (except -h!)
func TestArgParseAllFlags(t *testing.T) {
	Config := io.NewConf()
	Flags := ConstructFlags(Config)

	expectedFlags := []string{"PromptExec", "Verbose", "Output", "Tui", "Explain"}

	prompt := "This is an example prompt"
	outfile := "outfile"
	language := "powershell"
	cli_input := "fl -p -v -t -e -o " + outfile + " -l " + language + " " + prompt
	err := ArgParse(strings.Split(cli_input, " "), &Flags)
	if !testFlagsExpected(expectedFlags, Flags, t) {
		t.Fatalf(`ArgParse("%s") expects all and only the -p, -v, -o, -t flags. Actual: %+v`, cli_input, Flags)
	}
	if Flags.Prompt != prompt || err != nil {
		t.Fatalf(`ArgParse("%s") = "%s", %v. Expected '%s'`, cli_input, Flags.Prompt, err, prompt)
	}
	if Flags.Outfile != outfile {
		t.Fatalf(`ArgParse("%s") yields outfile = '%s'. Expected '%s'`, cli_input, Flags.Outfile, outfile)
	}
	if Flags.Language != language {
		t.Fatalf(`ArgParse("%s") yields outfile = '%s'. Expected '%s'`, cli_input, Flags.Outfile, language)
	}
}
