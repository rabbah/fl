package exec

/*
 * @NOTE: Exec() shell output may generate a \n. check for it.
 * these tests assume a bash environment
 */

import (
	"regexp"
	"strings"
	"testing"
)

// test cmd gen when quotes are involved
func TestQuotes(t *testing.T) {
	cmd_str := `echo 'this is a command with "quotes"'`
	expected_cmd := `echo 'this is a command with "quotes"'`
	expected_result := "this is a command with \"quotes\"\n"

	Cmd := Command(cmd_str)

	generated_cmd := Cmd.Cmd.String()
	match, _ := regexp.MatchString(expected_cmd, generated_cmd)
	if !match {
		t.Fatalf(`Command("%s") = "%s", expected "%s"`, cmd_str, generated_cmd, expected_cmd)
	}

	generated_result, err := Cmd.Exec()
	if generated_result != expected_result || err != nil {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, expected_cmd, generated_result, err, expected_result, nil)
	}
}

// test cmd gen when using <<<
func TestStdIn(t *testing.T) {
	cmd_str := "tr A-Z a-z <<< TeSTsTrIng"
	expected_cmd := "tr A-Z a-z <<< TeSTsTrIng"
	expected_result := "teststring\n"

	Cmd := Command(cmd_str)

	generated_cmd := Cmd.Cmd.String()
	match, _ := regexp.MatchString(expected_cmd, generated_cmd)
	if !match {
		t.Fatalf(`Command("%s") = "%s", expected "%s"`, cmd_str, generated_cmd, expected_cmd)
	}

	generated_result, err := Cmd.Exec()
	if generated_result != expected_result || err != nil {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, expected_cmd, generated_result, err, expected_result, nil)
	}
}

// test cmd gen when using <<<
func TestMultilineCmd(t *testing.T) {
	cmd_str := "mkdir test\n" +
		"ls\n" +
		"rmdir test"
	expected_cmd := "mkdir test\n" +
		"ls\n" +
		"rmdir test"
	// variable, depends on ls results!
	// however, should show as output entry of ls - hence \ntest\n
	expected_result := "\ntest\n"

	Cmd := Command(cmd_str)

	generated_cmd := Cmd.Cmd.String()
	match, _ := regexp.MatchString(expected_cmd, generated_cmd)
	if !match {
		t.Fatalf(`Command("%s") = "%s", expected "%s"`, cmd_str, generated_cmd, expected_cmd)
	}

	generated_result, err := Cmd.Exec()
	match, _ = regexp.MatchString(expected_result, generated_result)
	if !match || err != nil {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, expected_cmd, generated_result, err, expected_result, nil)
	}
}

/*
 * @TODO:
 * > >> |
 * ;
 * && ||
 * ~ $HOME
 */

// create a testfile using EXEC, check command with flags
func TestExecution(t *testing.T) {
	filename := "TESTFILE"
	cmd_str := "touch " + filename
	cmd_ls := "ls -a -l"
	cmd_cleanup := "rm " + filename

	// create file
	Cmd := Command(cmd_str)
	res, err := Cmd.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_str, res, err)
	}

	// check if file exists
	Cmd = Command(cmd_ls)
	res, err = Cmd.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_ls, res, err)
	}
	if !strings.Contains(string(res), filename) {
		t.Fatalf(`Testfile not properly created`)
	}

	// delete file
	Cmd = Command(cmd_cleanup)
	res, err = Cmd.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_cleanup, res, err)
	}

	// check if file was deleted
	Cmd = Command(cmd_ls)
	res, err = Cmd.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_ls, res, err)
	}
	if strings.Contains(string(res), filename) {
		t.Fatalf(`Testfile not properly deleted`)
	}
}
