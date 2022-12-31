package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rivo/uniseg"
)

type msgToot struct {
	ID            mastodon.ID
	Text          string
	Reply         *mastodon.Status
	Edit          *mastodon.Status
	MediaIDs      []mastodon.ID
	Sensitive     bool
	CWText        string
	ScheduledAt   *time.Time
	QuoteIncluded bool
	Visibility    string
	Language      string
}

type ComposeView struct {
	tutView    *TutView
	shared     *Shared
	View       *tview.Flex
	content    *tview.TextView
	input      *MediaInput
	info       *tview.TextView
	controls   *tview.Flex
	visibility *tview.DropDown
	lang       *tview.DropDown
	media      *MediaList
	msg        *msgToot
}

var visibilities = map[string]int{
	mastodon.VisibilityPublic:        0,
	mastodon.VisibilityUnlisted:      1,
	mastodon.VisibilityFollowersOnly: 2,
	mastodon.VisibilityDirectMessage: 3,
}
var visibilitiesStr = []string{
	mastodon.VisibilityPublic,
	mastodon.VisibilityUnlisted,
	mastodon.VisibilityFollowersOnly,
	mastodon.VisibilityDirectMessage,
}

func NewComposeView(tv *TutView) *ComposeView {
	cv := &ComposeView{
		tutView:    tv,
		shared:     tv.Shared,
		content:    NewTextView(tv.tut.Config),
		input:      NewMediaInput(tv),
		controls:   NewControlView(tv.tut.Config),
		info:       NewTextView(tv.tut.Config),
		visibility: NewDropDown(tv.tut.Config),
		lang:       NewDropDown(tv.tut.Config),
		media:      NewMediaList(tv),
	}
	cv.content.SetDynamicColors(true)
	cv.View = newComposeUI(cv)
	return cv
}

func newComposeUI(cv *ComposeView) *tview.Flex {
	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if cv.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(cv.tutView.Shared.Top.View, 1, 0, false)
	}
	r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(cv.content, 0, 2, false), 0, 2, false).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(cv.visibility, 1, 0, false).
			AddItem(cv.lang, 1, 0, false).
			AddItem(cv.info, 5, 0, false).
			AddItem(cv.media.View, 0, 1, false), 0, 1, false), 0, 1, false).
		AddItem(cv.input.View, 1, 0, false).
		AddItem(cv.controls, 1, 0, false).
		AddItem(cv.tutView.Shared.Bottom.View, 2, 0, false)
	return r
}

type ComposeControls uint

const (
	ComposeNormal ComposeControls = iota
	ComposeMedia
)

func (cv *ComposeView) msgLength() int {
	m := cv.msg
	charCount := uniseg.GraphemeClusterCount(m.Text)
	spoilerCount := uniseg.GraphemeClusterCount(m.CWText)
	totalCount := charCount
	if m.Sensitive {
		totalCount += spoilerCount
	}
	charsLeft := cv.tutView.tut.Config.General.CharLimit - totalCount
	return charsLeft
}

func (cv *ComposeView) SetControls(ctrl ComposeControls) {
	var items []Control
	switch ctrl {
	case ComposeNormal:
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposePost, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeEditText, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeVisibility, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeToggleContentWarning, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeEditCW, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeMediaFocus, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposePoll, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeLanguage, true))
		if cv.msg.Reply != nil {
			items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.ComposeIncludeQuote, true))
		}
	case ComposeMedia:
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaAdd, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaDelete, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.MediaEditDesc, true))
		items = append(items, NewControl(cv.tutView.tut.Config, cv.tutView.tut.Config.Input.GlobalBack, true))
	}
	cv.controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			cv.controls.AddItem(NewControlButton(cv.tutView, item), item.Len+1, 0, false)
		} else {
			cv.controls.AddItem(NewControlButton(cv.tutView, item), item.Len, 0, false)
		}
	}
}

func (cv *ComposeView) SetStatus(reply *mastodon.Status, edit *mastodon.Status) error {
	cv.tutView.PollView.Reset()
	cv.media.Reset()
	msg := &msgToot{}
	me := cv.tutView.tut.Client.Me
	visibility := mastodon.VisibilityPublic
	lang := ""
	if me.Source != nil && me.Source.Privacy != nil {
		visibility = *me.Source.Privacy
	}
	if me.Source != nil && me.Source.Language != nil {
		lang = *me.Source.Language
	}
	if reply != nil {
		if reply.Reblog != nil {
			reply = reply.Reblog
		}
		msg.Reply = reply
		if reply.Sensitive {
			msg.Sensitive = true
			msg.CWText = reply.SpoilerText
		}
		if visibilities[reply.Visibility] > visibilities[visibility] {
			visibility = reply.Visibility
		}
	}
	msg.Visibility = visibility
	msg.Language = lang
	cv.msg = msg
	cv.msg.Text = cv.getAccs()

	if edit != nil {
		source, err := cv.tutView.tut.Client.Client.GetStatusSource(context.Background(), edit.ID)
		if err != nil {
			cv.tutView.ShowError(
				fmt.Sprintf("Couldn't get status. Error: %v\n", err),
			)
			return err
		}
		msg := &msgToot{}
		msg.Edit = edit
		msg.ID = source.ID
		msg.Text = source.Text
		msg.CWText = source.SpoilerText
		for _, mid := range edit.MediaAttachments {
			msg.MediaIDs = append(msg.MediaIDs, mid.ID)
		}
		msg.Sensitive = edit.Sensitive
		msg.Visibility = edit.Visibility
		msg.Language = edit.Language
		if edit.Poll != nil {
			cv.tutView.PollView.AddPoll(edit.Poll)
		}
		if len(edit.MediaAttachments) > 0 {
			cv.media.AddFromEdit(edit)
		}

		cv.msg = msg
	}

	if cv.tutView.tut.Config.General.QuoteReply && edit == nil {
		cv.IncludeQuote()
	}
	cv.visibility.SetLabel("Visibility: ")
	index := 0
	for i, v := range visibilitiesStr {
		if cv.msg.Visibility == v {
			index = i
			break
		}
	}
	cv.visibility.SetOptions(visibilitiesStr, cv.visibilitySelected)
	cv.visibility.SetCurrentOption(index)
	cv.visibility.SetInputCapture(cv.visibilityInput)

	cv.lang.SetLabel("Lang: ")
	langStrs := []string{}
	for i, l := range util.Languages {
		if cv.msg.Language == l.Code {
			index = i
		}
		langStrs = append(langStrs, fmt.Sprintf("%s (%s)", l.Local, l.English))
	}
	cv.lang.SetOptions(langStrs, cv.langSelected)
	cv.lang.SetCurrentOption(index)

	cv.UpdateContent()
	cv.SetControls(ComposeNormal)
	return nil
}

func (cv *ComposeView) getAccs() string {
	if cv.msg.Reply == nil {
		return ""
	}
	s := cv.msg.Reply
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
	cv.UpdateContent()
}

func (cv *ComposeView) EditSpoiler() {
	text, err := OpenEditor(cv.tutView, cv.msg.CWText)
	if err != nil {
		cv.tutView.ShowError(
			fmt.Sprintf("Couldn't open editor. Error: %v", err),
		)
		return
	}
	cv.msg.CWText = text
	cv.UpdateContent()
}

func (cv *ComposeView) ToggleCW() {
	cv.msg.Sensitive = !cv.msg.Sensitive
	cv.UpdateContent()
}

func (cv *ComposeView) UpdateContent() {
	cv.info.SetText(fmt.Sprintf("Chars left: %d\nCW: %t\nHas poll: %t\n", cv.msgLength(), cv.msg.Sensitive, cv.tutView.PollView.HasPoll()))
	normal := config.ColorMark(cv.tutView.tut.Config.Style.Text)
	subtleColor := config.ColorMark(cv.tutView.tut.Config.Style.Subtle)
	warningColor := config.ColorMark(cv.tutView.tut.Config.Style.WarningText)

	var outputHead string
	var output string

	if cv.msg.Reply != nil {
		var acct string
		if cv.msg.Reply.Account.DisplayName != "" {
			acct = fmt.Sprintf("%s (%s)\n", cv.msg.Reply.Account.DisplayName, cv.msg.Reply.Account.Acct)
		} else {
			acct = fmt.Sprintf("%s\n", cv.msg.Reply.Account.Acct)
		}
		outputHead += subtleColor + "Replying to " + tview.Escape(acct) + "\n" + normal
	}
	if cv.msg.CWText != "" && !cv.msg.Sensitive {
		outputHead += warningColor + "You have entered content warning text, but haven't set an content warning. Do it by pressing " + tview.Escape("[T]") + "\n\n" + normal
	}

	if cv.msg.Sensitive && cv.msg.CWText == "" {
		outputHead += warningColor + "You have added an content warning, but haven't set any text above the hidden text. Do it by pressing " + tview.Escape("[C]") + "\n\n" + normal
	}

	if cv.msg.Sensitive && cv.msg.CWText != "" {
		outputHead += subtleColor + "Content warning\n\n" + normal
		outputHead += tview.Escape(cv.msg.CWText)
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
	s := cv.msg.Reply
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
	cv.UpdateContent()
}

func (cv *ComposeView) HasMedia() bool {
	return len(cv.media.Files) > 0
}

func (cv *ComposeView) visibilityInput(event *tcell.EventKey) *tcell.EventKey {
	if cv.tutView.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}
	if cv.tutView.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}
	if cv.tutView.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) ||
		cv.tutView.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		cv.exitVisibility()
		return nil
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

func (cv *ComposeView) langInput(event *tcell.EventKey) *tcell.EventKey {
	if cv.tutView.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}
	if cv.tutView.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}
	if cv.tutView.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) ||
		cv.tutView.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		cv.exitLang()
		return nil
	}
	return event
}

func (cv *ComposeView) exitLang() {
	cv.tutView.tut.App.SetInputCapture(cv.tutView.Input)
	cv.tutView.tut.App.SetFocus(cv.content)
}

func (cv *ComposeView) langSelected(s string, index int) {
	i, _ := cv.lang.GetCurrentOption()
	if i >= 0 && i < len(util.Languages) {
		cv.msg.Language = util.Languages[i].Code
	}
	cv.exitLang()
}

func (cv *ComposeView) FocusLang() {
	cv.tutView.tut.App.SetInputCapture(cv.langInput)
	cv.tutView.tut.App.SetFocus(cv.lang)
	ev := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	cv.tutView.tut.App.QueueEvent(ev)
}

func (cv *ComposeView) Post() {
	toot := cv.msg
	send := mastodon.Toot{
		Status: strings.TrimSpace(toot.Text),
	}
	if toot.Reply != nil {
		send.InReplyToID = toot.Reply.ID
	}
	if toot.Edit != nil && toot.Edit.InReplyToID != nil {
		send.InReplyToID = mastodon.ID(toot.Edit.InReplyToID.(string))
	}
	if toot.Sensitive {
		send.Sensitive = true
		send.SpoilerText = toot.CWText
	}

	if cv.HasMedia() {
		attachments := cv.media.Files
		for _, ap := range attachments {
			if ap.Remote {
				send.MediaIDs = append(send.MediaIDs, ap.ID)
				continue
			}
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
	}
	if cv.tutView.PollView.HasPoll() && !cv.HasMedia() {
		send.Poll = cv.tutView.PollView.GetPoll()
	}
	send.Visibility = cv.msg.Visibility
	send.Language = cv.msg.Language

	var err error
	var newPost *mastodon.Status
	if toot.Edit != nil {
		newPost, err = cv.tutView.tut.Client.Client.UpdateStatus(context.Background(), &send, toot.Edit.ID)
		if err == nil {
			item, itemErr := cv.tutView.GetCurrentItem()
			if itemErr != nil {
				return
			}
			if item.Type() != api.StatusType {
				return
			}
			s := item.Raw().(*mastodon.Status)
			*s = *newPost
			cv.tutView.RedrawContent()
		}
	} else {
		_, err = cv.tutView.tut.Client.Client.PostStatus(context.Background(), &send)
	}
	if err != nil {
		cv.tutView.ShowError(
			fmt.Sprintf("Couldn't post toot. Error: %v\n", err),
		)
		return
	}
	cv.tutView.SetPage(MainFocus)
}

type MediaList struct {
	tutView     *TutView
	View        *tview.Flex
	heading     *tview.TextView
	text        *tview.TextView
	list        *tview.List
	Files       []UploadFile
	scrollSleep *scrollSleep
}

func NewMediaList(tv *TutView) *MediaList {
	ml := &MediaList{
		tutView: tv,
		heading: NewTextView(tv.tut.Config),
		text:    NewTextView(tv.tut.Config),
		list:    NewList(tv.tut.Config, false),
	}
	ml.scrollSleep = NewScrollSleep(ml.Next, ml.Prev)
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
	Remote      bool
	ID          mastodon.ID
}

func (m *MediaList) AddFromEdit(edit *mastodon.Status) {
	m.Files = nil
	m.list.Clear()
	for i, ma := range edit.MediaAttachments {
		m.Files = append(m.Files, UploadFile{
			Description: ma.Description,
			Remote:      true,
			ID:          ma.ID,
		})
		m.list.AddItem(fmt.Sprintf("From edit: %d", i+1), "", 0, nil)
	}
	index := m.list.GetItemCount()
	if index > 0 {
		m.list.SetCurrentItem(index - 1)
	}
	m.Draw()
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
	if len(m.Files) != 0 && index > len(m.Files)-1 && m.Files[index].Description != "" {
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
	if len(m.Files) == 0 || index > len(m.Files)-1 {
		return
	}
	m.list.RemoveItem(index)
	m.list.SetCurrentItem(index)
	m.Files = append(m.Files[:index], m.Files[index+1:]...)
	m.Draw()
}

func (m *MediaList) EditDesc() {
	index := m.list.GetCurrentItem()
	if len(m.Files) == 0 || index > len(m.Files) {
		return
	}
	file := m.Files[index]
	if file.Remote {
		m.tutView.ShowError(
			"Can't edit desc of a file that's already uploaded",
		)
		return
	}
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
