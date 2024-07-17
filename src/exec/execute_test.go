package exec

import (
	"strings"
	"testing"
)

// create a testfile using EXEC, check command with flags
func TestExecution(t *testing.T) {
	filename := "TESTFILE"
	cmd := "touch " + filename
	cmd_ls := "ls -a -l"
	cmd_cleanup := "rm " + filename

	// create file
	res, err := Exec(cmd)
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd, res, err)
	}

	// check if file exists
	res, err = Exec(cmd_ls)
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_ls, res, err)
	}
	if !strings.Contains(string(res), filename) {
		t.Fatalf(`Testfile not properly created`)
	}

	// delete file
	res, err = Exec(cmd_cleanup)
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_cleanup, res, err)
	}

	// check if file was deleted
	res, err = Exec(cmd_ls)
	if err != nil {
		t.Fatalf(`Exec("%s") = "%s", %v. Expected no err`, cmd_ls, res, err)
	}
	if strings.Contains(string(res), filename) {
		t.Fatalf(`Testfile not properly deleted`)
	}
}
