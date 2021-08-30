package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const version string = "0.0.29"

func main() {
	newUser := false
	selectedUser := ""
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "example-config":
			CreateDefaultConfig("./config.example.ini")
			os.Exit(0)
		case "--new-user", "-n":
			newUser = true
		case "--user", "-u":
			if len(os.Args) > 2 {
				name := os.Args[2]
				selectedUser = strings.TrimSpace(name)
			} else {
				log.Fatalln("--user/-u must be followed by a user name. Like -u tut")
			}
		case "--help", "-h":
			fmt.Print("tut - a TUI for Mastodon with vim inspired keys.\n\n")
			fmt.Print("Usage:\n\n")
			fmt.Print("\tTo run the program you just have to write tut\n\n")

			fmt.Print("Commands:\n\n")
			fmt.Print("\texample-config - creates the default configuration file in the current directory and names it ./config.example.ini\n\n")

			fmt.Print("Flags:\n\n")
			fmt.Print("\t--help -h - prints this message\n")
			fmt.Print("\t--version -v - prints the version\n")
			fmt.Print("\t--new-user -n - add one more user to tut\n")
			fmt.Print("\t--user <name> -u <name> - login directly to user namde <name>\n")
			fmt.Print("\t\tDon't use a = between --user and the <name>\n")
			fmt.Print("\t\tIf two users are named the same. Use full name like tut@fosstodon.org\n\n")

			fmt.Print("Configuration:\n\n")
			fmt.Printf("\tThe config is located in XDG_CONFIG_HOME/tut/config.ini which usally equals to ~/.config/tut/config.ini.\n")
			fmt.Printf("\tThe program will generate the file the first time you run tut. The file has comments which exmplains what each configuration option does.\n\n")

			fmt.Print("Contact info for issues or questions:\n\n")
			fmt.Printf("\t@rasmus@mastodon.acc.sunet.se\n\trasmus@lindroth.xyz\n")
			fmt.Printf("\thttps://github.com/RasmusLindroth/tut\n")
			os.Exit(0)
		case "--version", "-v":
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
		Accounts:    &AccountData{},
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
		app.Accounts, err = GetAccounts(path)
		if err != nil {
			log.Fatalln(
				fmt.Sprintf("Couldn't access accounts.toml. Error: %v", err),
			)
		}
		if len(app.Accounts.Accounts) == 1 && !newUser {
			app.Login(0)
		}
	}

	if len(app.Accounts.Accounts) > 1 && !newUser {
		if selectedUser != "" {
			useHost := false
			found := false
			if strings.Contains(selectedUser, "@") {
				useHost = true
			}
			for i, acc := range app.Accounts.Accounts {
				accName := acc.Name
				if useHost {
					host := strings.TrimPrefix(acc.Server, "https://")
					host = strings.TrimPrefix(host, "http://")
					accName += "@" + host
				}
				if accName == selectedUser {
					app.Login(i)
					app.UI.LoggedIn()
					found = true
				}
			}
			if found == false {
				log.Fatalf("Couldn't find a user named %s. Try again", selectedUser)
			}
		} else {
			app.UI.SetFocus(UserSelectFocus)
		}
	} else if !app.HaveAccount || newUser {
		app.UI.SetFocus(AuthOverlayFocus)
	} else {
		app.UI.LoggedIn()
	}

	app.FileList = []string{}

	app.UI.Root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !app.HaveAccount {
			if app.UI.Focus == UserSelectFocus {
				app.UI.UserSelectOverlay.InputHandler(event)
				return nil
			} else {
				if event.Key() == tcell.KeyRune {
					switch event.Rune() {
					}
				}
				return event
			}
		}

		if app.UI.Focus == LinkOverlayFocus {
			app.UI.LinkOverlay.InputHandler(event)
			return nil
		}

		if app.UI.Focus == VisibilityOverlayFocus {
			app.UI.VisibilityOverlay.InputHandler(event)
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
				case 'i', 'I':
					app.UI.MessageBox.IncludeQuote()
					return nil
				case 'm', 'M':
					app.UI.SetFocus(MessageAttachmentFocus)
					return nil
				case 'v', 'V':
					app.UI.SetFocus(VisibilityOverlayFocus)
					return nil
				case 'q', 'Q':
					app.UI.SetFocus(LeftPaneFocus)
					return nil
				}
			} else {
				switch event.Key() {
				case tcell.KeyEsc:
					if app.UI.StatusView.lastList == NotificationPaneFocus {
						app.UI.SetFocus(NotificationPaneFocus)
					} else {
						app.UI.SetFocus(LeftPaneFocus)
					}
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
				case 'e', 'E':
					app.UI.MediaOverlay.EditDesc()
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

		if app.UI.Focus == LeftPaneFocus || app.UI.Focus == RightPaneFocus || app.UI.Focus == NotificationPaneFocus {
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

	app.UI.CmdBar.Input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		words := strings.Split(":blocking,:boosts,:bookmarks,:compose,:favorites,:favorited,:muting,:profile,:saved,:tag,:timeline,:tl,:user,:quit,:q", ",")
		if currentText == "" {
			return
		}

		if len(currentText) > 2 && currentText[:3] == ":tl" {
			words = strings.Split(":tl home,:tl notifications,:tl local,:tl federated,:tl direct,:tl favorited", ",")
		}
		if len(currentText) > 8 && currentText[:9] == ":timeline" {
			words = strings.Split(":timeline home,:timeline notifications,:timeline local,:timeline federated,:timeline direct,:timeline favorited", ",")
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
