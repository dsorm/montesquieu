package postgres

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/article/store"
)

// Postgres implementation of ArticleStore
type Store struct{}

func (p *Store) GetArticleNumber() uint64 {
	panic("implement me")
}

func (p *Store) Init(f func(), cfg store.ArticleStoreConfig) error {
	panic("implement me")
}

func (p *Store) LoadArticlesForIndex(page uint64) []article.Article {
	panic("implement me")
}

func (p *Store) GetArticleByID(ID string) (article.Article, bool) {
	panic("implement me")
}
