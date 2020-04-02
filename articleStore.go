package main

import "strconv"

// import "github.com/lib/pq"

type ArticleStore interface {

	/*
	 ArticleStore should be prepared for work upon returning nil from this function
	 Non-nil response means an error has occurred; error will be shown in console
	 If the first argument is nil, it means the store shouldn't monitor changes
	 If a function is passed, it should be called every time a change is detected
	 In the second parameter, *Config is passed, which includes parsed config.json,
	 mainly for database credentials if needed
	*/
	Init(f func(), cfg *Config) error

	/*
	 Should return articles for this index page
	 Index pages start at 0, the number of articles per page is defined in
	 Config.ArticlesPerPage
	*/
	LoadArticlesForIndex(page uint64) []Article

	/*
	 Should return the article by the unique ID, obviously the ID in Article will be
	 ignored, so it can be set to nil
	 If an article with the ID can't be found, the second return parameter should return
	 false, else if an article was found, return true
	*/
	GetArticleByID(ID string) (Article, bool)

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

// Postgres implementation
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

// Mock implementation
type MockStore struct {

	// save the cfg pointer for later use
	cfg *Config

	// stores articles sorted from most recent[0] to oldest[...]
	articlesByTimestamp []Article

	// stores articles indexed by their IDs
	articlesByID map[string]Article
}

func (ms *MockStore) LoadArticlesForIndex(page uint64) []Article {
	// return articles starting from
	starti := ms.cfg.ArticlesPerPage * page
	// and ending with these...
	endi := starti + ms.cfg.ArticlesPerPage

	if endi > ms.GetArticleNumber() {
		endi = ms.GetArticleNumber()
	}
	// and return a slice made from the selection
	return ms.articlesByTimestamp[starti:endi]
}

func (ms *MockStore) GetArticleByID(ID string) (Article, bool) {
	// val stores the value, if there's none, it simply stores a zeroed Article
	// exists stores boolean value meaning the existence of an article with the ID
	val, exists := ms.articlesByID[ID]
	return val, exists
}

func (ms *MockStore) GetArticleNumber() uint64 {
	num := len(ms.articlesByTimestamp)
	return uint64(num)
}

func (ms *MockStore) Init(f func(), cfg *Config) error {
	// doesn't implement notify at all, since MockStore cannot change contents at runtime

	// prepare the struct
	ms.cfg = cfg
	ms.articlesByTimestamp = make([]Article, 0, 0)
	ms.articlesByID = make(map[string]Article)

	// lets fill articles with some mock articles
	ms.articlesByTimestamp = append(ms.articlesByTimestamp, Article{
		timestamp: 1585828351,
		ID:        "welcome",
		Name:      "Welcome to your brand new Goblog installation!",
		Content:   "Thank you for choosing Goblog! You should consider <b>changing the config.json</b>, since now Goblog only displays mock content, and you won't be able to make articles until you change ArticleStore value to a real store.",
	})

	// lets generate another mock articles
	for i := 1; i < 11; i++ {
		ms.articlesByTimestamp = append(ms.articlesByTimestamp, Article{
			timestamp: ms.articlesByTimestamp[i-1].timestamp - 1,
			ID:        "article" + strconv.Itoa(i+1),
			Name:      "Article " + strconv.Itoa(i+1),
			Content:   "Lorem ipsum dolor sit amet",
		})
	}

	// make a copy that's sorted by ID
	for _, v := range ms.articlesByTimestamp {
		ms.articlesByID[v.ID] = v
	}

	// i don't think there's even a remote possibilty of error in this function
	return nil
}
