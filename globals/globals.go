package globals

import (
	"github.com/david-sorm/goblog/config"
	"html/template"
)

type BlogInformation struct {
	Name string
}

var Templates []*template.Template

var BlogInfo BlogInformation

var Cfg *config.Config

const (
	TemplateIndex = iota
	TemplateArticle
)
