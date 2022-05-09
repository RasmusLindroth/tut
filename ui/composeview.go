package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/gdamore/tcell/v2"
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
	Visibility    string
}

type ComposeView struct {
	tutView    *TutView
	shared     *Shared
	View       *tview.Flex
	content    *tview.TextView
	input      *MediaInput
	info       *tview.TextView
	controls   *tview.TextView
	visibility *tview.DropDown
	media      *MediaList
	msg        *msgToot
}

var visibilities = []string{mastodon.VisibilityPublic, mastodon.VisibilityUnlisted, mastodon.VisibilityFollowersOnly, mastodon.VisibilityDirectMessage}

func NewComposeView(tv *TutView) *ComposeView {
	cv := &ComposeView{
		tutView:    tv,
		shared:     tv.Shared,
		content:    NewTextView(tv.tut.Config),
		input:      NewMediaInput(tv),
		controls:   NewTextView(tv.tut.Config),
		info:       NewTextView(tv.tut.Config),
		visibility: NewDropDown(tv.tut.Config),
		media:      NewMediaList(tv),
	}
	cv.content.SetDynamicColors(true)
	cv.controls.SetDynamicColors(true)
	cv.View = newComposeUI(cv)
	return cv
}

func newComposeUI(cv *ComposeView) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(cv.tutView.Shared.Top.View, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(cv.content, 0, 2, false), 0, 2, false).
			AddItem(tview.NewBox(), 2, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(cv.visibility, 1, 0, false).
				AddItem(cv.info, 4, 0, false).
				AddItem(cv.media.View, 0, 1, false), 0, 1, false), 0, 1, false).
		AddItem(cv.input.View, 1, 0, false).
		AddItem(cv.controls, 1, 0, false).
		AddItem(cv.tutView.Shared.Bottom.View, 2, 0, false)
}

type ComposeControls uint

const (
	ComposeNormal ComposeControls = iota
	ComposeMedia
)

func (cv *ComposeView) msgLength() int {
	m := cv.msg
	charCount := uniseg.GraphemeClusterCount(m.Text)
	spoilerCount := uniseg.GraphemeClusterCount(m.SpoilerText)
	totalCount := charCount
	if m.Sensitive {
		totalCount += spoilerCount
	}
	charsLeft := cv.tutView.tut.Config.General.CharLimit - totalCount
	return charsLeft
}

func (cv *ComposeView) SetControls(ctrl ComposeControls) {
	var items []string
	switch ctrl {
	case ComposeNormal:
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposePost, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeEditText, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeVisibility, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeToggleContentWarning, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeEditSpoiler, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeMediaFocus, true))
		if cv.msg.Status != nil {
			items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeIncludeQuote, true))
		}
	case ComposeMedia:
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaAdd, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaDelete, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaEditDesc, true))
		items = append(items, config.ColorFromKey(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.GlobalBack, true))
	}
	res := strings.Join(items, " ")
	cv.controls.SetText(res)
}

func (cv *ComposeView) SetStatus(status *mastodon.Status) {
	msg := &msgToot{}
	if status != nil {
		if status.Reblog != nil {
			status = status.Reblog
		}
		msg.Status = status
		if status.Sensitive {
			msg.Sensitive = true
			msg.SpoilerText = status.SpoilerText
		}
		msg.Visibility = status.Visibility
	}
	cv.msg = msg
	cv.msg.Text = cv.getAccs()
	if cv.tutView.tut.Config.General.QuoteReply {
		cv.IncludeQuote()
	}
	cv.visibility.SetLabel("Visibility: ")
	index := 0
	for i, v := range visibilities {
		if msg.Visibility == v {
			index = i
			break
		}
	}
	cv.visibility.SetOptions(visibilities, cv.visibilitySelected)
	cv.visibility.SetCurrentOption(index)
	cv.visibility.SetInputCapture(cv.visibilityInput)
	cv.updateContent()
	cv.SetControls(ComposeNormal)
}

func (cv *ComposeView) getAccs() string {
	if cv.msg.Status == nil {
		return ""
	}
	s := cv.msg.Status
	var users []string
	if s.Account.Acct != cv.tutView.tut.Client.Me.Acct {
		users = append(users, "@"+s.Account.Acct)
	}
	for _, men := range s.Mentions {
		if men.Acct == cv.tutView.tut.Client.Me.Acct {
			continue
		}
		users = append(users, "@"+men.Acct)
	}
	t := strings.Join(users, " ")
	return t
}

func (cv *ComposeView) EditText() {
	text, err := OpenEditor(cv.tutView, cv.msg.Text)
	if err != nil {
		cv.tutView.ShowError(
			fmt.Sprintf("Couldn't open editor. Error: %v", err),
		)
		return
	}
	cv.msg.Text = text
	cv.updateContent()
}

func (cv *ComposeView) EditSpoiler() {
	text, err := OpenEditor(cv.tutView, cv.msg.SpoilerText)
	if err != nil {
		cv.tutView.ShowError(
			fmt.Sprintf("Couldn't open editor. Error: %v", err),
		)
		return
	}
	cv.msg.SpoilerText = text
	cv.updateContent()
}

func (cv *ComposeView) ToggleCW() {
	cv.msg.Sensitive = !cv.msg.Sensitive
	cv.updateContent()
}

func (cv *ComposeView) updateContent() {
	cv.info.SetText(fmt.Sprintf("Chars left: %d\nSpoiler: %t\n", cv.msgLength(), cv.msg.Sensitive))
	normal := config.ColorMark(cv.tutView.tut.Config.Style.Text)
	subtleColor := config.ColorMark(cv.tutView.tut.Config.Style.Subtle)
	warningColor := config.ColorMark(cv.tutView.tut.Config.Style.WarningText)

	var outputHead string
	var output string

	if cv.msg.Status != nil {
		var acct string
		if cv.msg.Status.Account.DisplayName != "" {
			acct = fmt.Sprintf("%s (%s)\n", cv.msg.Status.Account.DisplayName, cv.msg.Status.Account.Acct)
		} else {
			acct = fmt.Sprintf("%s\n", cv.msg.Status.Account.Acct)
		}
		outputHead += subtleColor + "Replying to " + tview.Escape(acct) + "\n" + normal
	}
	if cv.msg.SpoilerText != "" && !cv.msg.Sensitive {
		outputHead += warningColor + "You have entered spoiler text, but haven't set an content warning. Do it by pressing " + tview.Escape("[T]") + "\n\n" + normal
	}

	if cv.msg.Sensitive && cv.msg.SpoilerText == "" {
		outputHead += warningColor + "You have added an content warning, but haven't set any text above the hidden text. Do it by pressing " + tview.Escape("[C]") + "\n\n" + normal
	}

	if cv.msg.Sensitive && cv.msg.SpoilerText != "" {
		outputHead += subtleColor + "Content warning\n\n" + normal
		outputHead += tview.Escape(cv.msg.SpoilerText)
		outputHead += "\n\n" + subtleColor + "---hidden content below---\n\n" + normal
	}
	output = outputHead + normal + tview.Escape(cv.msg.Text)

	cv.content.SetText(output)
}

func (cv *ComposeView) IncludeQuote() {
	if cv.msg.QuoteIncluded {
		return
	}
	t := cv.msg.Text
	s := cv.msg.Status
	if s == nil {
		return
	}
	tootText, _ := util.CleanHTML(s.Content)

	t += "\n\n"
	for _, line := range strings.Split(tootText, "\n") {
		t += "> " + line + "\n"
	}
	t += "\n"
	cv.msg.Text = t
	cv.msg.QuoteIncluded = true
	cv.updateContent()
}

func (cv *ComposeView) visibilityInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k', 'K':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'q', 'Q':
			cv.exitVisibility()
			return nil
		}
	} else {
		switch event.Key() {
		case tcell.KeyEsc:
			cv.exitVisibility()
			return nil
		}
	}
	return event
}

func (cv *ComposeView) exitVisibility() {
	cv.tutView.tut.App.SetInputCapture(cv.tutView.Input)
	cv.tutView.tut.App.SetFocus(cv.content)
}

func (cv *ComposeView) visibilitySelected(s string, index int) {
	_, cv.msg.Visibility = cv.visibility.GetCurrentOption()
	cv.exitVisibility()
}

func (cv *ComposeView) FocusVisibility() {
	cv.tutView.tut.App.SetInputCapture(cv.visibilityInput)
	cv.tutView.tut.App.SetFocus(cv.visibility)
	ev := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	cv.tutView.tut.App.QueueEvent(ev)
}

func (cv *ComposeView) Post() {
	toot := cv.msg
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

	attachments := cv.media.Files
	for _, ap := range attachments {
		f, err := os.Open(ap.Path)
		if err != nil {
			cv.tutView.ShowError(
				fmt.Sprintf("Couldn't upload media. Error: %v\n", err),
			)
			f.Close()
			return
		}
		media := &mastodon.Media{
			File: f,
		}
		if ap.Description != "" {
			media.Description = ap.Description
		}
		a, err := cv.tutView.tut.Client.Client.UploadMediaFromMedia(context.Background(), media)
		if err != nil {
			cv.tutView.ShowError(
				fmt.Sprintf("Couldn't upload media. Error: %v\n", err),
			)
			f.Close()
			return
		}
		f.Close()
		send.MediaIDs = append(send.MediaIDs, a.ID)
	}
	send.Visibility = cv.msg.Visibility

	_, err := cv.tutView.tut.Client.Client.PostStatus(context.Background(), &send)
	if err != nil {
		cv.tutView.ShowError(
			fmt.Sprintf("Couldn't post toot. Error: %v\n", err),
		)
		return
	}
	cv.tutView.SetPage(MainFocus)
}

type MediaList struct {
	tutView *TutView
	View    *tview.Flex
	heading *tview.TextView
	text    *tview.TextView
	list    *tview.List
	Files   []UploadFile
}

func NewMediaList(tv *TutView) *MediaList {
	ml := &MediaList{
		tutView: tv,
		heading: NewTextView(tv.tut.Config),
		text:    NewTextView(tv.tut.Config),
		list:    NewList(tv.tut.Config),
	}
	ml.heading.SetText(fmt.Sprintf("Media files: %d", ml.list.GetItemCount()))
	ml.heading.SetBorderPadding(1, 1, 0, 0)
	ml.View = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ml.heading, 1, 0, false).
		AddItem(ml.text, 1, 0, false).
		AddItem(ml.list, 0, 1, false)
	return ml
}

type UploadFile struct {
	Path        string
	Description string
}

func (m *MediaList) Reset() {
	m.Files = nil
	m.list.Clear()
	m.Draw()
}

func (m *MediaList) AddFile(f string) {
	file := UploadFile{Path: f}
	m.Files = append(m.Files, file)
	m.list.AddItem(filepath.Base(f), "", 0, nil)
	index := m.list.GetItemCount()
	m.list.SetCurrentItem(index - 1)
	m.Draw()
}

func (m *MediaList) Draw() {
	topText := "File desc: "

	index := m.list.GetCurrentItem()
	if len(m.Files) != 0 && index < len(m.Files) && m.Files[index].Description != "" {
		topText += tview.Escape(m.Files[index].Description)
	}
	m.text.SetText(topText)
}

func (m *MediaList) SetFocus(reset bool) {
	if reset {
		m.tutView.ComposeView.input.View.SetText("")
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			pwd = ""
		} else {
			pwd = home
		}
	}
	if !strings.HasSuffix(pwd, "/") {
		pwd += "/"
	}
	m.tutView.ComposeView.input.View.SetText(pwd)
}

func (m *MediaList) Prev() {
	index := m.list.GetCurrentItem()
	if index-1 >= 0 {
		m.list.SetCurrentItem(index - 1)
	}
	m.Draw()
}

func (m *MediaList) Next() {
	index := m.list.GetCurrentItem()
	if index+1 < m.list.GetItemCount() {
		m.list.SetCurrentItem(index + 1)
	}
	m.Draw()
}

func (m *MediaList) Delete() {
	index := m.list.GetCurrentItem()
	if len(m.Files) == 0 || index > len(m.Files) {
		return
	}
	m.list.RemoveItem(index)
	m.Files = append(m.Files[:index], m.Files[index+1:]...)
	m.Draw()
}

func (m *MediaList) EditDesc() {
	index := m.list.GetCurrentItem()
	if len(m.Files) == 0 || index > len(m.Files) {
		return
	}
	file := m.Files[index]
	desc, err := OpenEditor(m.tutView, file.Description)
	if err != nil {
		m.tutView.ShowError(
			fmt.Sprintf("Couldn't edit description. Error: %v\n", err),
		)
		return
	}
	file.Description = desc
	m.Files[index] = file
	m.Draw()
}

type MediaInput struct {
	tutView              *TutView
	View                 *tview.InputField
	text                 string
	autocompleteIndex    int
	autocompleteList     []string
	isAutocompleteChange bool
}

func NewMediaInput(tv *TutView) *MediaInput {
	m := &MediaInput{
		tutView: tv,
		View:    NewInputField(tv.tut.Config),
	}
	m.View.SetChangedFunc(m.HandleChanges)
	return m
}

func (m *MediaInput) AddRune(r rune) {
	newText := m.View.GetText() + string(r)
	m.text = newText
	m.View.SetText(m.text)
	m.saveAutocompleteState()
}

func (m *MediaInput) HandleChanges(text string) {
	if m.isAutocompleteChange {
		m.isAutocompleteChange = false
		return
	}
	m.saveAutocompleteState()
}

func (m *MediaInput) saveAutocompleteState() {
	text := m.View.GetText()
	m.text = text
	m.autocompleteList = util.FindFiles(text)
	m.autocompleteIndex = 0
}

func (m *MediaInput) AutocompletePrev() {
	if len(m.autocompleteList) == 0 {
		return
	}
	index := m.autocompleteIndex - 1
	if index < 0 {
		index = len(m.autocompleteList) - 1
	}
	m.autocompleteIndex = index
	m.showAutocomplete()
}

func (m *MediaInput) AutocompleteTab() {
	if len(m.autocompleteList) == 0 {
		return
	}
	same := ""
	for i := 0; i < len(m.autocompleteList[0]); i++ {
		match := true
		c := m.autocompleteList[0][i]
		for _, item := range m.autocompleteList {
			if i >= len(item) || c != item[i] {
				match = false
				break
			}
		}
		if !match {
			break
		}
		same += string(c)
	}
	if same != m.text {
		m.text = same
		m.View.SetText(same)
		m.saveAutocompleteState()
	} else {
		m.AutocompleteNext()
	}
}

func (m *MediaInput) AutocompleteNext() {
	if len(m.autocompleteList) == 0 {
		return
	}
	index := m.autocompleteIndex + 1
	if index >= len(m.autocompleteList) {
		index = 0
	}
	m.autocompleteIndex = index
	m.showAutocomplete()
}

func (m *MediaInput) CheckDone() {
	path := m.View.GetText()
	if util.IsDir(path) {
		m.saveAutocompleteState()
		return
	}

	m.tutView.ComposeView.media.AddFile(path)
	m.tutView.ComposeView.media.SetFocus(true)
	m.tutView.SetPage(MediaFocus)
}

func (m *MediaInput) showAutocomplete() {
	m.isAutocompleteChange = true
	m.View.SetText(m.autocompleteList[m.autocompleteIndex])
	if len(m.autocompleteList) < 3 {
		m.saveAutocompleteState()
	}
}
