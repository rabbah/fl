package exec

/*
 * @NOTE: Exec() shell output may generate a \n. check for it.
 * these tests assume a bash environment
 */

import (
	"os"
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

// test cmd gen with multiple lines
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

// test that ~ and $HOME expand to the same thing (os.UserHomeDir(). also check not empty.)
func TestTildeAndEnvVarExpansion(t *testing.T) {
	cmd_str_tilde := "echo ~"
	cmd_str_env := "echo $HOME"
	expected, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf(`error executing os.UserHomeDir()`)
	}
	// exec outputs an extra \n
	expected = expected + "\n"

	Cmd_tilde := Command(cmd_str_tilde)
	Cmd_env := Command(cmd_str_env)
	// skip Command() test, unecessary for this test

	// check os.getenv not empty
	if match, _ := regexp.MatchString("/", expected); !match {
		t.Fatalf(`os.UserHomeDir() is empty, but it should have a value`)
	}
	// check tilde = $HOME = os.Getenv($HOME)
	exec_tilde, err := Cmd_tilde.Exec()
	if exec_tilde != expected {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, Cmd_tilde, exec_tilde, err, expected, nil)
	}
	exec_env, err := Cmd_env.Exec()
	if exec_env != expected {
		t.Fatalf(`Exec("%s") = ("%s","%v"), expected ("%s","%v")`, Cmd_env, exec_env, err, expected, nil)
	}
}

/*
 * @TODO:
 * > >> |
 * ;
 * && ||
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
