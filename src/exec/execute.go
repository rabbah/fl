package exec

import (
	"fmt"
	"os/exec"
	"strings"
)

// wrap os.exec struct for decoupling
type Exec struct {
	Cmd *exec.Cmd
}

func PromptExec() (userExec bool) {
	var userInput string
	fmt.Print("Would you like to execute the command? (y/n): ")
	fmt.Scanln(&userInput)
	userInput = strings.ToLower(userInput)
	if userInput == "y" || userInput == "yes" {
		return true
	}

	return false
}

func Command(result string) Exec {
	// convert to arr of values (exec requires a specific format)
	fullCmd := strings.Split(result, " ")
	cmd := fullCmd[0]
	args := []string{}

	if len(fullCmd) > 1 {
		args = fullCmd[1:]
	}

	out := exec.Command(cmd, args...)
	return Exec{Cmd: out}
}

func (ex Exec) Exec() (res string, err error) {
	var tmp []byte
	tmp, err = ex.Cmd.Output()
	return string(tmp), err
}
