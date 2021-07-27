package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
	"github.com/rivo/uniseg"
)

type msgToot struct {
	Text          string
	Status        *mastodon.Status
	MediaIDs      []mastodon.ID
	Sensitive     bool
	SpoilerText   string
	ScheduledAt   *time.Time
	QuoteIncluded bool
}

func VisibilityToText(s string) string {
	switch s {
	case mastodon.VisibilityPublic:
		return "Public"
	case mastodon.VisibilityUnlisted:
		return "Unlisted"
	case mastodon.VisibilityFollowersOnly:
		return "Followers"
	case mastodon.VisibilityDirectMessage:
		return "Direct"
	default:
		return "Public"
	}
}

func NewMessageBox(app *App) *MessageBox {
	m := &MessageBox{
		app:      app,
		Flex:     tview.NewFlex(),
		View:     tview.NewTextView(),
		Index:    0,
		Controls: tview.NewTextView(),
	}

	m.View.SetBackgroundColor(app.Config.Style.Background)
	m.View.SetTextColor(app.Config.Style.Text)
	m.View.SetDynamicColors(true)
	m.Controls.SetDynamicColors(true)
	m.Controls.SetBackgroundColor(app.Config.Style.Background)
	m.Controls.SetTextColor(app.Config.Style.Text)
	m.Flex.SetDrawFunc(app.Config.ClearContent)

	return m
}

type MessageBox struct {
	app         *App
	Flex        *tview.Flex
	View        *tview.TextView
	Controls    *tview.TextView
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
	visibility := mastodon.VisibilityPublic
	if status != nil && status.Visibility == mastodon.VisibilityDirectMessage {
		visibility = mastodon.VisibilityDirectMessage
	}
	m.app.UI.VisibilityOverlay.SetVisibilty(visibility)

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

	attachments := m.app.UI.MediaOverlay.Files
	for _, ap := range attachments {
		f, err := os.Open(ap.Path)
		if err != nil {
			m.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't upload media. Error: %v\n", err))
			return
		}
		media := &mastodon.Media{
			File: f,
		}
		if ap.Description != "" {
			media.Description = ap.Description
		}
		a, err := m.app.API.Client.UploadMediaFromMedia(context.Background(), media)
		if err != nil {
			m.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't upload media. Error: %v\n", err))
			return
		}
		f.Close()
		send.MediaIDs = append(send.MediaIDs, a.ID)
	}

	send.Visibility = m.app.UI.VisibilityOverlay.GetVisibility()

	_, err := m.app.API.Client.PostStatus(context.Background(), &send)
	if err != nil {
		m.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't post toot. Error: %v\n", err))
		return
	}
	m.app.UI.SetFocus(LeftPaneFocus)
}

func (m *MessageBox) TootLength() int {
	charCount := uniseg.GraphemeClusterCount(m.currentToot.Text)
	spoilerCount := uniseg.GraphemeClusterCount(m.currentToot.SpoilerText)
	totalCount := charCount
	if m.currentToot.Sensitive {
		totalCount += spoilerCount
	}
	charsLeft := m.app.Config.General.CharLimit - totalCount
	return charsLeft
}

func (m *MessageBox) Draw() {
	var items []string
	items = append(items, ColorKey(m.app.Config.Style, "", "P", "ost"))
	items = append(items, ColorKey(m.app.Config.Style, "", "E", "dit"))
	items = append(items, ColorKey(m.app.Config.Style, "", "V", "isibility"))
	items = append(items, ColorKey(m.app.Config.Style, "", "T", "oggle CW"))
	items = append(items, ColorKey(m.app.Config.Style, "", "C", "ontent warning text"))
	items = append(items, ColorKey(m.app.Config.Style, "", "M", "edia attachment"))
	items = append(items, ColorKey(m.app.Config.Style, "", "I", "nclude quote"))
	status := strings.Join(items, " ")
	m.Controls.SetText(status)

	var outputHead string
	var output string

	normal := ColorMark(m.app.Config.Style.Text)
	subtleColor := ColorMark(m.app.Config.Style.Subtle)
	warningColor := ColorMark(m.app.Config.Style.WarningText)

	charsLeft := m.TootLength()

	outputHead += subtleColor + VisibilityToText(m.app.UI.VisibilityOverlay.GetVisibility()) + ", "
	if charsLeft > 0 {
		outputHead += fmt.Sprintf("%d chars left", charsLeft) + "\n\n" + normal
	} else {
		outputHead += warningColor + fmt.Sprintf("%d chars left", charsLeft) + "\n\n" + normal
	}
	if m.currentToot.Status != nil {
		var acct string
		if m.currentToot.Status.Account.DisplayName != "" {
			acct = fmt.Sprintf("%s (%s)\n", m.currentToot.Status.Account.DisplayName, m.currentToot.Status.Account.Acct)
		} else {
			acct = fmt.Sprintf("%s\n", m.currentToot.Status.Account.Acct)
		}
		outputHead += subtleColor + "Replying to " + tview.Escape(acct) + "\n" + normal
	}
	if m.currentToot.SpoilerText != "" && !m.currentToot.Sensitive {
		outputHead += warningColor + "You have entered spoiler text, but haven't set an content warning. Do it by pressing " + tview.Escape("[T]") + "\n\n" + normal
	}

	if m.currentToot.Sensitive && m.currentToot.SpoilerText == "" {
		outputHead += warningColor + "You have added an content warning, but haven't set any text above the hidden text. Do it by pressing " + tview.Escape("[C]") + "\n\n" + normal
	}

	if m.currentToot.Sensitive && m.currentToot.SpoilerText != "" {
		outputHead += subtleColor + "Content warning\n\n" + normal
		outputHead += tview.Escape(m.currentToot.SpoilerText)
		outputHead += "\n\n" + subtleColor + "---hidden content below---\n\n" + normal
	}
	output = outputHead + normal + tview.Escape(m.currentToot.Text)

	m.View.SetText(output)
	m.View.ScrollToEnd()
	m.maxIndex, _ = m.View.GetScrollOffset()
	m.View.ScrollTo(m.Index, 0)
}

func (m *MessageBox) EditText() {
	t := m.currentToot.Text
	s := m.currentToot.Status
	m.currentToot.QuoteIncluded = false

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
		m.currentToot.Text = t

		if m.app.Config.General.QuoteReply {
			m.IncludeQuote()
		}
	}
	text, err := openEditor(m.app.UI.Root, m.currentToot.Text)
	if err != nil {
		m.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't edit toot. Error: %v\n", err))
		m.Draw()
		return
	}
	m.currentToot.Text = text
	m.Draw()
}

func (m *MessageBox) IncludeQuote() {
	if m.currentToot.QuoteIncluded {
		return
	}
	t := m.currentToot.Text
	s := m.currentToot.Status
	if s == nil {
		return
	}
	tootText, _ := cleanTootHTML(s.Content)

	t += "\n"
	for _, line := range strings.Split(tootText, "\n") {
		t += "> " + line + "\n"
	}
	t += "\n"
	m.currentToot.Text = t
	m.currentToot.QuoteIncluded = true
}

func (m *MessageBox) EditSpoiler() {
	text, err := openEditor(m.app.UI.Root, m.currentToot.SpoilerText)
	if err != nil {
		m.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't edit spoiler. Error: %v\n", err))
		m.Draw()
		return
	}
	m.currentToot.SpoilerText = text
	m.Draw()
}
