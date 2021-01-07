package handlers

import (
	"fmt"
	articlePkg "github.com/david-sorm/montesquieu/article"
	"github.com/david-sorm/montesquieu/globals"
	templates "github.com/david-sorm/montesquieu/template"
	"net/http"
	"strconv"
	"strings"
)

type ArticleView struct {
	BlogName string
	Article  articlePkg.Article
	RootURL  string
}

func HandleArticle(rw http.ResponseWriter, req *http.Request) {
	// split the uri (example: /articles/1 )
	split := strings.Split(req.RequestURI, "/")

	// make sure there are 3 splits
	if len(split) != 3 {
		Handle404(rw, req)
		return
	}

	convertInt, err := strconv.Atoi(split[2])
	if err != nil {
		fmt.Println("An error has happened while converting article id from request:", err)
	}

	// make sure article with the ID exists
	article, exists := globals.Cfg.Store.GetArticleByID(uint64(convertInt))
	if !exists {
		Handle404(rw, req)
		return
	}

	// respond
	articleView := ArticleView{
		BlogName: globals.Cfg.BlogName,
		Article:  article,
		RootURL:  "//" + req.Host + "/",
	}
	if err := templates.Store.Lookup("article.gohtml").Execute(rw, articleView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}
