package api

import (
	"context"
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/util"
)

type statusToggleFunc func(s *mastodon.Status) (*mastodon.Status, error)

func toggleHelper(s *mastodon.Status, comp bool, on, off statusToggleFunc) (*mastodon.Status, error) {
	if s == nil {
		return nil, fmt.Errorf("no status")
	}
	so := s
	reblogged := false
	if s.Reblog != nil {
		s = s.Reblog
		reblogged = true
	}

	var ns *mastodon.Status
	var err error
	if comp {
		ns, err = off(s)
	} else {
		ns, err = on(s)
	}
	if err != nil {
		return nil, err
	}
	if reblogged {
		so.Reblog = ns
		return so, nil
	}
	return ns, nil
}

func (ac *AccountClient) BoostToggle(s *mastodon.Status) (*mastodon.Status, error) {
	return toggleHelper(s,
		util.StatusOrReblog(s).Reblogged.(bool),
		ac.Boost, ac.Unboost,
	)
}

func (ac *AccountClient) Boost(s *mastodon.Status) (*mastodon.Status, error) {
	return ac.Client.Reblog(context.Background(), s.ID)
}

func (ac *AccountClient) Unboost(s *mastodon.Status) (*mastodon.Status, error) {
	return ac.Client.Unreblog(context.Background(), s.ID)
}

func (ac *AccountClient) FavoriteToogle(s *mastodon.Status) (*mastodon.Status, error) {
	return toggleHelper(s,
		util.StatusOrReblog(s).Favourited.(bool),
		ac.Favorite, ac.Unfavorite,
	)
}

func (ac *AccountClient) Favorite(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := ac.Client.Favourite(context.Background(), s.ID)
	return status, err
}

func (ac *AccountClient) Unfavorite(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := ac.Client.Unfavourite(context.Background(), s.ID)
	return status, err
}

func (ac *AccountClient) BookmarkToogle(s *mastodon.Status) (*mastodon.Status, error) {
	return toggleHelper(s,
		util.StatusOrReblog(s).Bookmarked.(bool),
		ac.Bookmark, ac.Unbookmark,
	)
}

func (ac *AccountClient) Bookmark(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := ac.Client.Bookmark(context.Background(), s.ID)
	return status, err
}

func (ac *AccountClient) Unbookmark(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := ac.Client.Unbookmark(context.Background(), s.ID)
	return status, err
}

func (ac *AccountClient) DeleteStatus(s *mastodon.Status) error {
	return ac.Client.DeleteStatus(context.Background(), util.StatusOrReblog(s).ID)
}

func (ac *AccountClient) GetStatus(id mastodon.ID) (*mastodon.Status, error) {
	return ac.Client.GetStatus(context.Background(), id)
}
