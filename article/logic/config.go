package logic

import (
	"github.com/david-sorm/goblog/article/store"
	"github.com/david-sorm/goblog/article/store/mock"
	"github.com/david-sorm/goblog/article/store/postgres"
)

func ParseArticleStore(str string) store.ArticleStore {
	if str == "mock" {
		store := mock.MockStore{}
		return &store
	}

	if str == "postgres" {
		store := postgres.PostgresStore{}
		return &store
	}

	return nil
}
