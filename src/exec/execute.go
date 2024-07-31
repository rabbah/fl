package exec

import (
	"fmt"
	"os/exec"
	"strings"
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

func Command(result string) Exec {
	out := exec.Command("bash", "-c", result)
	return Exec{Cmd: out}
}
