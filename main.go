package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	config := Config{
		Style: StyleConfig{
			Background:             tcell.ColorDefault,
			Text:                   tcell.ColorWhite,
			Subtle:                 tcell.ColorGray,
			WarningText:            tcell.NewRGBColor(249, 38, 114),
			TextSpecial1:           tcell.NewRGBColor(174, 129, 255),
			TextSpecial2:           tcell.NewRGBColor(166, 226, 46),
			TopBarBackground:       tcell.NewRGBColor(249, 38, 114),
			TopBarText:             tcell.ColorWhite,
			StatusBarBackground:    tcell.NewRGBColor(249, 38, 114),
			StatusBarText:          tcell.ColorWhite,
			ListSelectedBackground: tcell.NewRGBColor(249, 38, 114),
			ListSelectedText:       tcell.ColorWhite,
		},
	}

	err := CreateConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	path, exists, err := CheckConfig("accounts.toml")
	if err != nil {
		log.Fatalln(err)
	}
	app := &App{
		App:         tview.NewApplication(),
		API:         &API{},
		HaveAccount: false,
		Config:      &config,
	}

	clearContent := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for cx := x; cx < width+x; cx++ {
			for cy := y; cy < height+y; cy++ {
				screen.SetContent(cx, cy, ' ', nil, tcell.StyleDefault.Background(app.Config.Style.Background))
			}
		}
		y2 := y + height
		for cx := x + 1; cx < width+x; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
			screen.SetContent(cx, y2, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		}
		x2 := x + width
		for cy := y + 1; cy < height+y; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
			screen.SetContent(x2, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		}
		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndRight, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		screen.SetContent(x, y+height, tview.BoxDrawingsLightUpAndRight, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		screen.SetContent(x+width, y, tview.BoxDrawingsLightDownAndLeft, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		screen.SetContent(x+width, y+height, tview.BoxDrawingsLightUpAndLeft, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		return x + 1, y + 1, width - 1, height - 1
	}
	if exists {
		accounts, err := GetAccounts(path)
		if err != nil {
			log.Fatalln(err)
		}
		if len(accounts.Accounts) > 0 {
			a := accounts.Accounts[0]
			client, err := a.Login()
			if err == nil {
				app.API.Client = client
				app.HaveAccount = true
			}
		}
	}

	app.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	tview.Borders.HorizontalFocus = tview.BoxDrawingsLightHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsLightVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsLightDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsLightDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsLightUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsLightUpAndLeft

	top := tview.NewTextView()
	top.SetBackgroundColor(app.Config.Style.TopBarBackground)
	top.SetTextColor(app.Config.Style.TopBarText)

	app.UI = &UI{app: app, Top: top, Timeline: TimelineHome}

	app.UI.TootList = NewTootList(app, tview.NewList())
	app.UI.TootList.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.TootList.View.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	app.UI.TootList.View.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	app.UI.TootList.View.ShowSecondaryText(false)
	app.UI.TootList.View.SetHighlightFullLine(true)

	app.UI.TootList.View.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
		if app.HaveAccount {
			app.UI.StatusText.ShowToot(index)
		}
	})

	app.UI.StatusText = NewStatusText(app, tview.NewTextView(),
		NewControls(app, tview.NewTextView()), NewLinkOverlay(app),
	)
	app.UI.StatusText.View.SetWordWrap(true).SetDynamicColors(true)
	app.UI.StatusText.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.StatusText.View.SetTextColor(app.Config.Style.Text)
	app.UI.StatusText.Controls.View.SetDynamicColors(true)
	app.UI.StatusText.Controls.View.SetBackgroundColor(app.Config.Style.Background)

	app.UI.CmdBar = NewCmdBar(app,
		tview.NewInputField(),
	)
	app.UI.CmdBar.View.SetFieldBackgroundColor(app.Config.Style.Background)
	app.UI.CmdBar.View.SetFieldTextColor(app.Config.Style.Text)
	app.UI.Status = tview.NewTextView()
	app.UI.Status.SetBackgroundColor(app.Config.Style.StatusBarBackground)
	app.UI.Status.SetTextColor(app.Config.Style.StatusBarText)

	verticalLine := tview.NewBox().SetBackgroundColor(app.Config.Style.Background)
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(app.Config.Style.Subtle))
		}
		return 0, 0, 0, 0
	})

	app.UI.Pages = tview.NewPages()
	app.UI.Pages.SetBackgroundColor(app.Config.Style.Background)
	app.UI.Pages.AddPage("main",
		tview.NewFlex().
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(top, 1, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(app.UI.TootList.View, 0, 2, false).
					AddItem(verticalLine, 1, 0, false).
					AddItem(tview.NewBox().SetBackgroundColor(app.Config.Style.Background), 1, 0, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(app.UI.StatusText.View, 0, 9, false).
						AddItem(app.UI.StatusText.Controls.View, 1, 0, false),
						0, 4, false),
					0, 1, false).
				AddItem(app.UI.Status, 1, 1, false).
				AddItem(app.UI.CmdBar.View, 1, 0, false), 0, 1, false), true, true)

	flLinks := tview.NewFlex()
	flLinks.SetDrawFunc(clearContent)
	flToot := tview.NewFlex()
	flToot.SetDrawFunc(clearContent)
	modal := func(fl *tview.Flex, p tview.Primitive, c tview.Primitive) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(fl.SetDirection(tview.FlexRow).
					AddItem(p, 0, 9, true).
					AddItem(c, 2, 1, false), 0, 8, false).
				AddItem(nil, 0, 1, false), 0, 8, true).
			AddItem(nil, 0, 1, false)
	}

	authModal := func(f *tview.Flex, p tview.Primitive, c tview.Primitive) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(f.SetDirection(tview.FlexRow).
					AddItem(c, 4, 1, false).
					AddItem(p, 0, 9, true), 0, 9, true).
				AddItem(nil, 0, 1, false), 0, 6, true).
			AddItem(nil, 0, 1, false)
	}

	app.UI.MessageBox = NewMessageBox(app, tview.NewTextView(),
		NewControls(app, tview.NewTextView()),
	)

	app.UI.MessageBox.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MessageBox.View.SetTextColor(app.Config.Style.Text)
	app.UI.MessageBox.View.SetDynamicColors(true)
	app.UI.MessageBox.Controls.View.SetDynamicColors(true)
	app.UI.MessageBox.Controls.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MessageBox.Controls.View.SetTextColor(app.Config.Style.Text)
	app.UI.Pages.AddPage("toot",
		modal(flToot, app.UI.MessageBox.View, app.UI.MessageBox.Controls.View), true, false)

	app.UI.Pages.AddPage("links", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(app.UI.StatusText.LinkOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(app.UI.StatusText.LinkOverlay.List, 0, 10, true).
				AddItem(app.UI.StatusText.LinkOverlay.TextBottom, 1, 1, true), 0, 8, false).
			AddItem(nil, 0, 1, false), 0, 8, true).
		AddItem(nil, 0, 1, false), true, false)

	app.UI.StatusText.LinkOverlay.Flex.SetDrawFunc(clearContent)
	app.UI.StatusText.LinkOverlay.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	app.UI.StatusText.LinkOverlay.List.SetBackgroundColor(app.Config.Style.Background)
	app.UI.StatusText.LinkOverlay.List.SetMainTextColor(app.Config.Style.Text)
	app.UI.StatusText.LinkOverlay.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	app.UI.StatusText.LinkOverlay.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	app.UI.StatusText.LinkOverlay.List.ShowSecondaryText(false)
	app.UI.StatusText.LinkOverlay.List.SetHighlightFullLine(true)

	app.UI.AuthOverlay = NewAuthoverlay(app, tview.NewFlex(), tview.NewInputField(),
		NewControls(app, tview.NewTextView()))

	app.UI.Pages.AddPage("login",
		authModal(app.UI.AuthOverlay.Flex, app.UI.AuthOverlay.View, app.UI.AuthOverlay.Controls.View), true, false)
	app.UI.AuthOverlay.Flex.SetDrawFunc(clearContent)
	app.UI.AuthOverlay.Flex.SetBackgroundColor(app.Config.Style.Background)
	app.UI.AuthOverlay.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.AuthOverlay.View.SetFieldBackgroundColor(app.Config.Style.Background)
	app.UI.AuthOverlay.View.SetFieldTextColor(app.Config.Style.Text)
	app.UI.AuthOverlay.Controls.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.AuthOverlay.Controls.View.SetTextColor(app.Config.Style.Text)
	app.UI.AuthOverlay.Draw()

	app.UI.MediaOverlay = NewMediaView(app)
	app.UI.Pages.AddPage("media", tview.NewFlex().AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(app.UI.MediaOverlay.Flex.SetDirection(tview.FlexRow).
				AddItem(app.UI.MediaOverlay.TextTop, 1, 1, true).
				AddItem(app.UI.MediaOverlay.FileList, 0, 10, true).
				AddItem(app.UI.MediaOverlay.TextBottom, 1, 1, true).
				AddItem(app.UI.MediaOverlay.InputField.View, 2, 1, false), 0, 8, false).
			AddItem(nil, 0, 1, false), 0, 8, true).
		AddItem(nil, 0, 1, false), true, false)

	app.UI.MediaOverlay.FileList.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.FileList.SetMainTextColor(app.Config.Style.Text)
	app.UI.MediaOverlay.FileList.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	app.UI.MediaOverlay.FileList.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	app.UI.MediaOverlay.FileList.ShowSecondaryText(false)
	app.UI.MediaOverlay.FileList.SetHighlightFullLine(true)

	app.UI.MediaOverlay.Flex.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.TextTop.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.TextTop.SetTextColor(app.Config.Style.Text)
	app.UI.MediaOverlay.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.TextBottom.SetTextColor(app.Config.Style.Text)
	app.UI.MediaOverlay.InputField.View.SetBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.InputField.View.SetFieldBackgroundColor(app.Config.Style.Background)
	app.UI.MediaOverlay.InputField.View.SetFieldTextColor(app.Config.Style.Text)
	app.UI.MediaOverlay.Flex.SetDrawFunc(clearContent)

	if !app.HaveAccount {
		app.UI.SetFocus(AuthOverlayFocus)
	} else {
		app.UI.LoggedIn()
	}

	app.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !app.HaveAccount {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				}
			}
			return event
		}

		if app.UI.Focus == LinkOverlayFocus {
			app.UI.StatusText.LinkOverlay.InputHandler(event)
			return nil
		}

		if app.UI.Focus == CmdBarFocus {
			switch event.Key() {
			case tcell.KeyEsc:
				app.UI.CmdBar.View.SetText("")
				app.UI.CmdBar.View.Autocomplete().Blur()
				app.UI.SetFocus(LeftPaneFocus)
				return nil
			}
			return event
		}

		if app.UI.Focus == MessageFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'p', 'P':
					app.UI.MessageBox.Post()
					return nil
				case 'e', 'E':
					app.UI.MessageBox.EditText()
					return nil
				case 'c', 'C':
					app.UI.MessageBox.EditSpoiler()
					return nil
				case 't', 'T':
					app.UI.MessageBox.ToggleSpoiler()
					return nil
				case 'm', 'M':
					app.UI.SetFocus(MessageAttachmentFocus)
					return nil
				case 'q', 'Q':
					app.UI.SetFocus(LeftPaneFocus)
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyEsc:
					app.UI.SetFocus(LeftPaneFocus)
					return nil
				}
			}
			return event
		}

		if app.UI.Focus == MessageAttachmentFocus && app.UI.MediaOverlay.Focus == MediaFocusOverview {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'j', 'J':
					app.UI.MediaOverlay.Next()
				case 'k', 'K':
					app.UI.MediaOverlay.Prev()
				case 'd', 'D':
					app.UI.MediaOverlay.Delete()
				case 'a', 'A':
					app.UI.MediaOverlay.SetFocus(MediaFocusAdd)
				case 'q', 'Q':
					app.UI.SetFocus(MessageFocus)
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyUp:
					app.UI.MediaOverlay.Prev()
				case tcell.KeyDown:
					app.UI.MediaOverlay.Next()
				case tcell.KeyEsc:
					app.UI.SetFocus(MessageFocus)
					return nil
				}
			}
			return event
		}

		if app.UI.Focus == MessageAttachmentFocus && app.UI.MediaOverlay.Focus == MediaFocusAdd {
			if event.Key() == tcell.KeyRune {
				app.UI.MediaOverlay.InputField.AddRune(event.Rune())
				return nil
			}
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyDown:
				app.UI.MediaOverlay.InputField.AutocompleteNext()
				return nil
			case tcell.KeyBacktab, tcell.KeyUp:
				app.UI.MediaOverlay.InputField.AutocompletePrev()
				return nil
			case tcell.KeyEnter:
				app.UI.MediaOverlay.InputField.CheckDone()
				return nil
			case tcell.KeyEsc:
				app.UI.MediaOverlay.SetFocus(MediaFocusOverview)
			}
			return event
		}

		if app.UI.Focus == LeftPaneFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'v', 'V':
					app.UI.SetFocus(RightPaneFocus)
					return nil
				case 'k', 'K':
					app.UI.TootList.Prev()
					return nil
				case 'j', 'J':
					app.UI.TootList.Next()
					return nil
				case 'q', 'Q':
					app.App.Stop()
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyUp:
					app.UI.TootList.Prev()
					return nil
				case tcell.KeyDown:
					app.UI.TootList.Next()
					return nil
				case tcell.KeyEsc:
					app.UI.TootList.GoBack()
					return nil
				case tcell.KeyCtrlC:
					app.App.Stop()
					return nil
				}
			}
		}

		if app.UI.Focus == RightPaneFocus {
			if event.Key() != tcell.KeyRune {
				switch event.Key() {
				case tcell.KeyEsc:
					app.UI.SetFocus(LeftPaneFocus)
				}
			}
		}

		if app.UI.Focus == LeftPaneFocus || app.UI.Focus == RightPaneFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case ':':
					app.UI.CmdBar.View.SetText(":")
					app.UI.SetFocus(CmdBarFocus)
					return nil
				case 't', 'T':
					app.UI.ShowThread()
				case 's', 'S':
					app.UI.ShowSensetive()
				case 'c', 'C':
					app.UI.NewToot()
				case 'o', 'O':
					app.UI.ShowLinks()
				case 'r', 'R':
					app.UI.Reply()
				case 'm', 'M':
					app.UI.OpenMedia()
				case 'f', 'F':
					//TODO UPDATE TOOT IN LIST
					app.UI.FavoriteEvent()
				case 'b':
					//TODO UPDATE TOOT IN LIST
					app.UI.BoostEvent()
				case 'd':
					app.UI.DeleteStatus()
				}
			}
		}

		return event
	})

	app.UI.MediaOverlay.InputField.View.SetChangedFunc(
		app.UI.MediaOverlay.InputField.HandleChanges,
	)

	words := strings.Split(":q,:quit,:timeline", ",")
	app.UI.CmdBar.View.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if currentText == "" {
			return
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) < 1 {
			entries = nil
		}
		return
	})

	app.UI.CmdBar.View.SetDoneFunc(func(key tcell.Key) {
		input := app.UI.CmdBar.GetInput()
		parts := strings.Split(input, " ")
		if len(parts) == 0 {
			return
		}
		switch parts[0] {
		case ":q":
			fallthrough
		case ":quit":
			app.App.Stop()
		case ":timeline":
			if len(parts) < 2 {
				break
			}
			switch parts[1] {
			case "local":
				app.UI.SetTimeline(TimelineLocal)
				app.UI.SetFocus(LeftPaneFocus)
				app.UI.CmdBar.ClearInput()
			case "federated":
				app.UI.SetTimeline(TimelineFederated)
				app.UI.SetFocus(LeftPaneFocus)
				app.UI.CmdBar.ClearInput()
			case "direct":
				app.UI.SetTimeline(TimelineDirect)
				app.UI.SetFocus(LeftPaneFocus)
				app.UI.CmdBar.ClearInput()
			case "home":
				app.UI.SetTimeline(TimelineHome)
				app.UI.SetFocus(LeftPaneFocus)
				app.UI.CmdBar.ClearInput()
			}
		}
	})

	app.UI.AuthOverlay.View.SetDoneFunc(func(key tcell.Key) {
		app.UI.AuthOverlay.GotInput()
	})

	if err := app.App.SetRoot(app.UI.Pages, true).Run(); err != nil {
		panic(err)
	}
}
