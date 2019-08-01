package main

import (
	"os/exec"
	"testing"

	"github.com/logrusorgru/aurora"
)

func TestRun(t *testing.T) {
	_, err := run(
		exec.Command("go", "help"),
		".",
		"Couldn't run the Go command.",
		aurora.Bold("Command :"), "go", "help",
	)
	if err != nil {
		t.Errorf("Failed to run a command: %s.", err.Error())
	}
}
