package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type FocusAt uint

const (
	LeftPaneFocus FocusAt = iota
	RightPaneFocus
	NotificationPaneFocus
	CmdBarFocus
	MessageFocus
	MessageAttachmentFocus
	LinkOverlayFocus
	VisibilityOverlayFocus
	AuthOverlayFocus
	UserSelectFocus
)

func NewUI(app *App) *UI {
	ui := &UI{
		app:  app,
		Root: tview.NewApplication(),
	}

	return ui
}

func (ui *UI) Init() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    ui.app.Config.Style.StatusBarViewText, // main text color, selected text
		ContrastBackgroundColor:     ui.app.Config.Style.Background,
		MoreContrastBackgroundColor: ui.app.Config.Style.StatusBarBackground, //background color
		BorderColor:                 ui.app.Config.Style.Subtle,
		TitleColor:                  ui.app.Config.Style.Text,
		GraphicsColor:               ui.app.Config.Style.Text,
		PrimaryTextColor:            ui.app.Config.Style.StatusBarViewBackground, //backround color selected
		SecondaryTextColor:          ui.app.Config.Style.Text,
		TertiaryTextColor:           ui.app.Config.Style.Text,
		InverseTextColor:            ui.app.Config.Style.Text,
		ContrastSecondaryTextColor:  ui.app.Config.Style.Text,
	}
	ui.Top = NewTop(ui.app)
	ui.Pages = tview.NewPages()
	ui.Timeline = ui.app.Config.General.StartTimeline
	ui.CmdBar = NewCmdBar(ui.app)
	ui.StatusBar = NewStatusBar(ui.app)
	ui.MessageBox = NewMessageBox(ui.app)
	ui.LinkOverlay = NewLinkOverlay(ui.app)
	ui.VisibilityOverlay = NewVisibilityOverlay(ui.app)
	ui.AuthOverlay = NewAuthOverlay(ui.app)
	ui.UserSelectOverlay = NewUserSelectOverlay(ui.app)
	ui.MediaOverlay = NewMediaOverlay(ui.app)

	ui.Pages.SetBackgroundColor(ui.app.Config.Style.Background)

	verticalLine := tview.NewBox()
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		var s tcell.Style
		s = s.Background(ui.app.Config.Style.Background).Foreground(ui.app.Config.Style.Subtle)
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, s)
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
	ui.Pages.AddPage("visibility", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.VisibilityOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(ui.VisibilityOverlay.List, 0, 10, true).
				AddItem(ui.VisibilityOverlay.TextBottom, 1, 1, true), 0, 8, false).
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
	ui.Pages.AddPage("userselect",
		tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(ui.UserSelectOverlay.Flex.SetDirection(tview.FlexRow).
					AddItem(ui.UserSelectOverlay.Text, 2, 1, false).
					AddItem(ui.UserSelectOverlay.List, 0, 9, true), 0, 9, true).
				AddItem(nil, 0, 1, false), 0, 6, true).
			AddItem(nil, 0, 1, false),
		true, false)
	ui.Pages.AddPage("media", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.MediaOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(ui.MediaOverlay.TextTop, 2, 1, true).
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
	app               *App
	Root              *tview.Application
	Focus             FocusAt
	Top               *Top
	MessageBox        *MessageBox
	CmdBar            *CmdBar
	StatusBar         *StatusBar
	Pages             *tview.Pages
	LinkOverlay       *LinkOverlay
	VisibilityOverlay *VisibilityOverlay
	AuthOverlay       *AuthOverlay
	UserSelectOverlay *UserSelectOverlay
	MediaOverlay      *MediaView
	Timeline          TimelineType
	StatusView        *StatusView
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
		ui.FocusAt(ui.StatusView.text, "-- VIEW --")
	case CmdBarFocus:
		ui.FocusAt(ui.CmdBar.Input, "-- CMD --")
	case MessageFocus:
		ui.MessageBox.Draw()
		ui.Pages.ShowPage("toot")
		ui.Pages.HidePage("media")
		ui.Pages.HidePage("visibility")
		ui.Root.SetFocus(ui.MessageBox.View)
		ui.FocusAt(ui.MessageBox.View, "-- TOOT --")
	case MessageAttachmentFocus:
		ui.Pages.ShowPage("media")
	case LinkOverlayFocus:
		ui.Pages.ShowPage("links")
		ui.Root.SetFocus(ui.LinkOverlay.List)
		ui.FocusAt(ui.LinkOverlay.List, "-- LINK --")
	case VisibilityOverlayFocus:
		ui.VisibilityOverlay.Show()
		ui.Pages.ShowPage("visibility")
		ui.Root.SetFocus(ui.VisibilityOverlay.List)
		ui.FocusAt(ui.VisibilityOverlay.List, "-- VISIBILITY --")
	case AuthOverlayFocus:
		ui.Pages.ShowPage("login")
		ui.FocusAt(ui.AuthOverlay.Input, "-- LOGIN --")
	case UserSelectFocus:
		ui.UserSelectOverlay.Draw()
		ui.Pages.ShowPage("userselect")
		ui.FocusAt(ui.UserSelectOverlay.List, "-- SELECT USER --")
	case NotificationPaneFocus:
		ui.Pages.SwitchToPage("main")
		ui.FocusAt(nil, "-- NOTIFICATIONS --")

		ui.StatusView.notificationView.list.SetSelectedBackgroundColor(
			ui.app.Config.Style.ListSelectedBackground,
		)
		ui.StatusView.notificationView.list.SetSelectedTextColor(
			ui.app.Config.Style.ListSelectedText,
		)

		ui.StatusView.list.SetSelectedBackgroundColor(
			ui.app.Config.Style.StatusBarViewBackground,
		)
		ui.StatusView.list.SetSelectedTextColor(
			ui.app.Config.Style.StatusBarViewText,
		)
		ui.StatusView.notificationView.iconList.SetSelectedBackgroundColor(
			ui.app.Config.Style.ListSelectedBackground,
		)
		ui.StatusView.notificationView.iconList.SetSelectedTextColor(
			ui.app.Config.Style.ListSelectedText,
		)

		ui.StatusView.iconList.SetSelectedBackgroundColor(
			ui.app.Config.Style.StatusBarViewBackground,
		)
		ui.StatusView.iconList.SetSelectedTextColor(
			ui.app.Config.Style.StatusBarViewText,
		)
	default:
		ui.app.UI.StatusBar.Text.SetBackgroundColor(
			ui.app.Config.Style.StatusBarBackground,
		)
		ui.app.UI.StatusBar.Text.SetTextColor(
			ui.app.Config.Style.StatusBarText,
		)
		ui.StatusView.list.SetSelectedBackgroundColor(
			ui.app.Config.Style.ListSelectedBackground,
		)
		ui.StatusView.list.SetSelectedTextColor(
			ui.app.Config.Style.ListSelectedText,
		)
		ui.StatusView.iconList.SetSelectedBackgroundColor(
			ui.app.Config.Style.ListSelectedBackground,
		)
		ui.StatusView.iconList.SetSelectedTextColor(
			ui.app.Config.Style.ListSelectedText,
		)

		if ui.app.Config.General.NotificationFeed {
			ui.StatusView.notificationView.list.SetSelectedBackgroundColor(
				ui.app.Config.Style.StatusBarViewBackground,
			)
			ui.StatusView.notificationView.list.SetSelectedTextColor(
				ui.app.Config.Style.StatusBarViewText,
			)
			ui.StatusView.notificationView.iconList.SetSelectedBackgroundColor(
				ui.app.Config.Style.StatusBarViewBackground,
			)
			ui.StatusView.notificationView.iconList.SetSelectedTextColor(
				ui.app.Config.Style.StatusBarViewText,
			)
		}
		ui.Pages.SwitchToPage("main")
		ui.FocusAt(nil, "-- LIST --")
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
		ui.app.FileList = append(ui.app.FileList, files...)
	}
}

func (ui *UI) SetTopText(s string) {
	if s == "" {
		ui.Top.Text.SetText("tut")
	} else {
		ui.Top.Text.SetText(fmt.Sprintf("tut - %s - %s", s, ui.app.FullUsername))
	}
}

func (ui *UI) LoggedIn() {
	ui.StatusView = NewStatusView(ui.app, ui.Timeline)

	verticalLine := tview.NewBox()
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		var s tcell.Style
		s = s.Background(ui.app.Config.Style.Background).Foreground(ui.app.Config.Style.Subtle)
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, s)
		}
		return 0, 0, 0, 0
	})
	horizontalLine := tview.NewBox()
	horizontalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		var s tcell.Style
		s = s.Background(ui.app.Config.Style.Background).Foreground(ui.app.Config.Style.Subtle)
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, s)
		}
		return 0, 0, 0, 0
	})

	ui.Pages.RemovePage("main")
	mainText := tview.NewTextView()
	mainText.SetBackgroundColor(ui.app.Config.Style.Background)
	mainText.SetTextColor(ui.app.Config.Style.Subtle)
	mainText.SetText("")
	mainText.SetTextAlign(tview.AlignCenter)

	notificationText := tview.NewTextView()
	notificationText.SetBackgroundColor(ui.app.Config.Style.Background)
	notificationText.SetTextColor(ui.app.Config.Style.Subtle)
	notificationText.SetText("[N]otifications")
	notificationText.SetTextAlign(tview.AlignCenter)

	var listViewRow *tview.Flex
	var listViewColumn *tview.Flex
	lp := ui.app.Config.General.ListProportion
	cp := ui.app.Config.General.ContentProportion
	nf := 1
	if ui.app.Config.General.HideNotificationText {
		nf = 0
	}

	if ui.app.Config.General.NotificationFeed {
		listViewRow = tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.StatusView.GetLeftView(), 0, 1, false).
			AddItem(notificationText, 1, 0, false).
			AddItem(ui.StatusView.GetNotificationView(), 0, 1, false), 0, 1, false)

		listViewColumn = tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(mainText, nf, 0, false).
				AddItem(ui.StatusView.GetLeftView(), 0, 1, false), 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(notificationText, nf, 0, false).
				AddItem(ui.StatusView.GetNotificationView(), 0, 1, false), 0, 1, false), 0, 1, false)
	} else {
		listViewRow = tview.NewFlex().AddItem(ui.StatusView.GetLeftView(), 0, 1, false)
		listViewColumn = tview.NewFlex().AddItem(ui.StatusView.GetLeftView(), 0, 1, false)
	}

	var listViewChoice *tview.Flex
	if ui.app.Config.General.ListSplit == ListRow {
		listViewChoice = listViewRow
	} else {
		listViewChoice = listViewColumn
	}

	switch ui.app.Config.General.ListPlacement {
	case ListPlacementLeft:
		ui.Pages.AddPage("main",
			tview.NewFlex().
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(ui.Top.Text, 1, 0, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
						AddItem(listViewChoice, 0, lp, false).
						AddItem(verticalLine, 1, 0, false).
						AddItem(tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background), 1, 0, false).
						AddItem(ui.StatusView.GetRightView(), 0, cp, false),
						0, 1, false).
					AddItem(ui.StatusBar.Text, 1, 1, false).
					AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)
	case ListPlacementRight:
		ui.Pages.AddPage("main",
			tview.NewFlex().
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(ui.Top.Text, 1, 0, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
						AddItem(tview.NewBox().SetBackgroundColor(ui.app.Config.Style.Background), 1, 0, false).
						AddItem(ui.StatusView.GetRightView(), 0, cp, false).
						AddItem(verticalLine, 1, 0, false).
						AddItem(listViewChoice, 0, 1, false), 0, lp, false).
					AddItem(ui.StatusBar.Text, 1, 1, false).
					AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)
	case ListPlacementTop:
		ui.Pages.AddPage("main",
			tview.NewFlex().
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(ui.Top.Text, 1, 0, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(listViewChoice, 0, lp, false).
						AddItem(horizontalLine, 1, 0, false).
						AddItem(ui.StatusView.GetRightView(), 0, cp, false),
						0, 1, false).
					AddItem(ui.StatusBar.Text, 1, 1, false).
					AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)
	case ListPlacementBottom:
		ui.Pages.AddPage("main",
			tview.NewFlex().
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(ui.Top.Text, 1, 0, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(ui.StatusView.GetRightView(), 0, cp, false).
						AddItem(horizontalLine, 1, 0, false).
						AddItem(listViewChoice, 0, lp, false),
						0, 1, false).
					AddItem(ui.StatusBar.Text, 1, 1, false).
					AddItem(ui.CmdBar.Input, 1, 0, false), 0, 1, false), true, true)
	}
	ui.Pages.SendToBack("main")

	ui.SetFocus(LeftPaneFocus)

	me, err := ui.app.API.Client.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	ui.app.Me = me
	ui.StatusView.AddFeed(
		NewTimelineFeed(ui.app, ui.Timeline, nil),
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
