package exec

/*
 * @NOTE: Exec() shell output may generate a \n. check for it.
 * these tests assume a bash environment
 */

import (
	"os"
	"os/exec"
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

// test cmd gen when using > and >>
func TestRedirectAndOutput(t *testing.T) {
	cmd_str_redirect := "echo test1 > test.txt"
	cmd_str_append := "echo test2 >> test.txt"
	test_cmd_contents := "cat test.txt"
	expected_contents := "test1\ntest2\n"

	Cmd_redirect := Command(cmd_str_redirect)
	Cmd_append := Command(cmd_str_append)

	// redirect first and check for no error
	_, err := Cmd_redirect.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") returned err: %v`, cmd_str_redirect, err)
	}

	// check append for err
	_, err = Cmd_append.Exec()
	if err != nil {
		t.Fatalf(`Exec("%s") returned err: %v`, cmd_str_append, err)
	}

	// check that the contents were written
	actual_contents, _ := Command(test_cmd_contents).Exec()
	if actual_contents != expected_contents {
		exec.Command("rm", "test.txt").Run()
		t.Fatalf(
			"actual != expected\n\n" +
				"actual:\n" +
				actual_contents +
				"\nexpected:\n" +
				expected_contents +
				"\n\n",
		)
	}

	// clean up
	exec.Command("rm", "test.txt").Run()
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

// test cmd gen with "|"
func TestPipeCmd(t *testing.T) {
	cmd_str := "echo 'hello world\n" +
		"bye world' " +
		"| grep 'hello'"
	expected_cmd := "echo 'hello world\n" +
		"bye world' " +
		"| grep 'hello'"
	expected_result := "hello world\n"

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
func TestSemicolonCmd(t *testing.T) {
	cmd_str := "mkdir test;ls;rmdir test"
	expected_cmd := "mkdir test;ls;rmdir test"
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
		exec.Command("rm", filename).Run()
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_cleanup, res, err)
	}

	// check if file was deleted
	Cmd = Command(cmd_ls)
	res, err = Cmd.Exec()
	if err != nil {
		exec.Command("rm", filename).Run()
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_ls, res, err)
	}
	if strings.Contains(string(res), filename) {
		exec.Command("rm", filename).Run()
		t.Fatalf(`Testfile not properly deleted`)
	}
}
