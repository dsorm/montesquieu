package mock

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/users"
	"strconv"
	"sync"
)

// Store is a mock implementation of the Store interface
type Store struct {
	cfg store.StoreConfig

	m sync.Mutex

	// stores articles sorted from most recent[0] to oldest[...]
	articlesByTimestamp []article.Article

	// stores articles indexed by their IDs
	articlesByID map[string]article.Article

	users  []users.User
	admins []users.User
}

func (ms *Store) ListUsers(from uint64, to uint64) []users.User {
	ms.m.Lock()
	defer ms.m.Unlock()
	return ms.users[from:to]
}

func (ms *Store) GetUserID(user users.User) (users.User, bool) {
	ms.m.Lock()
	defer ms.m.Unlock()
	// find the user by login
	for _, v := range ms.users {
		if v.Login == user.Login {
			return v, true
		}
	}
	return users.User{}, false
}

func (ms *Store) GetUser(user users.User) users.User {
	ms.m.Lock()
	defer ms.m.Unlock()
	// find the user by ID
	for _, v := range ms.users {
		if v.ID == user.ID {
			return v
		}
	}
	return user
}

func (ms *Store) ListAuthors(from uint64, to uint64) []users.Author {
	ms.m.Lock()
	defer ms.m.Unlock()

	authors := make([]users.Author, 0, 0)
	for _, v := range ms.users[from:to] {
		authors = append(authors, users.Author{
			User:     v,
			AuthorID: v.ID,
		})
	}
	return authors
}

func (ms *Store) ListAdmins(from uint64, to uint64) []users.Admin {
	ms.m.Lock()
	defer ms.m.Unlock()

	admins := make([]users.Admin, 0, 0)
	for _, v := range ms.admins {
		admins = append(admins, users.Admin(v))
	}
	return admins[from:to]
}

func (ms *Store) Info() store.StoreInfo {
	return store.StoreInfo{
		Name:      "mock",
		Developer: "david-sorm",
	}
}

// too lazy to implement, and not needed
func (ms *Store) NewArticle(a article.Article) {
	return
}

func (ms *Store) EditArticle(a article.Article) {
	return
}

func (ms *Store) RemoveArticle(a article.Article) {
	return
}

func (ms *Store) AddUser(user users.User) {
	ms.m.Lock()
	ms.users = append(ms.users, user)
	ms.m.Unlock()
}

func (ms *Store) EditUser(user users.User) {
	ms.m.Lock()
	// find the user by ID
	for k, v := range ms.users {
		if v.ID == user.ID {
			ms.users[k] = user
		}
	}
	ms.m.Unlock()
}

func (ms *Store) RemoveUser(user users.User) {
	ms.m.Lock()
	// find the user by ID
	for k, v := range ms.users {
		if v.ID == user.ID {
			// copy the last user into the position of deleted user
			ms.users[k] = ms.users[len(ms.users)-1]

			// delete the last user
			ms.users = ms.users[:len(ms.users)-1]
		}
	}
	ms.m.Unlock()
}

// everyone is a an author since i'm way too lazy to implement this
// also user id == author id
func (ms *Store) GetAuthor(u users.User) users.Author {
	ms.m.Lock()
	defer ms.m.Unlock()

	u2 := ms.GetUser(u)
	return users.Author{
		User:     u2,
		AuthorID: u2.ID,
	}
}

func (ms *Store) AddAuthor(author users.Author) {
	return
}

func (ms *Store) LinkAuthor(author users.Author, user users.User) {
	return
}

func (ms *Store) RemoveAuthor(author users.Author) {
	return
}

func (ms *Store) PromoteToAdmin(admin users.Admin) {
	ms.m.Lock()
	// find the user by ID
	for _, v := range ms.users {
		if v.ID == admin.ID {
			ms.admins = append(ms.admins, v)
		}
	}
	ms.m.Unlock()
}

func (ms *Store) DemoteFromAdmin(admin users.Admin) {
	ms.m.Lock()
	// find the admin by ID
	for k, v := range ms.users {
		if v.ID == admin.ID {
			// copy the last admin into the position of deleted admin
			ms.admins[k] = ms.admins[len(ms.admins)-1]

			// delete the last admin
			ms.admins = ms.admins[:len(ms.admins)-1]
		}
	}
	ms.m.Unlock()
}

func (ms *Store) LoadArticlesSortedByLatest(from uint64, to uint64) []article.Article {
	/*
		// return articles starting from
		starti := ms.cfg.ArticlesPerIndexPage * page
		// and ending with these...
		endi := starti + ms.cfg.ArticlesPerIndexPage
		if endi > ms.GetArticleNumber() {
			endi = ms.GetArticleNumber()
		}
	*/

	return ms.articlesByTimestamp[from:to]
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

func (ms *Store) Init(_ func(), cfg store.StoreConfig) error {
	// doesn't implement notify at all, since MockStore cannot change contents at runtime

	// copy cfg
	ms.cfg = cfg

	// prepare the struct
	ms.articlesByTimestamp = make([]article.Article, 0, 0)
	ms.articlesByID = make(map[string]article.Article)
	ms.users = make([]users.User, 0, 0)
	ms.users = make([]users.User, 0, 0)

	// lets fill articles with some mock articles
	ms.articlesByTimestamp = append(ms.articlesByTimestamp, article.Article{
		Timestamp: 1585828351,
		ID:        "welcome",
		Name:      "Welcome to your brand new Goblog installation!",
		Content:   "Thank you for choosing Goblog! You should consider <b>changing the config.json</b>, since now Goblog only displays mock content, and you won't be able to make articles until you use a real Store.",
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
