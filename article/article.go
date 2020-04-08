package article

import (
	"html/template"
)

type Article struct {
	// Date of release, used for sorting articles on run page
	Timestamp uint64

	// Unique identifier used in URL
	ID string

	// title of article
	Name string

	// type template.HTML allows unescaped html
	Content template.HTML
}
