package main

import(
    "testing"
)

// test argparse when help flag is detected
func TestArgParseHelp(t *testing.T) {
    prompt := "This is an example prompt"
    msg := argParse()
    if prompt != msg {
        t.Fatalf(`prompt: '%s' not equal to msg: '%s'`, prompt, msg)
    }
}