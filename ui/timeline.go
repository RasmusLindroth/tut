package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/feed"
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
}

func NewTimeline(tv *TutView, update chan bool) *Timeline {
	tl := &Timeline{
		tutView: tv,
		Feeds:   []*FeedHolder{},
		update:  update,
	}
	var nf *Feed
	for _, f := range tv.tut.Config.General.Timelines {
		switch f.FeedType {
		case feed.TimelineHome:
			nf = NewHomeFeed(tv)
		case feed.Conversations:
			nf = NewConversationsFeed(tv)
		case feed.TimelineLocal:
			nf = NewLocalFeed(tv)
		case feed.TimelineFederated:
			nf = NewFederatedFeed(tv)
		case feed.Saved:
			nf = NewBookmarksFeed(tv)
		case feed.Favorited:
			nf = NewFavoritedFeed(tv)
		case feed.Notification:
			nf = NewNotificationFeed(tv)
		case feed.Lists:
			nf = NewListsFeed(tv)
		case feed.Tag:
			nf = NewTagFeed(tv, f.Subaction)
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
	case feed.Favorited:
		ct = "favorited"
	case feed.Notification:
		ct = "notifications"
	case feed.Tag:
		ct = fmt.Sprintf("tag #%s", name)
	case feed.Thread:
		ct = "thread feed"
	case feed.TimelineFederated:
		ct = "federated"
	case feed.TimelineHome:
		ct = "home"
	case feed.TimelineLocal:
		ct = "local"
	case feed.Saved:
		ct = "saved/bookmarked toots"
	case feed.User:
		ct = fmt.Sprintf("user %s", name)
	case feed.UserList:
		ct = fmt.Sprintf("user search %s", name)
	case feed.Conversations:
		ct = "direct"
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
	case feed.FollowRequests:
		ct = "follow requests"
	case feed.Blocking:
		ct = "blocking"
	case feed.Muting:
		ct = "muting"
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
		f.LoadNewer()
	}
	tl.DrawContent()
}

func (tl *Timeline) HomeItemFeed() {
	fh := tl.Feeds[tl.FeedFocusIndex]
	f := fh.Feeds[fh.FeedIndex]
	f.List.SetCurrentItem(0)
	f.LoadNewer()
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
