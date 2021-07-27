package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewUserSelectOverlay(app *App) *UserSelectOverlay {
	u := &UserSelectOverlay{
		app:  app,
		Flex: tview.NewFlex(),
		List: tview.NewList(),
		Text: tview.NewTextView(),
	}

	u.Flex.SetBackgroundColor(app.Config.Style.Background)
	u.List.SetMainTextColor(app.Config.Style.Text)
	u.List.SetBackgroundColor(app.Config.Style.Background)
	u.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	u.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	u.List.ShowSecondaryText(false)
	u.List.SetHighlightFullLine(true)
	u.Text.SetBackgroundColor(app.Config.Style.Background)
	u.Text.SetTextColor(app.Config.Style.Text)
	u.Flex.SetDrawFunc(app.Config.ClearContent)
	return u
}

type UserSelectOverlay struct {
	app  *App
	Flex *tview.Flex
	List *tview.List
	Text *tview.TextView
}

func (u *UserSelectOverlay) Prev() {
	index := u.List.GetCurrentItem()
	if index-1 >= 0 {
		u.List.SetCurrentItem(index - 1)
	}
}

func (u *UserSelectOverlay) Next() {
	index := u.List.GetCurrentItem()
	if index+1 < u.List.GetItemCount() {
		u.List.SetCurrentItem(index + 1)
	}
}
func (u *UserSelectOverlay) Done() {
	index := u.List.GetCurrentItem()
	u.app.Login(index)
	u.app.UI.LoggedIn()
}

func (u *UserSelectOverlay) InputHandler(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			u.Next()
		case 'k', 'K':
			u.Prev()
		case 'q', 'Q':
			u.app.UI.Root.Stop()
		}
	} else {
		switch event.Key() {
		case tcell.KeyEnter:
			u.Done()
		case tcell.KeyUp:
			u.Prev()
		case tcell.KeyDown:
			u.Next()
		}
	}
}

func (u *UserSelectOverlay) Draw() {
	u.Text.SetText("Select the user you want to use for this session by pressing Enter.")
	if len(u.app.Accounts.Accounts) > 0 {
		for i := 0; i < len(u.app.Accounts.Accounts); i++ {
			acc := u.app.Accounts.Accounts[i]
			u.List.AddItem(fmt.Sprintf("%s - %s", acc.Name, acc.Server), "", 0, nil)
		}
	}
}
