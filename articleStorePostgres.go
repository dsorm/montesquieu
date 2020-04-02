package main

// Postgres implementation of ArticleStore
type PostgresStore struct{}

func (p *PostgresStore) GetArticleNumber() uint64 {
	panic("implement me")
}

func (p *PostgresStore) Init(f func(), cfg *Config) error {
	panic("implement me")
}

func (p *PostgresStore) LoadArticlesForIndex(page uint64) []Article {
	panic("implement me")
}

func (p *PostgresStore) GetArticleByID(ID string) (Article, bool) {
	panic("implement me")
}
