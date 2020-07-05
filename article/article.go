package article

import (
	"html/template"
)

type Article struct {
	// title of article
	Name string

	// Unique identifier used in URL
	ID string

	// Unique identifier of the Author
	AuthorID string

	// Date of release, used for sorting articles on run page
	Timestamp uint64

	// type template.HTML allows unescaped html
	Content template.HTML
}
