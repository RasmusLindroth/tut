package ui

import (
	"github.com/RasmusLindroth/tut/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTextView(cnf *config.Config) *tview.TextView {
	tw := tview.NewTextView()
	tw.SetDynamicColors(true)
	return tw
}

func NewList(cnf *config.Config) *tview.List {
	l := tview.NewList()
	l.ShowSecondaryText(false)
	l.SetHighlightFullLine(true)
	return l
}

func NewDropDown(cnf *config.Config) *tview.DropDown {
	dd := tview.NewDropDown()
	return dd
}

func NewInputField(cnf *config.Config) *tview.InputField {
	i := tview.NewInputField()
	i.SetBackgroundColor(cnf.Style.Background)
	i.SetFieldBackgroundColor(cnf.Style.Background)
	i.SetFieldTextColor(cnf.Style.Text)
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
