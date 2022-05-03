package api

import "github.com/RasmusLindroth/go-mastodon"

type RequestData struct {
	MinID mastodon.ID
	MaxID mastodon.ID
}

type AccountClient struct {
	Client  *mastodon.Client
	Streams map[string]*Stream
	Me      *mastodon.Account
}

type User struct {
	Data     *mastodon.Account
	Relation *mastodon.Relationship
}
