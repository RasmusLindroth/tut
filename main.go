package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

const version string = "0.0.10"

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "example-config":
			CreateDefaultConfig("./config.example.ini")
			os.Exit(0)
		case "--help":
		case "-h":
			fmt.Print("tut - a TUI for Mastodon with vim inspired keys.\n\n")
			fmt.Print("Usage:\n\n")
			fmt.Print("\tTo run the program you just have to write tut\n\n")

			fmt.Print("Commands:\n\n")
			fmt.Print("\texample-config - creates the default configuration file in the current directory and names it ./config.example.ini\n\n")

			fmt.Print("Flags:\n\n")
			fmt.Print("\t--help -h - prints this message\n")
			fmt.Print("\t--version -v - prints the version\n\n")

			fmt.Print("Configuration:\n\n")
			fmt.Printf("\tThe config is located in XDG_CONFIG_HOME/tut/config.ini which usally equals to ~/.config/tut/config.ini.\n")
			fmt.Printf("\tThe program will generate the file the first time you run tut. The file has comments which exmplains what each configuration option does.\n\n")

			fmt.Print("Contact info for issues or questions:\n\n")
			fmt.Printf("\t@rasmus@mastodon.acc.sunet.se\n\trasmus@lindroth.xyz\n")
			fmt.Printf("\thttps://github.com/RasmusLindroth/tut\n")
			os.Exit(0)
		case "--version":
		case "-v":
			fmt.Printf("tut version %s\n\n", version)
			fmt.Printf("https://github.com/RasmusLindroth/tut\n")
			os.Exit(0)
		}
	}

	err := CreateConfigDir()
	if err != nil {
		log.Fatalln(
			fmt.Sprintf("Couldn't create or access the configuration dir. Error: %v", err),
		)
	}
	path, exists, err := CheckConfig("config.ini")
	if err != nil {
		log.Fatalln(
			fmt.Sprintf("Couldn't access config.ini. Error: %v", err),
		)
	}
	if !exists {
		err = CreateDefaultConfig(path)
		if err != nil {
			log.Fatalf("Couldn't create default config. Error: %v", err)
		}
	}
	config, err := ParseConfig(path)
	if err != nil {
		log.Fatalf("Couldn't open or parse the config. Error: %v", err)
	}

	app := &App{
		API:         &API{},
		HaveAccount: false,
		Config:      &config,
	}

	app.UI = NewUI(app)
	app.UI.Init()

	path, exists, err = CheckConfig("accounts.toml")
	if err != nil {
		log.Fatalln(
			fmt.Sprintf("Couldn't access accounts.toml. Error: %v", err),
		)
	}

	if exists {
		accounts, err := GetAccounts(path)
		if err != nil {
			log.Fatalln(
				fmt.Sprintf("Couldn't access accounts.toml. Error: %v", err),
			)
		}
		if len(accounts.Accounts) > 0 {
			a := accounts.Accounts[0]
			client, err := a.Login()
			if err == nil {
				app.API.SetClient(client)
				app.HaveAccount = true

				me, err := app.API.Client.GetAccountCurrentUser(context.Background())
				if err != nil {
					log.Fatalln(err)
				}
				app.Me = me
			}
		}
	}

	if !app.HaveAccount {
		app.UI.SetFocus(AuthOverlayFocus)
	} else {
		app.UI.LoggedIn()
	}

	app.FileList = []string{}

	app.UI.Root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !app.HaveAccount {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				}
			}
			return event
		}

		if app.UI.Focus == LinkOverlayFocus {
			app.UI.LinkOverlay.InputHandler(event)
			return nil
		}

		if app.UI.Focus == CmdBarFocus {
			switch event.Key() {
			case tcell.KeyEnter:
				app.UI.CmdBar.DoneFunc(tcell.KeyEnter)
			case tcell.KeyEsc:
				app.UI.CmdBar.ClearInput()
				app.UI.SetFocus(LeftPaneFocus)
				return nil
			}
			return event
		}

		if app.UI.Focus == MessageFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'p', 'P':
					app.UI.MessageBox.Post()
					return nil
				case 'e', 'E':
					app.UI.MessageBox.EditText()
					return nil
				case 'c', 'C':
					app.UI.MessageBox.EditSpoiler()
					return nil
				case 't', 'T':
					app.UI.MessageBox.ToggleSpoiler()
					return nil
				case 'm', 'M':
					app.UI.SetFocus(MessageAttachmentFocus)
					return nil
				case 'q', 'Q':
					app.UI.SetFocus(LeftPaneFocus)
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyEsc:
					app.UI.SetFocus(LeftPaneFocus)
					return nil
				}
			}
			return event
		}

		if app.UI.Focus == MessageAttachmentFocus && app.UI.MediaOverlay.Focus == MediaFocusOverview {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'j', 'J':
					app.UI.MediaOverlay.Next()
				case 'k', 'K':
					app.UI.MediaOverlay.Prev()
				case 'd', 'D':
					app.UI.MediaOverlay.Delete()
				case 'a', 'A':
					app.UI.MediaOverlay.SetFocus(MediaFocusAdd)
				case 'q', 'Q':
					app.UI.SetFocus(MessageFocus)
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyUp:
					app.UI.MediaOverlay.Prev()
				case tcell.KeyDown:
					app.UI.MediaOverlay.Next()
				case tcell.KeyEsc:
					app.UI.SetFocus(MessageFocus)
					return nil
				}
			}
			return event
		}

		if app.UI.Focus == MessageAttachmentFocus && app.UI.MediaOverlay.Focus == MediaFocusAdd {
			if event.Key() == tcell.KeyRune {
				app.UI.MediaOverlay.InputField.AddRune(event.Rune())
				return nil
			}
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyDown:
				app.UI.MediaOverlay.InputField.AutocompleteNext()
				return nil
			case tcell.KeyBacktab, tcell.KeyUp:
				app.UI.MediaOverlay.InputField.AutocompletePrev()
				return nil
			case tcell.KeyEnter:
				app.UI.MediaOverlay.InputField.CheckDone()
				return nil
			case tcell.KeyEsc:
				app.UI.MediaOverlay.SetFocus(MediaFocusOverview)
			}
			return event
		}

		if app.UI.Focus == LeftPaneFocus || app.UI.Focus == RightPaneFocus {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case ':':
					app.UI.CmdBar.ClearInput()
					app.UI.CmdBar.Input.SetText(":")
					app.UI.SetFocus(CmdBarFocus)
					return nil
				}
			}
			return app.UI.StatusView.Input(event)
		}

		return event
	})

	app.UI.MediaOverlay.InputField.View.SetChangedFunc(
		app.UI.MediaOverlay.InputField.HandleChanges,
	)

	app.UI.CmdBar.Input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		return event
	})

	app.UI.CmdBar.Input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		words := strings.Split(":blocking,:boosts,:compose,:favorites,:muting,:profile,:tag,:timeline,:tl,:user,:quit,:q", ",")
		if currentText == "" {
			return
		}

		if len(currentText) > 2 && currentText[:3] == ":tl" {
			words = strings.Split(":tl home,:tl notifications,:tl local,:tl federated,:tl direct", ",")
		}
		if len(currentText) > 8 && currentText[:9] == ":timeline" {
			words = strings.Split(":timeline home,:timeline notifications,:timeline local,:timeline federated,:timeline direct", ",")
		}

		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) < 1 {
			entries = nil
		}
		return
	})

	app.UI.AuthOverlay.Input.SetDoneFunc(func(key tcell.Key) {
		app.UI.AuthOverlay.GotInput()
	})

	if err := app.UI.Root.SetRoot(app.UI.Pages, true).Run(); err != nil {
		panic(err)
	}

	for _, f := range app.FileList {
		os.Remove(f)
	}
}
