package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type List struct {
}

func drawList(tut *Tut, data *mastodon.List, main *tview.TextView, controls *tview.TextView) {

	controlItem := config.ColorFromKey(tut.Config, tut.Config.Input.ListOpenFeed, true)

	main.SetText(fmt.Sprintf("List %s", tview.Escape(data.Title)))
	controls.SetText(controlItem)
}
