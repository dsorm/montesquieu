package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BlogInfo struct {
	Name string
}

var articles []Article

var templates []*template.Template

var blogInfo BlogInfo

var cfg *Config

const (
	templateIndex = iota
	templateArticle
)

func handleCss(rw http.ResponseWriter, req *http.Request) {

	// get the file
	file, err := os.Open("html" + req.RequestURI)
	defer file.Close()
	if err != nil {
		// TODO solve missing fonts
		// fmt.Println("Error while opening css: ", err.Error())
		rw.WriteHeader(404)
		return
	}

	// https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/

	// get file size
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// send the headers
	rw.Header().Set("Content-Type", "text/css")
	rw.Header().Set("Content-Length", fileSize)

	// send the file
	file.Seek(0, 0)
	io.Copy(rw, file)
}

func handleFonts(rw http.ResponseWriter, req *http.Request) {

	// get the file
	file, err := os.Open("html" + req.RequestURI)
	defer file.Close()
	if err != nil {
		fmt.Println("Error while opening css: ", err.Error())
		rw.WriteHeader(404)
		return
	}

	// send the content-type header
	// css
	if strings.HasSuffix(req.RequestURI, ".css") {
		rw.Header().Set("Content-Type", "text/css")
		// ttf
	} else if strings.HasSuffix(req.RequestURI, ".ttf") {
		rw.Header().Set("Content-Type", "font/ttf")
		// woff
	} else if strings.HasSuffix(req.RequestURI, ".woff") {
		rw.Header().Set("Content-Type", "font/woff")
		// woff2
	} else if strings.HasSuffix(req.RequestURI, ".woff2") {
		rw.Header().Set("Content-Type", "font/woff2")
	}

	// https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/

	// get file size
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// send the content-length header
	rw.Header().Set("Content-Length", fileSize)

	// send the file
	file.Seek(0, 0)
	io.Copy(rw, file)
}

// Generic 404 page, for use by other handlers in cases of invalid URL
// TODO make nicer 404 page
func handle404(rw http.ResponseWriter, _ *http.Request) {
	// better than nothing, i guess
	rw.WriteHeader(404)
	fmt.Fprintf(rw, "Error 404: Not Found\n")
}

func main() {
	fmt.Println("Goblog starting...")
	// get the config
	fmt.Println("Loading config...")
	var err error
	cfg, err = NewConfig()
	if err != nil {
		fmt.Printf("While verifying the config, some errors in config.json were found. Please fix them before running Goblog:\n%s", err.Error())
		return
	}

	// init
	fmt.Println("Initializing ArticleStore...")
	err = cfg.ArticleStore.Init(nil, cfg)
	if err != nil {
		fmt.Println("An error has happeneed while initializing ArticleStore: ", err.Error())
	}
	templates = make([]*template.Template, 0, 10)

	// prepare data for Views
	blogInfo = BlogInfo{
		Name: cfg.BlogName,
	}

	// parse and load all templates
	fmt.Println("Parsing templates...")
	templates = template.Must(template.ParseFiles("html/index.gohtml", "html/article.gohtml")).Templates()

	// register all controllers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/article/", handleArticle)
	mux.HandleFunc("/css/", handleCss)
	mux.HandleFunc("/fonts/", handleFonts)

	fmt.Println("All ok!")
	fmt.Println("Server starting at port", cfg.ListenOn)

	// start the web server
	if err := http.ListenAndServe(cfg.ListenOn, mux); err != nil {
		fmt.Println("Error while starting web server:", err.Error())
	}
}
