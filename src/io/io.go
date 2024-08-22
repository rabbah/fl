package io

import (
	"io"
	"os"
)

const (
	outfilePerms = 0666
)

// output passed string data (public)
func Output(outfile string, data string) (err error) {
	err = output(outfile, []byte(data))
	return err
}

// wrapper for Readall
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// output passed byte data (private)
func output(outfile string, data []byte) (err error) {
	err = os.WriteFile(outfile, data, outfilePerms)
	return err
}
