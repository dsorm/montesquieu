package mock

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/article/store"
	"strconv"
)

// Mock implementation of ArticleStore
type Store struct {
	cfg store.ArticleStoreConfig

	// stores articles sorted from most recent[0] to oldest[...]
	articlesByTimestamp []article.Article

	// stores articles indexed by their IDs
	articlesByID map[string]article.Article
}

func (ms *Store) LoadArticlesForIndex(page uint64) []article.Article {
	// return articles starting from
	starti := ms.cfg.ArticlesPerIndexPage * page
	// and ending with these...
	endi := starti + ms.cfg.ArticlesPerIndexPage
	if endi > ms.GetArticleNumber() {
		endi = ms.GetArticleNumber()
	}
	// and return a slice made from the selection
	return ms.articlesByTimestamp[starti:endi]
}

func (ms *Store) GetArticleByID(ID string) (article.Article, bool) {
	// val stores the value, if there's none, it simply stores a zeroed Article
	// exists stores boolean value meaning the existence of an article with the ID
	val, exists := ms.articlesByID[ID]
	return val, exists
}

func (ms *Store) GetArticleNumber() uint64 {
	num := len(ms.articlesByTimestamp)
	return uint64(num)
}

func (ms *Store) Init(_ func(), cfg store.ArticleStoreConfig) error {
	// doesn't implement notify at all, since MockStore cannot change contents at runtime

	// copy cfg
	ms.cfg = cfg

	// prepare the struct
	ms.articlesByTimestamp = make([]article.Article, 0, 0)
	ms.articlesByID = make(map[string]article.Article)

	// lets fill articles with some mock articles
	ms.articlesByTimestamp = append(ms.articlesByTimestamp, article.Article{
		Timestamp: 1585828351,
		ID:        "welcome",
		Name:      "Welcome to your brand new Goblog installation!",
		Content:   "Thank you for choosing Goblog! You should consider <b>changing the config.json</b>, since now Goblog only displays mock content, and you won't be able to make articles until you change ArticleStore value to a real store.",
	})

	// lets generate another mock articles
	for i := 1; i < 11; i++ {
		ms.articlesByTimestamp = append(ms.articlesByTimestamp, article.Article{
			Timestamp: ms.articlesByTimestamp[i-1].Timestamp - 1,
			ID:        "article" + strconv.Itoa(i+1),
			Name:      "Article " + strconv.Itoa(i+1),
			Content:   "Lorem ipsum dolor sit amet",
		})
	}

	// make a copy that's sorted by ID
	for _, v := range ms.articlesByTimestamp {
		ms.articlesByID[v.ID] = v
	}

	// i don't think there's even a remote possibility of error in this function
	return nil
}
