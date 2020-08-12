package handlers

import (
	"fmt"
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/globals"
	templates "github.com/david-sorm/goblog/template"
	"net/http"
	"strconv"
	"strings"
)

type IndexView struct {
	BlogName string

	// a list of articles which should be displayed on the page
	Articles []article.Article

	// the last page, if there's one, used for the last button
	// essentially Page - 1
	LastPage uint64

	// the current page
	Page uint64

	// the next page, if there's one, used for the next button
	// essentially Page + 1
	NextPage uint64

	// the biggest page
	MaxPage uint64
}

// makes sure that we have the correct number of pages
func countMaxPage(NumOfArticles uint64, ArticlesPerPage uint64) uint64 {
	tempFloat := float64(NumOfArticles)/float64(ArticlesPerPage) - 1.0
	tempInt := (NumOfArticles / ArticlesPerPage) - 1

	if float64(tempInt) == tempFloat {
		// if the result is round, it's pretty easy
		return tempInt
	} else {
		// if it's not, lets add another page for 'leftover' articles
		return tempInt + 1
	}
}

// executes
func HandleIndex(rw http.ResponseWriter, req *http.Request) {
	uri := req.URL.RequestURI()
	indexView := IndexView{
		BlogName: globals.Cfg.BlogName,

		// first page if we don't specify below
		Page: 0,

		// -1 since pages are zero-indexed
		MaxPage: countMaxPage(globals.Cfg.Store.GetArticleNumber(), globals.Cfg.ArticlesPerPage),
	}
	// get rid of the '/' at the beginning
	uri = strings.TrimPrefix(uri, "/")

	// if there's something more than just '', try to figure out whether we've got this page or not
	if len(uri) > 0 {
		uriNum, err := strconv.ParseUint(uri, 10, 64)
		if err == nil {

			// redirect page /0 to /, since it looks ugly
			if uriNum == 0 {
				http.Redirect(rw, req, "/", 301)
			}

			/*
			 if this is a valid number, larger or equal 0 and smaller or equal the biggest page,
			 set the page number to this number
			*/
			if uriNum >= 0 && uriNum <= indexView.MaxPage {
				indexView.Page = uriNum
			}
		} else {
			// if this is BS, send a 404
			Handle404(rw, req)
			return
		}
	}

	// for the buttons
	indexView.LastPage = indexView.Page - 1
	indexView.NextPage = indexView.Page + 1

	// calculate the articles
	// articles starting from
	starti := globals.Cfg.ArticlesPerPage * indexView.Page
	// and ending with these...
	endi := starti + globals.Cfg.ArticlesPerPage

	articleNum := globals.Cfg.Store.GetArticleNumber()
	if endi > articleNum {
		endi = articleNum
	}

	// insert the actual articles into page
	indexView.Articles = globals.Cfg.Store.LoadArticlesSortedByLatest(starti, endi)

	// execute template

	if err := templates.Store.Lookup("index.gohtml").Execute(rw, indexView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}
