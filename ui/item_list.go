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

	controlItem := config.ColorKey(tut.Config, "", "O", "pen")

	main.SetText(fmt.Sprintf("Press O or <Enter> to open list %s", tview.Escape(data.Title)))
	controls.SetText(controlItem)
}
