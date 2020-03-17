package main

import "github.com/rivo/tview"

func NewControls(app *App, view *tview.TextView) *Controls {
	return &Controls{
		app:  app,
		View: view,
	}
}

type Controls struct {
	app  *App
	View *tview.TextView
}
