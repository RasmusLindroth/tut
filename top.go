package main

import "github.com/rivo/tview"

func NewTop(app *App) *Top {
	t := &Top{
		Text: tview.NewTextView(),
	}

	t.Text.SetBackgroundColor(app.Config.Style.TopBarBackground)
	t.Text.SetTextColor(app.Config.Style.TopBarText)

	return t
}

type Top struct {
	Text *tview.TextView
}
