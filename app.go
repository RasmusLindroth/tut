package main

import (
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type App struct {
	App         *tview.Application
	UI          *UI
	Me          *mastodon.Account
	API         *API
	HaveAccount bool
}
