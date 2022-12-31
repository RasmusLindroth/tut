package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/auth"
	"github.com/rivo/tview"
)

type LoginView struct {
	tutView     *TutView
	accounts    *auth.AccountData
	View        tview.Primitive
	list        *tview.List
	scrollSleep *scrollSleep
}

func NewLoginView(tv *TutView, accs *auth.AccountData) *LoginView {
	tv.Shared.Top.SetText("select account")
	list := NewList(tv.tut.Config, false)
	for _, a := range accs.Accounts {
		list.AddItem(fmt.Sprintf("%s - %s", a.Name, a.Server), "", 0, nil)
	}

	v := tview.NewFlex().SetDirection(tview.FlexRow)
	if tv.tut.Config.General.TerminalTitle < 2 {
		v.AddItem(tv.Shared.Top.View, 1, 0, false)
	}
	v.AddItem(list, 0, 1, false).
		AddItem(tv.Shared.Bottom.View, 2, 0, false)

	lv := &LoginView{
		tutView:  tv,
		accounts: accs,
		View:     v,
		list:     list,
	}
	lv.scrollSleep = NewScrollSleep(lv.Next, lv.Prev)
	return lv
}

func (l *LoginView) Selected() {
	acc := l.accounts.Accounts[l.list.GetCurrentItem()]
	l.tutView.loggedIn(acc)
}

func (l *LoginView) Next() {
	listNext(l.list)
}

func (l *LoginView) Prev() {
	listPrev(l.list)
}
