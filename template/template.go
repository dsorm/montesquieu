package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

var Store *template.Template

// These are the required default templates for proper startup.
// If any template is missing, the server will panic.
var requiredTemplates = []string{
	"article.gohtml",
	"index.gohtml",
	"adminPanel.gohtml",
	"adminPanelHeader.gohtml",
	"adminPanelFooter.gohtml",
	"adminPanelUsers.gohtml",
	"adminPanelAuthors.gohtml",
	"adminPanelAdmins.gohtml",
	"adminPanelConfiguration.gohtml",
}

func checkRequiredTemplates() {
	for _, v := range requiredTemplates {
		// check every required template by searching it
		if Store.Lookup(v) == nil {
			fmt.Println("Template", v,
				"is missing. Please check if it's in the root of /html, or if the permissions are correct. Halting...")
			panic(nil)
		}
	}
	fmt.Println("All required templates present!")
}

// Prepares all templates, not used when unit testing
func Load() {

	// create a list of .gohtml template files
	var templateFiles []string
	templateFiles = make([]string, 0, 10)
	dirContent, err := ioutil.ReadDir("html")
	if err != nil {
		fmt.Println("Can't read contents of /html: ", err.Error())
	}

	// select only files ending with .gohtml
	fmt.Println("These template files are being loaded:")
	for _, v := range dirContent {
		if strings.HasSuffix(v.Name(), ".gohtml") {
			// don't forget the folder to make it a valid path
			templateFiles = append(templateFiles, "html/"+v.Name())
			fmt.Printf("%v; ", v.Name())
		}
	}
	fmt.Println()

	// parse all the selected files
	fmt.Println("Parsing templates...")
	Store, err = template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Println("Error while parsing gohtml templates from /html:", err.Error())
	}

	// this is probably total bs and will never happen, but go won't stop yelling at me
	// for "not using" the variable
	if Store == nil {
		fmt.Println("Template store is nil, halting...")
		panic(nil)
	}

	// we don't like nil pointer exceptions...
	checkRequiredTemplates()
}
