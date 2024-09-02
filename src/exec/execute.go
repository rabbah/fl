package exec

import (
	"os/exec"
)

// wrap os.exec struct for decoupling
type Exec struct {
	Cmd *exec.Cmd
}

func Command(result string) Exec {
	out := exec.Command("bash", "-c", result)
	return Exec{Cmd: out}
}

func (ex Exec) Exec() (res string, err error) {
	var tmp []byte
	tmp, err = ex.Cmd.Output()
	return string(tmp), err
}
