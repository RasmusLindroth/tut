package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/feed"
	"github.com/RasmusLindroth/tut/util"
)

func (tv *TutView) ComposeCommand() {
	tv.InitPost(nil)
}

func (tv *TutView) BlockingCommand() {
	tv.Timeline.AddFeed(
		NewBlocking(tv),
	)
}

func (tv *TutView) BookmarksCommand() {
	tv.Timeline.AddFeed(
		NewBookmarksFeed(tv),
	)
}
func (tv *TutView) FavoritedCommand() {
	tv.Timeline.AddFeed(
		NewFavoritedFeed(tv),
	)
}

func (tv *TutView) MutingCommand() {
	tv.Timeline.AddFeed(
		NewMuting(tv),
	)
}

func (tv *TutView) FollowRequestsCommand() {
	tv.Timeline.AddFeed(
		NewFollowRequests(tv),
	)
}

func (tv *TutView) LocalCommand() {
	tv.Timeline.AddFeed(
		NewLocalFeed(tv),
	)
}

func (tv *TutView) FederatedCommand() {
	tv.Timeline.AddFeed(
		NewFederatedFeed(tv),
	)
}

func (tv *TutView) DirectCommand() {
	tv.Timeline.AddFeed(
		NewConversationsFeed(tv),
	)
}

func (tv *TutView) HomeCommand() {
	tv.Timeline.AddFeed(
		NewHomeFeed(tv),
	)
}

func (tv *TutView) NotificationsCommand() {
	tv.Timeline.AddFeed(
		NewNotificationFeed(tv),
	)
}

func (tv *TutView) ListsCommand() {
	tv.Timeline.AddFeed(
		NewListsFeed(tv),
	)
}

func (tv *TutView) TagCommand(tag string) {
	tv.Timeline.AddFeed(
		NewTagFeed(tv, tag),
	)
}

func (tv *TutView) TagFollowCommand(tag string) {
	err := tv.tut.Client.FollowTag(tag)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't follow tag. Error: %v\n", err))
		return
	}
}

func (tv *TutView) TagUnfollowCommand(tag string) {
	err := tv.tut.Client.UnfollowTag(tag)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't unfollow tag. Error: %v\n", err))
		return
	}
}

func (tv *TutView) WindowCommand(index string) {
	i, err := strconv.Atoi(index)
	if err != nil {
		tv.ShowError(
			fmt.Sprintf("couldn't convert str to int. Error %v", err),
		)
		return
	}
	tv.FocusFeed(i)
}

func (tv *TutView) BoostsCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.StatusType {
		return
	}
	s := item.Raw().(*mastodon.Status)
	s = util.StatusOrReblog(s)
	tv.Timeline.AddFeed(
		NewBoosts(tv, s.ID),
	)
}

func (tv *TutView) FavoritesCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.StatusType {
		return
	}
	s := item.Raw().(*mastodon.Status)
	s = util.StatusOrReblog(s)
	tv.Timeline.AddFeed(
		NewFavoritesStatus(tv, s.ID),
	)
}

func (tv *TutView) FollowingCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.UserType && item.Type() != api.ProfileType {
		return
	}
	s := item.Raw().(*api.User)
	tv.Timeline.AddFeed(
		NewFollowing(tv, s.Data.ID),
	)
}

func (tv *TutView) FollowersCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.UserType && item.Type() != api.ProfileType {
		return
	}
	s := item.Raw().(*api.User)
	tv.Timeline.AddFeed(
		NewFollowers(tv, s.Data.ID),
	)
}

func (tv *TutView) HistoryCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.StatusType {
		return
	}
	tv.Timeline.AddFeed(
		NewHistoryFeed(tv, item),
	)
}

func (tv *TutView) ProfileCommand() {
	item, err := tv.tut.Client.GetUserByID(tv.tut.Client.Me.ID)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't load user. Error: %v\n", err))
		return
	}
	tv.Timeline.AddFeed(
		NewUserFeed(tv, item),
	)
}

func (tv *TutView) PreferencesCommand() {
	tv.SetPage(PreferenceFocus)
}

func (tv *TutView) ListPlacementCommand(lp config.ListPlacement) {
	tv.tut.Config.General.ListPlacement = lp
	tv.MainView.ForceUpdate()
}

func (tv *TutView) ListSplitCommand(ls config.ListSplit) {
	tv.tut.Config.General.ListSplit = ls
	tv.MainView.ForceUpdate()
}

func (tv *TutView) ProportionsCommand(lp string, cp string) {
	lpi, err := strconv.Atoi(lp)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't parse list proportion. Error: %v\n", err))
		return
	}
	cpi, err := strconv.Atoi(cp)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't parse content proportion. Error: %v\n", err))
		return
	}
	tv.tut.Config.General.ListProportion = lpi
	tv.tut.Config.General.ContentProportion = cpi
	tv.MainView.ForceUpdate()
}

func (tv *TutView) LoadNewerCommand() {
	f := tv.GetCurrentFeed()
	f.LoadNewer(true)
}

func (tv *TutView) ClearNotificationsCommand() {
	err := tv.tut.Client.ClearNotifications()
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't clear notifications. Error: %v\n", err))
		return
	}
	for _, tl := range tv.Timeline.Feeds {
		for _, f := range tl.Feeds {
			if f.Data.Type() == feed.Notification {
				f.Data.Clear()
			}
		}
	}
}
