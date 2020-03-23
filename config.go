package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

type Config struct {
	// Blog's name
	BlogName string

	// The port on which the server should listen on
	ListenOn string

	/*
	 Type of database
	 Currently only 'sqlite' is supported
	*/
	DBType string
}

func (cfg *Config) verifyConfig() string {
	str := ""

	// verify blogname
	if cfg.BlogName == "" {
		str += "BlogName can't be empty\n"
	}
	if utf8.RuneCountInString(cfg.BlogName) > 240 {
		str += "BlogName can't be longer than 240 characters\n"
	}

	// verify port
	if cfg.ListenOn == "" {
		str += "ListenOn can't be empty\n"
	}

	// verify database type
	if cfg.DBType == "" {
		str += "DBType can't be empty\n"
	}

	validType := false

	validType = cfg.DBType == "sqlite"

	if !validType {
		str += "DBType is invalid\n"
	}

	return str
}

func (cfg *Config) readConfig() {
	// open file
	file, err := os.Open("config.json")
	if err != nil {
		panic("Can't read config.json")
	}

	// read json, unmarshal and return
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Error while reading config.json")
	}
	if json.Unmarshal(bytes, &cfg) != nil {
		panic("The syntax of config.json is invalid")
	}
}

func (cfg *Config) createConfig() {
	// open file
	file, err := os.Create("config.json")
	if err != nil {
		panic("Failed to create config.json")
	}

	// default configuration file
	cfg.BlogName = "My blog"
	cfg.ListenOn = ":8080"
	cfg.DBType = "sqlite"

	// marshal json and save
	bytes, _ := json.MarshalIndent(cfg, "", "\t")
	_, err = file.Write(bytes)
	if err != nil {
		panic("Failed to write config.json")
	}

}

func NewConfig() *Config {
	cfg := &Config{}

	// check if config does exist, and if it doesn't, create a new one
	if file, err := os.Open("config.json"); err != nil {
		_ = file.Close()
		cfg.createConfig()
	}

	// read and verify the config
	cfg.readConfig()
	errors := cfg.verifyConfig()
	if len(errors) != 0 {
		panic(errors)
	}

	return cfg
}
