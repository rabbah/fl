package exec

import (
	"fmt"
	"strings"
	"os/exec"
)

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

func Exec(result string)(res string, err error) {
	// convert to arr of values (exec requires a specific format)
	fullCmd := strings.Split(result, " ")
	cmd := fullCmd[0]
	args := []string{}

	if len(fullCmd) > 1 {
		args=fullCmd[1:]
	}

	var out []byte
	out, err = exec.Command(cmd, args...).Output()

	return string(out), err
}