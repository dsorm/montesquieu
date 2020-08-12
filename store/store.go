package store

import (
	"github.com/david-sorm/goblog/article"
	"github.com/david-sorm/goblog/users"
)

// import "github.com/lib/pq"

// StoreConfig contains data passed to a Store implementation
type StoreConfig struct {
	Host                 string
	Database             string
	Username             string
	Password             string
	ArticlesPerIndexPage uint64
}

// StoreInfo should contain info about the store implementation, so goblog can
// properly register it
type StoreInfo struct {
	// should be a json-friendly and short name
	Name string

	// doesn't have to be json-friendly
	Developer string
}

/*
 Store is an interface meant to be implemented by a package which should do the
 actual work of managing and keeping the data
*/
type Store interface {

	// General

	// Info() has to return general info about the Store implementation itself
	Info() StoreInfo

	/*
	 Store should be prepared for work upon returning nil from this function
	 Non-nil response means an error has occurred; error will be shown in console
	 If the first argument is nil, it means the store shouldn't monitor changes
	 If a function is passed, it should be called every time a change is detected
	 The second parameter is a config that contains relevant parsed data from config
	 file
	*/
	Init(f func(), cfg StoreConfig) error

	// Articles

	/*
	 Should return a slice of articles sorted from latest.
	 'from' means how many articles from latest should be cut off from the start
	 (0 = don't cut off anything).
	 'to' means how many articles minus latest should be cut off to the end.
	 Example: LoadArticlesSortedByLatest(2,7) should load 5 articles, starting
	 with the 3rd most recent and article and ending with the 7th
	*/
	LoadArticlesSortedByLatest(from uint64, to uint64) []article.Article

	/*
	 Should return the article by the unique ID, obviously the ID in Article will
	 be ignored, so it can be set to nil.
	 If an article with the ID can't be found, the second return parameter should
	 return false, else if an article was found, return true
	*/
	GetArticleByID(ID string) (article.Article, bool)

	/*
	 Should return the total number of articles, used for determining how many
	 index pages we have
	*/
	GetArticleNumber() uint64

	// When called, the Store should make a new article in its database and save it.
	NewArticle(article.Article)

	// Store should look up the article by its ID and make corresponding changes
	EditArticle(article.Article)

	// The article should be looked up by its ID and deleted
	RemoveArticle(article.Article)

	// Users

	// Lists Users, sorts by ID
	ListUsers(from uint64, to uint64) []users.User

	// Gets user ID from login name
	// Returns whether a matching user was find using bool
	// True = Found, False = Not
	GetUserID(users.User) (users.User, bool)

	// Searches for a user by ID
	GetUser(users.User) users.User

	// Makes a new user
	AddUser(users.User)

	// Edits a user according to his ID
	EditUser(users.User)

	// Removes a user according to his ID
	RemoveUser(users.User)

	// Authors

	// Lists Authors, sorts by ID
	ListAuthors(from uint64, to uint64) []users.Author

	// Returns nil if the User is not an Author
	GetAuthor(users.User) users.Author

	// Adds an Author
	AddAuthor(users.Author)

	// Links a user to an Author
	// If User is nil, any link of an Author to a User should be deleted
	LinkAuthor(users.Author, users.User)

	// Removes an author
	RemoveAuthor(users.Author)

	// Admins

	// Lists Admins, sorts by ID
	ListAdmins(from uint64, to uint64) []users.Admin

	// Promotes a User to be an Admin
	PromoteToAdmin(users.Admin)

	// Demotes an Admin to a User only
	DemoteFromAdmin(users.Admin)
}

/*
 CachingStore should mostly have the same functionality as Store,
 only with the difference of Use(Store) and different internal logic
 (returning data from i's own cache instead of doing queries every time there's an
 article request, etc.)
*/
type CachingStore interface {
	Store

	/*
	 We use this method to pass the Store which should be used by this CachingStore
	 CachingStore should call Init() on the Store before it starts initialising itself.
	 Any errors that happened during the Init() of the Store should be returned
	 through CachingStore's Init()
	*/
	Use(Store)
}
