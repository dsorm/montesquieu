package handlers

import (
	"fmt"
	"net/http"
)
import templates "github.com/david-sorm/goblog/template"

func HandleAdminPanel(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanel.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelArticles(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanelArticles.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelUsers(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanelUsers.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelAuthors(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanelAuthors.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelAdmins(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanelAdmins.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func HandleAdminPanelConfiguration(rw http.ResponseWriter, req *http.Request) {
	if err := templates.Store.Lookup("adminPanelConfiguration.gohtml").Execute(rw, nil); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}
