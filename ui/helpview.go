package ui

import (
	"bytes"

	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type HelpView struct {
	tutView  *TutView
	shared   *Shared
	View     *tview.Flex
	content  *tview.TextView
	controls *tview.Flex
}

type HelpData struct {
	Style config.Style
}

func NewHelpView(tv *TutView) *HelpView {
	content := NewTextView(tv.tut.Config)
	controls := NewControlView(tv.tut.Config)
	hv := &HelpView{
		tutView:  tv,
		shared:   tv.Shared,
		content:  content,
		controls: controls,
	}
	hd := HelpData{Style: tv.tut.Config.Style}
	var output bytes.Buffer
	err := tv.tut.Config.Templates.Help.ExecuteTemplate(&output, "help.tmpl", hd)
	if err != nil {
		panic(err)
	}
	hv.content.SetText(output.String())
	var items []Control
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.GlobalBack, true))
	items = append(items, NewControl(tv.tut.Config, tv.tut.Config.Input.GlobalExit, true))
	for i, item := range items {
		if i < len(items)-1 {
			hv.controls.AddItem(NewControlButton(hv.tutView.tut.Config, item.Label), item.Len+1, 0, false)
		} else {
			hv.controls.AddItem(NewControlButton(hv.tutView.tut.Config, item.Label), item.Len, 0, false)
		}
	}
	hv.View = newHelpViewUI(hv)
	return hv
}

func newHelpViewUI(hv *HelpView) *tview.Flex {
	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if hv.tutView.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(hv.shared.Top.View, 1, 0, false)
	}
	r.AddItem(hv.content, 0, 1, false).
		AddItem(hv.controls, 1, 0, false).
		AddItem(hv.shared.Bottom.View, 2, 0, false)
	return r
}
