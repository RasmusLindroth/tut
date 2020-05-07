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

func NewAuthOverlay(app *App) *AuthOverlay {
	a := &AuthOverlay{
		app:      app,
		Flex:     tview.NewFlex(),
		Input:    tview.NewInputField(),
		Text:     tview.NewTextView(),
		authStep: authNoneStep,
	}

	a.Flex.SetBackgroundColor(app.Config.Style.Background)
	a.Input.SetBackgroundColor(app.Config.Style.Background)
	a.Input.SetFieldBackgroundColor(app.Config.Style.Background)
	a.Input.SetFieldTextColor(app.Config.Style.Text)
	a.Text.SetBackgroundColor(app.Config.Style.Background)
	a.Text.SetTextColor(app.Config.Style.Text)
	a.Flex.SetDrawFunc(app.Config.ClearContent)
	a.Draw()
	return a
}

type AuthOverlay struct {
	app         *App
	Flex        *tview.Flex
	Input       *tview.InputField
	Text        *tview.TextView
	authStep    authStep
	account     AccountRegister
	redirectURL string
}

func (a *AuthOverlay) GotInput() {
	input := strings.TrimSpace(a.Input.GetText())
	switch a.authStep {
	case authInstanceStep:
		if !(strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://")) {
			input = "https://" + input
		}

		_, err := TryInstance(input)
		if err != nil {
			a.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't connect to instance %s\n", input))
			return
		}

		acc, err := Authorize(input)
		if err != nil {
			a.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't authorize. Error: %v\n", err))
			return
		}
		a.account = acc
		openURL(a.app.Config.Media, acc.AuthURI)
		a.Input.SetText("")
		a.authStep = authCodeStep
		a.Draw()
	case authCodeStep:
		client, err := AuthorizationCode(a.account, input)
		if err != nil {
			a.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't verify the code. Error: %v\n", err))
			a.Input.SetText("")
			return
		}
		path, _, err := CheckConfig("accounts.toml")
		if err != nil {
			log.Fatalf("Couldn't open the account file for reading. Error: %v", err)
		}
		ad := AccountData{
			Accounts: []Account{
				{
					Server:       client.Config.Server,
					ClientID:     client.Config.ClientID,
					ClientSecret: client.Config.ClientSecret,
					AccessToken:  client.Config.AccessToken,
				},
			},
		}
		err = ad.Save(path)
		if err != nil {
			log.Fatalf("Couldn't save the account file. Error: %v", err)
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
		a.Input.SetText("")
		a.Draw()
		return
	case authInstanceStep:
		a.Input.SetLabel("Instance: ")
		a.Text.SetText("Enter the url of your instance. Will default to https://\nPress Enter when done")
	case authCodeStep:
		a.Text.SetText(fmt.Sprintf("The login URL has opened in your browser. If it didn't work open this URL\n%s", a.account.AuthURI))
		a.Input.SetLabel("Authorization code: ")
	}
}
