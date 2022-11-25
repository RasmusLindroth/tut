package auth

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/util"
)

func AddAccount(ad *AccountData) *mastodon.Client {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("You will have to log in to your Mastodon instance to be able")
	fmt.Println("to use Tut. The default protocol is https:// so you won't need")
	fmt.Println("it. E.g. write fosstodon.org and press <Enter>.")
	fmt.Println("--------------------------------------------------------------")

	var server string
	for {
		var err error
		fmt.Print("Instance: ")
		server, err = util.ReadLine(reader)
		if err != nil {
			log.Fatalln(err)
		}
		if !(strings.HasPrefix(server, "https://") || strings.HasPrefix(server, "http://")) {
			server = "https://" + server
		}
		client := mastodon.NewClient(&mastodon.Config{
			Server: server,
		})
		_, err = client.GetInstance(context.Background())
		if err != nil {
			fmt.Printf("\nCouldn't connect to instance %s:\n%s\nTry again or press ^C.\n", server, err)
			fmt.Println("--------------------------------------------------------------")
		} else {
			break
		}
	}
	srv, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:       server,
		ClientName:   "tut-tui",
		Scopes:       "read write follow",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
		Website:      "https://github.com/RasmusLindroth/tut",
	})
	if err != nil {
		fmt.Printf("Couldn't register the app. Error: %v\n\nExiting...\n", err)
		os.Exit(1)
	}

	util.OpenURL(srv.AuthURI)
	fmt.Println("You need to authorize Tut to use your account. Your browser")
	fmt.Println("should've opened. If not you can use the URL below.")
	fmt.Printf("\n%s\n\n", srv.AuthURI)

	var client *mastodon.Client
	for {
		var err error
		fmt.Print("Authorization code: ")
		code, err := util.ReadLine(reader)
		if err != nil {
			log.Fatalln(err)
		}
		client = mastodon.NewClient(&mastodon.Config{
			Server:       server,
			ClientID:     srv.ClientID,
			ClientSecret: srv.ClientSecret,
		})

		err = client.AuthenticateToken(context.Background(), code, "urn:ietf:wg:oauth:2.0:oob")
		if err != nil {
			fmt.Printf("\nError: %v\nTry again or press ^C.\n", err)
			fmt.Println("--------------------------------------------------------------")
		} else {
			break
		}
	}
	me, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		fmt.Printf("\nCouldn't get user. Error: %v\nExiting...\n", err)
		os.Exit(1)
	}
	acc := Account{
		Name:         me.Username,
		Server:       client.Config.Server,
		ClientID:     client.Config.ClientID,
		ClientSecret: client.Config.ClientSecret,
		AccessToken:  client.Config.AccessToken,
	}
	if ad == nil {
		ad = &AccountData{
			Accounts: []Account{acc},
		}
	} else {
		ad.Accounts = append(ad.Accounts, acc)
	}
	path, _, err := util.CheckConfig("accounts.toml")
	if err != nil {
		fmt.Printf("Couldn't open the account file for reading. Error: %v\n", err)
		os.Exit(1)
	}
	err = ad.Save(path)
	if err != nil {
		fmt.Printf("Couldn't update the account file. Error: %v\n", err)
		os.Exit(1)
	}
	return client
}
