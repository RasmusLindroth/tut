package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
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
	top.SetBackgroundColor(tcell.ColorGreen)

	app.UI = &UI{app: app, Top: top}

	app.UI.TootList = NewTootList(app, tview.NewTable())
	app.UI.TootList.View.SetSelectedStyle(tcell.ColorWhite, tcell.ColorRed, tcell.AttrMask(0))
	app.UI.TootList.View.SetSelectable(true, false)
	app.UI.TootList.View.SetBackgroundColor(tcell.ColorDefault)

	lo := NewLinkOverlay(app, tview.NewTextView(),
		NewControls(app, tview.NewTextView()),
	)
	lo.View.SetBorderPadding(0, 0, 0, 0)

	app.UI.StatusText = NewStatusText(app, tview.NewTextView(),
		NewControls(app, tview.NewTextView()), lo,
	)
	app.UI.StatusText.View.SetWordWrap(true).SetDynamicColors(true)
	app.UI.StatusText.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.StatusText.Controls.View.SetDynamicColors(true)
	app.UI.StatusText.Controls.View.SetBackgroundColor(tcell.ColorDefault)

	app.UI.CmdBar = NewCmdBar(app,
		tview.NewInputField(),
	)
	app.UI.CmdBar.View.SetFieldBackgroundColor(tcell.ColorDefault)
	app.UI.Status = tview.NewTextView()
	app.UI.Status.SetBackgroundColor(tcell.ColorBrown)

	verticalLine := tview.NewBox().SetBackgroundColor(tcell.ColorDefault)
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		}
		return 0, 0, 0, 0
	})

	app.UI.Pages = tview.NewPages()
	app.UI.Pages.SetBackgroundColor(tcell.ColorDefault)
	app.UI.Pages.AddPage("main",
		tview.NewFlex().
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(top, 1, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(app.UI.TootList.View, 0, 2, false).
					AddItem(verticalLine, 1, 0, false).
					AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 1, 0, false).
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
	//app.UI.MessageBox.View.SetBorder(true).SetTitle("New toot")
	app.UI.MessageBox.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.MessageBox.View.SetDynamicColors(true)
	app.UI.MessageBox.Controls.View.SetDynamicColors(true)
	app.UI.MessageBox.Controls.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.Pages.AddPage("toot",
		modal(flToot, app.UI.MessageBox.View, app.UI.MessageBox.Controls.View), true, false)

	//app.UI.StatusText.LinkOverlay.View.SetBorder(true).SetTitle("Follow link")
	app.UI.StatusText.LinkOverlay.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.StatusText.LinkOverlay.View.SetDynamicColors(true)
	app.UI.StatusText.LinkOverlay.Controls.View.SetDynamicColors(true)
	app.UI.StatusText.LinkOverlay.Controls.View.SetBackgroundColor(tcell.ColorDefault)

	links := modal(flLinks, app.UI.StatusText.LinkOverlay.View, app.UI.StatusText.LinkOverlay.Controls.View)

	app.UI.Pages.AddPage("links",
		links, true, false)

	app.UI.AuthOverlay = NewAuthoverlay(app, tview.NewFlex(), tview.NewInputField(),
		NewControls(app, tview.NewTextView()))

	app.UI.Pages.AddPage("login",
		authModal(app.UI.AuthOverlay.Flex, app.UI.AuthOverlay.View, app.UI.AuthOverlay.Controls.View), true, false)
	app.UI.AuthOverlay.Flex.SetDrawFunc(clearContent)
	app.UI.AuthOverlay.Flex.SetBackgroundColor(tcell.ColorDefault)
	app.UI.AuthOverlay.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.AuthOverlay.View.SetFieldBackgroundColor(tcell.ColorDefault)
	app.UI.AuthOverlay.View.SetFieldTextColor(tcell.ColorWhite)
	app.UI.AuthOverlay.Controls.View.SetBackgroundColor(tcell.ColorDefault)
	app.UI.AuthOverlay.Draw()

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
			if !app.UI.StatusText.LinkOverlay.Scroll {
				if event.Key() == tcell.KeyRune {
					switch event.Rune() {
					case 't':
						app.UI.StatusText.LinkOverlay.ActivateScroll()
					default:
						app.UI.StatusText.LinkOverlay.AddRune(event.Rune())
					}
				} else {
					switch event.Key() {
					case tcell.KeyEsc:
						if app.UI.StatusText.LinkOverlay.HasInput() {
							app.UI.StatusText.LinkOverlay.Clear()
						} else {
							app.UI.SetFocus(LeftPaneFocus)
							app.UI.StatusText.LinkOverlay.DisableScroll()
						}
						return nil
					}
				}
				return nil
			} else {
				switch event.Key() {
				case tcell.KeyEsc:
					app.UI.StatusText.LinkOverlay.DisableScroll()
					return nil
				}

			}

			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 't':
					app.UI.StatusText.LinkOverlay.DisableScroll()
					return nil
				case 'q':
					app.UI.SetFocus(LeftPaneFocus)
					app.UI.StatusText.LinkOverlay.DisableScroll()
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyEsc:
					app.UI.SetFocus(LeftPaneFocus)
					app.UI.StatusText.LinkOverlay.DisableScroll()
					return nil
				}
			}
			return event
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
				case 'p':
					app.UI.MessageBox.Post()
					return nil
				case 'e':
					app.UI.MessageBox.EditText()
					return nil
				case 'c':
					app.UI.MessageBox.EditSpoiler()
					return nil
				case 't':
					app.UI.MessageBox.ToggleSpoiler()
					return nil
				case 'q':
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

		if app.UI.Focus == LeftPaneFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'v':
					app.UI.SetFocus(RightPaneFocus)
					return nil
				case 'k':
					app.UI.TootList.Prev()
					return nil
				case 'j':
					app.UI.TootList.Next()
					return nil
				case 'q':
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
				case 't':
					app.UI.ShowThread()
				case 's':
					app.UI.ShowSensetive()
				case 'c':
					app.UI.NewToot()
				case 'o':
					app.UI.ShowLinks()
				case 'r':
					app.UI.Reply()
				case 'm':
					app.UI.OpenMedia()
				}
			}
		}

		return event
	})

	words := strings.Split(":q,:quit", ",")
	app.UI.CmdBar.View.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if currentText == "" {
			return
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) <= 1 {
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
		}

	})

	app.UI.AuthOverlay.View.SetDoneFunc(func(key tcell.Key) {
		app.UI.AuthOverlay.GotInput()
	})

	if err := app.App.SetRoot(app.UI.Pages, true).Run(); err != nil {
		panic(err)
	}
}
