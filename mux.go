package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	badger "github.com/dgraph-io/badger"
)

// errorf writes a swagger-compliant error response.
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

func addHandler(db *badger.DB, w http.ResponseWriter, r *http.Request) {
	var received AddURL

	if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			errorf(w, http.StatusBadRequest, "Body was not valid JSON: %v", err)
			return
		}
		errorf(w, http.StatusInternalServerError, "Could not get body: %v", err)
		return
	}

	repo, err := receiveURL(db, received.URL)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Couldn't properly handle the URL : %v", err)
		return
	}

	b, err := json.Marshal(repo)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not marshal JSON: %v", err)
		return
	}
	w.Write(b)
}
