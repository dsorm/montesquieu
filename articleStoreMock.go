package main

import "strconv"

// Mock implementation of ArticleStore
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
