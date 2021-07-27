package main

import (
	"context"
	"log"
	"strings"

	"github.com/mattn/go-mastodon"
)

type App struct {
	UI           *UI
	Me           *mastodon.Account
	API          *API
	Config       *Config
	FullUsername string
	HaveAccount  bool
	Accounts     *AccountData
	FileList     []string
}

func (a *App) Login(index int) {
	if index >= len(a.Accounts.Accounts) {
		log.Fatalln("Tried to login with an account that doesn't exist")
	}
	acc := a.Accounts.Accounts[index]
	client, err := acc.Login()
	if err == nil {
		a.API.SetClient(client)
		a.HaveAccount = true

		me, err := a.API.Client.GetAccountCurrentUser(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		a.Me = me
		if acc.Name == "" {
			a.Accounts.Accounts[index].Name = me.Username

			path, _, err := CheckConfig("accounts.toml")
			if err != nil {
				log.Fatalf("Couldn't open the account file for reading. Error: %v", err)
			}
			err = a.Accounts.Save(path)
			if err != nil {
				log.Fatalf("Couldn't update the account file. Error: %v", err)
			}
		}

		host := strings.TrimPrefix(acc.Server, "https://")
		host = strings.TrimPrefix(host, "http://")
		a.FullUsername = me.Username + "@" + host
	}
}
