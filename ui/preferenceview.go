package ui

import (
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var visibilitiesPrefStr = []string{
	mastodon.VisibilityPublic,
	mastodon.VisibilityUnlisted,
	mastodon.VisibilityFollowersOnly,
}

type preferences struct {
	displayname string
	bio         string
	fields      []mastodon.Field
	visibility  string
}

type PreferenceView struct {
	tutView     *TutView
	shared      *Shared
	View        *tview.Flex
	displayName *tview.TextView
	bio         *tview.TextView
	fields      *tview.List
	visibility  *tview.DropDown
	controls    *tview.Flex
	preferences *preferences
	fieldFocus  bool
}

func NewPreferenceView(tv *TutView) *PreferenceView {
	p := &PreferenceView{
		tutView:     tv,
		shared:      tv.Shared,
		displayName: NewTextView(tv.tut.Config),
		bio:         NewTextView(tv.tut.Config),
		fields:      NewList(tv.tut.Config, false),
		visibility:  NewDropDown(tv.tut.Config),
		controls:    NewControlView(tv.tut.Config),
		preferences: &preferences{},
	}
	p.View = preferenceViewUI(p)
	p.MainFocus()
	p.Update()

	return p
}

func preferenceViewUI(p *PreferenceView) *tview.Flex {
	p.visibility.SetLabel("Default toot visibility: ")
	p.visibility.SetOptions(visibilitiesPrefStr, p.visibilitySelected)

	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if p.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(p.shared.Top.View, 1, 0, false)
	}
	r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(p.displayName, 1, 0, false).
			AddItem(p.visibility, 2, 0, false).
			AddItem(p.fields, 0, 1, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(p.bio, 0, 1, false), 0, 1, false), 0, 1, false).
		AddItem(p.controls, 1, 0, false).
		AddItem(p.shared.Bottom.View, 2, 0, false)
	return r
}

func (p *PreferenceView) Update() {
	pf := &preferences{}
	me := p.tutView.tut.Client.Me
	pf.displayname = me.DisplayName
	if me.Source != nil {
		if me.Source.Note != nil {
			pf.bio = *me.Source.Note
		}
		if me.Source.Fields != nil {
			for _, f := range *me.Source.Fields {
				pf.fields = append(pf.fields, mastodon.Field{
					Name:  f.Name,
					Value: f.Value,
				})
			}
		}
		if me.Source.Privacy != nil {
			pf.visibility = *me.Source.Privacy
		}
	}
	p.preferences = pf
	p.update()
}

func (p *PreferenceView) update() {
	pf := p.preferences
	p.displayName.SetText(fmt.Sprintf("Display name: %s", tview.Escape(pf.displayname)))
	p.bio.SetText(fmt.Sprintf("Bio:\n%s", tview.Escape(pf.bio)))

	p.fields.Clear()
	for _, f := range pf.fields {
		p.fields.AddItem(fmt.Sprintf("%s: %s\n", tview.Escape(f.Name), tview.Escape(f.Value)), "", 0, nil)
	}

	index := 0
	for i, v := range visibilitiesPrefStr {
		if pf.visibility == v {
			index = i
			break
		}
	}
	p.visibility.SetCurrentOption(index)
}

func (p *PreferenceView) HasFieldFocus() bool {
	return p.fieldFocus
}

func (p *PreferenceView) FieldFocus() {
	p.fieldFocus = true

	var items []Control
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceFieldsAdd, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceFieldsEdit, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceFieldsDelete, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.GlobalBack, true))
	p.controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len+1, 0, false)
		} else {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len, 0, false)
		}
	}
	cnf := p.tutView.tut.Config
	p.fields.SetSelectedBackgroundColor(cnf.Style.ListSelectedBackground)
	p.fields.SetSelectedTextColor(cnf.Style.ListSelectedText)
}

func (p *PreferenceView) MainFocus() {
	p.fieldFocus = false

	var items []Control
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceName, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceVisibility, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceBio, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceFields, true))
	items = append(items, NewControl(p.tutView.tut.Config, p.tutView.tut.Config.Input.PreferenceSave, true))
	p.controls.Clear()
	for i, item := range items {
		if i < len(items)-1 {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len+1, 0, false)
		} else {
			p.controls.AddItem(NewControlButton(p.tutView, item), item.Len, 0, false)
		}
	}

	cnf := p.tutView.tut.Config
	p.fields.SetSelectedBackgroundColor(cnf.Style.Background)
	p.fields.SetSelectedTextColor(cnf.Style.Text)
}

func (p *PreferenceView) PrevField() {
	index := p.fields.GetCurrentItem()
	if index-1 >= 0 {
		p.fields.SetCurrentItem(index - 1)
	}
}

func (p *PreferenceView) NextField() {
	index := p.fields.GetCurrentItem()
	if index+1 < p.fields.GetItemCount() {
		p.fields.SetCurrentItem(index + 1)
	}
}

func (p *PreferenceView) AddField() {
	if p.fields.GetItemCount() > 3 {
		p.tutView.ShowError("You can have a maximum of four fields.")
		return
	}
	name, valid, err := OpenEditorLengthLimit(p.tutView, "name", 255)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't add name. Error: %v\n", err),
		)
		return
	}
	if !valid {
		p.tutView.ShowError("Name can't be empty.")
		return
	}
	value, valid, err := OpenEditorLengthLimit(p.tutView, "value", 255)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't add value. Error: %v\n", err),
		)
		return
	}
	if !valid {
		p.tutView.ShowError("Value can't be empty.")
		return
	}
	field := mastodon.Field{
		Name:  name,
		Value: value,
	}
	p.preferences.fields = append(p.preferences.fields, field)
	p.update()
	p.fields.SetCurrentItem(p.fields.GetItemCount() - 1)
}

func (p *PreferenceView) EditField() {
	if p.fields.GetItemCount() == 0 {
		return
	}
	index := p.fields.GetCurrentItem()
	if index < 0 || index >= len(p.preferences.fields) {
		return
	}
	curr := p.preferences.fields[index]
	name, valid, err := OpenEditorLengthLimit(p.tutView, curr.Name, 255)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't edit name. Error: %v\n", err),
		)
		return
	}
	if !valid {
		p.tutView.ShowError("Name can't be empty.")
		return
	}
	value, valid, err := OpenEditorLengthLimit(p.tutView, curr.Value, 255)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't edit value. Error: %v\n", err),
		)
		return
	}
	if !valid {
		p.tutView.ShowError("Value can't be empty.")
		return
	}
	field := mastodon.Field{
		Name:  name,
		Value: value,
	}
	p.preferences.fields[index] = field
	p.update()
}

func (p *PreferenceView) DeleteField() {
	if p.fields.GetItemCount() == 0 {
		return
	}
	index := p.fields.GetCurrentItem()
	if index < 0 || index >= len(p.preferences.fields) {
		return
	}
	p.fields.RemoveItem(index)
	p.preferences.fields = append(p.preferences.fields[:index], p.preferences.fields[index+1:]...)
	p.update()
}

func (p *PreferenceView) EditBio() {
	bio := p.preferences.bio
	text, _, err := OpenEditorLengthLimit(p.tutView, bio, 500)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't edit bio. Error: %v\n", err),
		)
		return
	}
	p.preferences.bio = text
	p.update()
}

func (p *PreferenceView) EditDisplayname() {
	dn := p.preferences.displayname
	text, _, err := OpenEditorLengthLimit(p.tutView, dn, 30)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't edit display name. Error: %v\n", err),
		)
		return
	}
	p.preferences.displayname = text
	p.update()
}

func (p *PreferenceView) visibilityInput(event *tcell.EventKey) *tcell.EventKey {
	if p.tutView.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}
	if p.tutView.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}
	if p.tutView.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) ||
		p.tutView.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		p.exitVisibility()
		return nil
	}
	return event
}

func (p *PreferenceView) exitVisibility() {
	p.tutView.tut.App.SetInputCapture(p.tutView.Input)
	p.tutView.tut.App.SetFocus(p.tutView.View)
}

func (p *PreferenceView) visibilitySelected(s string, index int) {
	_, p.preferences.visibility = p.visibility.GetCurrentOption()
	p.exitVisibility()
}

func (p *PreferenceView) FocusVisibility() {
	p.tutView.tut.App.SetInputCapture(p.visibilityInput)
	p.tutView.tut.App.SetFocus(p.visibility)
	ev := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	p.tutView.tut.App.QueueEvent(ev)
}

func (p *PreferenceView) Save() {
	og := &preferences{}
	me := p.tutView.tut.Client.Me
	og.displayname = me.DisplayName
	if me.Source != nil {
		if me.Source.Note != nil {
			og.bio = *me.Source.Note
		}
		if me.Source.Fields != nil {
			for _, f := range *me.Source.Fields {
				og.fields = append(og.fields, mastodon.Field{
					Name:  f.Name,
					Value: f.Value,
				})
			}
		}
		if me.Source.Privacy != nil {
			og.visibility = *me.Source.Privacy
		}
	}

	profile := mastodon.Profile{
		Source: &mastodon.AccountSource{},
	}
	if og.displayname != p.preferences.displayname {
		profile.DisplayName = &p.preferences.displayname
	}
	if og.bio != p.preferences.bio {
		profile.Note = &p.preferences.bio
	}
	if og.visibility != p.preferences.visibility {
		profile.Source.Privacy = &p.preferences.visibility
	}
	profile.Fields = &p.preferences.fields

	err := p.tutView.tut.Client.SavePreferences(&profile)
	if err != nil {
		p.tutView.ShowError(
			fmt.Sprintf("Couldn't update preferences. Error: %v\n", err),
		)
		return
	}
	p.tutView.SetPage(MainFocus)
}
