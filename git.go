// Operations and flow control on Git repositories managed by the application.

package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"golang.org/x/xerrors"
)

func addURL(db *badger.DB, link string) (repo Repo, err error) {
	repo, err = dbGet(db, link)
	if xerrors.Is(err, badger.ErrKeyNotFound) {
		return onNewRepo(db, link)
	} else if err != nil {
		return
	}
	return onOldRepo(db, repo)
}

// rmRepo completely removes any traces of a repository in the system.
func rmRepo(db *badger.DB, repo Repo) {

	if repo.UUID != "" {
		rm(repo.UUID)
	}

	if repo.URL != "" {
		ipfsKeyRmName(repo.URL)
		dbDelete(db, repo.URL)
	}

	if repo.IPFS != "" {
		ipfsClusterRm(repo.IPFS)
	}

	if repo.Key != "" {
		ipfsKeyRm(repo.Key)
	}

	if repo.IPNS != "" {
		ipfsKeyRm(repo.IPNS)
	}

}

func onAllRepos(db *badger.DB) {
	fmt.Println("Refreshing all repos...")
	fmt.Println()

	ch := dbList(db)

	var wg sync.WaitGroup
	for c := 1; c <= runtime.NumCPU(); c++ {
		wg.Add(1)
		go func(ch chan repoerr) {
			for re := range ch {
				if re.err != nil {
					fmt.Println("Couldn't refresh a repo.")
					fmt.Println(re.err.Error())
				}

				onOldRepo(db, re.repo)
			}
			wg.Done()
		}(ch)
	}

	wg.Wait()

	fmt.Println("All repos are refreshed.")
	fmt.Println()

	return
}

func onOldRepo(db *badger.DB, repo Repo) (Repo, error) {

	// Show old values
	fmt.Println(aurora.Bold("UUID :"), repo.UUID)
	fmt.Println(aurora.Bold("URL :"), aurora.Blue(repo.URL))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))
	fmt.Println(aurora.Bold("Key :"), repo.Key)
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))
	fmt.Println()

	// Pull
	_, err := gitPull(repo.UUID)
	if err != nil {
		return repo, err
	}

	// Remove old IPFS
	ipfsClusterRm(repo.IPFS)

	// Size
	size, err := dirSize(dirHome + dirGit + "/" + repo.UUID)
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		return repo, err
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// Add new IPFS
	out, err := ipfsClusterAdd(repo.URL, rmin, rmax, repo.UUID)
	if err != nil {
		return repo, err
	}

	repo.IPFS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))

	// IPNS
	out, err = ipfsNamePublish(repo.Key, repo.IPFS)
	if err != nil {
		return repo, err
	}

	repo.IPNS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))

	err = dbSet(db, repo)
	if err != nil {
		fmt.Println("Couldn't save the updated repo.")
		fmt.Println(err.Error())
		return repo, err
	}

	fmt.Println("Saved", aurora.Blue(repo.URL).String()+".")
	return repo, err
}

func onNewRepo(db *badger.DB, link string) (repo Repo, err error) {

	// URL
	repo.URL = strings.TrimSpace(link)
	fmt.Println(aurora.Bold("URL :"), aurora.Blue(repo.URL))
	if repo.URL == "" {
		rmRepo(db, repo)
		return repo, ErrNoURL
	}

	// UUID
	buuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Couldn't generate a new UUID.")
		fmt.Println(err.Error())
		rmRepo(db, repo)
		return
	}

	repo.UUID = strings.TrimSpace(buuid.String())
	fmt.Println(aurora.Bold("UUID :"), repo.UUID)

	// Clone
	_, err = gitClone(repo.URL, repo.UUID)
	if err != nil {
		rmRepo(db, repo)
		return
	}

	// Size
	size, err := dirSize(dirHome + dirGit + "/" + repo.UUID)
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		rmRepo(db, repo)
		return
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// IPFS-Cluster
	out, err := ipfsClusterAdd(repo.URL, rmin, rmax, repo.UUID)
	if err != nil {
		rmRepo(db, repo)
		return
	}

	repo.IPFS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))

	// Key
	out, err = ipfsKeyGen(repo.URL)
	if err != nil {
		rmRepo(db, repo)
		return
	}

	repo.Key = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("Key :"), repo.Key)

	// IPNS
	out, err = ipfsNamePublish(repo.Key, repo.IPFS)
	if err != nil {
		rmRepo(db, repo)
		return
	}

	repo.IPNS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))

	// Save
	err = dbSet(db, repo)
	if err != nil {
		fmt.Println("Couldn't save the newly created repo.")
		fmt.Println(err.Error())
		rmRepo(db, repo)
		return
	}

	fmt.Println("Saved", aurora.Blue(repo.URL).String()+".")
	return repo, err
}
