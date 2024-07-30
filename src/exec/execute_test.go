package exec

import (
	"regexp"
	"strings"
	"testing"
)

// test cmd gen when quotes are involved
func TestQuotes(t *testing.T) {
	cmd_str := `echo 'this is a command with "quotes"'`
	expected_cmd := `echo 'this is a command with "quotes"'`
	expected_result := `this is a command with "quotes"`

	Cmd := Command(cmd_str)

	generated_cmd := Cmd.Cmd.String()
	match, err := regexp.MatchString(expected_cmd, generated_cmd)
	if !match || err != nil {
		t.Fatalf(`Command("%s") = "%s", expected "%s"`, cmd_str, generated_cmd, expected_cmd)
	}

	generated_result, err := Cmd.Exec()
	if generated_result != expected_result || err != nil {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, expected_cmd, generated_result, err, expected_result, nil)
	}
}

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
