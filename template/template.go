package templates

import (
	"fmt"
	"html/template"
)

var Store []*template.Template

var templatePaths = []string{
	"html/index.gohtml",
	"html/article.gohtml",
}

/*
   used for more idiomatic and organized code
   (eg. templates.Store[template.Index] instead of templates.Store[0]),
   indexes have to be in sync with templatePaths
*/
const (
	Index = iota
	Article
)

// Prepares all templates, not used when unit testing
func Load() {

	// initialize the Store slice
	Store = make([]*template.Template, 0, 10)

	// load and parse templates into Store
	temp, err := template.ParseFiles(templatePaths...)
	if err != nil {
		fmt.Println("An error has happened while parsing templates:", err.Error())
		return
	}
	Store = temp.Templates()

}
