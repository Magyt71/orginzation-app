package main

import (
	"image/color"
	"sync"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var (
	GuiISOpen bool
	guimu     sync.Mutex
)

func GuiLoop() {
	for range showGuiCh {
		OpenWindow()
	}
}

func OpenWindow() {
	W := new(app.Window)
	W.Option(app.Title("Orginzation_app"), app.Size(800, 600))

	if err := Run(W); err != nil {
		Org.log("GUI Error: " + err.Error())
	}
}

func Run(Window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			title := material.H1(theme, "Hello, Gio")

			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}

			title.Color = maroon

			title.Alignment = text.Middle

			title.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
