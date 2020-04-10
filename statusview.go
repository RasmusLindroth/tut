package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

func NewStatusView(app *App, tl TimelineType) *StatusView {
	t := &StatusView{
		app:          app,
		list:         tview.NewList(),
		text:         tview.NewTextView(),
		controls:     tview.NewTextView(),
		focus:        LeftPaneFocus,
		loadingNewer: false,
		loadingOlder: false,
	}
	t.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.text, 0, 9, false).
		AddItem(t.controls, 1, 0, false)

	t.list.SetBackgroundColor(app.Config.Style.Background)
	t.list.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	t.list.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	t.list.ShowSecondaryText(false)
	t.list.SetHighlightFullLine(true)

	t.text.SetWordWrap(true).SetDynamicColors(true)
	t.text.SetBackgroundColor(app.Config.Style.Background)
	t.text.SetTextColor(app.Config.Style.Text)
	t.controls.SetDynamicColors(true)
	t.controls.SetBackgroundColor(app.Config.Style.Background)

	if app.Config.General.AutoLoadNewer {
		go func() {
			d := time.Second * time.Duration(app.Config.General.AutoLoadSeconds)
			ticker := time.NewTicker(d)
			for {
				select {
				case <-ticker.C:
					t.loadNewer()
				}
			}
		}()
	}
	return t
}

type StatusView struct {
	app          *App
	list         *tview.List
	flex         *tview.Flex
	text         *tview.TextView
	controls     *tview.TextView
	feeds        []Feed
	focus        FocusAt
	loadingNewer bool
	loadingOlder bool
}

func (t *StatusView) AddFeed(f Feed) {
	t.feeds = append(t.feeds, f)
	f.DrawList()
	t.list.SetCurrentItem(f.GetSavedIndex())
	f.DrawToot()
	t.drawDesc()
}

func (t *StatusView) RemoveLatestFeed() {
	t.feeds = t.feeds[:len(t.feeds)-1]
	feed := t.feeds[len(t.feeds)-1]
	feed.DrawList()
	t.list.SetCurrentItem(feed.GetSavedIndex())
	feed.DrawToot()
	t.drawDesc()
}

func (t *StatusView) GetLeftView() tview.Primitive {
	if len(t.feeds) > 0 {
		feed := t.feeds[len(t.feeds)-1]
		feed.DrawList()
		feed.DrawToot()
	}
	return t.list
}

func (t *StatusView) GetRightView() tview.Primitive {
	return t.flex
}

func (t *StatusView) GetTextWidth() int {
	_, _, width, _ := t.text.GetInnerRect()
	return width
}

func (t *StatusView) GetCurrentItem() int {
	return t.list.GetCurrentItem()
}

func (t *StatusView) GetCurrentStatus() *mastodon.Status {
	if len(t.feeds) == 0 {
		return nil
	}
	return t.feeds[len(t.feeds)-1].GetCurrentStatus()
}

func (t *StatusView) GetCurrentUser() *mastodon.Account {
	if len(t.feeds) == 0 {
		return nil
	}
	return t.feeds[len(t.feeds)-1].GetCurrentUser()
}

func (t *StatusView) ScrollToBeginning() {
	t.text.ScrollToBeginning()
}

func (t *StatusView) inputBoth(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'g':
			t.home()
		case 'G':
			t.end()
		case 'q', 'Q':
			if len(t.feeds) > 1 {
				t.RemoveLatestFeed()
			} else {
				t.app.UI.Root.Stop()
			}
		}
	} else {
		switch event.Key() {
		case tcell.KeyCtrlC:
			t.app.UI.Root.Stop()
		case tcell.KeyHome:
			t.home()
		case tcell.KeyEnd:
			t.end()
		}
	}
	if len(t.feeds) > 0 {
		feed := t.feeds[len(t.feeds)-1]
		feed.Input(event)
	}
}

func (t *StatusView) inputLeft(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'v', 'V':
			t.app.UI.FocusAt(t.text, "--VIEW--")
			t.focus = RightPaneFocus
		case 'k', 'K':
			t.prev()
		case 'j', 'J':
			t.next()
		}
	} else {
		switch event.Key() {
		case tcell.KeyUp:
			t.prev()
		case tcell.KeyDown:
			t.next()
		case tcell.KeyEsc:
			if len(t.feeds) > 1 {
				t.RemoveLatestFeed()
			}
		}
	}
}

func (t *StatusView) inputRight(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {

		}
	} else {
		switch event.Key() {
		case tcell.KeyEsc:
			t.app.UI.FocusAt(nil, "--LIST--")
			t.focus = LeftPaneFocus
		}
	}
}

func (t *StatusView) Input(event *tcell.EventKey) *tcell.EventKey {
	t.inputBoth(event)
	if len(t.feeds) == 0 {
		return event
	}

	if t.focus == LeftPaneFocus {
		t.inputLeft(event)
		return nil
	} else {
		t.inputRight(event)
	}

	return event
}

func (t *StatusView) SetList(items <-chan string) {
	t.list.Clear()
	for s := range items {
		t.list.AddItem(s, "", 0, nil)
	}
}
func (t *StatusView) SetText(text string) {
	t.text.SetText(text)
}

func (t *StatusView) SetControls(text string) {
	t.controls.SetText(text)
}

func (t *StatusView) drawDesc() {
	if len(t.feeds) == 0 {
		t.app.UI.SetTopText("")
		return
	}
	l := len(t.feeds)
	f := t.feeds[l-1]
	t.app.UI.SetTopText(
		fmt.Sprintf("%s (%d/%d)", f.GetDesc(), l, l),
	)
}

func (t *StatusView) prev() {
	current := t.list.GetCurrentItem()
	if current-1 >= 0 {
		current--
	}
	t.list.SetCurrentItem(current)
	t.feeds[len(t.feeds)-1].DrawToot()

	if current < 4 {
		t.loadNewer()
	}
}

func (t *StatusView) next() {
	t.list.SetCurrentItem(
		t.list.GetCurrentItem() + 1,
	)
	t.feeds[len(t.feeds)-1].DrawToot()

	count := t.list.GetItemCount()
	current := t.list.GetCurrentItem()
	if (count - current + 1) < 5 {
		t.loadOlder()
	}
}

func (t *StatusView) home() {
	t.list.SetCurrentItem(0)
	t.feeds[len(t.feeds)-1].DrawToot()
	t.loadOlder()
}

func (t *StatusView) end() {
	t.list.SetCurrentItem(-1)
	t.feeds[len(t.feeds)-1].DrawToot()
	t.loadOlder()
}

func (t *StatusView) loadNewer() {
	feedIndex := len(t.feeds) - 1
	if t.loadingNewer || feedIndex < 0 {
		return
	}
	t.loadingNewer = true
	go func() {
		new := t.feeds[feedIndex].LoadNewer()
		if new == 0 || feedIndex != len(t.feeds)-1 {
			t.loadingNewer = false
			return
		}
		t.app.UI.Root.QueueUpdateDraw(func() {
			index := t.list.GetCurrentItem()
			t.feeds[feedIndex].DrawList()
			newIndex := index + new
			if index == 0 && t.feeds[feedIndex].FeedType() == UserFeedType {
				newIndex = 0
			}
			t.list.SetCurrentItem(newIndex)
			t.loadingNewer = false
		})
	}()
}

func (t *StatusView) loadOlder() {
	feedIndex := len(t.feeds) - 1
	if t.loadingOlder || feedIndex < 0 {
		return
	}
	t.loadingOlder = true
	go func() {
		new := t.feeds[feedIndex].LoadOlder()
		if new == 0 || feedIndex != len(t.feeds)-1 {
			t.loadingOlder = false
			return
		}
		t.app.UI.Root.QueueUpdateDraw(func() {
			index := t.list.GetCurrentItem()
			t.feeds[feedIndex].DrawList()
			t.list.SetCurrentItem(index)
			t.loadingOlder = false
		})
	}()
}
