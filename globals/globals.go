package globals

import (
	"github.com/david-sorm/montesquieu/config"
)

type BlogInformation struct {
	Name string
}

var BlogInfo BlogInformation

var Cfg *config.Config
