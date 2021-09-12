package main

import "github.com/rivo/tview"

type NotificationView struct {
	app          *App
	list         *tview.List
	iconList     *tview.List
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

	nv.iconList = tview.NewList()
	nv.iconList.SetMainTextColor(app.Config.Style.Text)
	nv.iconList.SetBackgroundColor(app.Config.Style.Background)
	nv.iconList.SetSelectedTextColor(app.Config.Style.StatusBarViewText)
	nv.iconList.SetSelectedBackgroundColor(app.Config.Style.StatusBarViewBackground)
	nv.iconList.ShowSecondaryText(false)
	nv.iconList.SetHighlightFullLine(true)

	nv.feed = NewNotificationFeed(app, true)

	return nv
}

func (n *NotificationView) SetList(items <-chan ListItem) {
	n.list.Clear()
	n.iconList.Clear()
	for s := range items {
		n.list.AddItem(s.Text, "", 0, nil)
		n.iconList.AddItem(s.Icons, "", 0, nil)
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
			n.iconList.SetCurrentItem(newIndex)
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
			n.iconList.SetCurrentItem(index)
			n.loadingOlder = false
		})
	}()
}
