package main

import (
	"github.com/mattn/go-mastodon"
)

type App struct {
	UI          *UI
	Me          *mastodon.Account
	API         *API
	Config      *Config
	HaveAccount bool
}
