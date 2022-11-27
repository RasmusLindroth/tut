package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/rivo/tview"
)

type Tag struct {
}

func drawTag(tv *TutView, data *mastodon.Tag, main *tview.TextView, controls *tview.Flex) {
	controls.Clear()
	var items []Control
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.TagOpenFeed, true))
	if data.Following != nil && data.Following == true {
		items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.TagFollow, false))
	} else {
		items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.TagFollow, true))

	}
	controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			controls.AddItem(NewControlButton(tv, item), item.Len+1, 0, false)
		} else {
			controls.AddItem(NewControlButton(tv, item), item.Len, 0, false)
		}
	}
	if main != nil {
		out := fmt.Sprintf("#%s\n\n", tview.Escape(data.Name))
		for _, h := range data.History {
			i, err := strconv.ParseInt(h.Day, 10, 64)
			if err != nil {
				continue
			}
			tm := time.Unix(i, 0)
			out += fmt.Sprintf("%s: %s accounts and %s toots\n",
				tm.Format("2006-01-02"), h.Accounts, h.Uses)
		}
		main.SetText(out)
		main.ScrollToBeginning()
	}
}
