package logic

import (
	storePkg "github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/store/postgres"
)

func ParseStore(str string) storePkg.Store {
	if str == "postgres" {
		store := postgres.Store{}
		return &store
	}

	return nil
}
