package templates

import (
	"fmt"
	"github.com/david-sorm/montesquieu/globals"
	"github.com/radovskyb/watcher"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var Store *template.Template

// Hot Swap Templates should be run as a singleton, so that's the reason for this variable
var liveTemplatesRunning = false

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

// Hot Swap Templates are used for faster development of templates, because it auto-reloads templates when changes are
// detected instead of having to manually restart Montesquieu
func InitHotSwapTemplates() {
	// run as singleton
	if liveTemplatesRunning {
		return
	}
	fmt.Println("Setting up Hot Swap Templates...")

	// voodoo magic begins
	w := watcher.New()

	// send only one event at a time to the channel
	w.SetMaxEvents(1)

	// only watch files with ".gohtml" suffix
	regex := regexp.MustCompile("^.+\\.gohtml$")
	w.AddFilterHook(watcher.RegexFilterHook(regex, false))

	// set receivers for receiving events
	go func() {
		for {
			select {
			case _ = <-w.Event:
				t := time.Now().Format("15:04:05")
				fmt.Println("[Hot Swap Templates] Change detected @", t)
				Load()
			case err := <-w.Error:
				fmt.Println("An error has happened while watching files for Hot Swap Templates:", err.Error())
			case <-w.Closed:
				return
			}
		}
	}()

	// watch the html/ folder
	if err := w.Add("html/"); err != nil {
		fmt.Println("Failed to set up a watch for Hot Swap Templates:", err.Error())
	}

	// start watching asynchronously
	go func() {
		if err := w.Start(time.Millisecond * 500); err != nil {
			fmt.Println("Failed to set up a watch for Hot Swap Templates:", err.Error())
		}
	}()
	liveTemplatesRunning = true
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

	// set up watcher for Hot Swap Templates
	if globals.Cfg.HotSwapTemplates {
		InitHotSwapTemplates()
	}
}
