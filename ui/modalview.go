package ui

import (
	"github.com/rivo/tview"
)

type ModalView struct {
	tutView *TutView
	View    *tview.Modal
	res     chan bool
}

func NewModalView(tv *TutView) *ModalView {
	mv := &ModalView{
		tutView: tv,
		View:    tview.NewModal(),
		res:     make(chan bool, 1),
	}
	mv.View.SetText("Are you sure?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				mv.res <- true
			} else {
				mv.res <- false
			}
		})
	return mv
}

func (mv *ModalView) run(text string) (chan bool, func()) {
	mv.View.SetText(text)
	mv.tutView.SetPage(ModalFocus)
	return mv.res, func() {
		mv.tutView.tut.App.QueueUpdateDraw(func() {
			mv.tutView.PrevFocus()
		})
	}
}
func (mv *ModalView) Run(text string, fn func()) {
	if !mv.tutView.tut.Config.General.Confirmation {
		fn()
		return
	}
	r, f := mv.run(text)
	go func() {
		if <-r {
			fn()
		}
		f()
	}()
}

func (mv *ModalView) Stop(fn func()) {
	fn()
}
