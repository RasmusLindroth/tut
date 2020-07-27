package main

import "github.com/rivo/tview"

type NotificationView struct {
	app          *App
	list         *tview.List
	feed         Feed
	loadingNewer bool
	loadingOlder bool
}

func NewNotificationView(app *App) *NotificationView {
	nv := &NotificationView{
		app:          app,
		loadingNewer: false,
		loadingOlder: false,
	}

	nv.list = tview.NewList()
	nv.list.SetMainTextColor(app.Config.Style.Text)
	nv.list.SetBackgroundColor(app.Config.Style.Background)
	nv.list.SetSelectedTextColor(app.Config.Style.StatusBarViewText)
	nv.list.SetSelectedBackgroundColor(app.Config.Style.StatusBarViewBackground)
	nv.list.ShowSecondaryText(false)
	nv.list.SetHighlightFullLine(true)

	nv.feed = NewNotificationFeed(app, true)

	return nv
}

func (n *NotificationView) SetList(items <-chan string) {
	n.list.Clear()
	for s := range items {
		n.list.AddItem(s, "", 0, nil)
	}
}

func (n *NotificationView) loadNewer() {
	if n.loadingNewer {
		return
	}
	n.loadingNewer = true
	go func() {
		new := n.feed.LoadNewer()
		if new == 0 {
			n.loadingNewer = false
			return
		}
		n.app.UI.Root.QueueUpdateDraw(func() {
			index := n.list.GetCurrentItem()
			n.feed.DrawList()
			newIndex := index + new

			n.list.SetCurrentItem(newIndex)
			n.loadingNewer = false
		})
	}()
}

func (n *NotificationView) loadOlder() {
	if n.loadingOlder {
		return
	}
	n.loadingOlder = true
	go func() {
		new := n.feed.LoadOlder()
		if new == 0 {
			n.loadingOlder = false
			return
		}
		n.app.UI.Root.QueueUpdateDraw(func() {
			index := n.list.GetCurrentItem()
			n.feed.DrawList()
			n.list.SetCurrentItem(index)
			n.loadingOlder = false
		})
	}()
}
