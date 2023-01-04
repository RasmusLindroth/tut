package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rivo/uniseg"
)

type EditorView struct {
	tutView   *TutView
	shared    *Shared
	View      *tview.Flex
	editor    *tview.TextArea
	controls  *tview.Flex
	info      *tview.TextView
	limit     int
	prevPage  string
	prevFocus tview.Primitive
	prevInput func(event *tcell.EventKey) *tcell.EventKey
	exitFunc  func(string)
}

func NewEditorView(tv *TutView) *EditorView {
	e := &EditorView{
		tutView:  tv,
		shared:   tv.Shared,
		editor:   NewTextArea(tv.tut.Config),
		controls: NewControlView(tv.tut.Config),
		info:     NewTextView(tv.tut.Config),
	}
	e.View = editorViewUI(e)
	return e
}

func editorViewUI(e *EditorView) *tview.Flex {
	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if e.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(e.shared.Top.View, 1, 0, false)
	}
	if e.limit > 0 {
		r.AddItem(e.info, 1, 0, false).
			AddItem(e.editor, 0, 1, false).
			AddItem(e.controls, 1, 0, false).
			AddItem(e.shared.Bottom.View, 2, 0, false)
	} else {
		r.AddItem(e.editor, 0, 1, false).
			AddItem(e.controls, 1, 0, false).
			AddItem(e.shared.Bottom.View, 2, 0, false)
	}
	return r
}

func (e *EditorView) Init(text string, textLimit int, setReturn bool, exit func(string)) {
	e.editor.SetText(text, true)
	e.limit = textLimit
	*e.View = *editorViewUI(e)
	e.exitFunc = exit
	if setReturn {
		e.prevPage, _ = e.tutView.View.GetFrontPage()
		e.prevFocus = e.tutView.tut.App.GetFocus()
		e.prevInput = e.tutView.tut.App.GetInputCapture()
	}
	e.info.SetText("")
	e.editor.SetChangedFunc(e.updateInfo)
	e.updateInfo()
	e.tutView.View.HidePage(e.prevPage)
	e.tutView.View.ShowPage("editor")
	e.tutView.tut.App.SetInputCapture(e.tutView.InputEditorView)
	e.tutView.tut.App.SetFocus(e.editor)
}

func (e *EditorView) updateInfo() {
	if e.limit > 0 {
		content := e.editor.GetText()
		charCount := uniseg.GraphemeClusterCount(content)
		charsLeft := e.limit - charCount
		e.info.SetText(
			fmt.Sprintf("Chars left: %d", charsLeft),
		)
	}
}

func (e *EditorView) ExitTextAreaInput() {
	e.tutView.View.HidePage("editor")
	e.tutView.View.ShowPage(e.prevPage)
	e.tutView.tut.App.SetInputCapture(e.prevInput)
	e.tutView.tut.App.SetFocus(e.prevFocus)
	e.exitFunc(e.editor.GetText())
}
