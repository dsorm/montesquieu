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
