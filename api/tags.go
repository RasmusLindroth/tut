package api

import (
	"context"
	"errors"
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
