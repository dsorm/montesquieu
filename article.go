package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Article struct {
	Name string

	// type template.HTML allows unescaped html
	Content template.HTML
}

type ArticleView struct {
	Info    BlogInfo
	Article *Article
	RootURL string
}

func handleArticle(rw http.ResponseWriter, req *http.Request) {
	// split the uri (example: /articles/1 )
	split := strings.Split(req.RequestURI, "/")

	// make sure there are 3 splits
	if len(split) != 3 {
		handle404(rw, req)
		return
	}

	// convert to number
	i, err := strconv.Atoi(split[2])
	if err != nil {
		handle404(rw, req)
		return
	}

	// make sure article with the number exists
	if articles[i] == nil {
		handle404(rw, req)
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
