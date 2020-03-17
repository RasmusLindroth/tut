package main

import (
	"strings"

	"github.com/rivo/tview"
)

func NewCmdBar(app *App, view *tview.InputField) *CmdBar {
	return &CmdBar{
		app:  app,
		View: view,
	}
}

type CmdBar struct {
	app  *App
	View *tview.InputField
}

func (c *CmdBar) GetInput() string {
	return strings.TrimSpace(c.View.GetText())
}
