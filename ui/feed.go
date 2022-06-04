package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/feed"
	"github.com/gdamore/tcell/v2"
	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
)

type FeedList struct {
	Text        *tview.List
	Symbol      *tview.List
	stickyCount int
}

func (fl *FeedList) InFocus(style config.Style) {
	inFocus(fl.Text, style)
	inFocus(fl.Symbol, style)
}

func inFocus(l *tview.List, style config.Style) {
	l.SetBackgroundColor(style.Background)
	l.SetMainTextColor(style.Text)
	l.SetSelectedBackgroundColor(style.ListSelectedBackground)
	l.SetSelectedTextColor(style.ListSelectedText)
}

func (fl *FeedList) OutFocus(style config.Style) {
	outFocus(fl.Text, style)
	outFocus(fl.Symbol, style)
}

func outFocus(l *tview.List, style config.Style) {
	l.SetBackgroundColor(style.Background)
	l.SetMainTextColor(style.Text)
	l.SetSelectedBackgroundColor(style.StatusBarViewBackground)
	l.SetSelectedTextColor(style.StatusBarViewText)
}

type Feed struct {
	tutView   *TutView
	Data      *feed.Feed
	ListIndex int
	List      *FeedList
	Content   *FeedContent
}

func (f *Feed) ListInFocus() {
	f.List.InFocus(f.tutView.tut.Config.Style)
}

func (f *Feed) ListOutFocus() {
	f.List.OutFocus(f.tutView.tut.Config.Style)
}

func (f *Feed) LoadOlder() {
	f.Data.LoadOlder()
}

func (f *Feed) LoadNewer() {
	if f.Data.HasStream() {
		return
	}
	f.Data.LoadNewer()
}

func (f *Feed) Delete() {
	id := f.List.GetCurrentID()
	f.Data.Delete(id)
}

func (f *Feed) DrawContent() {
	id := f.List.GetCurrentID()
	for _, item := range f.Data.List() {
		if id != item.ID() {
			continue
		}
		DrawItem(f.tutView.tut, item, f.Content.Main, f.Content.Controls, f.Data.Type())
		f.tutView.ShouldSync()
	}
}

func (f *Feed) update() {
	for nft := range f.Data.Update {
		switch nft {
		case feed.DesktopNotificationFollower:
			if f.tutView.tut.Config.NotificationConfig.NotificationFollower {
				beeep.Notify("New follower", "", "")
			}
		case feed.DesktopNotificationFavorite:
			if f.tutView.tut.Config.NotificationConfig.NotificationFavorite {
				beeep.Notify("Favorited your toot", "", "")
			}
		case feed.DesktopNotificationMention:
			if f.tutView.tut.Config.NotificationConfig.NotificationMention {
				beeep.Notify("Mentioned you", "", "")
			}
		case feed.DesktopNotificationBoost:
			if f.tutView.tut.Config.NotificationConfig.NotificationBoost {
				beeep.Notify("Boosted your toot", "", "")
			}
		case feed.DesktopNotificationPoll:
			if f.tutView.tut.Config.NotificationConfig.NotificationPoll {
				beeep.Notify("Poll has ended", "", "")
			}
		case feed.DesktopNotificationPost:
			if f.tutView.tut.Config.NotificationConfig.NotificationPost {
				beeep.Notify("New post", "", "")
			}
		}
		f.tutView.tut.App.QueueUpdateDraw(func() {
			lLen := f.List.GetItemCount()
			curr := f.List.GetCurrentID()
			f.List.Clear()
			for _, item := range f.Data.List() {
				main, symbol := DrawListItem(f.tutView.tut.Config, item)
				f.List.AddItem(main, symbol, item.ID())
			}
			f.List.SetByID(curr)
			if lLen == 0 {
				f.DrawContent()
			}
		})
	}
}

func NewHomeFeed(tv *TutView) *Feed {
	f := feed.NewTimelineHome(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFederatedFeed(tv *TutView) *Feed {
	f := feed.NewTimelineFederated(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewLocalFeed(tv *TutView) *Feed {
	f := feed.NewTimelineLocal(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewNotificationFeed(tv *TutView) *Feed {
	f := feed.NewNotifications(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewThreadFeed(tv *TutView, item api.Item) *Feed {
	f := feed.NewThread(tv.tut.Client, item.Raw().(*mastodon.Status))
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	for i, s := range f.List() {
		main, symbol := DrawListItem(tv.tut.Config, s)
		fd.List.AddItem(main, symbol, s.ID())
		if s.Raw().(*mastodon.Status).ID == item.Raw().(*mastodon.Status).ID {
			fd.List.SetCurrentItem(i)
		}
	}
	fd.DrawContent()

	return fd
}

func NewConversationsFeed(tv *TutView) *Feed {
	f := feed.NewConversations(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewUserFeed(tv *TutView, item api.Item) *Feed {
	if item.Type() != api.UserType && item.Type() != api.ProfileType {
		panic("Can't open user. Wrong type.\n")
	}
	u := item.Raw().(*api.User)
	f := feed.NewUserProfile(tv.tut.Client, u)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewUserSearchFeed(tv *TutView, search string) *Feed {
	f := feed.NewUserSearch(tv.tut.Client, search)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	for _, s := range f.List() {
		main, symbol := DrawListItem(tv.tut.Config, s)
		fd.List.AddItem(main, symbol, s.ID())
	}
	fd.DrawContent()

	return fd
}

func NewTagFeed(tv *TutView, search string) *Feed {
	f := feed.NewTag(tv.tut.Client, search)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}
func NewListsFeed(tv *TutView) *Feed {
	f := feed.NewListList(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewListFeed(tv *TutView, l *mastodon.List) *Feed {
	f := feed.NewList(tv.tut.Client, l)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFavoritedFeed(tv *TutView) *Feed {
	f := feed.NewFavorites(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}

	go fd.update()
	return fd
}

func NewBookmarksFeed(tv *TutView) *Feed {
	f := feed.NewBookmarks(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFavoritesStatus(tv *TutView, id mastodon.ID) *Feed {
	f := feed.NewFavoritesStatus(tv.tut.Client, id)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewBoosts(tv *TutView, id mastodon.ID) *Feed {
	f := feed.NewBoosts(tv.tut.Client, id)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFollowers(tv *TutView, id mastodon.ID) *Feed {
	f := feed.NewFollowers(tv.tut.Client, id)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFollowing(tv *TutView, id mastodon.ID) *Feed {
	f := feed.NewFollowing(tv.tut.Client, id)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewBlocking(tv *TutView) *Feed {
	f := feed.NewBlocking(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewMuting(tv *TutView) *Feed {
	f := feed.NewMuting(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFollowRequests(tv *TutView) *Feed {
	f := feed.NewFollowRequests(tv.tut.Client)
	f.LoadNewer()
	fd := &Feed{
		tutView:   tv,
		Data:      f,
		ListIndex: 0,
		List:      NewFeedList(tv.tut, f.StickyCount()),
		Content:   NewFeedContent(tv.tut),
	}
	go fd.update()

	return fd
}

func NewFeedList(t *Tut, stickyCount int) *FeedList {
	fl := &FeedList{
		Text:        NewList(t.Config),
		Symbol:      NewList(t.Config),
		stickyCount: stickyCount,
	}
	return fl
}

func (fl *FeedList) AddItem(text string, symbols string, id uint) {
	fl.Text.AddItem(text, fmt.Sprintf("%d", id), 0, nil)
	fl.Symbol.AddItem(symbols, fmt.Sprintf("%d", id), 0, nil)
}

func (fl *FeedList) Next() (loadOlder bool) {
	ni := fl.Text.GetCurrentItem() + 1
	if ni >= fl.Text.GetItemCount() {
		ni = fl.Text.GetItemCount() - 1
		if ni < 0 {
			ni = 0
		}
	}
	fl.Text.SetCurrentItem(ni)
	fl.Symbol.SetCurrentItem(ni)
	return fl.Text.GetItemCount()-(ni+1) < 5
}

func (fl *FeedList) Prev() (loadNewer bool) {
	ni := fl.Text.GetCurrentItem() - 1
	if ni < 0 {
		ni = 0
	}
	fl.Text.SetCurrentItem(ni)
	fl.Symbol.SetCurrentItem(ni)
	return ni-fl.stickyCount < 4
}

func (fl *FeedList) Clear() {
	fl.Text.Clear()
	fl.Symbol.Clear()
}

func (fl *FeedList) GetItemCount() int {
	return fl.Text.GetItemCount()
}

func (fl *FeedList) SetCurrentItem(index int) {
	fl.Text.SetCurrentItem(index)
	fl.Symbol.SetCurrentItem(index)
}

func (fl *FeedList) GetCurrentID() uint {
	if fl.GetItemCount() == 0 {
		return 0
	}
	i := fl.Text.GetCurrentItem()
	_, sec := fl.Text.GetItemText(i)
	id, err := strconv.ParseUint(sec, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

func (fl *FeedList) SetByID(id uint) {
	if fl.Text.GetItemCount() == 0 {
		return
	}
	s := fmt.Sprintf("%d", id)
	items := fl.Text.FindItems("", s, false, false)
	for _, i := range items {
		_, sec := fl.Text.GetItemText(i)
		if sec == s {
			fl.Text.SetCurrentItem(i)
			fl.Symbol.SetCurrentItem(i)
			break
		}
	}
}

type FeedContent struct {
	Main     *tview.TextView
	Controls *tview.TextView
}

func NewFeedContent(t *Tut) *FeedContent {
	m := NewTextView(t.Config)
	m.SetWordWrap(true)

	if t.Config.General.MaxWidth > 0 {
		mw := t.Config.General.MaxWidth
		m.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			rWidth := width
			if rWidth > mw {
				rWidth = mw
			}
			return x, y, rWidth, height
		})
	}
	c := NewTextView(t.Config)
	c.SetDynamicColors(true)
	fc := &FeedContent{
		Main:     m,
		Controls: c,
	}
	return fc
}
