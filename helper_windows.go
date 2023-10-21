package main

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/energye/systray"
)

func spawn(exe string, args ...string) *exec.Cmd {
	cmd := exec.Command(exe, args...)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	return cmd
}

func createTray() {
	systray.Run(systemTray, func() {
		os.Exit(0)
	})
}

func quitTray() {
	systray.Quit()
}
