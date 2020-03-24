package main

import "github.com/rivo/tview"

func NewStatusBar(app *App) *StatusBar {
	s := &StatusBar{
		Text: tview.NewTextView(),
	}

	s.Text.SetBackgroundColor(app.Config.Style.StatusBarBackground)
	s.Text.SetTextColor(app.Config.Style.StatusBarText)

	return s
}

type StatusBar struct {
	Text *tview.TextView
}

func (s *StatusBar) SetText(t string) {
	s.Text.SetText(t)
}
