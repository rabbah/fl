package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	home, _     = os.UserHomeDir()
	confPath, _ = filepath.Abs(home + "/" + ".flconfig")

	new_conf_statement = `
    PLEASE READ:
    As a safety precaution, the default behavior of this program is to disable automatic execution of
    generated commands and never execute a command. To change these default behaviors, read'fl -h' to
    view usage.
	Default config path: '$HOME/.flconfig'

`
)

type Config struct {
	Autoexec bool
	Language string
}

func NewConf() Config {
	return Config{
		Autoexec: false,
		Language: "Unix/Bash",
	}
}

// create config for new user w/o a config
func initializeConf() (err error) {
	Config := NewConf()

	err = Config.SaveConf()

	return err
}

// save edited config as-is
func (Config Config) SaveConf() (err error) {
	// marshall
	b, err := json.Marshal(Config)
	output(confPath, b)

	return err
}

func ReadConf() (config Config, err error) {

	// check file exists, create if not
	exists, err := pathExists(confPath)
	if err != nil {
		return config, errors.New("could not 'stat' " + confPath)
	} else if !exists {
		err = initializeConf()
		if err != nil {
			return config, errors.New("could not initialize " + confPath)
		} else {
			// assume this is the user's first execution
			fmt.Print(new_conf_statement)
		}
	}

	// open file
	file, err := os.Open(confPath)
	if err != nil {
		return config, errors.New("could not open " + confPath)
	}
	defer file.Close()

	// unmarshall
	decoder := json.NewDecoder(file)
	config = Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return config, errors.New("could not unmarshall " + confPath)
	}

	return config, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
