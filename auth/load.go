package auth

import (
	"log"

	"github.com/RasmusLindroth/tut/util"
)

func StartAuth(newUser bool) *AccountData {
	path, exists, err := util.CheckConfig("accounts.toml")
	if err != nil {
		log.Fatalf("Couldn't open the account file for reading. Error: %v", err)
	}
	var accs *AccountData
	if exists {
		accs, err = GetAccounts(path)
	}
	if err != nil || accs == nil || len(accs.Accounts) == 0 || newUser {
		if err == nil && accs != nil {
			AddAccount(accs)
		} else {
			AddAccount(nil)
		}
		return StartAuth(false)
	}
	return accs
}
