package main

// Repo is a user-submitted Git repository.
type Repo struct {
	UUID string // 1. Generate a UUID
	URL  string // 2. Download the repo
	IPFS string // 3. Add it to IPFS
	Key  string // 4. Generate a key
	IPNS string // 5. Add it to IPNS
}
