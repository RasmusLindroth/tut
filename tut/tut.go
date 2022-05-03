package tut

import (
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type Tut struct {
	//Config
	Client *api.AccountClient
	App    *tview.Application
	Config *config.Config
}
