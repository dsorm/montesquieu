package handlers

import (
	"fmt"
	"github.com/david-sorm/goblog/globals"
	"net/http"
)
import templates "github.com/david-sorm/goblog/template"

func HandleAdminPanel(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanel.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelArticles(rw http.ResponseWriter, req *http.Request) {
	data := globals.Cfg.Store.LoadArticlesSortedByLatest(0, 100)
	if err := templates.Store.Lookup("adminPanelArticles.gohtml").Execute(rw, data); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelUsers(rw http.ResponseWriter, req *http.Request) {
	data := globals.Cfg.Store.ListUsers(0, 100)
	if err := templates.Store.Lookup("adminPanelUsers.gohtml").Execute(rw, data); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelAuthors(rw http.ResponseWriter, req *http.Request) {
	data := globals.Cfg.Store.ListAuthors(0, 100)
	if err := templates.Store.Lookup("adminPanelAuthors.gohtml").Execute(rw, data); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelAdmins(rw http.ResponseWriter, req *http.Request) {
	data := globals.Cfg.Store.ListAdmins(0, 100)
	if err := templates.Store.Lookup("adminPanelAdmins.gohtml").Execute(rw, data); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelConfiguration(rw http.ResponseWriter, req *http.Request) {
	cachingEngine := globals.Cfg.CachingStore != nil
	data := struct {
		BlogName        string
		ArticlesPerPage uint64
		Store           string
		StoreHost       string
		StoreDB         string
		StoreUser       string
		CachingStore    bool
		ListenOn        string
	}{
		globals.Cfg.BlogName,
		globals.Cfg.ArticlesPerPage,
		"something",
		globals.Cfg.StoreHost,
		globals.Cfg.StoreDB,
		globals.Cfg.StoreUser,
		cachingEngine,
		globals.Cfg.ListenOn,
	}
	if err := templates.Store.Lookup("adminPanelConfiguration.gohtml").Execute(rw, data); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}
