package main

import (
	"fmt"

	badger "github.com/dgraph-io/badger"
	"golang.org/x/xerrors"
)

func getFromURL(db *badger.DB, link string, ch chan Repo, oldRepo func(*badger.DB, Repo) (Repo, error), newRepo func(*badger.DB, string) (Repo, error)) (err error) {
	err = db.View(func(txn *badger.Txn) (err error) {

		// Select from URL
		item, err := txn.Get([]byte(link))
		if err == badger.ErrKeyNotFound {
			repo, err := newRepo(db, link)
			if err != nil {
				return xerrors.Errorf("Couldn't add a new repo : %w", err)
			}
			ch <- repo
		} else if err != nil {
			return xerrors.Errorf("Couldn't select an URL : %w", err)
		}

		// Get value
		err = item.Value(func(bytes []byte) (err error) {

			// Decode
			repo, err := decodeRepo(bytes)
			if err != nil {
				return xerrors.Errorf("Couldn't decode a repo : %w", err)
			}

			// Execute the callback only if there's no error.
			repo, err = oldRepo(db, repo)
			if err != nil {
				return xerrors.Errorf("Couldn't refresh an old repo : %w", err)
			}
			ch <- repo

			// Unreachable code
			return nil
		})

		if err != nil {
			return xerrors.Errorf("Couldn't get an URL's value : %w", err)
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("Couldn't create a view from an URL : %w", err)
	}

	return nil
}

func (repo Repo) save(db *badger.DB) (err error) {
	err = db.Update(func(txn *badger.Txn) error {

		// Encode
		bytes, err := repo.encode()
		if err != nil {
			return xerrors.Errorf("Couldn't encode a repo : %w", err)
		}

		// Save
		err = txn.Set([]byte(repo.URL), bytes)
		if err != nil {
			return xerrors.Errorf("Couldn't save a repo : %w", err)
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("Couldn't create an update transaction : %w", err)
	}

	return nil
}

func getRepos(db *badger.DB, ch chan Repo) error {

	// Create a view
	err := db.View(func(txn *badger.Txn) error {

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		defer close(ch)

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			err := item.Value(func(v []byte) error {

				repo, err := decodeRepo(v)
				if err != nil {
					return xerrors.Errorf("Couldn't decode a repo : %w", err)
				}

				ch <- repo

				return nil
			})

			if err != nil {
				fmt.Println("Couldn't get a value.")
				fmt.Println(err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return xerrors.Errorf("Couldn't create a view : %w", err)
	}

	return nil
}
