package io

import (
	"os"
	"testing"
)

func TestOutput(t *testing.T) {
	outfile := "testfile"
	outdata := "Hello World!"
	var actualdata []byte

	err := Output(outfile, outdata)
	// check error on output
	if err != nil {
		t.Fatalf(`Output("%s", "%s") = "%v". Expected <nil>`, outfile, outdata, err)
	}
	// check file was created
	if _, err = os.Stat(outfile); err != nil {
		t.Fatalf(`File was not created`)
	}
	// check contents are equal as expected
	actualdata, err = os.ReadFile(outfile)
	if outdata != string(actualdata) || err != nil {
		t.Fatalf(`ReadFile('%s') = %s, %v. Expected %s, nil`, outfile, actualdata, err, outdata)
	}

	// remove the created file on success
	os.Remove(outfile)
}
