package main

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"golang.org/x/xerrors"
)

func receiveURL(db *badger.DB, url string) {
	err := getFromURL(db, url, onOldRepo)
	if xerrors.Is(err, badger.ErrKeyNotFound) {
		onNewRepo(db, url)
	} else if err != nil {
		fmt.Println(err.Error())
	}
}

func onAllRepos(db *badger.DB) {
	fmt.Println("Refreshing all repos...")
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

	return
}

func onOldRepo(db *badger.DB, repo Repo) {

	// Show old values
	fmt.Println(aurora.Bold("UUID :"), repo.UUID)
	fmt.Println(aurora.Bold("URL :"), aurora.Blue(repo.URL))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))
	fmt.Println(aurora.Bold("Key :"), repo.Key)
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))
	fmt.Println()

	// Pull
	out, err := exec.Command("git", "-C", repo.UUID, "pull").Output()
	if err != nil {
		fmt.Println("Couldn't clone the repository.")
		fmt.Println(aurora.Bold("Command :"), "git", "-C", repo.UUID, "pull")

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}
	// fmt.Println(string(out))

	// Remove old IPFS
	out, err = exec.Command("ipfs-cluster-ctl", "pin", "rm", repo.IPFS).Output()
	if err != nil {
		fmt.Println("Couldn't add the repository to IPFS.")
		fmt.Println(aurora.Bold("Command :"), "ipfs-cluster-ctl", "pin", "rm", aurora.Cyan(repo.IPFS))

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	// Size
	size, err := dirSize(repo.UUID)
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		return
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// Add new IPFS
	out, err = exec.Command("ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", repo.URL, "--replication-min", rmin, "--replication-max", rmax, repo.UUID).Output()
	if err != nil {
		fmt.Println("Couldn't add the repository to IPFS.")
		fmt.Println(aurora.Bold("Command :"), "ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", aurora.Blue(repo.URL), "--replication-min", aurora.Bold(rmin), "--replication-max", aurora.Bold(rmax), repo.UUID)

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	repo.IPFS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(repo.IPFS))

	// IPNS
	out, err = exec.Command("ipfs", "name", "publish", "--key", repo.Key, "--quieter", "/ipfs/"+repo.IPFS).Output()
	if err != nil {
		fmt.Println("Couldn't generate a new key.")
		fmt.Println(aurora.Bold("Command :"), "ipfs", "name", "publish", "--key", repo.Key, "--quieter", aurora.Cyan("/ipfs/"+repo.IPFS))

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	repo.IPNS = strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(repo.IPNS))

	err = repo.save(db)
	if err != nil {
		fmt.Println("Couldn't save the newly created repo.")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Saved", aurora.Blue(repo.URL).String()+".")
}

func onNewRepo(db *badger.DB, link string) {

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
	out, err := exec.Command("git", "clone", link, uuid).Output()
	if err != nil {
		fmt.Println("Couldn't clone the repository.")
		fmt.Println(aurora.Bold("Command :"), "git clone", aurora.Blue(link), uuid)

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}
	// fmt.Println(string(out))

	// Size
	size, err := dirSize(uuid)
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		return
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// IPFS-Cluster
	out, err = exec.Command("ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", link, "--replication-min", rmin, "--replication-max", rmax, uuid).Output()
	if err != nil {
		fmt.Println("Couldn't add the repository to IPFS.")
		fmt.Println(aurora.Bold("Command :"), "ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", aurora.Blue(link), "--replication-min", aurora.Bold(rmin), "--replication-max", aurora.Bold(rmax), uuid)

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	ipfs := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(ipfs))

	// Key
	escaped := url.PathEscape(link)
	out, err = exec.Command("ipfs", "key", "gen", "--type", "ed25519", escaped).Output()
	if err != nil {
		fmt.Println("Couldn't generate a new key.")
		fmt.Println(aurora.Bold("Command :"), "ipfs", "key", "gen", "--type", "ed25519", aurora.Blue(escaped))

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	key := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("Key :"), key)

	// IPNS
	out, err = exec.Command("ipfs", "name", "publish", "--key", key, "--quieter", "/ipfs/"+ipfs).Output()
	if err != nil {
		fmt.Println("Couldn't generate a new key.")
		fmt.Println(aurora.Bold("Command :"), "ipfs", "name", "publish", "--key", key, "--quieter", aurora.Cyan("/ipfs/"+ipfs))

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	ipns := strings.TrimSpace(string(out))
	fmt.Println(aurora.Bold("IPNS :"), aurora.Cyan(ipns))

	repo := Repo{
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
}
