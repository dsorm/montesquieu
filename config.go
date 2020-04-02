package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"unicode/utf8"
)

// "parsed" config that's served to the app
type Config struct {
	// Blog's name
	BlogName string

	// The port on which the server should listen on
	ListenOn string

	/*
	 How many articles should be displayed on one index page
	 While technically there's a limit of 2^64 articles on page
	 if you really hate your ArticleStore (and internet connections,
	 browsers, etc.), we recommend sticking to something like 5
	 articles
	*/
	ArticlesPerPage uint64

	/*
	 Type of database
	 Currently only 'mock' and `postgres` is supported
	*/
	ArticleStore ArticleStore

	/*
	 Login info for ArticleStore driver, if needed
	 Postgres: requires all filled out
	 Mock: ignores all the fields
	*/
	ArticleStoreHost     string
	ArticleStoreDB       string
	ArticleStoreUser     string
	ArticleStorePassword string

	/*
	 Type of caching engine used between the app and the store
	 Currently only 'internal' or 'off' is supported
	*/
	CachingEngine ArticleStore
}

// "unparsed" config that's served from and to the user
type ConfigFile struct {
	BlogName             string
	ArticlesPerPage      string
	ListenOn             string
	ArticleStore         string
	ArticleStoreHost     string
	ArticleStoreDB       string
	ArticleStoreUser     string
	ArticleStorePassword string
	CachingEngine        string
}

// parses ConfigFile from user into Config for the app
// it's assumed that ConfigFile is verified and correct
func (cfg *ConfigFile) parseFile() *Config {
	parsedCfg := &Config{
		BlogName: cfg.BlogName,
		ListenOn: cfg.ListenOn,
		//ArticlesPerPage:	  0,
		//ArticleStore:       nil,
		ArticleStoreHost:     cfg.ArticleStoreHost,
		ArticleStoreDB:       cfg.ArticleStoreDB,
		ArticleStoreUser:     cfg.ArticleStoreUser,
		ArticleStorePassword: cfg.ArticleStorePassword,
		//CachingEngine:      nil,
	}

	if cfg.ArticleStore == "mock" {
		store := MockStore{}
		parsedCfg.ArticleStore = &store
	}

	if cfg.ArticleStore == "postgres" {
		store := PostgresStore{}
		parsedCfg.ArticleStore = &store
	}

	// convert ArticlesPerPage to int
	preconvert, _ := strconv.ParseInt(cfg.ArticlesPerPage, 10, 64)
	parsedCfg.ArticlesPerPage = uint64(preconvert)

	// TODO deal with CachingEngine
	return parsedCfg
}

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
