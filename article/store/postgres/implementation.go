package postgres

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/article/store"
)

// Postgres implementation of ArticleStore
type PostgresStore struct{}

func (p *PostgresStore) GetArticleNumber() uint64 {
	panic("implement me")
}

func (p *PostgresStore) Init(f func(), cfg store.ArticleStoreConfig) error {
	panic("implement me")
}

func (p *PostgresStore) LoadArticlesForIndex(page uint64) []article.Article {
	panic("implement me")
}

func (p *PostgresStore) GetArticleByID(ID string) (article.Article, bool) {
	panic("implement me")
}
