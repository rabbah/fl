package utils

import (
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Failed to initialize clipboard.")
		os.Exit(128)
	}
}

func Clip(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
