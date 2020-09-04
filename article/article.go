package article

import (
	"html/template"
)

type Article struct {
	// title of article
	Title string

	// Unique identifier used internally
	ID uint64

	// Unique identifier of the Author
	AuthorID uint64

	// Date of release, used for sorting articles on run page
	// TODO change the type to time.Time
	Timestamp uint64

	// type template.HTML allows unescaped html
	Content template.HTML
}
