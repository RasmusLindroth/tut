package api

import (
	"context"

	"github.com/RasmusLindroth/go-mastodon"
)

type TimelineType uint

func (ac *AccountClient) getStatusSimilar(fn func() ([]*mastodon.Status, error), filter string) ([]Item, error) {
	var items []Item
	statuses, err := fn()
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		item := NewStatusItem(s, ac.Filters, filter, false)
		items = append(items, item)
	}
	return items, nil
}

func (ac *AccountClient) getUserSimilar(fn func() ([]*mastodon.Account, error)) ([]Item, error) {
	var items []Item
	users, err := fn()
	if err != nil {
		return items, err
	}
	ids := []string{}
	for _, u := range users {
		ids = append(ids, string(u.ID))
	}
	rel, err := ac.Client.GetAccountRelationships(context.Background(), ids)
	if err != nil {
		return items, err
	}
	for _, u := range users {
		for _, r := range rel {
			if u.ID == r.ID {
				items = append(items, NewUserItem(&User{
					Data:     u,
					Relation: r,
				}, false))
				break
			}
		}
	}
	return items, nil
}

func (ac *AccountClient) GetTimeline(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetTimelineHome(context.Background(), pg)
	}
	return ac.getStatusSimilar(fn, "home")
}

func (ac *AccountClient) GetTimelineFederated(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetTimelinePublic(context.Background(), false, pg)
	}
	return ac.getStatusSimilar(fn, "public")
}

func (ac *AccountClient) GetTimelineLocal(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetTimelinePublic(context.Background(), true, pg)
	}
	return ac.getStatusSimilar(fn, "public")
}

func (ac *AccountClient) GetNotifications(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	notifications, err := ac.Client.GetNotifications(context.Background(), pg)
	if err != nil {
		return items, err
	}
	ids := []string{}
	for _, n := range notifications {
		ids = append(ids, string(n.Account.ID))
	}
	rel, err := ac.Client.GetAccountRelationships(context.Background(), ids)
	if err != nil {
		return items, err
	}
	for _, n := range notifications {
		for _, r := range rel {
			if n.Account.ID == r.ID {
				item := NewNotificationItem(n, &User{
					Data: &n.Account, Relation: r,
				}, ac.Filters)
				items = append(items, item)
				break
			}
		}
	}
	return items, nil
}

func (ac *AccountClient) GetHistory(status *mastodon.Status) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetStatusHistory(context.Background(), status.ID)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusHistoryItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetThread(status *mastodon.Status) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetStatusContext(context.Background(), status.ID)
	if err != nil {
		return items, err
	}
	for _, s := range statuses.Ancestors {
		items = append(items, NewStatusItem(s, ac.Filters, "thread", false))
	}
	items = append(items, NewStatusItem(status, ac.Filters, "thread", false))
	for _, s := range statuses.Descendants {
		items = append(items, NewStatusItem(s, ac.Filters, "thread", false))
	}
	return items, nil
}

func (ac *AccountClient) GetFavorites(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetFavourites(context.Background(), pg)
	}
	return ac.getStatusSimilar(fn, "home")
}

func (ac *AccountClient) GetBookmarks(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetBookmarks(context.Background(), pg)
	}
	return ac.getStatusSimilar(fn, "home")
}

func (ac *AccountClient) GetConversations(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	conversations, err := ac.Client.GetConversations(context.Background(), pg)
	if err != nil {
		return items, err
	}
	for _, c := range conversations {
		item := NewStatusItem(c.LastStatus, ac.Filters, "thread", false)
		items = append(items, item)
	}
	return items, nil
}

func (ac *AccountClient) GetUsers(search string) ([]Item, error) {
	var items []Item
	users, err := ac.Client.AccountsSearch(context.Background(), search, 10)
	if err != nil {
		return items, err
	}
	ids := []string{}
	for _, u := range users {
		ids = append(ids, string(u.ID))
	}
	rel, err := ac.Client.GetAccountRelationships(context.Background(), ids)
	if err != nil {
		return items, err
	}
	for _, u := range users {
		for _, r := range rel {
			if u.ID == r.ID {
				items = append(items, NewUserItem(&User{
					Data:     u,
					Relation: r,
				}, false))
				break
			}
		}
	}
	return items, nil
}

func (ac *AccountClient) GetBoostsStatus(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetRebloggedBy(context.Background(), id, pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetFavoritesStatus(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetFavouritedBy(context.Background(), id, pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetFollowers(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetAccountFollowers(context.Background(), id, pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetFollowing(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetAccountFollowing(context.Background(), id, pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetBlocking(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetBlocks(context.Background(), pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetMuting(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetMutes(context.Background(), pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetFollowRequests(pg *mastodon.Pagination) ([]Item, error) {
	fn := func() ([]*mastodon.Account, error) {
		return ac.Client.GetFollowRequests(context.Background(), pg)
	}
	return ac.getUserSimilar(fn)
}

func (ac *AccountClient) GetUser(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetAccountStatuses(context.Background(), id, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		item := NewStatusItem(s, ac.Filters, "account", false)
		items = append(items, item)
	}
	return items, nil
}

func (ac *AccountClient) GetUserPinned(id mastodon.ID) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetAccountPinnedStatuses(context.Background(), id)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		item := NewStatusItem(s, ac.Filters, "account", true)
		items = append(items, item)
	}
	return items, nil
}

func (ac *AccountClient) GetLists() ([]Item, error) {
	var items []Item
	lists, err := ac.Client.GetLists(context.Background())
	if err != nil {
		return items, err
	}
	for _, l := range lists {
		items = append(items, NewListsItem(l))
	}
	return items, nil
}

func (ac *AccountClient) GetListStatuses(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetTimelineList(context.Background(), id, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		item := NewStatusItem(s, ac.Filters, "home", false)
		items = append(items, item)
	}
	return items, nil
}

func (ac *AccountClient) GetTag(pg *mastodon.Pagination, search string) ([]Item, error) {
	fn := func() ([]*mastodon.Status, error) {
		return ac.Client.GetTimelineHashtag(context.Background(), search, false, pg)
	}
	return ac.getStatusSimilar(fn, "public")
}
