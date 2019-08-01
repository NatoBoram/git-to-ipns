// API

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	badger "github.com/dgraph-io/badger"
	"github.com/gorilla/mux"
)

// errorf writes a swagger-compliant error response.
//
// This function was written by Google under the MIT license.
func errorf(w http.ResponseWriter, code int, format string, a ...interface{}) {
	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	out.Code = code
	out.Message = fmt.Sprintf(format, a...)

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "Could not format JSON for original message."}`, 500)
		return
	}

	http.Error(w, string(b), code)
}

func reposGetHandler(db *badger.DB, w http.ResponseWriter, r *http.Request) {

	// Put all repos in an array
	ch := dbList(db)
	var repos []Repo
	for repoerr := range ch {
		if repoerr.err != nil {
			fmt.Println(repoerr.err.Error())
			continue
		}

		repos = append(repos, repoerr.repo)
	}

	// Repos to JSON
	b, err := json.Marshal(repos)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not marshal JSON: %v", err)
		return
	}

	// Output the list of repos
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func reposPostHandler(db *badger.DB, w http.ResponseWriter, r *http.Request) {
	var received PostRepos

	// Decode body
	if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			errorf(w, http.StatusBadRequest, "Body was not valid JSON: %v", err)
			return
		}

		errorf(w, http.StatusInternalServerError, "Could not get body: %v", err)
		return
	}

	// Add the received URL
	repo, err := addURL(db, received.URL)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Couldn't handle the URL : %v", err)
		return
	}

	// Repo to JSON
	b, err := json.Marshal(repo)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not marshal JSON: %v", err)
		return
	}

	// Output the newly created repo
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func repoGetHandler(db *badger.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	repo, err := dbGet(db, link)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Couldn't get the specified repository : %v", err)
		return
	}

	// Repo to JSON
	b, err := json.Marshal(repo)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not marshal JSON: %v", err)
		return
	}

	// Output the specified repo
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func repoDeleteHandler(db *badger.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	repo, err := dbGet(db, link)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Couldn't find the specified repository : %v", err)
		return
	}

	rmRepo(db, repo)

	// Output `200`
	w.WriteHeader(http.StatusOK)
}
