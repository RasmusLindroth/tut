package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
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
		lastList:     LeftPaneFocus,
		loadingNewer: false,
		loadingOlder: false,
	}
	if t.app.Config.General.NotificationFeed {
		t.notificationView = NewNotificationView(app)
	}

	t.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.text, 0, 9, false).
		AddItem(t.controls, 1, 0, false)

	t.list.SetMainTextColor(app.Config.Style.Text)
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
		if app.Config.General.NotificationFeed {
			go func() {
				d := time.Second * time.Duration(app.Config.General.AutoLoadSeconds)
				ticker := time.NewTicker(d)
				for {
					select {
					case <-ticker.C:
						t.notificationView.loadNewer()
					}
				}
			}()
		}
	}
	return t
}

type StatusView struct {
	app              *App
	list             *tview.List
	flex             *tview.Flex
	text             *tview.TextView
	controls         *tview.TextView
	feeds            []Feed
	focus            FocusAt
	lastList         FocusAt
	loadingNewer     bool
	loadingOlder     bool
	notificationView *NotificationView
}

func (t *StatusView) AddFeed(f Feed) {
	t.feeds = append(t.feeds, f)
	f.DrawList()
	t.list.SetCurrentItem(f.GetSavedIndex())
	f.DrawToot()
	t.drawDesc()

	if t.lastList == NotificationPaneFocus {
		t.app.UI.SetFocus(LeftPaneFocus)
		t.focus = LeftPaneFocus
		t.lastList = LeftPaneFocus
	}
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

func (t *StatusView) GetNotificationView() tview.Primitive {
	if t.notificationView != nil {
		t.notificationView.feed.DrawList()
		return t.notificationView.list
	}
	return nil
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
	if len(t.feeds) > 0 && t.lastList == LeftPaneFocus {
		feed := t.feeds[len(t.feeds)-1]
		feed.Input(event)
	} else if t.lastList == NotificationPaneFocus {
		t.notificationView.feed.Input(event)
	}
}

func (t *StatusView) inputBack(q bool) {
	if t.app.UI.Focus == LeftPaneFocus && len(t.feeds) > 1 {
		t.RemoveLatestFeed()
	} else if t.app.UI.Focus == LeftPaneFocus && q {
		t.app.UI.Root.Stop()
	} else if t.app.UI.Focus == NotificationPaneFocus {
		t.app.UI.SetFocus(LeftPaneFocus)
		t.focus = LeftPaneFocus
		t.lastList = LeftPaneFocus
		t.feeds[len(t.feeds)-1].DrawToot()
	}
}

func (t *StatusView) inputLeft(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'v', 'V':
			t.app.UI.SetFocus(RightPaneFocus)
			t.focus = RightPaneFocus
			t.app.UI.StatusBar.Text.SetBackgroundColor(
				t.app.Config.Style.StatusBarViewBackground,
			)
			t.app.UI.StatusBar.Text.SetTextColor(
				t.app.Config.Style.StatusBarViewText,
			)
		case 'k', 'K':
			t.prev()
		case 'j', 'J':
			t.next()
		case 'n', 'N':
			if t.app.Config.General.NotificationFeed {
				t.app.UI.SetFocus(NotificationPaneFocus)
				t.focus = NotificationPaneFocus
				t.lastList = NotificationPaneFocus
				t.notificationView.feed.DrawToot()
			}
		case 'q', 'Q':
			t.inputBack(true)
		}
	} else {
		switch event.Key() {
		case tcell.KeyUp:
			t.prev()
		case tcell.KeyDown:
			t.next()
		case tcell.KeyPgUp, tcell.KeyCtrlB:
			t.pgup()
		case tcell.KeyPgDn, tcell.KeyCtrlF:
			t.pgdown()
		case tcell.KeyEsc:
			t.inputBack(false)
		}
	}
}

func (t *StatusView) inputRightQuit() {
	if t.lastList == LeftPaneFocus {
		t.app.UI.SetFocus(LeftPaneFocus)
		t.focus = LeftPaneFocus
	} else if t.lastList == NotificationPaneFocus {
		t.app.UI.SetFocus(NotificationPaneFocus)
		t.focus = NotificationPaneFocus
	}
	t.app.UI.StatusBar.Text.SetBackgroundColor(
		t.app.Config.Style.StatusBarBackground,
	)
	t.app.UI.StatusBar.Text.SetTextColor(
		t.app.Config.Style.StatusBarText,
	)
}

func (t *StatusView) inputRight(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'q', 'Q':
			t.inputRightQuit()
		}
	} else {
		switch event.Key() {
		case tcell.KeyEsc:
			t.inputRightQuit()
		}
	}
}

func (t *StatusView) Input(event *tcell.EventKey) *tcell.EventKey {
	t.inputBoth(event)
	if len(t.feeds) == 0 {
		return event
	}

	switch t.focus {
	case LeftPaneFocus:
		t.inputLeft(event)
		return nil
	case NotificationPaneFocus:
		t.inputLeft(event)
		return nil
	default:
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
	var current int
	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		current = t.GetCurrentItem()
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		current = t.notificationView.list.GetCurrentItem()
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	if current-1 >= 0 {
		current--
	}
	list.SetCurrentItem(current)
	feed.DrawToot()

	if current < 4 {
		switch t.app.UI.Focus {
		case LeftPaneFocus:
			t.loadNewer()
		case NotificationPaneFocus:
			t.notificationView.loadNewer()
		}
	}
}

func (t *StatusView) next() {
	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	list.SetCurrentItem(
		list.GetCurrentItem() + 1,
	)
	feed.DrawToot()

	count := list.GetItemCount()
	current := list.GetCurrentItem()
	if (count - current + 1) < 5 {
		switch t.app.UI.Focus {
		case LeftPaneFocus:
			t.loadOlder()
		case NotificationPaneFocus:
			t.notificationView.loadOlder()
		}
	}
}

func (t *StatusView) pgdown() {
	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	_, _, _, height := list.GetInnerRect()
	i := list.GetCurrentItem() + height - 1
	list.SetCurrentItem(i)
	feed.DrawToot()

	count := list.GetItemCount()
	current := list.GetCurrentItem()
	if (count - current + 1) < 5 {
		switch t.app.UI.Focus {
		case LeftPaneFocus:
			t.loadOlder()
		case NotificationPaneFocus:
			t.notificationView.loadOlder()
		}
	}
}

func (t *StatusView) pgup() {
	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	_, _, _, height := list.GetInnerRect()
	i := list.GetCurrentItem() - height + 1
	if i < 0 {
		i = 0
	}
	list.SetCurrentItem(i)
	feed.DrawToot()

	current := list.GetCurrentItem()
	if current < 4 {
		switch t.app.UI.Focus {
		case LeftPaneFocus:
			t.loadNewer()
		case NotificationPaneFocus:
			t.notificationView.loadNewer()
		}
	}
}

func (t *StatusView) home() {
	if t.focus == RightPaneFocus {
		t.text.ScrollToBeginning()
		return
	}

	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	list.SetCurrentItem(0)
	feed.DrawToot()

	switch t.app.UI.Focus {
	case LeftPaneFocus:
		t.loadNewer()
	case NotificationPaneFocus:
		t.notificationView.loadNewer()
	}
}

func (t *StatusView) end() {
	if t.focus == RightPaneFocus {
		t.text.ScrollToEnd()
		return
	}

	var list *tview.List
	var feed Feed
	if t.app.UI.Focus == LeftPaneFocus {
		list = t.list
		feed = t.feeds[len(t.feeds)-1]
	} else {
		list = t.notificationView.list
		feed = t.notificationView.feed
	}

	list.SetCurrentItem(-1)
	feed.DrawToot()

	switch t.app.UI.Focus {
	case LeftPaneFocus:
		t.loadOlder()
	case NotificationPaneFocus:
		t.notificationView.loadOlder()
	}
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
