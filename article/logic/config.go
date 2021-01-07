package logic

import (
	storePkg "github.com/david-sorm/montesquieu/store"
	"github.com/david-sorm/montesquieu/store/postgres"
)

func ParseStore(str string) storePkg.Store {
	if str == "mock" {
		store := mock.Store{}
		return &store
	}

	if str == "postgres" {
		store := postgres.Store{}
		return &store
	}

	return nil
}
