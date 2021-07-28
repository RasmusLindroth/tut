package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

func NewVisibilityOverlay(app *App) *VisibilityOverlay {
	v := &VisibilityOverlay{
		app:        app,
		Flex:       tview.NewFlex(),
		TextBottom: tview.NewTextView(),
		List:       tview.NewList(),
	}

	v.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	v.TextBottom.SetDynamicColors(true)
	v.List.SetBackgroundColor(app.Config.Style.Background)
	v.List.SetMainTextColor(app.Config.Style.Text)
	v.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	v.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	v.List.ShowSecondaryText(false)
	v.List.SetHighlightFullLine(true)
	v.Flex.SetDrawFunc(app.Config.ClearContent)
	v.TextBottom.SetText(ColorKey(app.Config, "", "Enter", ""))
	return v
}

type VisibilityOverlay struct {
	app        *App
	Flex       *tview.Flex
	TextBottom *tview.TextView
	List       *tview.List
	Selected   int
}

func (v *VisibilityOverlay) SetVisibilty(s string) {
	v.List.Clear()
	visibilities := []string{
		mastodon.VisibilityPublic,
		mastodon.VisibilityFollowersOnly,
		mastodon.VisibilityUnlisted,
		mastodon.VisibilityDirectMessage,
	}

	selected := 0
	for i, item := range visibilities {
		if s == item {
			selected = i
		}
		v.List.AddItem(
			VisibilityToText(item),
			"", 0, nil)
	}
	v.List.SetCurrentItem(selected)
	v.Selected = selected
}

func (v *VisibilityOverlay) Show() {
	v.List.SetCurrentItem(v.Selected)
}

func (v *VisibilityOverlay) Prev() {
	index := v.List.GetCurrentItem()
	if index-1 >= 0 {
		v.List.SetCurrentItem(index - 1)
	}
}

func (v *VisibilityOverlay) Next() {
	index := v.List.GetCurrentItem()
	if index+1 < v.List.GetItemCount() {
		v.List.SetCurrentItem(index + 1)
	}
}

func (v *VisibilityOverlay) SetVisibilityIndex() {
	index := v.List.GetCurrentItem()
	v.Selected = index
}

func (v *VisibilityOverlay) GetVisibility() string {
	visibilities := []string{
		mastodon.VisibilityPublic,
		mastodon.VisibilityFollowersOnly,
		mastodon.VisibilityUnlisted,
		mastodon.VisibilityDirectMessage,
	}
	return visibilities[v.Selected]
}

func (v *VisibilityOverlay) InputHandler(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			v.Next()
		case 'k', 'K':
			v.Prev()
		case 'q', 'Q':
			v.app.UI.SetFocus(MessageFocus)
		}
	} else {
		switch event.Key() {
		case tcell.KeyEnter:
			v.SetVisibilityIndex()
			v.app.UI.SetFocus(MessageFocus)
		case tcell.KeyUp:
			v.Prev()
		case tcell.KeyDown:
			v.Next()
		case tcell.KeyEsc:
			v.app.UI.SetFocus(MessageFocus)
		}
	}
}
