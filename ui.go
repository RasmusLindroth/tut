package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type FocusAt uint

const (
	LeftPaneFocus FocusAt = iota
	RightPaneFocus
	CmdBarFocus
	MessageFocus
	LinkOverlayFocus
	AuthOverlayFocus
)

type UI struct {
	app         *App
	Focus       FocusAt
	Top         *tview.TextView
	StatusText  *StatusText
	TootList    *TootList
	MessageBox  *MessageBox
	CmdBar      *CmdBar
	Status      *tview.TextView
	Pages       *tview.Pages
	AuthOverlay *AuthOverlay
	Timeline    TimelineType
}

func (ui *UI) SetFocus(f FocusAt) {
	ui.Focus = f
	switch f {
	case RightPaneFocus:
		ui.Status.SetText("-- VIEW --")
		ui.app.App.SetFocus(ui.StatusText.View)
	case CmdBarFocus:
		ui.Status.SetText("-- CMD --")
		ui.app.App.SetFocus(ui.CmdBar.View)
	case MessageFocus:
		ui.Status.SetText("-- TOOT --")
		ui.Pages.ShowPage("toot")
		ui.app.App.SetFocus(ui.MessageBox.View)
	case LinkOverlayFocus:
		ui.Status.SetText("-- LINK --")
		ui.Pages.ShowPage("links")
		ui.app.App.SetFocus(ui.StatusText.LinkOverlay.View)
	case AuthOverlayFocus:
		ui.Status.SetText("-- LOGIN --")
		ui.Pages.ShowPage("login")
		ui.app.App.SetFocus(ui.StatusText.app.UI.AuthOverlay.View)
	default:
		ui.Status.SetText("-- LIST --")
		ui.app.App.SetFocus(ui.Pages)
		ui.Pages.HidePage("toot")
		ui.Pages.HidePage("links")
		ui.Pages.HidePage("login")
	}
}

func (ui *UI) SetTimeline(tl TimelineType) {
	ui.Timeline = tl
	statuses, err := ui.app.API.GetStatuses(tl)
	if err != nil {
		log.Fatalln(err)
	}
	ui.TootList.SetFeedStatuses(statuses)
}

func (ui *UI) ShowThread() {
	status, err := ui.TootList.GetStatus(ui.TootList.Index)
	if err != nil {
		log.Fatalln(err)
	}

	thread, index, err := ui.app.API.GetThread(status)
	if err != nil {
		log.Fatalln(err)
	}

	ui.TootList.SetThread(thread, index)
	ui.TootList.FocusThread()
	ui.SetFocus(LeftPaneFocus)
	ui.TootList.Draw()
}

func (ui *UI) ShowSensetive() {
	ui.StatusText.ShowTootOptions(ui.TootList.GetIndex(), true)
}

func (ui *UI) NewToot() {
	ui.app.App.SetFocus(ui.MessageBox.View)
	ui.MessageBox.NewToot()
	ui.MessageBox.Draw()
	ui.SetFocus(MessageFocus)
}

func (ui *UI) Reply() {
	status, err := ui.TootList.GetStatus(ui.TootList.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	if status.Reblog != nil {
		status = status.Reblog
	}
	ui.MessageBox.Reply(status)
	ui.MessageBox.Draw()
	ui.SetFocus(MessageFocus)
}

func (ui *UI) ShowLinks() {
	ui.StatusText.LinkOverlay.Draw()
	ui.SetFocus(LinkOverlayFocus)
}

func (ui *UI) OpenMedia() {
	status, err := ui.TootList.GetStatus(ui.TootList.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	if status.Reblog != nil {
		status = status.Reblog
	}

	if len(status.MediaAttachments) == 0 {
		//TODO show error that there's no media
		return
	}

	mediaGroup := make(map[string][]mastodon.Attachment)
	for _, m := range status.MediaAttachments {
		mediaGroup[m.Type] = append(mediaGroup[m.Type], m)
	}

	for key := range mediaGroup {
		var files []string
		for _, m := range mediaGroup[key] {
			f, err := downloadFile(m.URL)
			if err != nil {
				continue
			}
			files = append(files, f)
		}
		go openMedia(files)
	}
}

func (ui *UI) LoggedIn() {
	ui.SetFocus(LeftPaneFocus)
	fmt.Fprint(ui.Top, "tut\n")

	me, err := ui.app.API.Client.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	ui.app.Me = me

	ui.SetTimeline(ui.Timeline)
}

func (ui *UI) LoadNewer(status *mastodon.Status) int {
	statuses, _, err := ui.app.API.GetStatusesNewer(ui.Timeline, status)
	if err != nil {
		log.Fatalln(err)
	}
	ui.app.UI.TootList.PrependFeedStatuses(statuses)
	return len(statuses)
}

func (ui *UI) LoadOlder(status *mastodon.Status) int {
	statuses, _, err := ui.app.API.GetStatusesOlder(ui.Timeline, status)
	if err != nil {
		log.Fatalln(err)
	}
	ui.app.UI.TootList.AppendFeedStatuses(statuses)
	return len(statuses)
}

func (ui *UI) FavoriteEvent() {
	status, err := ui.TootList.GetStatus(ui.TootList.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	if status.Favourited == true {
		err = ui.app.API.Unfavorite(status)
	} else {
		err = ui.app.API.Favorite(status)
	}
}

func (ui *UI) BoostEvent() {
	status, err := ui.TootList.GetStatus(ui.TootList.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	if status.Reblogged == true {
		err = ui.app.API.Unboost(status)
	} else {
		err = ui.app.API.Boost(status)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func (ui *UI) DeleteStatus() {
	status, err := ui.TootList.GetStatus(ui.TootList.GetIndex())
	if err != nil {
		log.Fatalln(err)
	}
	err = ui.app.API.DeleteStatus(status)
	if err != nil {
		log.Fatalln(err)
	}
}

func clearContent(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	for cx := x; cx < width+x; cx++ {
		for cy := y; cy < height+y; cy++ {
			screen.SetContent(cx, cy, ' ', nil, tcell.StyleDefault.Background(tcell.ColorDefault))
		}
	}
	y2 := y + height
	for cx := x + 1; cx < width+x; cx++ {
		screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		screen.SetContent(cx, y2, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}
	x2 := x + width
	for cy := y + 1; cy < height+y; cy++ {
		screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		screen.SetContent(x2, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}
	screen.SetContent(x, y, tview.BoxDrawingsLightDownAndRight, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	screen.SetContent(x, y+height, tview.BoxDrawingsLightUpAndRight, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	screen.SetContent(x+width, y, tview.BoxDrawingsLightDownAndLeft, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	screen.SetContent(x+width, y+height, tview.BoxDrawingsLightUpAndLeft, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	return x + 1, y + 1, width - 1, height - 1
}
