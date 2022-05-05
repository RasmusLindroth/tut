package api

import (
	"context"

	"github.com/RasmusLindroth/go-mastodon"
)

type TimelineType uint

func (ac *AccountClient) GetTimeline(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetTimelineHome(context.Background(), pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetTimelineFederated(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetTimelinePublic(context.Background(), false, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetTimelineLocal(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetTimelinePublic(context.Background(), true, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
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
				items = append(items, NewNotificationItem(n, &User{
					Data:     &n.Account,
					Relation: r,
				}))
				break
			}
		}
	}
	return items, nil
}

func (ac *AccountClient) GetThread(status *mastodon.Status) ([]Item, int, error) {
	var items []Item
	statuses, err := ac.Client.GetStatusContext(context.Background(), status.ID)
	if err != nil {
		return items, 0, err
	}
	for _, s := range statuses.Ancestors {
		items = append(items, NewStatusItem(s))
	}
	items = append(items, NewStatusItem(status))
	for _, s := range statuses.Descendants {
		items = append(items, NewStatusItem(s))
	}
	return items, len(statuses.Ancestors), nil
}

func (ac *AccountClient) GetFavorites(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetFavourites(context.Background(), pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetBookmarks(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetBookmarks(context.Background(), pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetConversations(pg *mastodon.Pagination) ([]Item, error) {
	var items []Item
	conversations, err := ac.Client.GetConversations(context.Background(), pg)
	if err != nil {
		return items, err
	}
	for _, c := range conversations {
		items = append(items, NewStatusItem(c.LastStatus))
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

func (ac *AccountClient) GetUser(pg *mastodon.Pagination, id mastodon.ID) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetAccountStatuses(context.Background(), id, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
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
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}

func (ac *AccountClient) GetTag(pg *mastodon.Pagination, search string) ([]Item, error) {
	var items []Item
	statuses, err := ac.Client.GetTimelineHashtag(context.Background(), search, false, pg)
	if err != nil {
		return items, err
	}
	for _, s := range statuses {
		items = append(items, NewStatusItem(s))
	}
	return items, nil
}
