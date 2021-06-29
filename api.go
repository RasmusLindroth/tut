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
	TimelineBookmarked
)

type UserListType uint

const (
	UserListSearch UserListType = iota
	UserListBoosts
	UserListFavorites
	UserListFollowers
	UserListFollowing
	UserListBlocking
	UserListMuting
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

	var pgMin = mastodon.ID("")
	var pgMax = mastodon.ID("")
	if pg != nil {
		pgMin = pg.MinID
		pgMax = pg.MaxID
	}

	switch tl {
	case TimelineHome:
		statuses, err = api.Client.GetTimelineHome(context.Background(), pg)
	case TimelineDirect:
		var conv []*mastodon.Conversation
		conv, err = api.Client.GetConversations(context.Background(), pg)
		var cStatuses []*mastodon.Status
		for _, c := range conv {
			cStatuses = append(cStatuses, c.LastStatus)
		}
		statuses = cStatuses
	case TimelineLocal:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), true, pg)
	case TimelineFederated:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), false, pg)
	case TimelineBookmarked:
		statuses, err = api.Client.GetBookmarks(context.Background(), pg)
	default:
		err = errors.New("No timeline selected")
	}

	if err != nil {
		return statuses, err
	}

	if pg != nil && len(statuses) > 0 {
		if pgMin != "" && statuses[0].ID == pgMin {
			return []*mastodon.Status{}, nil
		} else if pgMax != "" && statuses[len(statuses)-1].ID == pgMax {
			return []*mastodon.Status{}, nil
		}
	}

	return statuses, err
}

func (api *API) GetStatuses(tl TimelineType) ([]*mastodon.Status, error) {
	return api.getStatuses(tl, nil)
}

func (api *API) GetStatusesOlder(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	return api.getStatuses(tl, pg)
}

func (api *API) GetStatusesNewer(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	return api.getStatuses(tl, pg)
}

func (api *API) GetTags(tag string) ([]*mastodon.Status, error) {
	return api.Client.GetTimelineHashtag(context.Background(), tag, false, nil)
}

func (api *API) GetTagsOlder(tag string, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	return api.Client.GetTimelineHashtag(context.Background(), tag, false, pg)
}

func (api *API) GetTagsNewer(tag string, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	return api.Client.GetTimelineHashtag(context.Background(), tag, false, pg)
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

func (api *API) GetUserStatusesOlder(u mastodon.Account, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	return api.Client.GetAccountStatuses(context.Background(), u.ID, pg)
}

func (api *API) GetUserStatusesNewer(u mastodon.Account, s *mastodon.Status) ([]*mastodon.Status, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	return api.Client.GetAccountStatuses(context.Background(), u.ID, pg)
}

func (api *API) GetNotifications() ([]*mastodon.Notification, error) {
	return api.Client.GetNotifications(context.Background(), nil)
}

func (api *API) GetNotificationsOlder(n *mastodon.Notification) ([]*mastodon.Notification, error) {
	pg := &mastodon.Pagination{
		MaxID: n.ID,
	}

	return api.Client.GetNotifications(context.Background(), pg)
}

func (api *API) GetNotificationsNewer(n *mastodon.Notification) ([]*mastodon.Notification, error) {
	pg := &mastodon.Pagination{
		MinID: n.ID,
	}

	return api.Client.GetNotifications(context.Background(), pg)
}

type UserData struct {
	User         *mastodon.Account
	Relationship *mastodon.Relationship
}

func (api *API) GetUsers(s string) ([]*UserData, error) {
	var ud []*UserData
	users, err := api.Client.AccountsSearch(context.Background(), s, 10)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		r, err := api.UserRelation(*u)
		if err != nil {
			return ud, err
		}
		ud = append(ud, &UserData{User: u, Relationship: r})
	}

	return ud, nil
}

func (api *API) getUserList(t UserListType, id string, pg *mastodon.Pagination) ([]*UserData, error) {

	var ud []*UserData
	var users []*mastodon.Account
	var err error
	var pgMin = mastodon.ID("")
	var pgMax = mastodon.ID("")
	if pg != nil {
		pgMin = pg.MinID
		pgMax = pg.MinID
	}

	switch t {
	case UserListSearch:
		users, err = api.Client.AccountsSearch(context.Background(), id, 10)
	case UserListBoosts:
		users, err = api.Client.GetRebloggedBy(context.Background(), mastodon.ID(id), pg)
	case UserListFavorites:
		users, err = api.Client.GetFavouritedBy(context.Background(), mastodon.ID(id), pg)
	case UserListFollowers:
		users, err = api.Client.GetAccountFollowers(context.Background(), mastodon.ID(id), pg)
	case UserListFollowing:
		users, err = api.Client.GetAccountFollowing(context.Background(), mastodon.ID(id), pg)
	case UserListBlocking:
		users, err = api.Client.GetBlocks(context.Background(), pg)
	case UserListMuting:
		users, err = api.Client.GetMutes(context.Background(), pg)
	}

	if err != nil {
		return ud, err
	}

	if pg != nil && len(users) > 0 {
		if pgMin != "" && users[0].ID == pgMin {
			return ud, nil
		} else if pgMax != "" && users[len(users)-1].ID == pgMax {
			return ud, nil
		}
	}

	for _, u := range users {
		r, err := api.UserRelation(*u)
		if err != nil {
			return ud, err
		}
		ud = append(ud, &UserData{User: u, Relationship: r})
	}
	return ud, nil
}

func (api *API) GetUserList(t UserListType, id string) ([]*UserData, error) {
	return api.getUserList(t, id, nil)
}

func (api *API) GetUserListOlder(t UserListType, id string, user *mastodon.Account) ([]*UserData, error) {
	if t == UserListSearch {
		return []*UserData{}, nil
	}
	pg := &mastodon.Pagination{
		MaxID: user.ID,
	}
	return api.getUserList(t, id, pg)
}

func (api *API) GetUserListNewer(t UserListType, id string, user *mastodon.Account) ([]*UserData, error) {
	if t == UserListSearch {
		return []*UserData{}, nil
	}
	pg := &mastodon.Pagination{
		MinID: user.ID,
	}
	return api.getUserList(t, id, pg)
}

func (api *API) GetUserByID(id mastodon.ID) (*mastodon.Account, error) {
	a, err := api.Client.GetAccount(context.Background(), id)
	return a, err
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

func (api *API) BookmarkToogle(s *mastodon.Status) (*mastodon.Status, error) {
	if s == nil {
		return nil, fmt.Errorf("No status")
	}

	if s.Bookmarked == true {
		return api.Unbookmark(s)
	}
	return api.Bookmark(s)
}

func (api *API) Bookmark(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Bookmark(context.Background(), s.ID)
	return status, err
}

func (api *API) Unbookmark(s *mastodon.Status) (*mastodon.Status, error) {
	status, err := api.Client.Unbookmark(context.Background(), s.ID)
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
