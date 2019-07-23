package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	badger "github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
)

func main() {

	// License
	fmt.Println("")
	fmt.Println(aurora.Bold("Gi :"), "Git to IPFS")
	fmt.Println("Copyright © 2019 Nato Boram")
	fmt.Println("This program is free software : you can redistribute it and/or modify it under the terms of the " + aurora.Underline("GNU General Public License").String() + " as published by the " + aurora.Underline("Free Software Foundation").String() + ", either version 3 of the License, or (at your option) any later version. This program is distributed in the hope that it will be useful, but " + aurora.Bold("without any warranty").String() + " ; without even the implied warranty of " + aurora.Italic("merchantability").String() + " or " + aurora.Italic("fitness for a particular purpose").String() + ". See the " + aurora.Underline("GNU General Public License").String() + " for more details. You should have received a copy of the " + aurora.Underline("GNU General Public License").String() + " along with this program. If not, see " + aurora.Blue("http://www.gnu.org/licenses/").String() + ".")
	fmt.Println(aurora.Bold("Contact :"), aurora.Blue("https://gitlab.com/NatoBoram/git-to-ipfs"))
	fmt.Println("")

	// Git
	err := initGit()
	if err != nil {
		return
	}

	// IPFS
	err = initIPFS()
	if err != nil {
		return
	}

	receiveURL("git@gitlab.com:NatoBoram/git-to-ipfs.git")
}

func initGit() (err error) {

	// Check for Git
	path, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("Git is not installed.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println(aurora.Bold("Git :"), aurora.Blue(path))

	fmt.Println("")
	return
}

func initIPFS() (err error) {

	// Check for IPFS
	path, err := exec.LookPath("ipfs")
	if err != nil {
		fmt.Println("IPFS is not installed.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println(aurora.Bold("IPFS :"), aurora.Blue(path))

	// Enable sharding
	exec.Command("ipfs", "config", "--json", "Experimental.ShardingEnabled", "true").Run()

	// Check for IPFS Cluster Service
	path, err = exec.LookPath("ipfs-cluster-service")
	if err != nil {
		fmt.Println("IPFS Cluster Service is not installed.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println(aurora.Bold("IPFS Cluster Service :"), aurora.Blue(path))

	// Check for IPFS Cluster Control
	path, err = exec.LookPath("ipfs-cluster-ctl")
	if err != nil {
		fmt.Println("IPFS Cluster Control is not installed.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println(aurora.Bold("IPFS Cluster Control :"), aurora.Blue(path))

	fmt.Println("")
	return
}

func initBager() {

	// Options
	options := badger.DefaultOptions("~/.config/Gi")

	db, err := badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func receiveURL(url string) {

	// UUID
	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Couldn't generate a new UUID.")
		fmt.Println(err.Error())
		return
	}

	fmt.Println(aurora.Bold("URL :"), aurora.Blue(url))
	fmt.Println(aurora.Bold("UUID :"), uuid.String())

	// Clone
	out, err := exec.Command("git", "clone", url, uuid.String()).Output()
	if err != nil {
		fmt.Println("Couldn't clone the repository.")
		fmt.Println(aurora.Bold("Command :"), "git clone", aurora.Blue(url), uuid.String())

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
	size, err := dirSize(uuid.String())
	if err != nil {
		fmt.Println("Couldn't get the size of the git repository.")
		fmt.Println(err.Error())
		return
	}

	rmin := rmin(size)
	rmax := rmax(size)

	// IPFS-Cluster
	out, err = exec.Command("ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--name", url, "--replication-min", rmin, "--replication-max", rmax, uuid.String()).Output()
	if err != nil {
		fmt.Println("Couldn't add the repository to IPFS.")
		fmt.Println(aurora.Bold("Command :"), "ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--name", url, "--replication-min", rmin, "--replication-max", rmax, uuid.String())

		// Log the error from the command
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}

	hash := string(out)
	fmt.Println(aurora.Bold("IPFS :"), aurora.Cyan(hash))

}

func dirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func rmin(size int64) string {
	return strconv.FormatInt(1, 10)
}

func rmax(size int64) string {
	return strconv.FormatInt(size/(speed*seconds)+1, 10)
}
