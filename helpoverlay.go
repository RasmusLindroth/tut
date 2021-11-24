package main

import (
	"bytes"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewHelpOverlay(app *App) *HelpOverlay {
	h := &HelpOverlay{
		app:        app,
		Flex:       tview.NewFlex(),
		TextMain:   tview.NewTextView(),
		TextBottom: tview.NewTextView(),
	}

	h.TextMain.SetBackgroundColor(app.Config.Style.Background)
	h.TextMain.SetDynamicColors(true)
	h.TextBottom.SetBackgroundColor(app.Config.Style.Background)
	h.TextBottom.SetDynamicColors(true)
	h.TextBottom.SetText(ColorKey(app.Config, "", "Q", "uit"))
	h.Flex.SetDrawFunc(app.Config.ClearContent)

	hd := HelpData{
		Style: app.Config.Style,
	}
	var output bytes.Buffer
	err := app.Config.Templates.HelpTemplate.ExecuteTemplate(&output, "help.tmpl", hd)
	if err != nil {
		panic(err)
	}
	h.TextMain.SetText(output.String())

	return h
}

type HelpData struct {
	Style StyleConfig
}

type HelpOverlay struct {
	app        *App
	Flex       *tview.Flex
	TextMain   *tview.TextView
	TextBottom *tview.TextView
}

func (h *HelpOverlay) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'q', 'Q':
			h.app.UI.StatusView.giveBackFocus()
			return nil
		}
	} else {
		switch event.Key() {
		case tcell.KeyEsc:
			h.app.UI.StatusView.giveBackFocus()
			return nil
		}
	}
	return event
}
