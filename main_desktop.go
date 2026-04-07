//go:build windows || linux || darwin

package main

import (
	"fmt"
	"os"

	"gioui.org/app"
	"github.com/getlantern/systray"
)

var (
	Org       *Organizer
	showGuiCh = make(chan struct{}, 1)
)

func main() {

	go systray.Run(OnReady, OnExit)

	app.Main()
}

func OnReady() {

	Org = NewOrganizer()
	Org.log("Starting")

	go GuiLoop()

	showGuiCh <- struct{}{}

	systray.SetTitle("Orginzation App")
	systray.SetTooltip("Orginzation App")

	IconRead, err := os.ReadFile("image/organization-chart.ico")
	if err != nil {
		Org.log(fmt.Sprintf("there an error reading the file : %v", err))
	}
	systray.SetIcon(IconRead)

	mStutes := systray.AddMenuItem("it's working", "")
	mStutes.Disable()

	systray.AddSeparator()

	mOpen := systray.AddMenuItem("Program settings", "Open Gio interface")
	mQuit := systray.AddMenuItem("Complete shutdown,", "Stop program and monitor")

	Org.Start()

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				Org.log("The Gio interface will be available here soon...")
				select {
				case showGuiCh <- struct{}{}:
				default:
				}

			case <-mQuit.ClickedCh:
				Org.log("the app is closing")
				systray.Quit()
			}
		}
	}()
}

func OnExit() {

	Org.log("Stopped")

	if Org != nil && Org.Watcher != nil {
		Org.Watcher.Close()
		Org.log("the watcher is closed")
	}
	os.Exit(0)
}
