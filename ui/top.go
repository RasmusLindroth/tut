package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

type Top struct {
	TutView *TutView
	View    *tview.TextView
}

func NewTop(tv *TutView) *Top {
	t := &Top{
		TutView: tv,
		View:    NewTextView(tv.tut.Config),
	}
	t.View.SetBackgroundColor(tv.tut.Config.Style.TopBarBackground)
	t.View.SetTextColor(tv.tut.Config.Style.TopBarText)

	return t
}

func (t *Top) SetText(s string) {
	if s == "" {
		t.View.SetText("tut")
	} else {
		t.View.SetText(fmt.Sprintf("tut - %s", s))
	}
}
