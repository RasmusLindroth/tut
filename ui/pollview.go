package ui

import (
	"fmt"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var durations = []string{
	"5 minutes",
	"30 minutes",
	"1 hour",
	"6 hours",
	"1 day",
	"3 days",
	"7 days",
}
var durationsTime = map[string]int64{
	"5 minutes":  60 * 50,
	"30 minutes": 60 * 30,
	"1 hour":     60 * 60,
	"6 hours":    60 * 60 * 6,
	"1 day":      60 * 60 * 24,
	"3 days":     60 * 60 * 24 * 3,
	"7 days":     60 * 60 * 24 * 7,
}

type PollView struct {
	tutView     *TutView
	shared      *Shared
	View        *tview.Flex
	info        *tview.TextView
	expiration  *tview.DropDown
	controls    *tview.Flex
	list        *tview.List
	poll        *mastodon.TootPoll
	scrollSleep *scrollSleep
}

func NewPollView(tv *TutView) *PollView {
	p := &PollView{
		tutView:    tv,
		shared:     tv.Shared,
		info:       NewTextView(tv.tut.Config),
		expiration: NewDropDown(tv.tut.Config),
		controls:   NewControlView(tv.tut.Config),
		list:       NewList(tv.tut.Config, false),
	}
	p.scrollSleep = NewScrollSleep(p.Next, p.Prev)
	p.Reset()
	p.View = pollViewUI(p)

	return p
}

func pollViewUI(p *PollView) *tview.Flex {
	var items []Control
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PollAdd, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PollEdit, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PollDelete, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PollMultiToggle, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PollExpiration, true))
	p.controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len+1, 0, false)
		} else {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len, 0, false)
		}
	}
	p.expiration.SetLabel("Expiration: ")
	p.expiration.SetOptions(durations, p.expirationSelected)
	p.expiration.SetCurrentOption(4)

	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if p.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(p.shared.Top.View, 1, 0, false)
	}
	r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(p.list, 0, 10, false), 0, 2, false).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(p.expiration, 1, 0, false).
			AddItem(p.info, 3, 0, false), 0, 1, false), 0, 1, false).
		AddItem(p.controls, 1, 0, false).
		AddItem(p.shared.Bottom.View, 2, 0, false)
	return r
}

func (p *PollView) Reset() {
	p.poll = &mastodon.TootPoll{
		Options:          []string{},
		ExpiresInSeconds: durationsTime[durations[4]],
		Multiple:         false,
		HideTotals:       false,
	}
	p.list.Clear()
	p.redrawInfo()
}

func (p *PollView) AddPoll(np *mastodon.Poll) {
	p.poll = &mastodon.TootPoll{
		Options:          []string{},
		ExpiresInSeconds: durationsTime[durations[4]],
		Multiple:         false,
		HideTotals:       false,
	}
	for _, opt := range np.Options {
		p.poll.Options = append(p.poll.Options, opt.Title)
		p.list.AddItem(opt.Title, "", 0, nil)
	}
	p.poll.Multiple = np.Multiple
	diff := time.Until(np.ExpiresAt)
	p.poll.ExpiresInSeconds = int64(diff.Seconds())
	p.redrawInfo()
}

func (p *PollView) HasPoll() bool {
	return p.list.GetItemCount() > 1
}

func (p *PollView) GetPoll() *mastodon.TootPoll {
	options := []string{}
	for i := 0; i < p.list.GetItemCount(); i++ {
		m, _ := p.list.GetItemText(i)
		options = append(options, m)
	}
	return &mastodon.TootPoll{
		Options:          options,
		ExpiresInSeconds: p.poll.ExpiresInSeconds,
		Multiple:         p.poll.Multiple,
		HideTotals:       false,
	}
}

func (p *PollView) redrawInfo() {
	content := fmt.Sprintf("Multiple answers: %v", p.poll.Multiple)
	p.info.SetText(content)
}

func (p *PollView) Prev() {
	index := p.list.GetCurrentItem()
	if index-1 >= 0 {
		p.list.SetCurrentItem(index - 1)
	}
}

func (p *PollView) Next() {
	index := p.list.GetCurrentItem()
	if index+1 < p.list.GetItemCount() {
		p.list.SetCurrentItem(index + 1)
	}
}

func (p *PollView) Add() {
	if p.list.GetItemCount() > 3 {
		p.tutView.ShowError("You can only have a maximum of 4 options.")
		return
	}
	text, valid, err := OpenEditorLengthLimit(p.tutView, "", 25)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't open editor. Error: %v", err),
		)
		return
	}
	if !valid {
		return
	}
	p.list.AddItem(text, "", 0, nil)
	p.list.SetCurrentItem(
		p.list.GetItemCount() - 1,
	)
}

func (p *PollView) Edit() {
	if p.list.GetItemCount() == 0 {
		return
	}
	text, _ := p.list.GetItemText(p.list.GetCurrentItem())
	text, valid, err := OpenEditorLengthLimit(p.tutView, text, 25)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't open editor. Error: %v", err),
		)
		return
	}
	if !valid {
		return
	}
	p.list.SetItemText(p.list.GetCurrentItem(), text, "")
}

func (p *PollView) Delete() {
	if p.list.GetItemCount() == 0 {
		return
	}
	item := p.list.GetCurrentItem()
	p.list.RemoveItem(item)
	if p.list.GetCurrentItem() < 0 && p.list.GetItemCount() > 0 {
		p.list.SetCurrentItem(0)
	}
}

func (p *PollView) ToggleMultiple() {
	p.poll.Multiple = !p.poll.Multiple
	p.redrawInfo()
}
func (p *PollView) expirationSelected(s string, index int) {
	_, v := p.expiration.GetCurrentOption()
	for k, dur := range durationsTime {
		if v == k {
			p.poll.ExpiresInSeconds = dur
			break
		}
	}
	p.exitExpiration()
}

func (p *PollView) expirationInput(event *tcell.EventKey) *tcell.EventKey {
	if p.tutView.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}
	if p.tutView.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}
	if p.tutView.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) ||
		p.tutView.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		p.exitExpiration()
		return nil
	}
	return event
}

func (p *PollView) FocusExpiration() {
	p.tutView.tut.App.SetInputCapture(p.expirationInput)
	p.tutView.tut.App.SetFocus(p.expiration)
	ev := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	p.tutView.tut.App.QueueEvent(ev)
}

func (p *PollView) exitExpiration() {
	p.tutView.tut.App.SetInputCapture(p.tutView.Input)
	p.tutView.tut.App.SetFocus(p.tutView.View)
}
