package postgres

import (
	"context"
	"fmt"
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/users"
	"html/template"
	"strconv"
)

// Postgres implementation of Store
type Store struct {
	ArticlesPerIndexPage uint64
}

func (p *Store) Info() store.StoreInfo {
	panic("implement me")
}

func (p *Store) ListUsers(from uint64, to uint64) []users.User {
	panic("implement me")
}

func (p *Store) GetUserID(user users.User) (users.User, bool) {
	panic("implement me")
}

func (p *Store) GetUser(user users.User) users.User {
	panic("implement me")
}

func (p *Store) ListAuthors(from uint64, to uint64) []users.Author {
	panic("implement me")
}

func (p *Store) ListAdmins(from uint64, to uint64) []users.Admin {
	panic("implement me")
}

func (p *Store) LoadArticlesSortedByLatest(from uint64, to uint64) []article.Article {
	rows, err := pool.Query(context.Background(), stmtArticles, from, to)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return []article.Article{}
	}

	articles := make([]article.Article, 0, p.ArticlesPerIndexPage)
	var title string
	var articleId int32
	var authorId int32
	var htmlPreview string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &articleId, &authorId, &htmlPreview, &timestamp)
		articles = append(articles, article.Article{
			Name:      title,
			ID:        strconv.Itoa(int(articleId)),
			AuthorID:  strconv.Itoa(int(authorId)),
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlPreview),
		})
	}

	return articles
}

func (p *Store) NewArticle(a article.Article) {
	panic("implement me")
}

func (p *Store) EditArticle(a article.Article) {
	panic("implement me")
}

func (p *Store) RemoveArticle(a article.Article) {
	panic("implement me")
}

func (p *Store) AddUser(user users.User) {
	panic("implement me")
}

func (p *Store) EditUser(user users.User) {
	panic("implement me")
}

func (p *Store) RemoveUser(user users.User) {
	panic("implement me")
}

func (p *Store) GetAuthor(user users.User) users.Author {
	panic("implement me")
}

func (p *Store) AddAuthor(author users.Author) {
	panic("implement me")
}

func (p *Store) LinkAuthor(author users.Author, user users.User) {
	panic("implement me")
}

func (p *Store) RemoveAuthor(author users.Author) {
	panic("implement me")
}

func (p *Store) PromoteToAdmin(admin users.Admin) {
	panic("implement me")
}

func (p *Store) DemoteFromAdmin(admin users.Admin) {
	panic("implement me")
}

func (p *Store) GetArticleNumber() uint64 {
	rows, err := pool.Query(context.Background(), "select count(*) from "+prefix+".articles;")
	if err != nil {
		fmt.Println("An error has happened while getting the number of articles from Postgres:", err.Error())
	}

	count := uint64(5)
	rows.Next()
	rows.Scan(&count)

	return count
}

func (p *Store) Init(f func(), cfg store.StoreConfig) error {
	p.ArticlesPerIndexPage = cfg.ArticlesPerIndexPage
	err := dbInit(cfg.Host, cfg.Database, cfg.Username, cfg.Password)

	if err != nil {
		return err
	}
	return nil
}

func (p *Store) LoadArticlesForIndex(page uint64) []article.Article {
	// return articles starting from
	offset := p.ArticlesPerIndexPage * page
	limit := p.ArticlesPerIndexPage

	rows, err := pool.Query(context.Background(), stmtArticles, offset, limit)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return []article.Article{}
	}

	articles := make([]article.Article, 0, p.ArticlesPerIndexPage)
	var title string
	var articleId int32
	var authorId int32
	var htmlPreview string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &articleId, &authorId, &htmlPreview, &timestamp)
		articles = append(articles, article.Article{
			Name:      title,
			ID:        strconv.Itoa(int(articleId)),
			AuthorID:  strconv.Itoa(int(authorId)),
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlPreview),
		})
	}

	return articles

}

func (p *Store) GetArticleByID(ID string) (article.Article, bool) {
	idNum, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return article.Article{}, false
	}

	rows, err := pool.Query(context.Background(),
		"select title, author_id, html_content, timestamp from "+prefix+".articles where article_id = $1",
		idNum)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return article.Article{}, false
	}

	var title string
	var authorId int32
	var htmlContent string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &authorId, &htmlContent, &timestamp)
		return article.Article{
			Name:      title,
			ID:        string(ID),
			AuthorID:  strconv.Itoa(int(authorId)),
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlContent),
		}, true
	}

	return article.Article{}, false
}
