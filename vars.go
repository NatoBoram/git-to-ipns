// Global state of the application.

package main

import "errors"

// Path
var (
	dirHome string
)

// Errors
var (
	ErrNoURL = errors.New("url is empty")
)

// PublicGateways is a list of public IPFS gateways that can be used with `ipfs swarm connect`.
var PublicGateways = []string{
	pgPermaweb,
	pgLibp2p,
}
