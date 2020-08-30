package postgres

import (
	"context"
	"fmt"
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/users"
	"html/template"
)

// Postgres implementation of Store
type Store struct {
	ArticlesPerIndexPage uint64
}

// The comments are here to please code quality analysis tools.
// They are absolutely useless.

// IsAdmin implements Store's IsAdmin function
func (p *Store) IsAdmin(id uint64) bool {
	rows, err := pool.Query(context.Background(), stmtIsAdmin, id)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while checking if the user is an admin:", err)
	}

	var count uint8
	for rows.Next() {
		rows.Scan(&count)
	}
	return count == 1
}

// Info implements Store's Info function
func (p *Store) Info() store.StoreInfo {
	return store.StoreInfo{
		Name:      "postgres",
		Developer: "david-sorm",
	}
}

// ListUsers implements Store's ListUsers function
func (p *Store) ListUsers(from uint64, to uint64) []users.User {
	rows, err := pool.Query(context.Background(), stmtListUsers, from, to)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while listing users:", err)
	}

	us := make([]users.User, 0, 0)
	for rows.Next() {
		u := users.User{}
		rows.Scan(&u.ID, &u.DisplayName, &u.Login)
		us = append(us, u)
	}
	return us
}

// GetUserID implements Store's GetUserID function
func (p *Store) GetUserID(login string) (uint64, bool) {
	rows, err := pool.Query(context.Background(), stmtGetUserID, login)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while getting user's id:", err)
	}

	if rows.Next() {
		// if there are any rows
		var id uint64
		rows.Scan(&id)
		return id, true
	} else {
		// if there aren't any rows
		return 0, false
	}
}

// GetUser implements Store's GetUser function
func (p *Store) GetUser(id uint64) users.User {
	rows, err := pool.Query(context.Background(), stmtGetUser, id)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while getting a user:", err)
	}

	u := users.User{}
	for rows.Next() {
		rows.Scan(&u.ID, &u.DisplayName, &u.Login, &u.Password)
	}
	return u
}

// ListAuthors implements Store's ListAuthors function
func (p *Store) ListAuthors(from uint64, to uint64) []users.Author {
	rows, err := pool.Query(context.Background(), stmtListAuthors, from, to)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while listing authors:", err)
	}

	authors := make([]users.Author, 0, 0)
	for rows.Next() {
		a := users.Author{}
		rows.Scan(&a.ID, &a.DisplayName, &a.Login, &a.AuthorID, &a.AuthorName)
		authors = append(authors, a)
	}
	return authors
}

// ListAdmins implements Store's ListAdmins function
func (p *Store) ListAdmins(from uint64, to uint64) []users.User {
	rows, err := pool.Query(context.Background(), stmtListAdmins, from, to)

	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while listing admins:", err)
	}

	admins := make([]users.User, 0, 0)
	for rows.Next() {
		u := users.User{}
		rows.Scan(&u.ID, &u.DisplayName, &u.Login)
		admins = append(admins, u)
	}
	return admins
}

// LoadArticlesSortedByLatest implements Store's LoadArticlesSortedByLatest function
func (p *Store) LoadArticlesSortedByLatest(from uint64, to uint64) []article.Article {
	rows, err := pool.Query(context.Background(), stmtLoadArticlesSortedByNewest, from, to)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return []article.Article{}
	}

	articles := make([]article.Article, 0, p.ArticlesPerIndexPage)
	var title string
	var articleId int32
	var authorId uint64
	var htmlPreview string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &articleId, &authorId, &htmlPreview, &timestamp)
		articles = append(articles, article.Article{
			Title:     title,
			ID:        uint64(articleId),
			AuthorID:  authorId,
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlPreview),
		})
	}

	return articles
}

// AddArticle implements Store's AddArticle function
func (p *Store) AddArticle(title string, authorId uint64, timestamp uint64,
	content template.HTML) {
	doExec(stmtNewArticle, "adding an article", title, authorId, content,
		content, timestamp)
}

// EditArticle implements Store's EditArticle function
func (p *Store) EditArticle(a article.Article) {
	doExec(stmtEditArticle, "editing an article", a.Title, a.AuthorID,
		string(a.Content), string(a.Content), a.Timestamp, a.ID)
}

// RemoveArticle implements Store's RemoveArticle function
func (p *Store) RemoveArticle(id uint64) {
	doExec(stmtRemoveArticle, "removing an article", id)
}

// AddUser implements Store's AddUser function
func (p *Store) AddUser(displayName string, login string, password string) {
	doExec(stmtAddUser, "adding a new user", displayName, login, password)
}

// EditUser implements Store's EditUser function
func (p *Store) EditUser(user users.User) {
	doExec(stmtEditUser, "editing a user", user.DisplayName, user.Login,
		user.Password, user.ID)
}

// RemoveUser implements Store's RemoveUser function
func (p *Store) RemoveUser(id uint64) {
	doExec(stmtRemoveUser, "removing a user", id)
}

// GetAuthor implements Store's GetAuthor function
func (p *Store) GetAuthor(userId uint64) users.Author {
	rows, err := pool.Query(context.Background(), stmtGetAuthor, userId)
	if err != nil {
		fmt.Println("[Postgres Store] An error has happened while getting an author:", err)
	}

	author := users.Author{}
	for rows.Next() {
		rows.Scan(&author.AuthorID, &author.AuthorName)
	}

	return author
}

// AddAuthor implements Store's AddAuthor function
func (p *Store) AddAuthor(userId uint64, authorName string) {
	doExec(stmtAddAuthor, "adding an author", userId, authorName)
}

// LinkAuthor implements Store's LinkAuthor function
func (p *Store) LinkAuthor(authorId uint64, userId uint64) {
	doExec(stmtLinkAuthor, "linking a user to an author", userId, authorId)
}

// RemoveAuthor implements Store's RemoveAuthor function
func (p *Store) RemoveAuthor(authorId uint64) {
	doExec(stmtRemoveAuthor, "removing an author", authorId)
}

// PromoteToAdmin implements Store's PromoteToAdmin function
func (p *Store) PromoteToAdmin(userId uint64) {
	doExec(stmtPromoteToAdmin, "promoting a user to an admin", userId)
}

// DemoteFromAdmin implements Store's DemoteFromAdmin function
func (p *Store) DemoteFromAdmin(userId uint64) {
	doExec(stmtDemoteFromAdmin, "demoting a user from an admin", userId)
}

// GetArticleNumber implements Store's GetArticleNumber function
func (p *Store) GetArticleNumber() uint64 {
	rows, err := pool.Query(context.Background(), stmtArticleNumber)
	if err != nil {
		fmt.Println("An error has happened while getting the number of articles from Postgres:", err.Error())
	}

	count := uint64(0)
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println("An error has happened while getting the article number: ", err.Error())
		}
	}
	return count

}

// Init implements Store's Init function
func (p *Store) Init(f func(), cfg store.StoreConfig) error {
	p.ArticlesPerIndexPage = cfg.ArticlesPerIndexPage
	err := dbInit(cfg.Host, cfg.Database, cfg.Username, cfg.Password)

	if err != nil {
		return err
	}
	return nil
}

// LoadArticlesForIndex implements Store's LoadArticlesForIndex function
func (p *Store) LoadArticlesForIndex(page uint64) []article.Article {
	// return articles starting from
	offset := p.ArticlesPerIndexPage * page
	limit := p.ArticlesPerIndexPage

	rows, err := pool.Query(context.Background(), stmtLoadArticlesSortedByNewest, offset, limit)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return []article.Article{}
	}

	articles := make([]article.Article, 0, p.ArticlesPerIndexPage)
	var title string
	var articleId uint64
	var authorId uint64
	var htmlPreview string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &articleId, &authorId, &htmlPreview, &timestamp)
		articles = append(articles, article.Article{
			Title:     title,
			ID:        articleId,
			AuthorID:  authorId,
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlPreview),
		})
	}

	return articles

}

// GetArticleByID implements Store's GetArticleByID function
func (p *Store) GetArticleByID(id uint64) (article.Article, bool) {
	rows, err := pool.Query(context.Background(), stmtGetArticleByID, id)
	if err != nil {
		fmt.Println("An error has happened while loading articles for index:", err.Error())
		return article.Article{}, false
	}

	var title string
	var authorId uint64
	var htmlContent string
	var timestamp int64

	for rows.Next() {
		rows.Scan(&title, &authorId, &htmlContent, &timestamp)
		return article.Article{
			Title:     title,
			ID:        id,
			AuthorID:  authorId,
			Timestamp: uint64(timestamp),
			Content:   template.HTML(htmlContent),
		}, true
	}

	return article.Article{}, false
}

// doExec is a helper function that helps prevent code duplication when doing
// simple pgx exec queries
func doExec(stmt string, activity string, arguments ...interface{}) {
	ct, err := pool.Exec(context.Background(), stmt, arguments...)
	if err != nil || ct.RowsAffected() == 0 {
		fmt.Println("[Postgres Store] An error has happened while "+activity+": ", err)
	}
}
