// High level operations on the database.

package main

import (
	"runtime"

	badger "github.com/dgraph-io/badger"
)

func dbGet(db *badger.DB, link string) (repo Repo, err error) {
	ch := make(chan Repo, 1)
	err = badgerGet(db, link, ch)
	if err != nil {
		return
	}

	repo = <-ch
	return
}

func dbList(db *badger.DB) (ch chan repoerr) {
	ch = make(chan repoerr, runtime.NumCPU())
	go badgerList(db, ch)
	return ch
}

func dbSet(db *badger.DB, repo Repo) error {
	return badgerSet(db, repo)
}

func dbDelete(db *badger.DB, link string) error {
	return badgerDelete(db, link)
}
