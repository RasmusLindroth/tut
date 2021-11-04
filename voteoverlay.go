package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

func NewVoteOverlay(app *App) *VoteOverlay {
	v := &VoteOverlay{
		app:        app,
		Flex:       tview.NewFlex(),
		TextTop:    tview.NewTextView(),
		TextBottom: tview.NewTextView(),
		List:       tview.NewList(),
	}

	v.TextTop.SetBackgroundColor(app.Config.Style.Background)
	v.TextTop.SetTextColor(app.Config.Style.Text)
	v.TextTop.SetDynamicColors(true)
	v.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	v.TextBottom.SetDynamicColors(true)
	v.List.SetBackgroundColor(app.Config.Style.Background)
	v.List.SetMainTextColor(app.Config.Style.Text)
	v.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	v.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	v.List.ShowSecondaryText(false)
	v.List.SetHighlightFullLine(true)
	v.Flex.SetDrawFunc(app.Config.ClearContent)
	var items []string
	items = append(items, ColorKey(app.Config, "Select ", "Space/Enter", ""))
	items = append(items, ColorKey(app.Config, "", "V", "ote"))
	v.TextBottom.SetText(strings.Join(items, " "))
	return v
}

type VoteOverlay struct {
	app        *App
	Flex       *tview.Flex
	TextTop    *tview.TextView
	TextBottom *tview.TextView
	List       *tview.List
	poll       *mastodon.Poll
	selected   []int
}

func (v *VoteOverlay) SetPoll(poll *mastodon.Poll) {
	v.poll = poll
	v.selected = []int{}
	v.List.Clear()
	if v.poll.Multiple {
		v.TextTop.SetText(
			tview.Escape("You can select multiple options. Press [v] to vote when you're finished selecting"),
		)
	} else {
		v.TextTop.SetText(
			tview.Escape("You can only select ONE option. Press [v] to vote when you're finished selecting"),
		)
	}
	for _, o := range poll.Options {
		v.List.AddItem(tview.Escape(o.Title), "", 0, nil)
	}
}

func (v *VoteOverlay) Prev() {
	index := v.List.GetCurrentItem()
	if index-1 >= 0 {
		v.List.SetCurrentItem(index - 1)
	}
}

func (v *VoteOverlay) Next() {
	index := v.List.GetCurrentItem()
	if index+1 < v.List.GetItemCount() {
		v.List.SetCurrentItem(index + 1)
	}
}

func (v *VoteOverlay) ToggleSelect() {
	index := v.List.GetCurrentItem()
	inSelected := false
	for _, value := range v.selected {
		if index == value {
			inSelected = true
			break
		}
	}
	if inSelected {
		v.Unselect()
	} else {
		v.Select()
	}
}

func (v *VoteOverlay) Select() {
	if !v.poll.Multiple && len(v.selected) > 0 {
		return
	}
	index := v.List.GetCurrentItem()
	inSelected := false
	for _, value := range v.selected {
		if index == value {
			inSelected = true
			break
		}
	}
	if inSelected {
		return
	}
	v.selected = append(v.selected, index)
	v.List.SetItemText(index,
		tview.Escape(fmt.Sprintf("[x] %s", v.poll.Options[index].Title)),
		"")
}

func (v *VoteOverlay) Unselect() {
	index := v.List.GetCurrentItem()
	sel := []int{}
	for _, value := range v.selected {
		if value == index {
			continue
		}
		sel = append(sel, value)
	}
	v.selected = sel
	v.List.SetItemText(index,
		tview.Escape(v.poll.Options[index].Title),
		"")
}

func (v *VoteOverlay) Vote() {
	if len(v.selected) == 0 {
		return
	}
	p, err := v.app.API.Vote(v.poll, v.selected...)
	if err != nil {
		v.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't vote. Error: %v\n", err))
		return
	}
	v.app.UI.StatusView.RedrawPoll(p)
	if v.app.UI.StatusView.lastList == NotificationPaneFocus {
		v.app.UI.SetFocus(NotificationPaneFocus)
	} else {
		v.app.UI.SetFocus(LeftPaneFocus)
	}
}

func (v *VoteOverlay) InputHandler(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			v.Next()
		case 'k', 'K':
			v.Prev()
		case 'v', 'V':
			v.Vote()
		case ' ':
			v.ToggleSelect()
		case 'q', 'Q':
			if v.app.UI.StatusView.lastList == NotificationPaneFocus {
				v.app.UI.SetFocus(NotificationPaneFocus)
			} else {
				v.app.UI.SetFocus(LeftPaneFocus)
			}
		}
	} else {
		switch event.Key() {
		case tcell.KeyEnter:
			v.ToggleSelect()
		case tcell.KeyUp:
			v.Prev()
		case tcell.KeyDown:
			v.Next()
		case tcell.KeyEsc:
			if v.app.UI.StatusView.lastList == NotificationPaneFocus {
				v.app.UI.SetFocus(NotificationPaneFocus)
			} else {
				v.app.UI.SetFocus(LeftPaneFocus)
			}
		}
	}
}
