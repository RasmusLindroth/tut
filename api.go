package main

import (
	"github.com/mattn/go-mastodon"
)

type API struct {
	Client *mastodon.Client
}

type AccountRegister struct {
	Account
	AuthURI string
}

func (api *API) SetClient(c *mastodon.Client) {
	api.Client = c
}
