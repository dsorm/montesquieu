package config

import (
	cfgLogic "github.com/david-sorm/montesquieu/article/logic"
	"github.com/david-sorm/montesquieu/store"
	"strconv"
	"strings"
)

// TODO get rid of old-fashioned config parsing from json and use envs instead

// "parsed" config that's served to the app
type Config struct {
	// Blog's name
	BlogName string

	// The port on which the server should listen on
	ListenOn string

	/*
	 How many articles should be displayed on one index page
	 While technically there's a limit of 2^64 articles on page
	 if you really hate your Store (and internet connections,
	 browsers, etc.), we recommend sticking to something like 5
	 articles
	*/
	ArticlesPerPage uint64

	/*
	 Type of database
	 Currently only `postgres` is supported
	*/
	Store store.Store

	/*
	 Login info for Store driver, if needed
	 Postgres: requires all except StorePort filled out
	*/
	StoreHost     string
	StoreDB       string
	StoreUser     string
	StorePassword string
	StorePort     string

	/*
	 Type of caching engine used between the app and the store
	 Currently only 'internal' or 'off' is supported
	*/
	CachingStore store.Store

	/*
	 For template-development purposes only, reloads templates without restarting
	 Recommended setting for production use: off
	*/
	HotSwapTemplates bool
}

// "unparsed" config that's served from and to the user
type file struct {
	BlogName         string
	ArticlesPerPage  string
	ListenOn         string
	Store            string
	StoreHost        string
	StoreDB          string
	StoreUser        string
	StorePassword    string
	StorePort        string
	CachingStore     string
	HotSwapTemplates string
}

// parses ConfigFile from user into Config for the app
// it's assumed that ConfigFile is verified and correct
func (cfg *file) parseFile() *Config {
	parsedCfg := &Config{
		BlogName: cfg.BlogName,
		ListenOn: cfg.ListenOn,
		//ArticlesPerPage:	  0,
		//Store:       nil,
		StoreHost:     cfg.StoreHost,
		StoreDB:       cfg.StoreDB,
		StoreUser:     cfg.StoreUser,
		StorePassword: cfg.StorePassword,
		StorePort:     cfg.StorePort,
		//CachingStore:      nil,
		HotSwapTemplates: strings.ToLower(cfg.HotSwapTemplates) == "yes",
	}

	// convert ArticlesPerPage to int
	preconvert, _ := strconv.ParseInt(cfg.ArticlesPerPage, 10, 64)
	parsedCfg.ArticlesPerPage = uint64(preconvert)

	parsedCfg.Store = cfgLogic.ParseStore(cfg.Store)

	// TODO deal with CachingStore
	return parsedCfg
}
