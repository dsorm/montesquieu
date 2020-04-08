package config

import (
	cfgLogic "github.com/david-sorm/goblog/article/logic"
	"github.com/david-sorm/goblog/article/store"
	"strconv"
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
	ArticleStore store.ArticleStore

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
	CachingEngine store.ArticleStore
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

	// convert ArticlesPerPage to int
	preconvert, _ := strconv.ParseInt(cfg.ArticlesPerPage, 10, 64)
	parsedCfg.ArticlesPerPage = uint64(preconvert)

	parsedCfg.ArticleStore = cfgLogic.ParseArticleStore(cfg.ArticleStore)

	// TODO deal with CachingEngine
	return parsedCfg
}
