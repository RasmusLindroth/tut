package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/rivo/tview"
)

type List struct {
}

func drawList(tv *TutView, data *mastodon.List, main *tview.TextView, controls *tview.Flex) {
	btn := NewControl(tv.tut.Config, tv.tut.Config.Input.ListOpenFeed, true)
	controls.Clear()
	controls.AddItem(NewControlButton(tv, btn), btn.Len, 0, false)

	main.SetText(fmt.Sprintf("List %s", tview.Escape(data.Title)))
}
