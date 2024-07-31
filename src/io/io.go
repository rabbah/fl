package io

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	outfilePerms = 0666
	confPath     = ".flconfig"
)

func Output(outfile string, data string) (err error) {
	err = output(outfile, []byte(data))
	return err
}

func output(outfile string, data []byte) (err error) {
	err = os.WriteFile(outfile, data, outfilePerms)
	return err
}

type confEntry uint

const (
	autoexec confEntry = iota
)

type Config struct {
	Autoexec bool
}

func (Config Config) SaveConf() (err error) {
	// marshall
	b, err := json.Marshal(Config)
	output(confPath, b)

	return err
}

func initConf() Config {
	return Config{
		Autoexec: false,
	}
}

func initializeConf() (err error) {
	Config := initConf()

	err = Config.SaveConf()

	return err
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

func ReadConf() (config Config, err error) {

	// check file exists
	exists, err := pathExists(confPath)
	if err != nil {
		return config, errors.New("could not 'stat' " + confPath)
	} else if !exists {
		err = initializeConf()
		if err != nil {
			return config, errors.New("could not initialize " + confPath)
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
