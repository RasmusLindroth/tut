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
		NewBlocking(tv, &config.Timeline{
			FeedType: config.Blocking,
		}),
	)
}

func (tv *TutView) BookmarksCommand() {
	tv.Timeline.AddFeed(
		NewBookmarksFeed(tv, &config.Timeline{
			FeedType: config.Saved,
		}))
}
func (tv *TutView) FavoritedCommand() {
	tv.Timeline.AddFeed(
		NewFavoritedFeed(tv, &config.Timeline{
			FeedType: config.Favorited,
		}),
	)
}

func (tv *TutView) MutingCommand() {
	tv.Timeline.AddFeed(
		NewMuting(tv, &config.Timeline{
			FeedType: config.Muting,
		}),
	)
}

func (tv *TutView) FollowRequestsCommand() {
	tv.Timeline.AddFeed(
		NewFollowRequests(tv, &config.Timeline{
			FeedType: config.FollowRequests,
		}),
	)
}

func (tv *TutView) LocalCommand() {
	tv.Timeline.AddFeed(
		NewLocalFeed(tv, &config.Timeline{
			FeedType: config.TimelineLocal,
		}),
	)
}

func (tv *TutView) FederatedCommand() {
	tv.Timeline.AddFeed(
		NewFederatedFeed(tv, &config.Timeline{
			FeedType: config.TimelineFederated,
		}),
	)
}

func (tv *TutView) SpecialCommand(hideBoosts, hideReplies bool) {
	tv.Timeline.AddFeed(
		NewHomeSpecialFeed(tv, &config.Timeline{
			FeedType:    config.TimelineHomeSpecial,
			HideBoosts:  hideBoosts,
			HideReplies: hideReplies,
		}),
	)
}

func (tv *TutView) DirectCommand() {
	tv.Timeline.AddFeed(
		NewConversationsFeed(tv, &config.Timeline{
			FeedType: config.Conversations,
		}),
	)
}

func (tv *TutView) HomeCommand() {
	tv.Timeline.AddFeed(
		NewHomeFeed(tv, &config.Timeline{
			FeedType: config.TimelineHome,
		}),
	)
}

func (tv *TutView) NotificationsCommand() {
	tv.Timeline.AddFeed(
		NewNotificationFeed(tv, &config.Timeline{
			FeedType: config.Notifications,
		}),
	)
}

func (tv *TutView) MentionsCommand() {
	tv.Timeline.AddFeed(
		NewNotificatioMentionsFeed(tv, &config.Timeline{
			FeedType: config.Mentions,
		}),
	)
}

func (tv *TutView) ListsCommand() {
	tv.Timeline.AddFeed(
		NewListsFeed(tv, &config.Timeline{
			FeedType: config.Lists,
		}),
	)
}

func (tv *TutView) TagCommand(tag string) {
	tv.Timeline.AddFeed(
		NewTagFeed(tv, &config.Timeline{
			FeedType:  config.Tag,
			Subaction: tag,
		}),
	)
}

func (tv *TutView) TagsCommand() {
	tv.Timeline.AddFeed(
		NewTagsFeed(tv, &config.Timeline{
			FeedType: config.Tags,
		}),
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
	hideBoosts := false
	hideReplies := false
	name := ""
	if len(parts) > 1 {
		tfStr := []string{"true", "false"}
		name = parts[1]
		if slices.Contains(tfStr, name) {
			name = ""
		}
		if len(parts) > 2 && slices.Contains(tfStr, parts[len(parts)-2]) &&
			slices.Contains(tfStr, parts[len(parts)-1]) {
			hideBoosts = parts[len(parts)-2] == "true"
			hideReplies = parts[len(parts)-1] == "true"
		} else if len(parts) > 1 && slices.Contains(tfStr, parts[len(parts)-1]) {
			hideReplies = parts[len(parts)-1] == "true"
		} else if name == "" {
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
		hideBoosts = false
		hideReplies = false
	case "special-boosts":
		ft = config.TimelineHomeSpecial
		hideBoosts = false
		hideReplies = true
	case "special-replies":
		ft = config.TimelineHomeSpecial
		hideBoosts = true
		hideReplies = false
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
	found := tv.Timeline.FindAndGoTo(ft, data, hideBoosts, hideReplies)
	if found {
		return
	}
	timelineData := config.Timeline{
		Name:        name,
		FeedType:    ft,
		Subaction:   data,
		HideBoosts:  hideBoosts,
		HideReplies: hideReplies,
	}
	nf := CreateFeed(tv, &timelineData)
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
		NewBoosts(tv, s.ID, &config.Timeline{
			FeedType:    config.Boosts,
			HideBoosts:  false,
			HideReplies: true,
		}),
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
		NewFavoritesStatus(tv, s.ID, &config.Timeline{
			FeedType: config.Favorites,
		}),
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
		NewFollowing(tv, s.Data.ID, &config.Timeline{
			FeedType: config.Following,
		}),
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
		NewFollowers(tv, s.Data.ID, &config.Timeline{
			FeedType: config.Followers,
		}),
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
		NewHistoryFeed(tv, item, &config.Timeline{
			FeedType: config.History,
		}),
	)
}

func (tv *TutView) ProfileCommand() {
	item, err := tv.tut.Client.GetUserByID(tv.tut.Client.Me.ID)
	if err != nil {
		tv.ShowError(fmt.Sprintf("Couldn't load user. Error: %v\n", err))
		return
	}
	tv.Timeline.AddFeed(
		NewUserFeed(tv, item, &config.Timeline{
			FeedType: config.User,
		}),
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
