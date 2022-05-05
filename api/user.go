package api

import (
	"context"
	"fmt"

	"github.com/RasmusLindroth/go-mastodon"
)

func (ac *AccountClient) GetUserByID(id mastodon.ID) (Item, error) {
	var item Item
	acc, err := ac.Client.GetAccount(context.Background(), id)
	if err != nil {
		return nil, err
	}
	rel, err := ac.Client.GetAccountRelationships(context.Background(), []string{string(acc.ID)})
	if err != nil {
		return nil, err
	}
	if len(rel) == 0 {
		return nil, fmt.Errorf("couldn't find user relationship")
	}
	item = NewUserItem(&User{
		Data:     acc,
		Relation: rel[0],
	}, false)
	return item, nil
}

func (ac *AccountClient) FollowToggle(u *User) (*mastodon.Relationship, error) {
	if u.Relation.Following {
		return ac.UnfollowUser(u.Data)
	}
	return ac.FollowUser(u.Data)
}

func (ac *AccountClient) FollowUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountFollow(context.Background(), u.ID)
}

func (ac *AccountClient) UnfollowUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountUnfollow(context.Background(), u.ID)
}

func (ac *AccountClient) BlockToggle(u *User) (*mastodon.Relationship, error) {
	if u.Relation.Blocking {
		return ac.UnblockUser(u.Data)
	}
	return ac.BlockUser(u.Data)
}

func (ac *AccountClient) BlockUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountBlock(context.Background(), u.ID)
}

func (ac *AccountClient) UnblockUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountUnblock(context.Background(), u.ID)
}

func (ac *AccountClient) MuteToggle(u *User) (*mastodon.Relationship, error) {
	if u.Relation.Blocking {
		return ac.UnmuteUser(u.Data)
	}
	return ac.MuteUser(u.Data)
}

func (ac *AccountClient) MuteUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountMute(context.Background(), u.ID)
}

func (ac *AccountClient) UnmuteUser(u *mastodon.Account) (*mastodon.Relationship, error) {
	return ac.Client.AccountUnmute(context.Background(), u.ID)
}
