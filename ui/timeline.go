package ui

import (
	"fmt"
	"os"

	"github.com/RasmusLindroth/tut/feed"
)

type Timeline struct {
	tutView       *TutView
	Feeds         []*Feed
	FeedIndex     int
	Notifications *Feed
	update        chan bool
}

func NewTimeline(tv *TutView, update chan bool) *Timeline {
	tl := &Timeline{
		tutView:       tv,
		Feeds:         []*Feed{},
		FeedIndex:     0,
		Notifications: nil,
		update:        update,
	}
	var nf *Feed
	switch tv.tut.Config.General.StartTimeline {
	case feed.TimelineFederated:
		nf = NewFederatedFeed(tv)
	case feed.TimelineLocal:
		nf = NewLocalFeed(tv)
	case feed.Conversations:
		nf = NewConversationsFeed(tv)
	default:
		nf = NewHomeFeed(tv)
	}
	tl.Feeds = append(tl.Feeds, nf)
	tl.Notifications = NewNotificationFeed(tv)
	tl.Notifications.ListOutFocus()

	return tl
}

func (tl *Timeline) AddFeed(f *Feed) {
	tl.tutView.FocusFeed()
	tl.Feeds = append(tl.Feeds, f)
	tl.FeedIndex = tl.FeedIndex + 1
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) RemoveCurrent(quit bool) {
	if len(tl.Feeds) == 1 && !quit {
		return
	}
	if len(tl.Feeds) == 1 && quit {
		tl.tutView.tut.App.Stop()
		os.Exit(0)
	}
	tl.Feeds[tl.FeedIndex].Data.Close()
	tl.Feeds = append(tl.Feeds[:tl.FeedIndex], tl.Feeds[tl.FeedIndex+1:]...)
	ni := tl.FeedIndex - 1
	if ni < 0 {
		ni = 0
	}
	tl.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) NextFeed() {
	l := len(tl.Feeds)
	ni := tl.FeedIndex + 1
	if ni >= l {
		ni = l - 1
	}
	tl.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) PrevFeed() {
	ni := tl.FeedIndex - 1
	if ni < 0 {
		ni = 0
	}
	tl.FeedIndex = ni
	tl.tutView.Shared.Top.SetText(tl.GetTitle())
	tl.update <- true
}

func (tl *Timeline) DrawContent(main bool) {
	var f *Feed
	if main {
		f = tl.Feeds[tl.FeedIndex]
	} else {
		f = tl.Notifications
	}
	f.DrawContent()
}

func (tl *Timeline) GetFeedList() *FeedList {
	return tl.Feeds[tl.FeedIndex].List
}

func (tl *Timeline) GetFeedContent(main bool) *FeedContent {
	if main {
		return tl.Feeds[tl.FeedIndex].Content
	} else {
		return tl.Notifications.Content
	}
}

func (tl *Timeline) GetTitle() string {
	index := tl.FeedIndex
	total := len(tl.Feeds)
	current := tl.Feeds[index].Data.Type()
	name := tl.Feeds[index].Data.Name()
	ct := ""
	switch current {
	case feed.Favorited:
		ct = "favorited"
	case feed.Notification:
		ct = "notifications"
	case feed.Tag:
		ct = fmt.Sprintf("tag #%s", name)
	case feed.Thread:
		ct = "thread feed"
	case feed.TimelineFederated:
		ct = "timeline federated"
	case feed.TimelineHome:
		ct = "timeline home"
	case feed.TimelineLocal:
		ct = "timeline local"
	case feed.Saved:
		ct = "saved/bookmarked toots"
	case feed.User:
		ct = "timeline user"
	case feed.UserList:
		ct = fmt.Sprintf("user search %s", name)
	case feed.Conversations:
		ct = "timeline direct"
	case feed.Lists:
		ct = "lists"
	case feed.List:
		ct = fmt.Sprintf("list named %s", name)
	case feed.Boosts:
		ct = "boosts"
	case feed.Favorites:
		ct = "favorites"
	case feed.Followers:
		ct = "followers"
	case feed.Following:
		ct = "following"
	case feed.Blocking:
		ct = "blocking"
	case feed.Muting:
		ct = "muting"
	}
	return fmt.Sprintf("%s (%d/%d)", ct, index+1, total)
}

func (tl *Timeline) ScrollUp() {
	f := tl.Feeds[tl.FeedIndex]
	row, _ := f.Content.Main.GetScrollOffset()
	if row > 0 {
		row = row - 1
	}
	f.Content.Main.ScrollTo(row, 0)
}

func (tl *Timeline) ScrollDown() {
	f := tl.Feeds[tl.FeedIndex]
	row, _ := f.Content.Main.GetScrollOffset()
	f.Content.Main.ScrollTo(row+1, 0)
}

func (tl *Timeline) NextItemFeed(mainFocus bool) {
	var f *Feed
	if mainFocus {
		f = tl.Feeds[tl.FeedIndex]
	} else {
		f = tl.Notifications
	}
	loadMore := f.List.Next()
	if loadMore {
		f.LoadOlder()
	}
	tl.DrawContent(mainFocus)
}
func (tl *Timeline) PrevItemFeed(mainFocus bool) {
	var f *Feed
	if mainFocus {
		f = tl.Feeds[tl.FeedIndex]
	} else {
		f = tl.Notifications
	}
	loadMore := f.List.Prev()
	if loadMore {
		f.LoadNewer()
	}
	tl.DrawContent(mainFocus)
}

func (tl *Timeline) HomeItemFeed(mainFocus bool) {
	var f *Feed
	if mainFocus {
		f = tl.Feeds[tl.FeedIndex]
	} else {
		f = tl.Notifications
	}
	f.List.SetCurrentItem(0)
	f.LoadNewer()
	tl.DrawContent(mainFocus)
}

func (tl *Timeline) EndItemFeed(mainFocus bool) {
	var f *Feed
	if mainFocus {
		f = tl.Feeds[tl.FeedIndex]
	} else {
		f = tl.Notifications
	}
	ni := f.List.GetItemCount() - 1
	if ni < 0 {
		return
	}
	f.List.SetCurrentItem(ni)
	f.LoadOlder()
	tl.DrawContent(mainFocus)
}
