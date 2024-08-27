package utils

import (
	"fmt"

	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Failed to initialize clipboard.")
		panic(err)
	}
}

func Clip(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
