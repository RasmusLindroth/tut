package ui

import (
	"fmt"
	"strings"

	"github.com/RasmusLindroth/tut/config"
)

type FeedHolder struct {
	Name      string
	Feeds     []*Feed
	FeedIndex int
}

type Timeline struct {
	tutView        *TutView
	Feeds          []*FeedHolder
	FeedFocusIndex int
	update         chan bool
	scrollSleep    *scrollSleep
}

func NewTimeline(tv *TutView, update chan bool) *Timeline {
	tl := &Timeline{
		tutView: tv,
		Feeds:   []*FeedHolder{},
		update:  update,
	}
	tl.scrollSleep = NewScrollSleep(tl.NextItemFeed, tl.PrevItemFeed)
	var nf *Feed
	for _, f := range tv.tut.Config.General.Timelines {
		switch f.FeedType {
		case config.TimelineHome:
			nf = NewHomeFeed(tv, f.ShowBoosts, f.ShowReplies)
		case config.TimelineHomeSpecial:
			nf = NewHomeSpecialFeed(tv, f.ShowBoosts, f.ShowReplies)
		case config.Conversations:
			nf = NewConversationsFeed(tv)
		case config.TimelineLocal:
			nf = NewLocalFeed(tv, f.ShowBoosts, f.ShowReplies)
		case config.TimelineFederated:
			nf = NewFederatedFeed(tv, f.ShowBoosts, f.ShowReplies)
		case config.Saved:
			nf = NewBookmarksFeed(tv)
		case config.Favorited:
			nf = NewFavoritedFeed(tv)
		case config.Notifications:
			nf = NewNotificationFeed(tv, f.ShowBoosts, f.ShowReplies)
		case config.Lists:
			nf = NewListsFeed(tv)
		case config.Tag:
			nf = NewTagFeed(tv, f.Subaction, f.ShowBoosts, f.ShowReplies)
		default:
			fmt.Println("Invalid feed")
			tl.tutView.CleanExit(1)
		}
		tl.Feeds = append(tl.Feeds, &FeedHolder{
			Feeds: []*Feed{nf},
			Name:  f.Name,
		})
	}
	for i := 1; i < len(tl.Feeds); i++ {
		for _, f := range tl.Feeds[i].Feeds {
			f.ListOutFocus()
		}
	}

	return tl
}

func (tl *Timeline) AddFeed(f *Feed) {
	fh := tl.Feeds[tl.FeedFocusIndex]
	fh.Feeds = append(fh.Feeds, f)
	fh.FeedIndex = fh.FeedIndex + 1
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) RemoveCurrent(quit bool) bool {
	if len(tl.Feeds[tl.FeedFocusIndex].Feeds) == 1 && !quit {
		return true
	}
	if len(tl.Feeds[tl.FeedFocusIndex].Feeds) == 1 && quit {
		tl.tutView.tut.App.Stop()
		tl.tutView.CleanExit(0)
	}

	f := tl.Feeds[tl.FeedFocusIndex]
	f.Feeds[f.FeedIndex].Data.Close()
	f.Feeds = append(f.Feeds[:f.FeedIndex], f.Feeds[f.FeedIndex+1:]...)
	ni := f.FeedIndex - 1
	if ni < 0 {
		ni = 0
	}
	f.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
	return false
}

func (tl *Timeline) NextFeed() {
	f := tl.Feeds[tl.FeedFocusIndex]
	l := len(f.Feeds)
	ni := f.FeedIndex + 1
	if ni >= l {
		ni = l - 1
	}
	f.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) PrevFeed() {
	f := tl.Feeds[tl.FeedFocusIndex]
	ni := f.FeedIndex - 1
	if ni < 0 {
		ni = 0
	}
	f.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) DrawContent() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	f.DrawContent()
}

func (fh *FeedHolder) GetFeedList() *FeedList {
	return fh.Feeds[fh.FeedIndex].List
}

func (tl *Timeline) GetFeedContent() *FeedContent {
	fh := tl.Feeds[tl.FeedFocusIndex]
	return fh.Feeds[fh.FeedIndex].Content
}

func (tl *Timeline) GetTitle() string {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	index := fh.FeedIndex
	total := len(fh.Feeds)
	current := f.Data.Type()
	name := f.Data.Name()
	ct := ""
	switch current {
	case config.Favorited:
		ct = "favorited"
	case config.Notifications:
		ct = "notifications"
	case config.Tag:
		parts := strings.Split(name, " ")
		for i, p := range parts {
			parts[i] = fmt.Sprintf("#%s", p)
		}
		ct = fmt.Sprintf("tag %s", strings.Join(parts, " "))
	case config.Thread:
		ct = "thread feed"
	case config.History:
		ct = "history feed"
	case config.TimelineFederated:
		ct = "federated"
	case config.TimelineHome:
		ct = "home"
	case config.TimelineHomeSpecial:
		ct = "special"
	case config.TimelineLocal:
		ct = "local"
	case config.Saved:
		ct = "saved/bookmarked toots"
	case config.User:
		ct = fmt.Sprintf("user %s", name)
	case config.UserList:
		ct = fmt.Sprintf("user search %s", name)
	case config.Conversations:
		ct = "direct"
	case config.Lists:
		ct = "lists"
	case config.List:
		ct = fmt.Sprintf("list named %s", name)
	case config.Boosts:
		ct = "boosts"
	case config.Favorites:
		ct = "favorites"
	case config.Followers:
		ct = "followers"
	case config.Following:
		ct = "following"
	case config.FollowRequests:
		ct = "follow requests"
	case config.Blocking:
		ct = "blocking"
	case config.ListUsersAdd:
		ct = fmt.Sprintf("Add users to %s", name)
	case config.ListUsersIn:
		ct = fmt.Sprintf("Delete users from %s", name)
	}
	return fmt.Sprintf("%s (%d/%d)", ct, index+1, total)
}

func (tl *Timeline) ScrollUp() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	row, _ := f.Content.Main.GetScrollOffset()
	if row > 0 {
		row = row - 1
	}
	f.Content.Main.ScrollTo(row, 0)
}

func (tl *Timeline) ScrollDown() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	row, _ := f.Content.Main.GetScrollOffset()
	f.Content.Main.ScrollTo(row+1, 0)
}

func (tl *Timeline) NextItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	loadMore := f.List.Next()
	if loadMore {
		f.LoadOlder()
	}
	tl.DrawContent()
}
func (tl *Timeline) PrevItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	loadMore := f.List.Prev()
	if loadMore {
		f.LoadNewer(false)
	}
	tl.DrawContent()
}

func (tl *Timeline) SetItemFeedIndex(index int) {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	loadOlder, loadNewer := f.List.Set(index)
	if loadOlder {
		f.LoadOlder()
	}
	if loadNewer {
		f.LoadNewer(false)
	}
	tl.DrawContent()
}

func (tl *Timeline) HomeItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	f.List.SetCurrentItem(0)
	f.LoadNewer(false)
	tl.DrawContent()
}

func (tl *Timeline) DeleteItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	f.List.GetCurrentID()

	tl.DrawContent()
}

func (tl *Timeline) EndItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	ni := f.List.GetItemCount() - 1
	if ni < 0 {
		return
	}
	f.List.SetCurrentItem(ni)
	f.LoadOlder()
	tl.DrawContent()
}
