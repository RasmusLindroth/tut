package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/rivo/tview"
)

type List struct {
}

func drawList(tut *Tut, data *mastodon.List, main *tview.TextView, controls *tview.Flex) {
	btn := NewControl(tut.Config, tut.Config.Input.ListOpenFeed, true)
	controls.AddItem(NewControlButton(tut.Config, btn.Label), btn.Len, 0, false)

	main.SetText(fmt.Sprintf("List %s", tview.Escape(data.Title)))
}
