package main

import (
	"bytes"
	"encoding/gob"
)

// String

func encodeString(s string) (b []byte, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewEncoder(buffer).Encode(s)
	b = buffer.Bytes()
	return
}

func decodeString(b []byte) (s string, err error) {
	buffer := bytes.NewBuffer(b)
	err = gob.NewDecoder(buffer).Decode(&s)
	return
}

// Repo

func (repo Repo) encode() (b []byte, err error) {
	buffer := &bytes.Buffer{}
	err = gob.NewEncoder(buffer).Encode(repo)
	b = buffer.Bytes()
	return
}

func decodeRepo(b []byte) (repo Repo, err error) {
	buffer := bytes.NewBuffer(b)
	err = gob.NewDecoder(buffer).Decode(&repo)
	return
}
