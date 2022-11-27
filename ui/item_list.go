package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/rivo/tview"
)

type List struct {
}

func drawList(tv *TutView, data *mastodon.List, main *tview.TextView, controls *tview.Flex) {
	controls.Clear()
	var items []Control
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.ListOpenFeed, true))
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.ListUserList, true))
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.ListUserAdd, true))
	controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			controls.AddItem(NewControlButton(tv, item), item.Len+1, 0, false)
		} else {
			controls.AddItem(NewControlButton(tv, item), item.Len, 0, false)
		}
	}

	if main != nil {
		main.SetText(fmt.Sprintf("List %s", tview.Escape(data.Title)))
	}
}
