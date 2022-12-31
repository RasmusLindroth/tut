package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/rivo/tview"
)

type VoteView struct {
	tutView     *TutView
	shared      *Shared
	View        *tview.Flex
	textTop     *tview.TextView
	controls    *tview.Flex
	list        *tview.List
	poll        *mastodon.Poll
	selected    []int
	scrollSleep *scrollSleep
}

func NewVoteView(tv *TutView) *VoteView {
	v := &VoteView{
		tutView:  tv,
		shared:   tv.Shared,
		textTop:  NewTextView(tv.tut.Config),
		controls: NewControlView(tv.tut.Config),
		list:     NewList(tv.tut.Config, false),
	}
	v.scrollSleep = NewScrollSleep(v.Next, v.Prev)
	v.View = voteViewUI(v)

	return v
}

func voteViewUI(v *VoteView) *tview.Flex {
	var items []Control
	items = append(items, NewControl(v.tutView.tut.Config, v.tutView.tut.Config.Input.VoteSelect, true))
	v.controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			v.controls.AddItem(NewControlButton(v.tutView, item), item.Len+1, 0, false)
		} else {
			v.controls.AddItem(NewControlButton(v.tutView, item), item.Len, 0, false)
		}
	}

	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if v.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(v.shared.Top.View, 1, 0, false)
	}
	r.AddItem(v.textTop, 3, 0, false).
		AddItem(v.list, 0, 10, false).
		AddItem(v.controls, 1, 0, false).
		AddItem(v.shared.Bottom.View, 2, 0, false)
	return r
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
		v.tutView.ShowError(
			fmt.Sprintf("Couldn't vote. Error: %v\n", err),
		)
		return
	}
	v.tutView.FocusMainNoHistory()
	v.tutView.RedrawPoll(p)
}
