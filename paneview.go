package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type PaneView interface {
	GetLeftView() tview.Primitive
	GetRightView() tview.Primitive
	Input(event *tcell.EventKey) *tcell.EventKey
}
