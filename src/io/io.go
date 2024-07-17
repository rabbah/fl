package io

import (
	"os"
)

const (
	outfilePerms = 0666
)

func Output(outfile string, data string) (err error) {
	err = os.WriteFile(outfile, []byte(data), outfilePerms)
	return err
}
