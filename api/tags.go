package api

import (
	"context"
	"errors"

	"github.com/RasmusLindroth/go-mastodon"
)

func (ac *AccountClient) FollowTag(tag string) error {
	t, err := ac.Client.TagFollow(context.Background(), tag)
	if err != nil {
		return err
	}
	if t.Following == nil {
		return errors.New("following is set to nil")
	}
	if t.Following == false {
		return errors.New("following is still set to false")
	}
	return nil
}

func (ac *AccountClient) UnfollowTag(tag string) error {
	t, err := ac.Client.TagUnfollow(context.Background(), tag)
	if err != nil {
		return err
	}

	if t.Following == nil {
		return errors.New("following is set to nil")
	}
	if t.Following == true {
		return errors.New("following is still set to true")
	}
	return nil
}

func (ac *AccountClient) TagToggleFollow(tag *mastodon.Tag) (*mastodon.Tag, error) {
	var t *mastodon.Tag
	var err error
	switch tag.Following {
	case true:
		t, err = ac.Client.TagUnfollow(context.Background(), tag.Name)
	case false:
		t, err = ac.Client.TagFollow(context.Background(), tag.Name)
	default:
		t, err = ac.Client.TagFollow(context.Background(), tag.Name)
	}
	return t, err
}
