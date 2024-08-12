package helpers

import (
	"errors"
	"fl/io"
	"regexp"
	"strings"
)

/**********************
 * Private Flag Handlers
 *********************/

func confHandlerAutoexec(Config *io.Config, arg string) {
	// rewrite autoexec with right side of '='
	value := strings.Split(arg, "=")[1]
	if strings.ToLower(value) == "true" {
		Config.Autoexec = true
	} else {
		Config.Autoexec = false
	}
}

func confHandlerLanguage(Config *io.Config, arg string) {
	// rewrite autoexec with right side of '='
	value := strings.Split(arg, "=")[1]
	Config.Language = strings.ToLower(value)
}

/**********************
 * ConfParse
 *********************/

var (
	regex_autoexecute = regexp.MustCompile("--autoexecute=")
	regex_language    = regexp.MustCompile("--language=")
)

// check if this is a config command - follow format 'fl config <CONFIGCMD>'
func ConfParse(args []string, Config *io.Config) (confCmd bool, err error) {
	if len(args) > 1 && args[1] == "conf" {
		if regex_autoexecute.MatchString(args[2]) {
			confHandlerAutoexec(Config, args[2])
		} else if regex_language.MatchString(args[2]) {
			confHandlerLanguage(Config, args[2])
		} else {
			return true, errors.New("config param not found")
		}
		return true, nil
	} else {
		return false, nil
	}
}
