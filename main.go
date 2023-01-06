package main

import (
	"strings"

	"github.com/RasmusLindroth/tut/auth"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/ui"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

const version = "1.0.35"

var tutViews []*ui.TutView

func main() {
	util.SetTerminalTitle("tut")
	util.MakeDirs()
	newUser, selectedUser, cnfPath, cnfDir := ui.CliView(version)
	accs := auth.StartAuth(newUser)

	app := tview.NewApplication()
	cnf := config.Load(cnfPath, cnfDir)

	if cnf.General.MouseSupport {
		app.EnableMouse(true)
	}
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    cnf.Style.Background,              // background
		ContrastBackgroundColor:     cnf.Style.Text,                    //background for button, checkbox, form, modal
		MoreContrastBackgroundColor: cnf.Style.Text,                    //background for dropdown
		BorderColor:                 cnf.Style.Background,              //border
		TitleColor:                  cnf.Style.Text,                    //titles
		GraphicsColor:               cnf.Style.Text,                    //borders
		PrimaryTextColor:            cnf.Style.StatusBarViewBackground, //backround color selected
		SecondaryTextColor:          cnf.Style.Text,                    //text
		TertiaryTextColor:           cnf.Style.Text,                    //list secondary
		InverseTextColor:            cnf.Style.Text,                    //label activated
		ContrastSecondaryTextColor:  cnf.Style.Text,                    //foreground on input and prefix on dropdown
	}
	ui.SetVars(cnf, app, accs)
	users := strings.Fields(selectedUser)
	if len(users) > 0 {
		for _, user := range strings.Fields(selectedUser) {
			ui.NewTutView(user)
		}
	} else {
		ui.NewTutView(selectedUser)
	}
	ui.DoneAdding()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
