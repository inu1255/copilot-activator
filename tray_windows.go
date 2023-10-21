package main

import (
	_ "embed"
	"fmt"

	"github.com/energye/systray"
)

//go:embed build/windows/icon.ico
var icon []byte

func systemTray() {
	systray.SetIcon(icon) // read the icon from a file
	systray.SetTooltip(title)

	test := systray.AddMenuItem("test", "Show The Window")
	show := systray.AddMenuItem("Show", "Show The Window")
	// systray.AddSeparator()
	exit := systray.AddMenuItem("Exit", "Quit The Program")

	test.Click(func() {
		fmt.Println(isNormal())
	})
	show.Click(func() { showWindow() })
	exit.Click(func() {
		systray.Quit()
	})

	systray.SetOnClick(func(menu systray.IMenu) {
		toggleWindow()
	})
	systray.SetOnRClick(func(menu systray.IMenu) { _ = menu.ShowMenu() })
}
