// Constants.

package main

import "os"

// Relative paths for the whole application.
// Use them with `dirHome`.
const (
	dirConfig = "/.config/gipns"
	dirBadger = dirConfig + "/badger"
	dirGit    = dirConfig + "/git"
)

// Permissions
const (
	permPrivateDirectory os.FileMode = 0700
	permPrivateFile      os.FileMode = 0600
)

// Calculations
const (
	speed   = 10 * 1024 * 1024
	seconds = 60
)

// Public Gateways
const (
	pgPermaweb = "/dnsaddr/permaweb.io"
	pgLibp2p   = "/dnsaddr/bootstrap.libp2p.io"
)
