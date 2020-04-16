package run

import (
	"fmt"
	"github.com/david-sorm/goblog/article/store"
	"github.com/david-sorm/goblog/config"
	"github.com/david-sorm/goblog/globals"
	"github.com/david-sorm/goblog/handlers"
	templates "github.com/david-sorm/goblog/template"
	"net/http"
)

func Main() {
	fmt.Println("Goblog starting...")
	// get the config
	fmt.Println("Loading config...")
	var err error
	globals.Cfg, err = config.NewConfig()
	if err != nil {
		fmt.Printf("While verifying the config, some errors in config.json were found. Please fix them before running Goblog:\n%s", err.Error())
		return
	}

	// init
	fmt.Println("Initializing ArticleStore...")
	asCfg := store.ArticleStoreConfig{
		Host:                 globals.Cfg.ArticleStoreHost,
		Database:             globals.Cfg.ArticleStoreDB,
		Username:             globals.Cfg.ArticleStoreUser,
		Password:             globals.Cfg.ArticleStorePassword,
		ArticlesPerIndexPage: globals.Cfg.ArticlesPerPage,
	}
	err = globals.Cfg.ArticleStore.Init(func() {}, asCfg)
	if err != nil {
		fmt.Println("An error has happened while initializing ArticleStore: ", err.Error())
	}

	// prepare data for Views
	globals.BlogInfo = globals.BlogInformation{
		Name: globals.Cfg.BlogName,
	}

	// parse and load all templates
	fmt.Println("Parsing templates...")
	templates.Load()

	// register all controllers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/article/", handlers.HandleArticle)
	mux.HandleFunc("/css/", handlers.HandleCss)
	mux.HandleFunc("/fonts/", handlers.HandleFonts)

	fmt.Println("All ok!")
	fmt.Println("Server starting at port", globals.Cfg.ListenOn)

	// start the web server
	if err := http.ListenAndServe(globals.Cfg.ListenOn, mux); err != nil {
		fmt.Println("Error while starting web server:", err.Error())
	}
}
