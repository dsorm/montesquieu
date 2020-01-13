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

type Article struct {
	Name string
	// type template.HTML allows unescaped html
	Content template.HTML
}

type IndexView struct {
	Info     *BlogInfo
	Articles *[]*Article
}

type ArticleView struct {
	Info    *BlogInfo
	Article *Article
	RootURL string
}

var articles []*Article

var templates []*template.Template

var blogInfo *BlogInfo

var indexView *IndexView

const (
	templateIndex = iota
	templateArticle
)

func handleIndex(rw http.ResponseWriter, req *http.Request) {
	if err := templates[templateIndex].Execute(rw, indexView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func handleArticle(rw http.ResponseWriter, req *http.Request) {
	// split the uri (example: /articles/1 )
	split := strings.Split(req.RequestURI, "/")

	// make sure there are 3 splits
	if len(split) != 3 {
		rw.WriteHeader(404)
		return
	}

	// convert to number
	i, err := strconv.Atoi(split[2])
	if err != nil {
		rw.WriteHeader(404)
		return
	}

	// make sure article with the number exists
	if articles[i] == nil {
		rw.WriteHeader(404)
		return
	}

	// respond
	articleView := ArticleView{
		Info:    blogInfo,
		Article: articles[i],
		RootURL: "//" + req.Host + "/",
	}
	if err := templates[templateArticle].Execute(rw, articleView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

func handleCss(rw http.ResponseWriter, req *http.Request) {

	// get the file
	file, err := os.Open("html" + req.RequestURI)
	defer file.Close()
	if err != nil {
		fmt.Println("Error while opening css: ", err.Error())
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

func main() {

	// init
	articles = make([]*Article, 0, 10)
	templates = make([]*template.Template, 0, 10)

	// prepare data for Views
	blogInfo = &BlogInfo{
		Name: "Loremum Ipsium",
	}

	indexView = &IndexView{
		Info:     blogInfo,
		Articles: &articles,
	}

	// some garbage articles for testing
	articles = append(articles, &Article{
		Name:    "Lorem ipsum",
		Content: "... dolor sir amet :)",
	})
	articles = append(articles, &Article{
		Name:    "Go Templating Engine",
		Content: `<b>Go’s html/template package</b> provides a <i>rich templating language</i> for HTML templates. It is mostly used in web applications to display data in a structured way in a client’s browser. One great benefit of Go’s templating language is the automatic escaping of data. There is no need to worry about XSS attacks as Go parses the HTML template and escapes all inputs before displaying it to the browser.`,
	})

	// parse and load all templates
	templates = template.Must(template.ParseFiles("html/index.gohtml", "html/article.gohtml")).Templates()

	// register all controllers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/article/", handleArticle)
	mux.HandleFunc("/css/", handleCss)
	mux.HandleFunc("/fonts/", handleFonts)

	// start the web server
	if err := http.ListenAndServe(":80", mux); err != nil {
		fmt.Println("Error while starting web server:", err.Error())
	}
}
