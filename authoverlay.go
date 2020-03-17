package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/tview"
)

type authStep int

const (
	authNoneStep authStep = iota
	authInstanceStep
	authCodeStep
)

func NewAuthoverlay(app *App, flex *tview.Flex, view *tview.InputField, controls *Controls) *AuthOverlay {
	return &AuthOverlay{
		app:      app,
		Flex:     flex,
		View:     view,
		Controls: controls,
		authStep: authNoneStep,
	}
}

type AuthOverlay struct {
	app         *App
	Flex        *tview.Flex
	View        *tview.InputField
	Controls    *Controls
	authStep    authStep
	account     AccountRegister
	redirectURL string
}

func (a *AuthOverlay) GotInput() {
	input := strings.TrimSpace(a.View.GetText())
	switch a.authStep {
	case authInstanceStep:
		if !(strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://")) {
			input = "https://" + input
		}

		_, err := TryInstance(input)
		if err != nil {
			log.Fatalf("Couldn'n connect to instance %s\n", input)
		}

		acc, err := Authorize(input)
		if err != nil {
			log.Fatalln(err)
		}
		a.account = acc
		openURL(acc.AuthURI)
		a.View.SetText("")
		a.authStep = authCodeStep
		a.Draw()
	case authCodeStep:
		client, err := AuthorizationCode(a.account, input)
		if err != nil {
			log.Fatalln(err)
		}
		path, _, err := CheckConfig("accounts.toml")
		if err != nil {
			log.Fatalln(err)
		}
		ad := AccountData{
			Accounts: []Account{
				Account{
					Server:       client.Config.Server,
					ClientID:     client.Config.ClientID,
					ClientSecret: client.Config.ClientSecret,
					AccessToken:  client.Config.AccessToken,
				},
			},
		}
		err = ad.Save(path)
		if err != nil {
			log.Fatalln(err)
		}
		a.app.API.SetClient(client)
		a.app.HaveAccount = true
		a.app.UI.LoggedIn()
	}
}

func (a *AuthOverlay) Draw() {
	switch a.authStep {
	case authNoneStep:
		a.authStep = authInstanceStep
		a.View.SetText("")
		a.Draw()
		return
	case authInstanceStep:
		a.View.SetLabel("Instance: ")
		a.Controls.View.SetText("Enter the url of your instance. Will default to https://\nPress Enter when done")
	case authCodeStep:
		a.Controls.View.SetText(fmt.Sprintf("The login URL has opened in your browser. If it didn't work open this URL\n%s", a.account.AuthURI))
		a.View.SetLabel("Authorization code: ")
	}
}
