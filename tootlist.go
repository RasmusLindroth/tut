package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type TootListFocus int

const (
	TootListFeedFocus TootListFocus = iota
	TootListThreadFocus
)

type TootList struct {
	app            *App
	Index          int
	Statuses       []*mastodon.Status
	Thread         []*mastodon.Status
	ThreadIndex    int
	List           *tview.List
	Focus          TootListFocus
	loadingFeedOld bool
	loadingFeedNew bool
}

func NewTootList(app *App) *TootList {
	t := &TootList{
		app:   app,
		Index: 0,
		Focus: TootListFeedFocus,
		List:  tview.NewList(),
	}
	t.List.SetBackgroundColor(app.Config.Style.Background)
	t.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	t.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	t.List.ShowSecondaryText(false)
	t.List.SetHighlightFullLine(true)

	t.List.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
		if app.HaveAccount {
			app.UI.TootView.ShowToot(index)
		}
	})

	return t
}

func (t *TootList) GetStatuses() []*mastodon.Status {
	if t.Focus == TootListThreadFocus {
		return t.GetThread()
	}
	return t.GetFeed()
}

func (t *TootList) GetStatus(index int) (*mastodon.Status, error) {
	if t.Focus == TootListThreadFocus {
		return t.GetThreadStatus(index)
	}
	return t.GetFeedStatus(index)
}

func (t *TootList) SetFeedStatuses(s []*mastodon.Status) {
	t.Statuses = s
	t.Draw()
}

func (t *TootList) PrependFeedStatuses(s []*mastodon.Status) {
	t.Statuses = append(s, t.Statuses...)
	t.SetFeedIndex(
		t.GetFeedIndex() + len(s),
	)
	t.List.SetCurrentItem(t.GetFeedIndex())
}

func (t *TootList) AppendFeedStatuses(s []*mastodon.Status) {
	t.Statuses = append(t.Statuses, s...)
}

func (t *TootList) GetFeed() []*mastodon.Status {
	return t.Statuses
}

func (t *TootList) GetFeedStatus(index int) (*mastodon.Status, error) {
	statuses := t.GetFeed()
	if index < len(statuses) {
		return statuses[index], nil
	}
	return nil, fmt.Errorf("no status with that index")
}

func (t *TootList) GetIndex() int {
	if t.Focus == TootListThreadFocus {
		return t.GetThreadIndex()
	}
	return t.GetFeedIndex()
}

func (t *TootList) SetIndex(index int) {
	switch t.Focus {
	case TootListFeedFocus:
		t.SetFeedIndex(index)
	case TootListThreadFocus:
		t.SetThreadIndex(index)
	}
}

func (t *TootList) GetFeedIndex() int {
	return t.Index
}

func (t *TootList) SetFeedIndex(index int) {
	t.Index = index
}

func (t *TootList) GetThreadIndex() int {
	return t.ThreadIndex
}

func (t *TootList) SetThreadIndex(index int) {
	t.ThreadIndex = index
}

func (t *TootList) Prev() {
	index := t.GetIndex()
	statuses := t.GetStatuses()

	if index-1 > -1 {
		index--
	}

	if index < 5 && t.Focus == TootListFeedFocus {
		go func() {
			if t.loadingFeedNew {
				return
			}
			t.loadingFeedNew = true
			t.app.UI.LoadNewer(statuses[0])
			t.app.UI.Root.QueueUpdateDraw(func() {
				t.Draw()
				t.loadingFeedNew = false
			})
		}()
	}
	t.SetIndex(index)
	t.List.SetCurrentItem(index)
}

func (t *TootList) Next() {
	index := t.GetIndex()
	statuses := t.GetStatuses()

	if index+1 < len(statuses) {
		index++
	}

	if (len(statuses)-index) < 10 && t.Focus == TootListFeedFocus {
		go func() {
			if t.loadingFeedOld || len(statuses) == 0 {
				return
			}
			t.loadingFeedOld = true
			t.app.UI.LoadOlder(statuses[len(statuses)-1])
			t.app.UI.Root.QueueUpdateDraw(func() {
				t.Draw()
				t.loadingFeedOld = false
			})
		}()
	}
	t.SetIndex(index)
	t.List.SetCurrentItem(index)
}

func (t *TootList) Draw() {
	t.List.Clear()

	var statuses []*mastodon.Status
	var index int

	switch t.Focus {
	case TootListFeedFocus:
		statuses = t.GetFeed()
		index = t.GetFeedIndex()
	case TootListThreadFocus:
		statuses = t.GetThread()
		index = t.GetThreadIndex()
	}
	if len(statuses) == 0 {
		return
	}

	today := time.Now()
	ty, tm, td := today.Date()
	currRow := 0
	for _, s := range statuses {
		sLocal := s.CreatedAt.Local()
		sy, sm, sd := sLocal.Date()
		format := "2006-01-02 15:04"
		if ty == sy && tm == sm && td == sd {
			format = "15:04"
		}
		content := fmt.Sprintf("%s %s", sLocal.Format(format), s.Account.Acct)
		t.List.InsertItem(currRow, content, "", 0, nil)
		currRow++
	}
	t.List.SetCurrentItem(index)
	t.app.UI.TootView.ShowToot(index)
}

func (t *TootList) SetThread(s []*mastodon.Status, index int) {
	t.Thread = s
	t.SetThreadIndex(index)
}

func (t *TootList) GetThread() []*mastodon.Status {
	return t.Thread
}

func (t *TootList) GetThreadStatus(index int) (*mastodon.Status, error) {
	statuses := t.GetThread()
	if index < len(statuses) {
		return statuses[index], nil
	}
	return nil, fmt.Errorf("no status with that index")
}

func (t *TootList) FocusFeed() {
	t.Focus = TootListFeedFocus
}

func (t *TootList) FocusThread() {
	t.Focus = TootListThreadFocus
}

func (t *TootList) GoBack() {
	t.Focus = TootListFeedFocus
	t.Draw()
}

func (t *TootList) Reply() {
	status, err := t.GetStatus(t.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	if status.Reblog != nil {
		status = status.Reblog
	}

	users := []string{"@" + status.Account.Acct}
	for _, m := range status.Mentions {
		users = append(users, "@"+m.Acct)
	}
}
