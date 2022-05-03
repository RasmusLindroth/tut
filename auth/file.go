package auth

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/RasmusLindroth/tut/util"
	"github.com/pelletier/go-toml/v2"
)

func GetSecret(s string) string {
	var err error
	if strings.HasPrefix(s, "!CMD!") {
		s, err = util.CmdToString(s)
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
