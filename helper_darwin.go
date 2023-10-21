package main

import (
	"os/exec"
)

func spawn(exe string, args ...string) *exec.Cmd {
	cmd := exec.Command(exe, args...)
	return cmd
}

func createTray() {
}

func quitTray() {
}
