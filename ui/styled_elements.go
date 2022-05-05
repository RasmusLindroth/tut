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

func NewList(cnf *config.Config) *tview.List {
	l := tview.NewList()
	l.ShowSecondaryText(false)
	l.SetHighlightFullLine(true)
	l.SetBackgroundColor(cnf.Style.Background)
	l.SetMainTextColor(cnf.Style.Text)
	l.SetSelectedBackgroundColor(cnf.Style.ListSelectedBackground)
	l.SetSelectedTextColor(cnf.Style.ListSelectedText)
	return l
}

func NewDropDown(cnf *config.Config) *tview.DropDown {
	dd := tview.NewDropDown()
	dd.SetBackgroundColor(cnf.Style.Background)
	dd.SetFieldBackgroundColor(cnf.Style.Background)
	dd.SetFieldTextColor(cnf.Style.Text)

	selected := tcell.Style{}.
		Background(cnf.Style.ListSelectedBackground).
		Foreground(cnf.Style.ListSelectedText)
	unselected := tcell.Style{}.
		Background(cnf.Style.StatusBarViewBackground).
		Foreground(cnf.Style.StatusBarViewText)
	dd.SetListStyles(selected, unselected)
	return dd
}

func NewInputField(cnf *config.Config) *tview.InputField {
	i := tview.NewInputField()
	i.SetBackgroundColor(cnf.Style.Background)
	i.SetFieldBackgroundColor(cnf.Style.Background)

	selected := tcell.Style{}.
		Background(cnf.Style.ListSelectedBackground).
		Foreground(cnf.Style.ListSelectedText)
	unselected := tcell.Style{}.
		Background(cnf.Style.StatusBarViewBackground).
		Foreground(cnf.Style.StatusBarViewText)

	i.SetAutocompleteStyles(
		cnf.Style.Background,
		selected, unselected)
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
