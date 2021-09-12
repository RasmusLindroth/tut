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
	TimelineFavorited
	TimelineList
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

func (api *API) getStatuses(tl TimelineType, listInfo *ListInfo, pg *mastodon.Pagination) ([]*mastodon.Status, mastodon.ID, mastodon.ID, error) {
	var err error
	var statuses []*mastodon.Status

	if pg == nil {
		pg = &mastodon.Pagination{}
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
	case TimelineFavorited:
		statuses, err = api.Client.GetFavourites(context.Background(), pg)
	case TimelineList:
		if listInfo == nil {
			err = errors.New("no list id")
			return statuses, "", "", err
		}
		statuses, err = api.Client.GetTimelineList(context.Background(), listInfo.id, pg)
	default:
		err = errors.New("no timeline selected")
	}

	if err != nil {
		return statuses, "", "", err
	}

	min := mastodon.ID("")
	max := mastodon.ID("")
	if pg != nil {
		min = pg.MinID
		max = pg.MaxID
		if min == "" {
			min = "-1"
		}
		if max == "" {
			max = "-1"
		}
	}
	return statuses, min, max, err
}

func (api *API) GetStatuses(t *TimelineFeed) ([]*mastodon.Status, error) {
	statuses, pgmin, pgmax, err := api.getStatuses(t.timelineType, t.listInfo, nil)
	switch t.timelineType {
	case TimelineBookmarked, TimelineFavorited:
		if err == nil {
			t.linkPrev = pgmin
			t.linkNext = pgmax
		}
	}
	return statuses, err
}

func (api *API) GetStatusesOlder(t *TimelineFeed) ([]*mastodon.Status, error) {
	if len(t.statuses) == 0 {
		return api.GetStatuses(t)
	}

	switch t.timelineType {
	case TimelineBookmarked, TimelineFavorited:
		if t.linkNext == "-1" {
			return []*mastodon.Status{}, nil
		}
		pg := &mastodon.Pagination{
			MaxID: t.linkNext,
		}
		statuses, _, max, err := api.getStatuses(t.timelineType, t.listInfo, pg)
		if err == nil {
			t.linkNext = max
		}
		return statuses, err
	default:
		pg := &mastodon.Pagination{
			MaxID: t.statuses[len(t.statuses)-1].ID,
		}
		statuses, _, _, err := api.getStatuses(t.timelineType, t.listInfo, pg)
		return statuses, err
	}
}

func (api *API) GetStatusesNewer(t *TimelineFeed) ([]*mastodon.Status, error) {
	if len(t.statuses) == 0 {
		return api.GetStatuses(t)
	}

	switch t.timelineType {
	case TimelineBookmarked, TimelineFavorited:
		if t.linkPrev == "-1" {
			return []*mastodon.Status{}, nil
		}
		pg := &mastodon.Pagination{
			MinID: mastodon.ID(t.linkPrev),
		}
		statuses, min, _, err := api.getStatuses(t.timelineType, t.listInfo, pg)
		if err == nil {
			t.linkPrev = min
		}
		return statuses, err
	default:
		pg := &mastodon.Pagination{
			MinID: t.statuses[0].ID,
		}
		statuses, _, _, err := api.getStatuses(t.timelineType, t.listInfo, pg)
		return statuses, err
	}
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

func (api *API) getNotifications(pg *mastodon.Pagination) ([]*Notification, error) {
	var notifications []*Notification

	mnot, err := api.Client.GetNotifications(context.Background(), pg)
	if err != nil {
		return []*Notification{}, err
	}

	for _, np := range mnot {
		var r *mastodon.Relationship
		if np.Type == "follow" {
			r, err = api.GetRelation(&np.Account)
			if err != nil {
				return notifications, err
			}
		}
		notifications = append(notifications, &Notification{N: np, R: r})
	}

	return notifications, err
}

func (api *API) GetNotifications() ([]*Notification, error) {
	return api.getNotifications(nil)
}

func (api *API) GetNotificationsOlder(n *Notification) ([]*Notification, error) {
	pg := &mastodon.Pagination{
		MaxID: n.N.ID,
	}
	return api.getNotifications(pg)
}

func (api *API) GetNotificationsNewer(n *Notification) ([]*Notification, error) {
	pg := &mastodon.Pagination{
		MinID: n.N.ID,
	}
	return api.getNotifications(pg)
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
		r, err := api.GetRelation(u)
		if err != nil {
			return ud, err
		}
		ud = append(ud, &UserData{User: u, Relationship: r})
	}

	return ud, nil
}

func (api *API) GetRelation(u *mastodon.Account) (*mastodon.Relationship, error) {
	return api.UserRelation(*u)
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

func (api *API) GetLists() ([]*mastodon.List, error) {
	return api.Client.GetLists(context.Background())
}

func (api *API) BoostToggle(s *mastodon.Status) (*mastodon.Status, error) {
	if s == nil {
		return nil, fmt.Errorf("no status")
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
		return nil, fmt.Errorf("no status")
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
		return nil, fmt.Errorf("no status")
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
