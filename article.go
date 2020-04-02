package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Article struct {
	// Date of release, used for sorting articles on main page
	timestamp uint64

	// Unique identifier used in URL
	ID string

	// title of article
	Name string

	// type template.HTML allows unescaped html
	Content template.HTML
}

type ArticleView struct {
	Info    BlogInfo
	Article Article
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

	// make sure article with the ID exists
	article, exists := cfg.ArticleStore.GetArticleByID(split[2])
	if !exists {
		handle404(rw, req)
		return
	}

	// respond
	articleView := ArticleView{
		Info:    blogInfo,
		Article: article,
		RootURL: "//" + req.Host + "/",
	}
	if err := templates[templateArticle].Execute(rw, articleView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}
