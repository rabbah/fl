package helpers

import (
	"errors"
	"fl/auth"
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

func confHandlerFlid(Config *io.Config, arg string) {
	switch arg {
	case "reset":
		Config.FLID = ""
	case "restore":
		auth.RestoreFLID()
	}
}

/**********************
 * ConfParse
 *********************/

var (
	regex_autoexecute = regexp.MustCompile("--autoexecute=")
	regex_language    = regexp.MustCompile("--language=")
	regex_flid        = regexp.MustCompile("flid")
)

// check if this is a config command - follow format 'fl config <CONFIGCMD>'
func ConfParse(args []string, Config *io.Config) (confCmd bool, err error) {
	if len(args) > 1 && args[1] == "conf" {
		if regex_autoexecute.MatchString(args[2]) {
			confHandlerAutoexec(Config, args[2])
		} else if regex_language.MatchString(args[2]) {
			confHandlerLanguage(Config, args[2])
		} else if regex_flid.MatchString(args[2]) {
			if len(args) > 2 {
				confHandlerFlid(Config, args[3])
			}
		} else {
			return true, errors.New("config param not found")
		}
		return true, nil
	} else {
		return false, nil
	}
}
