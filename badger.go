package main

import (
	badger "github.com/dgraph-io/badger"
	"golang.org/x/xerrors"
)

func getFromURL(db *badger.DB, url string, callback func(repo Repo, err error)) (err error) {
	err = db.View(func(txn *badger.Txn) (err error) {

		// Select from URL
		item, err := txn.Get([]byte(url))
		if err != nil {
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
			callback(repo, err)

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
