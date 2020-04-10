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
	MessageAttachmentFocus
	LinkOverlayFocus
	AuthOverlayFocus
)

func NewUI(app *App) *UI {
	ui := &UI{
		app:  app,
		Root: tview.NewApplication(),
	}

	return ui
}

func (ui *UI) Init() {
	ui.Top = NewTop(ui.app)
	ui.Pages = tview.NewPages()
	ui.Timeline = ui.app.Config.General.StartTimeline
	ui.CmdBar = NewCmdBar(ui.app)
	ui.StatusBar = NewStatusBar(ui.app)
	ui.MessageBox = NewMessageBox(ui.app)
	ui.LinkOverlay = NewLinkOverlay(ui.app)
	ui.AuthOverlay = NewAuthOverlay(ui.app)
	ui.MediaOverlay = NewMediaOverlay(ui.app)

	ui.Pages.SetBackgroundColor(ui.app.Config.Style.Background)

	verticalLine := tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background)
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(ui.app.Config.Style.Subtle))
		}
		return 0, 0, 0, 0
	})
	ui.SetTopText("")
	ui.Pages.AddPage("main",
		tview.NewFlex().
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(ui.Top.Text, 1, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background), 0, 2, false).
					AddItem(verticalLine, 1, 0, false).
					AddItem(tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background), 1, 0, false).
					AddItem(tview.NewTextView().SetBackgroundColor(ui.app.Config.Style.Background),
						0, 4, false),
					0, 1, false).
				AddItem(ui.StatusBar.Text, 1, 1, false).
				AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)

	ui.Pages.AddPage("toot", tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.MessageBox.Flex.SetDirection(tview.FlexRow).
				AddItem(ui.MessageBox.View, 0, 9, true).
				AddItem(ui.MessageBox.Controls, 2, 1, false), 0, 8, false).
			AddItem(nil, 0, 1, false), 0, 8, true).
		AddItem(nil, 0, 1, false), true, false)

	ui.Pages.AddPage("links", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.LinkOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(ui.LinkOverlay.List, 0, 10, true).
				AddItem(ui.LinkOverlay.TextBottom, 1, 1, true), 0, 8, false).
			AddItem(nil, 0, 1, false), 0, 8, true).
		AddItem(nil, 0, 1, false), true, false)
	ui.Pages.AddPage("login",
		tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(ui.AuthOverlay.Flex.SetDirection(tview.FlexRow).
					AddItem(ui.AuthOverlay.Text, 4, 1, false).
					AddItem(ui.AuthOverlay.Input, 0, 9, true), 0, 9, true).
				AddItem(nil, 0, 1, false), 0, 6, true).
			AddItem(nil, 0, 1, false),
		true, false)

	ui.Pages.AddPage("media", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.MediaOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(ui.MediaOverlay.TextTop, 1, 1, true).
				AddItem(ui.MediaOverlay.FileList, 0, 10, true).
				AddItem(ui.MediaOverlay.TextBottom, 1, 1, true).
				AddItem(ui.MediaOverlay.InputField.View, 2, 1, false), 0, 8, false).
			AddItem(nil, 0, 1, false), 0, 8, true).
		AddItem(nil, 0, 1, false), true, false)

	ui.Root.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})
}

type UI struct {
	app          *App
	Root         *tview.Application
	Focus        FocusAt
	Top          *Top
	MessageBox   *MessageBox
	CmdBar       *CmdBar
	StatusBar    *StatusBar
	Pages        *tview.Pages
	LinkOverlay  *LinkOverlay
	AuthOverlay  *AuthOverlay
	MediaOverlay *MediaView
	Timeline     TimelineType
	StatusView   *StatusView
}

func (ui *UI) FocusAt(p tview.Primitive, s string) {
	if p == nil {
		ui.Root.SetFocus(ui.Pages)
	} else {
		ui.Root.SetFocus(p)
	}
	if s != "" {
		ui.StatusBar.SetText(s)
	}
}

func (ui *UI) SetFocus(f FocusAt) {
	ui.Focus = f
	switch f {
	case RightPaneFocus:
		ui.StatusBar.SetText("-- VIEW --")
	case CmdBarFocus:
		ui.StatusBar.SetText("-- CMD --")
		ui.Root.SetFocus(ui.CmdBar.Input)
	case MessageFocus:
		ui.StatusBar.SetText("-- TOOT --")
		ui.Pages.ShowPage("toot")
		ui.Pages.HidePage("media")
		ui.Root.SetFocus(ui.MessageBox.View)
	case MessageAttachmentFocus:
		ui.Pages.ShowPage("media")
	case LinkOverlayFocus:
		ui.StatusBar.SetText("-- LINK --")
		ui.Pages.ShowPage("links")
		ui.Root.SetFocus(ui.LinkOverlay.List)
	case AuthOverlayFocus:
		ui.StatusBar.SetText("-- LOGIN --")
		ui.Pages.ShowPage("login")
		ui.Root.SetFocus(ui.AuthOverlay.Input)
	default:
		ui.StatusBar.SetText("-- LIST --")
		ui.Root.SetFocus(ui.Pages)
		ui.Pages.HidePage("toot")
		ui.Pages.HidePage("media")
		ui.Pages.HidePage("links")
		ui.Pages.HidePage("login")
	}
}

func (ui *UI) NewToot() {
	ui.Root.SetFocus(ui.MessageBox.View)
	ui.MediaOverlay.Reset()
	ui.MessageBox.NewToot()
	ui.MessageBox.Draw()
	ui.SetFocus(MessageFocus)
}

func (ui *UI) Reply(status *mastodon.Status) {
	if status.Reblog != nil {
		status = status.Reblog
	}
	ui.MediaOverlay.Reset()
	ui.MessageBox.Reply(status)
	ui.MessageBox.Draw()
	ui.SetFocus(MessageFocus)
}

func (ui *UI) ShowLinks() {
	ui.SetFocus(LinkOverlayFocus)
}

func (ui *UI) OpenMedia(status *mastodon.Status) {
	if status.Reblog != nil {
		status = status.Reblog
	}

	if len(status.MediaAttachments) == 0 {
		return
	}

	mediaGroup := make(map[string][]mastodon.Attachment)
	for _, m := range status.MediaAttachments {
		mediaGroup[m.Type] = append(mediaGroup[m.Type], m)
	}

	for key := range mediaGroup {
		var files []string
		for _, m := range mediaGroup[key] {
			//'image', 'video', 'gifv', 'audio' or 'unknown'
			f, err := downloadFile(m.URL)
			if err != nil {
				continue
			}
			files = append(files, f)
		}
		go openMediaType(ui.app.Config.Media, files, key)
	}
}

func (ui *UI) SetTopText(s string) {
	if s == "" {
		ui.Top.Text.SetText("tut")
	} else {
		ui.Top.Text.SetText(fmt.Sprintf("tut - %s", s))
	}
}

func (ui *UI) LoggedIn() {
	ui.StatusView = NewStatusView(ui.app, ui.Timeline)

	verticalLine := tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background)
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(ui.app.Config.Style.Subtle))
		}
		return 0, 0, 0, 0
	})
	ui.Pages.RemovePage("main")
	ui.Pages.AddPage("main",
		tview.NewFlex().
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(ui.Top.Text, 1, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(ui.StatusView.GetLeftView(), 0, 2, false).
					AddItem(verticalLine, 1, 0, false).
					AddItem(tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background), 1, 0, false).
					AddItem(ui.StatusView.GetRightView(),
						0, 4, false),
					0, 1, false).
				AddItem(ui.StatusBar.Text, 1, 1, false).
				AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)
	ui.Pages.SendToBack("main")

	ui.SetFocus(LeftPaneFocus)

	me, err := ui.app.API.Client.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	ui.app.Me = me
	ui.StatusView.AddFeed(
		NewTimelineFeed(ui.app, ui.Timeline),
	)
}

func (conf *Config) ClearContent(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	for cx := x; cx < width+x; cx++ {
		for cy := y; cy < height+y; cy++ {
			screen.SetContent(cx, cy, ' ', nil, tcell.StyleDefault.Background(conf.Style.Background))
		}
	}
	y2 := y + height
	for cx := x + 1; cx < width+x; cx++ {
		screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
		screen.SetContent(cx, y2, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	}
	x2 := x + width
	for cy := y + 1; cy < height+y; cy++ {
		screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
		screen.SetContent(x2, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	}
	screen.SetContent(x, y, tview.BoxDrawingsLightDownAndRight, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	screen.SetContent(x, y+height, tview.BoxDrawingsLightUpAndRight, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	screen.SetContent(x+width, y, tview.BoxDrawingsLightDownAndLeft, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	screen.SetContent(x+width, y+height, tview.BoxDrawingsLightUpAndLeft, nil, tcell.StyleDefault.Foreground(conf.Style.Subtle))
	return x + 1, y + 1, width - 1, height - 1
}
