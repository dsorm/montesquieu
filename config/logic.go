package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"unicode/utf8"
)

// verifies config from user; if the string is not null, there's at least one error in the config
func (cfg *ConfigFile) verifyConfig() string {
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
	if cfg.ArticleStore == "" {
		str += "ArticleStore can't be empty\n"
	} else {
		validType := false

		if cfg.ArticleStore == "mock" {
			validType = true
		}

		if !validType {
			str += "ArticleStore is invalid\n"
		}
	}

	// verify caching engine
	if cfg.CachingEngine == "" {
		str += "CachingEngine can't be empty\n"
	} else {
		validType := false

		if cfg.CachingEngine == "internal" {
			validType = true
		}

		if cfg.CachingEngine == "off" {
			validType = true
		}

		if !validType {
			str += "CachingEngine is invalid\n"
		}
	}

	return str
}

func (cfg *ConfigFile) readConfig() {
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

func (cfg *ConfigFile) createConfig() {
	// open file
	file, err := os.Create("config.json")
	if err != nil {
		panic("Failed to create config.json")
	}

	// default configuration file
	cfg.BlogName = "My blog"
	cfg.ListenOn = ":8080"
	cfg.ArticleStore = "mock"
	cfg.CachingEngine = "off"
	cfg.ArticlesPerPage = "5"

	// marshal json and save
	bytes, _ := json.MarshalIndent(cfg, "", "\t")
	_, err = file.Write(bytes)
	if err != nil {
		panic("Failed to write config.json")
	}

}

func NewConfig() (*Config, error) {
	cfg := &ConfigFile{}

	// check if config does exist, and if it doesn't, create a new one
	if file, err := os.Open("config.json"); err != nil {
		_ = file.Close()
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
