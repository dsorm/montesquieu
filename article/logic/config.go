package logic

import (
	storePkg "github.com/david-sorm/goblog/article/store"
	"github.com/david-sorm/goblog/article/store/mock"
	"github.com/david-sorm/goblog/article/store/postgres"
)

func ParseArticleStore(str string) storePkg.ArticleStore {
	if str == "mock" {
		store := mock.Store{}
		return &store
	}

	if str == "postgres" {
		store := postgres.Store{}
		return &store
	}

	return nil
}
