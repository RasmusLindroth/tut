package auth

type Account struct {
	Name         string
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

type AccountData struct {
	Accounts []Account `yaml:"accounts"`
}
