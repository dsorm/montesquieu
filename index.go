package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type IndexView struct {
	Info *BlogInfo

	// a list of articles which should be displayed on the page
	Articles []*Article

	// the last page, if there's one, used for the last button
	// essentially Page - 1
	LastPage int

	// the current page
	Page int

	// the next page, if there's one, used for the next button
	// essentially Page + 1
	NextPage int

	// the biggest page
	MaxPage int
}

// executes
func handleIndex(rw http.ResponseWriter, req *http.Request) {
	uri := req.URL.RequestURI()
	indexView := IndexView{
		Info: &blogInfo,

		// first page if we don't specify below
		Page: 1,

		// mock, too
		MaxPage: 10,
	}
	// get rid of the '/' at the beginning
	uri = strings.TrimPrefix(uri, "/")

	// if there's something more than just '', try to figure out whether we've got this page or not
	if len(uri) > 0 {
		uriNum, err := strconv.Atoi(uri)
		if err == nil {
			/*
			 if this is a valid number, larger than 0 and smaller or equal the biggest page,
			 set the page number to this number
			*/
			if uriNum > 0 && uriNum <= indexView.MaxPage {
				indexView.Page = uriNum
			}
		} else {
			// if this is BS, send a 404
			handle404(rw, req)
			return
		}
	}

	// for the buttons
	indexView.LastPage = indexView.Page - 1
	indexView.NextPage = indexView.Page + 1

	// insert the actual articles into page
	indexView.Articles = prepareArticles(articles, indexView.Page)

	// TODO delete this debug code
	bytes, _ := json.MarshalIndent(indexView, "", "\t")
	fmt.Println(string(bytes))

	// execute template
	if err := templates[templateIndex].Execute(rw, indexView); err != nil {
		fmt.Println("Error while parsing template:", err.Error())
	}
}

// TODO don't just display the same articles on every page
func prepareArticles(a []*Article, page int) []*Article {
	return a
}
