package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
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
