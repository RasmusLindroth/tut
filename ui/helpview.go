package ui

import (
	"bytes"
	"strings"

	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type HelpView struct {
	tutView  *TutView
	shared   *Shared
	View     *tview.Flex
	content  *tview.TextView
	controls *tview.TextView
}

type HelpData struct {
	Style config.Style
}

func NewHelpView(tv *TutView) *HelpView {
	content := NewTextView(tv.tut.Config)
	controls := NewTextView(tv.tut.Config)
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
	var items []string
	items = append(items, config.ColorFromKey(tv.tut.Config, tv.tut.Config.Input.GlobalBack, true))
	items = append(items, config.ColorFromKey(tv.tut.Config, tv.tut.Config.Input.GlobalExit, true))
	res := strings.Join(items, " ")
	hv.controls.SetText(res)
	hv.View = newHelpViewUI(hv)
	return hv
}

func newHelpViewUI(hv *HelpView) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(hv.shared.Top.View, 1, 0, false).
		AddItem(hv.content, 0, 1, false).
		AddItem(hv.controls, 1, 0, false).
		AddItem(hv.shared.Bottom.View, 2, 0, false)
}
