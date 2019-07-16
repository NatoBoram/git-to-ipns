package main

import "os"

// Paths
const (
	rootFolder   = "./Gi"
	databasePath = rootFolder + "/database.json"
)

// Permissions
const (
	permPrivateDirectory os.FileMode = 0700
	permPrivateFile      os.FileMode = 0600
)

// Calculation
const (
	speed   = 10 * 1024 * 1024
	seconds = 60
)
