// Low-level operations on the database.

package main

import (
	"fmt"

	badger "github.com/dgraph-io/badger"
)

func badgerSet(db *badger.DB, repo Repo) error {
	return db.Update(func(txn *badger.Txn) error {

		// Encode
		bytes, err := repo.encode()
		if err != nil {
			return err
		}

		// Save
		err = txn.Set([]byte(repo.URL), bytes)
		return err
	})
}

// badgerList populates a channel with every single keys inside the Badger database.
func badgerList(db *badger.DB, ch chan repoerr) error {
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				repo, err := decodeRepo(v)
				ch <- repoerr{
					repo: repo,
					err:  err,
				}

				return err
			})

			// Continue even if there's an error
			if err != nil {
				fmt.Println("Couldn't get the value of an item.")
				fmt.Println(err.Error())
			}
		}

		return nil
	})

	close(ch)
	return err
}

// Get a repository from its URL.
//
// Can return `badger.ErrKeyNotFound`.
func badgerGet(db *badger.DB, link string, ch chan Repo) error {
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(link))
		if err != nil {
			return err
		}

		return item.Value(func(bytes []byte) error {
			repo, err := decodeRepo(bytes)
			ch <- repo
			return err
		})
	})

	close(ch)
	return err
}

func badgerDelete(db *badger.DB, link string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(link))
	})
}
