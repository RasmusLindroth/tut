package api

import (
	"context"

	"github.com/RasmusLindroth/go-mastodon"
)

func (ac *AccountClient) Vote(poll *mastodon.Poll, choices ...int) (*mastodon.Poll, error) {
	return ac.Client.PollVote(context.Background(), poll.ID, choices...)
}
