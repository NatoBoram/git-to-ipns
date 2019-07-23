package main

import (
	"fmt"
	"os/exec"
	"strings"

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

func onOldRepo(repo Repo, err error) {

	fmt.Println(
		"onOldRepo",
		repo.UUID,
		repo.URL,
		repo.IPFS,
		repo.Key,
		repo.IPNS,
	)
}

func onNewRepo(db *badger.DB, url string) {

	// URL
	url = strings.TrimSpace(url)
	fmt.Println(aurora.Bold("URL :"), aurora.Blue(url))

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
	out, err := exec.Command("git", "clone", url, uuid).Output()
	if err != nil {
		fmt.Println("Couldn't clone the repository.")
		fmt.Println(aurora.Bold("Command :"), "git clone", aurora.Blue(url), uuid)

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}
	fmt.Println(string(out))

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
	out, err = exec.Command("ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", url, "--replication-min", rmin, "--replication-max", rmax, uuid).Output()
	if err != nil {
		fmt.Println("Couldn't add the repository to IPFS.")
		fmt.Println(aurora.Bold("Command :"), "ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", aurora.Blue(url), "--replication-min", aurora.Bold(rmin), "--replication-max", aurora.Bold(rmax), uuid)

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
	out, err = exec.Command("ipfs", "key", "gen", "--type", "ed25519", url).Output()
	if err != nil {
		fmt.Println("Couldn't generate a new key.")
		fmt.Println(aurora.Bold("Command :"), "ipfs", "key", "gen", "--type", "ed25519", aurora.Blue(url))

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
		URL:  url,
		IPFS: ipfs,
		Key:  key,
		IPNS: ipns,
	}

	err = repo.save(db)
	if err != nil {
		fmt.Println("Couldn't save the newly created repo.")
		fmt.Println(err.Error())
	}
}
