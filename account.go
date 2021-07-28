package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/pelletier/go-toml/v2"
)

func GetSecret(s string) string {
	var err error
	if strings.HasPrefix(s, "!CMD!") {
		s, err = CmdToString(s)
		if err != nil {
			log.Fatalf("Couldn't run CMD on auth-file. Error; %v", err)
		}
	}
	return s
}

func GetAccounts(filepath string) (*AccountData, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return &AccountData{}, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return &AccountData{}, err
	}
	accounts := &AccountData{}
	err = toml.Unmarshal(data, accounts)

	for i, acc := range accounts.Accounts {
		accounts.Accounts[i].ClientID = GetSecret(acc.ClientID)
		accounts.Accounts[i].ClientSecret = GetSecret(acc.ClientSecret)
		accounts.Accounts[i].AccessToken = GetSecret(acc.AccessToken)
	}

	return accounts, err
}

type AccountData struct {
	Accounts []Account `yaml:"accounts"`
}

func (ad *AccountData) Save(filepath string) error {
	marshaled, err := toml.Marshal(ad)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(marshaled)
	return err
}

type Account struct {
	Name         string
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

func (a *Account) Login() (*mastodon.Client, error) {
	config := &mastodon.Config{
		Server:       a.Server,
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		AccessToken:  a.AccessToken,
	}
	client := mastodon.NewClient(config)
	_, err := client.GetAccountCurrentUser(context.Background())

	return client, err
}

func TryInstance(server string) (*mastodon.Instance, error) {
	client := mastodon.NewClient(&mastodon.Config{
		Server: server,
	})
	inst, err := client.GetInstance(context.Background())
	return inst, err
}

func Authorize(server string) (AccountRegister, error) {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:       server,
		ClientName:   "tut-tui",
		Scopes:       "read write follow",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
		Website:      "https://github.com/RasmusLindroth/tut",
	})
	if err != nil {
		return AccountRegister{}, err
	}

	acc := AccountRegister{
		Account: Account{
			Server:       server,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
		},
		AuthURI: app.AuthURI,
	}

	return acc, nil
}

func AuthorizationCode(acc AccountRegister, code string) (*mastodon.Client, error) {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       acc.Account.Server,
		ClientID:     acc.Account.ClientID,
		ClientSecret: acc.Account.ClientSecret,
	})

	err := client.AuthenticateToken(context.Background(), code, "urn:ietf:wg:oauth:2.0:oob")
	return client, err
}
