package store

import "github.com/david-sorm/goblog/article"

// import "github.com/lib/pq"

type ArticleStoreConfig struct {
	Host                 string
	Database             string
	Username             string
	Password             string
	ArticlesPerIndexPage uint64
}

type ArticleStore interface {

	/*
	 ArticleStore should be prepared for work upon returning nil from this function
	 Non-nil response means an error has occurred; error will be shown in console
	 If the first argument is nil, it means the store shouldn't monitor changes
	 If a function is passed, it should be called every time a change is detected
	 The second parameter is a config that contains relevant parsed data from config
	 file
	*/
	Init(f func(), cfg ArticleStoreConfig) error

	/*
	 Should return articles for this index page
	 Index pages start at 0, the number of articles per page is defined in
	 Config.ArticlesPerPage
	*/
	LoadArticlesForIndex(page uint64) []article.Article

	/*
	 Should return the article by the unique ID, obviously the ID in Article will be
	 ignored, so it can be set to nil
	 If an article with the ID can't be found, the second return parameter should return
	 false, else if an article was found, return true
	*/
	GetArticleByID(ID string) (article.Article, bool)

	/*
	 Should return the total number of articles, used for determining how many
	 index pages we have
	*/
	GetArticleNumber() uint64

	// TODO functionality for adding articles

}

type ArticleCachingEngine interface {
	/*
		 ArticleCachingEngine should mostly have the same functionality as ArticleStore,
		 only with the difference of Use(ArticleStore) and different internal logic
		(returning data from i's own cache instead of doing queries every time there's an
		article request, etc.)
	*/
	ArticleStore
	/*
	 This tells the ArticleCachingEngine which ArticleStore it should use
	 ArticleCachingEngine should take care of init etc., we just pass a new instance
	 of ArticleStore which the user has selected in config
	 We call this before Init(), so the cache engine should first Init() it's own
	 ArticleStore before initializing itself
	 If any errors happen while initialising ArticleStore, the error should be passed
	 from Init() called afterward
	*/
	Use(ArticleStore)
}
