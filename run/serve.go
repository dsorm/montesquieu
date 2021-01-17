package run

import (
	"fmt"
	"github.com/david-sorm/montesquieu/config"
	"github.com/david-sorm/montesquieu/globals"
	"github.com/david-sorm/montesquieu/handlers"
	"github.com/david-sorm/montesquieu/store"
	templates "github.com/david-sorm/montesquieu/template"
	"net/http"
)

func Main() {
	fmt.Println("Montesquieu starting...")
	// get the config
	fmt.Println("Loading config...")
	var err error
	globals.Cfg, err = config.NewConfig()
	if err != nil {
		fmt.Printf("While verifying the config, some errors in config.json were found. Please fix them before running Montesquieu:\n%s", err.Error())
		return
	}

	// init
	fmt.Println("Initializing Store...")
	asCfg := store.StoreConfig{
		Host:                 globals.Cfg.StoreHost,
		Database:             globals.Cfg.StoreDB,
		Username:             globals.Cfg.StoreUser,
		Password:             globals.Cfg.StorePassword,
		Port:                 globals.Cfg.StorePort,
		ArticlesPerIndexPage: globals.Cfg.ArticlesPerPage,
	}
	err = globals.Cfg.Store.Init(func() {}, asCfg)
	if err != nil {
		fmt.Println("An error has happened while initializing Store: ", err.Error())
	}

	// prepare data for Views
	globals.BlogInfo = globals.BlogInformation{
		Name: globals.Cfg.BlogName,
	}

	// parse and load all templates in html/
	templates.Load()

	// make FileServer controllers for handling fully static content
	handleCss := http.FileServer(http.Dir("html/css"))
	handleFonts := http.FileServer(http.Dir("html/fonts"))
	handleJs := http.FileServer(http.Dir("html/js"))

	// register all controllers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/article/", handlers.HandleArticle)
	mux.HandleFunc("/admin/panel", handlers.HandleAdminPanel)
	mux.HandleFunc("/admin/panel/articles", handlers.HandleAdminPanelArticles)
	mux.HandleFunc("/admin/panel/users", handlers.HandleAdminPanelUsers)
	mux.HandleFunc("/admin/panel/authors", handlers.HandleAdminPanelAuthors)
	mux.HandleFunc("/admin/panel/admins", handlers.HandleAdminPanelAdmins)
	mux.HandleFunc("/admin/panel/configuration", handlers.HandleAdminPanelConfiguration)
	mux.HandleFunc("/login", handlers.HandleLogin)

	// http.StripPrefix is needed for FileServer handlers so the paths work correctly
	mux.Handle("/css/", http.StripPrefix("/css/", handleCss))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", handleFonts))
	mux.Handle("/js/", http.StripPrefix("/js/", handleJs))

	fmt.Println("Server starting at port", globals.Cfg.ListenOn)

	// start the web server
	if err := http.ListenAndServe(globals.Cfg.ListenOn, mux); err != nil {
		fmt.Println("Error while starting web server:", err.Error())
	}
}
