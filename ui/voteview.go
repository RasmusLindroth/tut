package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type VoteView struct {
	tutView  *TutView
	shared   *Shared
	View     *tview.Flex
	textTop  *tview.TextView
	controls *tview.TextView
	list     *tview.List
	poll     *mastodon.Poll
	selected []int
}

func NewVoteView(tv *TutView) *VoteView {
	v := &VoteView{
		tutView:  tv,
		shared:   tv.Shared,
		textTop:  NewTextView(tv.tut.Config),
		controls: NewTextView(tv.tut.Config),
		list:     NewList(tv.tut.Config),
	}
	v.View = voteViewUI(v)

	return v
}

func voteViewUI(v *VoteView) *tview.Flex {
	var items []string
	items = append(items, config.ColorKey(v.tutView.tut.Config, "Select ", "Space/Enter", ""))
	items = append(items, config.ColorKey(v.tutView.tut.Config, "", "V", "ote"))
	v.controls.SetText(strings.Join(items, " "))

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(v.shared.Top.View, 1, 0, false).
		AddItem(v.textTop, 3, 0, false).
		AddItem(v.list, 0, 10, false).
		AddItem(v.controls, 1, 0, false).
		AddItem(v.shared.Bottom.View, 2, 0, false)
}

func (v *VoteView) SetPoll(poll *mastodon.Poll) {
	v.poll = poll
	v.selected = []int{}
	v.list.Clear()
	if v.poll.Multiple {
		v.textTop.SetText(
			tview.Escape("You can select multiple options. Press [v] to vote when you're finished selecting"),
		)
	} else {
		v.textTop.SetText(
			tview.Escape("You can only select ONE option. Press [v] to vote when you're finished selecting"),
		)
	}
	for _, o := range poll.Options {
		v.list.AddItem(tview.Escape(o.Title), "", 0, nil)
	}
}

func (v *VoteView) Prev() {
	index := v.list.GetCurrentItem()
	if index-1 >= 0 {
		v.list.SetCurrentItem(index - 1)
	}
}

func (v *VoteView) Next() {
	index := v.list.GetCurrentItem()
	if index+1 < v.list.GetItemCount() {
		v.list.SetCurrentItem(index + 1)
	}
}

func (v *VoteView) ToggleSelect() {
	index := v.list.GetCurrentItem()
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

func (v *VoteView) Select() {
	if !v.poll.Multiple && len(v.selected) > 0 {
		return
	}
	index := v.list.GetCurrentItem()
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
	v.list.SetItemText(index,
		tview.Escape(fmt.Sprintf("[x] %s", v.poll.Options[index].Title)),
		"")
}

func (v *VoteView) Unselect() {
	index := v.list.GetCurrentItem()
	sel := []int{}
	for _, value := range v.selected {
		if value == index {
			continue
		}
		sel = append(sel, value)
	}
	v.selected = sel
	v.list.SetItemText(index,
		tview.Escape(v.poll.Options[index].Title),
		"")
}

func (v *VoteView) Vote() {
	if len(v.selected) == 0 {
		return
	}
	p, err := v.tutView.tut.Client.Vote(v.poll, v.selected...)
	if err != nil {
		fmt.Printf("Couldn't vote. Error: %v\n", err)
		os.Exit(1)
		return
	}
	v.tutView.FocusMainNoHistory()
	v.tutView.RedrawPoll(p)
}

func (v *VoteView) InputHandler(event *tcell.EventKey) {
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
			//v.app.UI.StatusView.giveBackFocus()
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
			//v.app.UI.StatusView.giveBackFocus()
		}
	}
}
