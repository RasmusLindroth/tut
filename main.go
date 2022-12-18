package main

import (
	"github.com/RasmusLindroth/tut/auth"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/ui"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

const version = "1.0.27"

func main() {
	util.SetTerminalTitle("tut")
	util.MakeDirs()
	newUser, selectedUser, cnfPath, cnfDir := ui.CliView(version)
	accs := auth.StartAuth(newUser)

	app := tview.NewApplication()
	t := &ui.Tut{
		App:    app,
		Config: config.Load(cnfPath, cnfDir),
	}
	if t.Config.General.MouseSupport {
		app.EnableMouse(true)
	}
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    t.Config.Style.Background,              // background
		ContrastBackgroundColor:     t.Config.Style.Text,                    //background for button, checkbox, form, modal
		MoreContrastBackgroundColor: t.Config.Style.Text,                    //background for dropdown
		BorderColor:                 t.Config.Style.Background,              //border
		TitleColor:                  t.Config.Style.Text,                    //titles
		GraphicsColor:               t.Config.Style.Text,                    //borders
		PrimaryTextColor:            t.Config.Style.StatusBarViewBackground, //backround color selected
		SecondaryTextColor:          t.Config.Style.Text,                    //text
		TertiaryTextColor:           t.Config.Style.Text,                    //list secondary
		InverseTextColor:            t.Config.Style.Text,                    //label activated
		ContrastSecondaryTextColor:  t.Config.Style.Text,                    //foreground on input and prefix on dropdown
	}
	main := ui.NewTutView(t, accs, selectedUser)
	app.SetInputCapture(main.Input)
	if t.Config.General.MouseSupport {
		app.SetMouseCapture(main.MouseInput)
	}
	if err := app.SetRoot(main.View, true).Run(); err != nil {
		panic(err)
	}
}
