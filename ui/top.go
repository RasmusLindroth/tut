package ui

import (
	"fmt"
	"net/url"

	"github.com/RasmusLindroth/tut/util"
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
	if t.TutView.tut.Client != nil {
		acct := t.TutView.tut.Client.Me
		us := acct.Acct
		u, err := url.Parse(acct.URL)
		if err == nil {
			us = fmt.Sprintf("%s@%s", us, u.Host)
		}
		if s == "" {
			t.setText(fmt.Sprintf("tut - %s", us))
		} else {
			t.setText(fmt.Sprintf("tut - %s - %s", s, us))
		}
	} else {
		if s == "" {
			t.setText("tut")
		} else {
			t.setText(fmt.Sprintf("tut - %s", s))
		}
	}
}

func (t *Top) setText(s string) {
	t.View.SetText(s)
	if t.TutView.tut.Config.General.TerminalTitle > 0 {
		util.SetTerminalTitle(s)
	}
}
