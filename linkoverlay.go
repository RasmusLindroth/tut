package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewLinkOverlay(app *App) *LinkOverlay {
	l := &LinkOverlay{
		app:        app,
		Flex:       tview.NewFlex(),
		TextBottom: tview.NewTextView(),
		List:       tview.NewList(),
	}

	l.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	l.TextBottom.SetDynamicColors(true)
	l.List.SetBackgroundColor(app.Config.Style.Background)
	l.List.SetMainTextColor(app.Config.Style.Text)
	l.List.SetSelectedBackgroundColor(app.Config.Style.ListSelectedBackground)
	l.List.SetSelectedTextColor(app.Config.Style.ListSelectedText)
	l.List.ShowSecondaryText(false)
	l.List.SetHighlightFullLine(true)
	l.Flex.SetDrawFunc(app.Config.ClearContent)
	l.TextBottom.SetText(ColorKey(app.Config.Style, "", "O", "pen"))
	return l
}

type LinkOverlay struct {
	app        *App
	Flex       *tview.Flex
	TextBottom *tview.TextView
	List       *tview.List
	urls       []URL
}

func (l *LinkOverlay) SetURLs(urls []URL) {
	l.urls = urls
	l.List.Clear()
	for _, url := range urls {
		l.List.AddItem(url.Text, "", 0, nil)
	}
}

func (l *LinkOverlay) Prev() {
	index := l.List.GetCurrentItem()
	if index-1 >= 0 {
		l.List.SetCurrentItem(index - 1)
	}
}

func (l *LinkOverlay) Next() {
	index := l.List.GetCurrentItem()
	if index+1 < l.List.GetItemCount() {
		l.List.SetCurrentItem(index + 1)
	}
}

func (l *LinkOverlay) Open() {
	index := l.List.GetCurrentItem()
	if len(l.urls) == 0 || index >= len(l.urls) {
		return
	}
	openURL(l.urls[index].URL)
}

func (l *LinkOverlay) InputHandler(event *tcell.EventKey) {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			l.Next()
		case 'k', 'K':
			l.Prev()
		case 'o', 'O':
			l.Open()
		case 'q', 'Q':
			l.app.UI.SetFocus(LeftPaneFocus)
		}
	} else {
		switch event.Key() {
		case tcell.KeyUp:
			l.Prev()
		case tcell.KeyDown:
			l.Next()
		case tcell.KeyEsc:
			l.app.UI.SetFocus(LeftPaneFocus)
		}
	}
}
