package main

import (
	"testing"
)

func TestDirSize(t *testing.T) {
	_, err := dirSize(".")
	if err != nil {
		t.Errorf("dirSize failed: %s.", err.Error())
	}
}

// Replication minimum for a cluster consisting of a single peer.
// If the cluster consisted of many, many peers, then the minimum should be 3 to allow re-pinning.
func TestRMin(t *testing.T) {
	tests := []struct {
		input  int64
		output string
	}{
		{0, "1"},                      // Test value     |   0 MiB | 1
		{512 * 1024 * 1024, "1"},      // Test value     | 512 MiB | 1
		{1 * 1024 * 1024 * 1024, "1"}, // Test value     |   1 GiB | 1
		{321756774, "1"},              // Smallest build | 307 MiB | 1
		{499107798, "1"},              // Average build  | 476 MiB | 1
		{755192969, "1"},              // Largest build  | 720 MiB | 1
	}

	for _, test := range tests {
		output := rmin(test.input)
		if output != test.output {
			t.Errorf("Replication Minimum was incorrect, got: %s, want: %s.", output, test.output)
		}
	}
}

// Replication maximum for a cluster consisting of a single peer.
// The goal is to download the file under 60 seconds at 10 Mbps.
//
// This is the calculation for an average file.
//
// | Value | MiB | o         |
// | :---- | --: | --------: |
// | File  | 512 | 536870912 |
// | Speed |  10 |  10485760 |
//
// | Unit        | Value |
// | :---------- | ----: |
// | Time        | 51.2s |
// | Replication | 0.85x |
//
// The replication factor is truncated then is added +1.
func TestRMax(t *testing.T) {
	tests := []struct {
		input  int64
		output string
	}{
		{0, "1"},                      // Test value     |   0 MiB | 1
		{512 * 1024 * 1024, "1"},      // Test value     | 512 MiB | 1
		{1 * 1024 * 1024 * 1024, "2"}, // Test value     |   1 GiB | 2
		{321756774, "1"},              // Smallest build | 307 MiB | 1
		{499107798, "1"},              // Average build  | 476 MiB | 1
		{755192969, "2"},              // Largest build  | 720 MiB | 2
	}

	for _, test := range tests {
		output := rmax(test.input)

		if output != test.output {
			t.Errorf("Replication Maximum was incorrect, got: %s, want: %s.", output, test.output)
		}
	}
}
