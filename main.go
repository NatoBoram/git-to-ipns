package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"time"

	rice "github.com/GeertJohan/go.rice"
	badger "github.com/dgraph-io/badger"
	"github.com/gorilla/mux"
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

	// User
	path, err := initUser()
	if err != nil {
		return
	}

	// Git
	err = initGit()
	if err != nil {
		return
	}

	// IPFS
	err = initIPFS()
	if err != nil {
		return
	}

	// Badger
	db, err := initBager(path)
	if err != nil {
		return
	}
	defer db.Close()

	// Refresh repos
	go func() {
		for {
			onAllRepos(db)
			time.Sleep(24 * time.Hour)
		}
	}()

	// receiveURL(db, "git@gitlab.com:NatoBoram/git-to-ipfs.git")

	// Router
	initMux(db)

	// Listen to CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {

			// Close database
			err = db.Close()
			if err != nil {
				fmt.Println("Couldn't save the database.")
				log.Fatalln(err.Error())
			}

			os.Exit(0)
		}
	}()

	// Wait
	<-make(chan struct{})
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

func initBager(homeDir string) (db *badger.DB, err error) {

	// Options
	options := badger.DefaultOptions(homeDir + rootDir)

	db, err = badger.Open(options)
	if err != nil {
		fmt.Println("Couldn't open a Badger Database.")
		fmt.Println(err.Error())
	}

	fmt.Println()
	return db, err
}

func initUser() (string, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Couldn't get the current user.")
		fmt.Println(err.Error())
	}

	return usr.HomeDir, err
}

func initMux(db *badger.DB) {
	r := mux.NewRouter()

	r.HandleFunc("/api/add/", func(w http.ResponseWriter, r *http.Request) { addHandler(db, w, r) }).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("web").HTTPBox()))

	log.Fatal(http.ListenAndServe(":62458", r))
}
