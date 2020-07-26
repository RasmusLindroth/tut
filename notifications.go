package main

import "github.com/rivo/tview"

type NotificationView struct {
	app          *App
	list         *tview.List
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
	nv.list.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	nv.list.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	nv.list.ShowSecondaryText(false)
	nv.list.SetHighlightFullLine(true)

	return nv
}
