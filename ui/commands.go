package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"golang.org/x/exp/slices"
)

func (tv *TutView) ComposeCommand() {
	tv.InitPost(nil, nil)
}

func (tv *TutView) EditCommand() {
	item, itemErr := tv.GetCurrentItem()
	if itemErr != nil {
		return
	}
	if item.Type() != api.StatusType {
		return
	}
	s := item.Raw().(*mastodon.Status)
	s = util.StatusOrReblog(s)
	if tv.tut.Client.Me.ID != s.Account.ID {
		return
	}
	tv.InitPost(nil, s)
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
		NewLocalFeed(tv, true, true),
	)
}

func (tv *TutView) FederatedCommand() {
	tv.Timeline.AddFeed(
		NewFederatedFeed(tv, true, true),
	)
}

func (tv *TutView) SpecialCommand(boosts, replies bool) {
	tv.Timeline.AddFeed(
		NewHomeSpecialFeed(tv, boosts, replies),
	)
}

func (tv *TutView) DirectCommand() {
	tv.Timeline.AddFeed(
		NewConversationsFeed(tv),
	)
}

func (tv *TutView) HomeCommand() {
	tv.Timeline.AddFeed(
		NewHomeFeed(tv, true, true),
	)
}

func (tv *TutView) NotificationsCommand() {
	tv.Timeline.AddFeed(
		NewNotificationFeed(tv, true, true),
	)
}

func (tv *TutView) MentionsCommand() {
	tv.Timeline.AddFeed(
		NewNotificatioMentionsFeed(tv, true, true),
	)
}

func (tv *TutView) ListsCommand() {
	tv.Timeline.AddFeed(
		NewListsFeed(tv),
	)
}

func (tv *TutView) TagCommand(tag string) {
	tv.Timeline.AddFeed(
		NewTagFeed(tv, tag, true, true),
	)
}

func (tv *TutView) TagsCommand() {
	tv.Timeline.AddFeed(
		NewTagsFeed(tv),
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

func (tv *TutView) MoveWindowLeft() {
	tv.Timeline.MoveCurrentWindowLeft()
}

func (tv *TutView) MoveWindowRight() {
	tv.Timeline.MoveCurrentWindowRight()
}

func (tv *TutView) MoveWindowHome() {
	tv.Timeline.MoveCurrentWindowHome()
}

func (tv *TutView) MoveWindowEnd() {
	tv.Timeline.MoveCurrentWindowEnd()
}

func (tv *TutView) SwitchCommand(s string) {
	ft := config.InvalidFeed

	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	cmd := parts[0]
	var subaction string
	if strings.Contains(parts[0], " ") {
		p := strings.Split(cmd, " ")
		cmd = p[0]
		subaction = strings.Join(p[1:], " ")
	}
	showBoosts := true
	showReplies := true
	name := ""
	if len(parts) > 1 {
		tfStr := []string{"true", "false"}
		name = parts[1]
		if slices.Contains(tfStr, name) {
			name = ""
		}
		if len(parts) > 2 && slices.Contains(tfStr, parts[len(parts)-2]) &&
			slices.Contains(tfStr, parts[len(parts)-1]) {
			showBoosts = parts[len(parts)-2] == "true"
			showReplies = parts[len(parts)-1] == "true"
		} else if len(parts) > 1 && slices.Contains(tfStr, parts[len(parts)-1]) {
			showBoosts = parts[len(parts)-1] == "true"
		} else {
			fmt.Printf("switch is invalid . Check this for errors: switch %s\n", s)
			os.Exit(1)
		}
	}
	var data string
	switch cmd {
	case "home":
		ft = config.TimelineHome
	case "direct":
		ft = config.Conversations
	case "local":
		ft = config.TimelineLocal
	case "federated":
		ft = config.TimelineFederated
	case "special":
		ft = config.TimelineHomeSpecial
	case "special-all":
		ft = config.TimelineHomeSpecial
		showBoosts = true
		showReplies = true
	case "special-boosts":
		ft = config.TimelineHomeSpecial
		showBoosts = true
		showReplies = false
	case "special-replies":
		ft = config.TimelineHomeSpecial
		showBoosts = false
		showReplies = true
	case "bookmarks", "saved":
		ft = config.Saved
	case "favorited":
		ft = config.Favorited
	case "notifications":
		ft = config.Notifications
	case "lists":
		ft = config.Lists
	case "tag":
		ft = config.Tag
		data = subaction
	case "blocking":
		ft = config.Blocking
	case "muting":
		ft = config.Muting
	case "tags":
		ft = config.Tags
	case "mentions":
		ft = config.Mentions
	}
	found := tv.Timeline.FindAndGoTo(ft, data, showBoosts, showReplies)
	if found {
		return
	}
	nf := CreateFeed(tv, ft, data, showBoosts, showReplies)
	tv.Timeline.Feeds = append(tv.Timeline.Feeds, &FeedHolder{
		Feeds: []*Feed{nf},
		Name:  name,
	})
	tv.FocusFeed(len(tv.Timeline.Feeds) - 1)
	tv.Shared.Top.SetText(tv.Timeline.GetTitle())
	tv.Timeline.update <- true
}

func (tv *TutView) CloseWindowCommand() {
	tv.Timeline.CloseCurrentWindow()
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
			if f.Data.Type() == config.Notifications {
				f.Data.Clear()
			}
		}
	}
}

func (tv *TutView) ToggleStickToTop() {
	tv.tut.Config.General.StickToTop = !tv.tut.Config.General.StickToTop
}

func (tv *TutView) RefetchCommand() {
	item, itemErr := tv.GetCurrentItem()
	f := tv.GetCurrentFeed()
	if itemErr != nil {
		return
	}
	update := item.Refetch(tv.tut.Client)
	if update {
		f.DrawContent()
	}
}
