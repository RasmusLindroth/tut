package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type msgToot struct {
	Text        string
	Status      *mastodon.Status
	MediaIDs    []mastodon.ID
	Sensitive   bool
	SpoilerText string
	Visibility  string
	ScheduledAt *time.Time
}

func NewMessageBox(app *App, view *tview.TextView, controls *Controls) *MessageBox {
	return &MessageBox{
		app:      app,
		View:     view,
		Index:    0,
		Controls: controls,
	}
}

type MessageBox struct {
	app         *App
	View        *tview.TextView
	Controls    *Controls
	Index       int
	maxIndex    int
	currentToot msgToot
}

func (m *MessageBox) NewToot() {
	m.composeToot(nil)
}

func (m *MessageBox) Reply(status *mastodon.Status) {
	m.composeToot(status)
}

func (m *MessageBox) ToggleSpoiler() {
	m.currentToot.Sensitive = !m.currentToot.Sensitive
	m.Draw()
}

func (m *MessageBox) composeToot(status *mastodon.Status) {
	m.Index = 0
	mt := msgToot{}
	if status != nil {
		if status.Reblog != nil {
			status = status.Reblog
		}
		mt.Status = status
	}
	m.currentToot = mt
}

func (m *MessageBox) Up() {
	if m.Index-1 > -1 {
		m.Index--
	}
	m.View.ScrollTo(m.Index, 0)
}

func (m *MessageBox) Down() {
	m.Index++
	if m.Index > m.maxIndex {
		m.Index = m.maxIndex
	}
	m.View.ScrollTo(m.Index, 0)
}

func (m *MessageBox) Post() {
	toot := m.currentToot
	send := mastodon.Toot{
		Status: strings.TrimSpace(toot.Text),
	}
	if toot.Status != nil {
		send.InReplyToID = toot.Status.ID
	}
	if toot.Sensitive {
		send.Sensitive = true
		send.SpoilerText = toot.SpoilerText
	}

	_, err := m.app.API.Client.PostStatus(context.Background(), &send)
	if err != nil {
		log.Fatalln(err)
	}
	m.app.UI.SetFocus(LeftPaneFocus)
}

func (m *MessageBox) Draw() {
	info := "\n[P]ost [E]dit text, [T]oggle CW, [C]ontent warning text [M]edia attachment"
	status := tview.Escape(info)
	m.Controls.View.SetText(status)

	var outputHead string
	var output string

	if m.currentToot.Status != nil {
		var acct string
		if m.currentToot.Status.Account.DisplayName != "" {
			acct = fmt.Sprintf("%s (%s)\n", m.currentToot.Status.Account.DisplayName, m.currentToot.Status.Account.Acct)
		} else {
			acct = fmt.Sprintf("%s\n", m.currentToot.Status.Account.Acct)
		}
		outputHead += "[gray]Replying to " + tview.Escape(acct) + "\n"
	}

	if m.currentToot.SpoilerText != "" && !m.currentToot.Sensitive {
		outputHead += "[red]You have entered spoiler text, but haven't set an content warning. Do it by pressing " + tview.Escape("[T]") + "\n\n"
	}

	if m.currentToot.Sensitive && m.currentToot.SpoilerText == "" {
		outputHead += "[red]You have added an content warning, but haven't set any text above the hidden text. Do it by pressing " + tview.Escape("[C]") + "\n\n"
	}

	if m.currentToot.Sensitive && m.currentToot.SpoilerText != "" {
		outputHead += "[gray]Content warning\n\n"
		outputHead += tview.Escape(m.currentToot.SpoilerText)
		outputHead += "\n\n[gray]---hidden content below---\n\n"
	}

	output = outputHead + tview.Escape(m.currentToot.Text)

	m.View.SetText(output)
	m.View.ScrollToEnd()
	m.maxIndex, _ = m.View.GetScrollOffset()
	m.View.ScrollTo(m.Index, 0)
}

func (m *MessageBox) EditText() {
	t := m.currentToot.Text
	s := m.currentToot.Status
	if t == "" && s != nil {
		var users []string
		if s.Account.Acct != m.app.Me.Acct {
			users = append(users, "@"+s.Account.Acct)
		}
		for _, men := range s.Mentions {
			if men.Acct == m.app.Me.Acct {
				continue
			}
			users = append(users, "@"+men.Acct)
		}
		t = strings.Join(users, " ")
	}
	text, err := openEditor(m.app.App, t)
	if err != nil {
		log.Fatalln(err)
	}
	m.currentToot.Text = text
	m.Draw()
}

func (m *MessageBox) EditSpoiler() {
	text, err := openEditor(m.app.App, m.currentToot.SpoilerText)
	if err != nil {
		log.Fatalln(err)
	}
	m.currentToot.SpoilerText = text
	m.Draw()
}
