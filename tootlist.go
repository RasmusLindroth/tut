package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type tootListFocus int

const (
	feedFocus tootListFocus = iota
	threadFocus
)

type TootList struct {
	app            *App
	Index          int
	Statuses       []*mastodon.Status
	Thread         []*mastodon.Status
	ThreadIndex    int
	View           *tview.List
	focus          tootListFocus
	loadingFeedOld bool
	loadingFeedNew bool
}

func NewTootList(app *App, viewList *tview.List) *TootList {
	return &TootList{
		app:   app,
		Index: 0,
		focus: feedFocus,
		View:  viewList,
	}
}

func (t *TootList) GetStatuses() []*mastodon.Status {
	if t.focus == threadFocus {
		return t.GetThread()
	}
	return t.GetFeed()
}

func (t *TootList) GetStatus(index int) (*mastodon.Status, error) {
	if t.focus == threadFocus {
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
	t.View.SetCurrentItem(t.GetFeedIndex())
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
	if t.focus == threadFocus {
		return t.GetThreadIndex()
	}
	return t.GetFeedIndex()
}

func (t *TootList) SetIndex(index int) {
	switch t.focus {
	case feedFocus:
		t.SetFeedIndex(index)
	case threadFocus:
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

	if index < 5 && t.focus == feedFocus {
		go func() {
			if t.loadingFeedNew {
				return
			}
			t.loadingFeedNew = true
			t.app.UI.LoadNewer(statuses[0])
			t.app.App.QueueUpdateDraw(func() {
				t.Draw()
				t.loadingFeedNew = false
			})
		}()
	}
	t.SetIndex(index)
	t.View.SetCurrentItem(index)
}

func (t *TootList) Next() {
	index := t.GetIndex()
	statuses := t.GetStatuses()

	if index+1 < len(statuses) {
		index++
	}

	if (len(statuses)-index) < 10 && t.focus == feedFocus {
		go func() {
			if t.loadingFeedOld || len(statuses) == 0 {
				return
			}
			t.loadingFeedOld = true
			t.app.UI.LoadOlder(statuses[len(statuses)-1])
			t.app.App.QueueUpdateDraw(func() {
				t.Draw()
				t.loadingFeedOld = false
			})
		}()
	}
	t.SetIndex(index)
	t.View.SetCurrentItem(index)
}

func (t *TootList) Draw() {
	t.View.Clear()

	var statuses []*mastodon.Status
	var index int

	switch t.focus {
	case feedFocus:
		statuses = t.GetFeed()
		index = t.GetFeedIndex()
	case threadFocus:
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
		t.View.InsertItem(currRow, content, "", 0, nil)
		currRow++
	}
	t.View.SetCurrentItem(index)
	t.app.UI.StatusText.ShowToot(index)
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
	t.focus = feedFocus
}

func (t *TootList) FocusThread() {
	t.focus = threadFocus
}

func (t *TootList) GoBack() {
	t.focus = feedFocus
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
