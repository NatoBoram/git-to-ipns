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

func receiveURL(db *badger.DB, url string) (repo Repo, err error) {
	ch := make(chan Repo, 1)
	err = getFromURL(db, url, ch, onOldRepo, onNewRepo)
	if err != nil {
		close(ch)
		return repo, xerrors.Errorf("Couldn't get a repo from its URL : %w", err)
	}
	repo = <-ch
	return
}

func onAllRepos(db *badger.DB) {
	fmt.Println("Refreshing all repos...")
	fmt.Println()

	ch := make(chan Repo, runtime.NumCPU())

	go func() {
		err := getRepos(db, ch)
		if err != nil {
			fmt.Println("Couldn't get all known repos.")
			fmt.Println(err.Error())
			close(ch)
		}
	}()

	var wg sync.WaitGroup
	for c := 1; c <= runtime.NumCPU(); c++ {
		wg.Add(1)
		go func(ch chan Repo) {
			for repo := range ch {
				onOldRepo(db, repo)
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
	out, err := gitPull(repo.UUID)

	// Remove old IPFS
	out, err = ipfsClusterRm(repo.IPFS)

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
	out, err = ipfsClusterAdd(repo.URL, rmin, rmax, repo.UUID)
	repo.IPFS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))

	// IPNS
	out, err = ipfsNamePublish(repo.Key, repo.IPFS)
	repo.IPNS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))

	err = repo.save(db)
	if err != nil {
		fmt.Println("Couldn't save the newly created repo.")
		fmt.Println(err.Error())
		return repo, err
	}

	fmt.Println("Saved", aurora.Blue(repo.URL).String()+".")
	return repo, err
}

func onNewRepo(db *badger.DB, link string) (repo Repo, err error) {

	// URL
	link = strings.TrimSpace(link)
	fmt.Println(aurora.Bold("URL :"), aurora.Blue(link))

	// UUID
	buuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Couldn't generate a new UUID.")
		fmt.Println(err.Error())
		return
	}

	uuid := strings.TrimSpace(buuid.String())
	fmt.Println(aurora.Bold("UUID :"), uuid)

	// Clone
	out, err := gitClone(link, uuid)
	if err != nil {
		return
	}

	// Size
	size, err := dirSize(dirHome + dirGit + "/" + repo.UUID)
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		return
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// IPFS-Cluster
	out, err = ipfsClusterAdd(link, rmin, rmax, uuid)
	if err != nil {
		return
	}

	ipfs := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(ipfs))

	// Key
	out, err = ipfsKeyGen(link)
	if err != nil {
		return
	}

	key := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("Key :"), key)

	// IPNS
	out, err = ipfsNamePublish(key, ipfs)
	if err != nil {
		return
	}

	ipns := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(ipns))

	repo = Repo{
		UUID: uuid,
		URL:  link,
		IPFS: ipfs,
		Key:  key,
		IPNS: ipns,
	}

	err = repo.save(db)
	if err != nil {
		fmt.Println("Couldn't save the newly created repo.")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Saved", aurora.Blue(link).String()+".")
	return repo, err
}
