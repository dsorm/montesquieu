package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// verifies config from user; if the string is not null, there's at least one error in the config
func (cfg *file) verifyConfig() string {
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

	// verify articles per page
	if num, err := strconv.Atoi(cfg.ArticlesPerPage); err != nil || num <= 0 {
		str += "ArticlesPerPage has to be a valid positive integer\n"
	}

	// verify database type
	validType := false
	switch cfg.Store {
	case "":
		str += "Store can't be empty\n"
	case "postgres":
		validType = true
	}

	if !validType {
		str += "Store is invalid\n"
	}

	// verify caching engine
	if cfg.CachingStore == "" {
		str += "CachingStore can't be empty\n"
	} else {
		validType := false

		if cfg.CachingStore == "internal" {
			validType = true
		}

		if cfg.CachingStore == "off" {
			validType = true
		}

		if !validType {
			str += "CachingStore is invalid\n"
		}
	}

	// verify live templates
	if cfg.HotSwapTemplates == "" {
		str += "HotSwapTemplates can't be empty"
	} else if !(strings.ToLower(cfg.HotSwapTemplates) == "yes" || strings.ToLower(cfg.HotSwapTemplates) == "no") {
		str += "HotSwapTemplates can only be either 'yes' or 'no'"
	}

	return str
}

func (cfg *file) readConfig() {
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

func (cfg *file) createConfig() {
	// open file
	file, err := os.Create("config.json")
	if err != nil {
		panic("Failed to create config.json")
	}

	// default configuration file
	cfg.BlogName = "My blog"
	cfg.ListenOn = ":8080"
	cfg.Store = "postgres"
	cfg.CachingStore = "off"
	cfg.ArticlesPerPage = "5"

	// marshal json and save
	bytes, _ := json.MarshalIndent(cfg, "", "\t")
	_, err = file.Write(bytes)
	if err != nil {
		panic("Failed to write config.json")
	}

}

//noinspection ALL
func NewConfig() (*Config, error) {
	cfg := &file{}

	// check if config does exist, and if it doesn't, create a new one
	if file, err := os.Open("config.json"); err != nil {
		err = file.Close()
		if err != nil {
			fmt.Println("Error while opening config")
		}
		cfg.createConfig()
	}

	// read and verify the config
	cfg.readConfig()
	errs := cfg.verifyConfig()

	// if any errors were found, lets return the errors
	if len(errs) != 0 {
		return nil, errors.New(errs)
	}

	// parse and return config
	return cfg.parseFile(), nil
}
