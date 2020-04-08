package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-mastodon"
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
	mentions   []mastodon.Mention
	tags       []mastodon.Tag
}

func (l *LinkOverlay) SetLinks(urls []URL, status *mastodon.Status) {
	l.urls = []URL{}
	l.mentions = []mastodon.Mention{}
	l.tags = []mastodon.Tag{}

	if urls != nil {
		l.urls = urls
	}
	if status != nil {
		l.mentions = status.Mentions
		l.tags = status.Tags
	}

	l.List.Clear()
	for _, url := range urls {
		l.List.AddItem(url.Text, "", 0, nil)
	}
	for _, mention := range l.mentions {
		l.List.AddItem(mention.Acct, "", 0, nil)
	}
	for _, tag := range l.tags {
		l.List.AddItem("#"+tag.Name, "", 0, nil)
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
	total := len(l.urls) + len(l.mentions) + len(l.tags)
	if total == 0 || index >= total {
		return
	}
	if index < len(l.urls) {
		openURL(l.urls[index].URL)
		return
	}
	mIndex := index - len(l.urls)
	if mIndex < len(l.mentions) {
		u, err := l.app.API.GetUserByID(l.mentions[mIndex].ID)
		if err != nil {
			l.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load user. Error: %v\n", err))
			return
		}
		l.app.UI.StatusView.AddFeed(
			NewUserFeed(l.app, *u),
		)
		l.app.UI.SetFocus(LeftPaneFocus)
		return
	}
	tIndex := index - len(l.mentions) - len(l.urls)
	if tIndex < len(l.tags) {
		l.app.UI.StatusView.AddFeed(
			NewTagFeed(l.app, l.tags[tIndex].Name),
		)
		l.app.UI.SetFocus(LeftPaneFocus)
	}
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
		case tcell.KeyEnter:
			l.Open()
		case tcell.KeyUp:
			l.Prev()
		case tcell.KeyDown:
			l.Next()
		case tcell.KeyEsc:
			l.app.UI.SetFocus(LeftPaneFocus)
		}
	}
}
