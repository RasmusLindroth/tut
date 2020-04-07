package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-mastodon"
)

type TimelineType uint

const (
	TimelineHome TimelineType = iota
	TimelineDirect
	TimelineLocal
	TimelineFederated
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

func (api *API) getStatuses(tl TimelineType, pg *mastodon.Pagination) ([]*mastodon.Status, error) {
	var statuses []*mastodon.Status
	var err error

	switch tl {
	case TimelineHome:
		statuses, err = api.Client.GetTimelineHome(context.Background(), pg)
	case TimelineDirect:
		statuses, err = api.Client.GetTimelineDirect(context.Background(), pg)
	case TimelineLocal:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), true, pg)
	case TimelineFederated:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), false, pg)
	default:
		err = errors.New("No timeline selected")
	}

	return statuses, err
}

func (api *API) GetStatuses(tl TimelineType) ([]*mastodon.Status, error) {
	return api.getStatuses(tl, nil)
}

func (api *API) GetStatusesOlder(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	statuses, err := api.getStatuses(tl, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MinID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetStatusesNewer(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	statuses, err := api.getStatuses(tl, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MaxID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetThread(s *mastodon.Status) ([]*mastodon.Status, int, error) {
	cont, err := api.Client.GetStatusContext(context.Background(), s.ID)
	if err != nil {
		return nil, 0, err
	}
	thread := cont.Ancestors
	thread = append(thread, s)
	thread = append(thread, cont.Descendants...)
	return thread, len(cont.Ancestors), nil
}

func (api *API) GetUserStatuses(u mastodon.Account) ([]*mastodon.Status, error) {
	return api.Client.GetAccountStatuses(context.Background(), u.ID, nil)
}

func (api *API) GetUserStatusesOlder(u mastodon.Account, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	statuses, err := api.Client.GetAccountStatuses(context.Background(), u.ID, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MinID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetUserStatusesNewer(u mastodon.Account, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	statuses, err := api.Client.GetAccountStatuses(context.Background(), u.ID, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MaxID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetNotifications() ([]*mastodon.Notification, error) {
	return api.Client.GetNotifications(context.Background(), nil)
}

func (api *API) GetNotificationsOlder(n *mastodon.Notification) ([]*mastodon.Notification, bool, error) {
	pg := &mastodon.Pagination{
		MaxID: n.ID,
	}

	statuses, err := api.Client.GetNotifications(context.Background(), pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MinID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetNotificationsNewer(n *mastodon.Notification) ([]*mastodon.Notification, bool, error) {
	pg := &mastodon.Pagination{
		MinID: n.ID,
	}

	statuses, err := api.Client.GetNotifications(context.Background(), pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MaxID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) BoostToggle(s *mastodon.Status) (*mastodon.Status, error) {
	if s == nil {
		return nil, fmt.Errorf("No status")
	}

	if s.Reblogged == true {
		return api.Unboost(s)
	}
	return api.Boost(s)
}

func (api *API) Boost(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Reblog(context.Background(), s.ID)
	return status, err
}

func (api *API) Unboost(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Unreblog(context.Background(), s.ID)
	return status, err
}

func (api *API) FavoriteToogle(s *mastodon.Status) (*mastodon.Status, error) {
	if s == nil {
		return nil, fmt.Errorf("No status")
	}

	if s.Favourited == true {
		return api.Unfavorite(s)
	}
	return api.Favorite(s)
}

func (api *API) Favorite(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Favourite(context.Background(), s.ID)
	return status, err
}

func (api *API) Unfavorite(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Unfavourite(context.Background(), s.ID)
	return status, err
}

func (api *API) DeleteStatus(s *mastodon.Status) error {
	//TODO: check user here?
	return api.Client.DeleteStatus(context.Background(), s.ID)
}

func (api *API) UserRelation(u mastodon.Account) (*mastodon.Relationship, error) {
	relations, err := api.Client.GetAccountRelationships(context.Background(), []string{string(u.ID)})

	if err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}
	return relations[0], nil
}

func (api *API) FollowToggle(u mastodon.Account) (*mastodon.Relationship, error) {
	relation, err := api.UserRelation(u)
	if err != nil {
		return nil, err
	}
	if relation.Following {
		return api.UnfollowUser(u)
	}
	return api.FollowUser(u)
}

func (api *API) FollowUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountFollow(context.Background(), u.ID)
}

func (api *API) UnfollowUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountUnfollow(context.Background(), u.ID)
}

func (api *API) BlockToggle(u mastodon.Account) (*mastodon.Relationship, error) {
	relation, err := api.UserRelation(u)
	if err != nil {
		return nil, err
	}
	if relation.Blocking {
		return api.UnblockUser(u)
	}
	return api.BlockUser(u)
}

func (api *API) BlockUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountBlock(context.Background(), u.ID)
}

func (api *API) UnblockUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountUnblock(context.Background(), u.ID)
}

func (api *API) MuteToggle(u mastodon.Account) (*mastodon.Relationship, error) {
	relation, err := api.UserRelation(u)
	if err != nil {
		return nil, err
	}
	if relation.Blocking {
		return api.UnmuteUser(u)
	}
	return api.MuteUser(u)
}

func (api *API) MuteUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountMute(context.Background(), u.ID)
}

func (api *API) UnmuteUser(u mastodon.Account) (*mastodon.Relationship, error) {
	return api.Client.AccountUnmute(context.Background(), u.ID)
}
