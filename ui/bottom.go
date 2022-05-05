package ui

import "github.com/rivo/tview"

type Bottom struct {
	View      tview.Primitive
	StatusBar *StatusBar
	Cmd       *CmdBar
}

func NewBottom(tv *TutView) *Bottom {
	b := &Bottom{
		StatusBar: NewStatusBar(tv),
		Cmd:       NewCmdBar(tv),
	}
	view := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(b.StatusBar.View, 1, 0, false).
		AddItem(b.Cmd.View, 1, 0, false)

	b.View = view
	return b
}
