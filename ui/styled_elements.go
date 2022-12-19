package ui

import (
	"github.com/RasmusLindroth/tut/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewModal(cnf *config.Config) *tview.Modal {
	m := tview.NewModal()
	m.SetTextColor(cnf.Style.Text)
	m.SetBackgroundColor(cnf.Style.Background)
	m.SetBorderColor(cnf.Style.Background)
	m.SetBorder(false)
	m.SetButtonBackgroundColor(cnf.Style.ButtonColorTwo)
	m.SetButtonTextColor(cnf.Style.ButtonColorOne)
	tview.Styles.BorderColor = cnf.Style.Background
	return m
}

func NewTextView(cnf *config.Config) *tview.TextView {
	tw := tview.NewTextView()
	tw.SetBackgroundColor(cnf.Style.Background)
	tw.SetTextColor(cnf.Style.Text)
	tw.SetDynamicColors(true)
	return tw
}

func NewControlView(cnf *config.Config) *tview.Flex {
	f := tview.NewFlex().SetDirection(tview.FlexColumn)
	f.SetBackgroundColor(cnf.Style.Background)
	return f
}

func NewControlButton(tv *TutView, control Control) *tview.Button {
	btn := tview.NewButton(control.Label)
	style := tcell.Style{}
	style = style.Foreground(tv.tut.Config.Style.Text)
	style = style.Background(tv.tut.Config.Style.Background)
	btn.SetActivatedStyle(style)
	btn.SetStyle(style)
	btn.SetBackgroundColor(tv.tut.Config.Style.Background)
	btn.SetBackgroundColorActivated(tv.tut.Config.Style.Background)
	btn.SetLabelColor(tv.tut.Config.Style.Background)
	btn.SetLabelColorActivated(tv.tut.Config.Style.Background)
	btn.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if !btn.InRect(event.Position()) {
			return action, event
		}
		if action != tview.MouseLeftClick {
			return action, event
		}
		tv.tut.App.QueueEvent(control.Click())
		return action, nil
	})
	return btn
}

func NewList(cnf *config.Config, is_feed bool) *tview.List {
	l := tview.NewList()
	l.ShowSecondaryText(false)
	l.SetHighlightFullLine(true)
	l.SetBackgroundColor(cnf.Style.Background)
	l.SetMainTextColor(cnf.Style.Text)
	if is_feed && cnf.Style.ListSelectedBoldUnderline == 1 {
		l.SetSelectedBackgroundColor(cnf.Style.Background)
		s := tcell.Style.Attributes(tcell.Style{}, tcell.AttrBold|tcell.AttrUnderline)
		l.SetSelectedStyle(s)
	} else {
		l.SetSelectedBackgroundColor(cnf.Style.ListSelectedBackground)
	}
	l.SetSelectedTextColor(cnf.Style.ListSelectedText)
	return l
}

func NewDropDown(cnf *config.Config) *tview.DropDown {
	dd := tview.NewDropDown()
	dd.SetBackgroundColor(cnf.Style.Background)
	dd.SetFieldBackgroundColor(cnf.Style.Background)
	dd.SetFieldTextColor(cnf.Style.Text)

	selected := tcell.Style{}.
		Background(cnf.Style.AutocompleteSelectedBackground).
		Foreground(cnf.Style.AutocompleteSelectedText)
	unselected := tcell.Style{}.
		Background(cnf.Style.AutocompleteBackground).
		Foreground(cnf.Style.AutocompleteText)
	dd.SetListStyles(unselected, selected)
	return dd
}

func NewInputField(cnf *config.Config) *tview.InputField {
	i := tview.NewInputField()
	i.SetBackgroundColor(cnf.Style.Background)
	i.SetFieldBackgroundColor(cnf.Style.Background)
	i.SetFieldTextColor((cnf.Style.CommandText))

	selected := tcell.Style{}.
		Background(cnf.Style.AutocompleteSelectedBackground).
		Foreground(cnf.Style.AutocompleteSelectedText)
	unselected := tcell.Style{}.
		Background(cnf.Style.AutocompleteBackground).
		Foreground(cnf.Style.AutocompleteText)

	i.SetAutocompleteStyles(
		cnf.Style.AutocompleteBackground,
		unselected, selected)
	return i
}

func NewVerticalLine(cnf *config.Config) *tview.Box {
	verticalLine := tview.NewBox()
	verticalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		var s tcell.Style
		s = s.Background(cnf.Style.Background).Foreground(cnf.Style.Subtle)
		for cy := y; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, s)
		}
		return 0, 0, 0, 0
	})
	return verticalLine
}

func NewHorizontalLine(cnf *config.Config) *tview.Box {
	horizontalLine := tview.NewBox()
	horizontalLine.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		var s tcell.Style
		s = s.Background(cnf.Style.Background).Foreground(cnf.Style.Subtle)
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, s)
		}
		return 0, 0, 0, 0
	})
	return horizontalLine
}
