package api

import (
	"strings"
	"sync"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"golang.org/x/exp/slices"
)

var id uint = 0
var idMux sync.Mutex

func newID() uint {
	idMux.Lock()
	defer idMux.Unlock()
	id = id + 1
	return id
}

type Item interface {
	ID() uint
	Type() MastodonType
	ToggleCW()
	ShowCW() bool
	Raw() interface{}
	URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int)
	Filtered(config.FeedType) (bool, string, string, bool)
	ForceViewFilter()
	Pinned() bool
	Refetch(*AccountClient) bool
}

type filtered struct {
	InUse   bool
	Filters []filter
}

type filter struct {
	Values []string
	Where  []string
	Type   string
}

func getUrlsStatus(status *mastodon.Status) ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	if status.Reblog != nil {
		status = status.Reblog
	}
	_, urls := util.CleanHTML(status.Content)
	if status.Sensitive {
		_, u := util.CleanHTML(status.SpoilerText)
		urls = append(urls, u...)
	}

	realUrls := []util.URL{}
	for _, url := range urls {
		isNotMention := true
		for _, mention := range status.Mentions {
			if mention.URL == url.URL {
				isNotMention = false
			}
		}
		if isNotMention {
			realUrls = append(realUrls, url)
		}
	}

	length := len(realUrls) + len(status.Mentions) + len(status.Tags)
	return realUrls, status.Mentions, status.Tags, length
}
func getUrlsUser(user *mastodon.Account) ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	var urls []util.URL
	user.Note, urls = util.CleanHTML(user.Note)
	for _, f := range user.Fields {
		_, fu := util.CleanHTML(f.Value)
		urls = append(urls, fu...)
	}

	return urls, []mastodon.Mention{}, []mastodon.Tag{}, len(urls)
}

func NewStatusItem(item *mastodon.Status, pinned bool) (sitem Item) {
	filtered := filtered{InUse: false}
	if item == nil {
		return &StatusItem{id: newID(), item: item, showSpoiler: false, filtered: filtered, pinned: pinned}
	}
	s := util.StatusOrReblog(item)
	for _, f := range s.Filtered {
		filtered.InUse = true
		filtered.Filters = append(filtered.Filters,
			filter{
				Type:   f.Filter.FilterAction,
				Values: f.KeywordMatches,
				Where:  f.Filter.Context,
			})
	}
	sitem = &StatusItem{id: newID(), item: item, showSpoiler: false, filtered: filtered, pinned: pinned}
	return sitem
}

func NewStatusItemID(item *mastodon.Status, pinned bool, id uint) (sitem Item) {
	sitem = NewStatusItem(item, pinned)
	sitem.(*StatusItem).id = id
	return sitem
}

type StatusItem struct {
	id          uint
	item        *mastodon.Status
	showSpoiler bool
	forceView   bool
	filtered    filtered
	pinned      bool
}

func (s *StatusItem) ID() uint {
	return s.id
}

func (s *StatusItem) Type() MastodonType {
	return StatusType
}

func (s *StatusItem) ToggleCW() {
	s.showSpoiler = !s.showSpoiler
}

func (s *StatusItem) ShowCW() bool {
	return s.showSpoiler
}

func (s *StatusItem) Raw() interface{} {
	return s.item
}

func (s *StatusItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return getUrlsStatus(s.item)
}

func (s *StatusItem) Filtered(tl config.FeedType) (bool, string, string, bool) {
	if !s.filtered.InUse || s.forceView {
		return false, "", "", true
	}
	words := []string{}
	t := ""
	for _, f := range s.filtered.Filters {
		used := false
		for _, w := range f.Where {
			switch w {
			case "home":
				if tl == config.TimelineHome || tl == config.List || tl == config.TimelineHomeSpecial {
					used = true
				}
			case "thread":
				if tl == config.Thread || tl == config.Conversations {
					used = true
				}
			case "notifications":
				if tl == config.Notifications || tl == config.Mentions {
					used = true
				}
			case "account":
				if tl == config.User {
					used = true
				}
			case "public":
				where := []config.FeedType{
					config.TimelineFederated,
					config.TimelineLocal,
					config.Tag,
				}
				if slices.Contains(where, tl) {
					used = true
				}
			}
			if used {
				words = append(words, f.Values...)
				if t == "" || t == "warn" {
					t = f.Type
				}
				break
			}
		}
	}
	return len(words) > 0, t, strings.Join(words, ", "), s.forceView
}

func (s *StatusItem) ForceViewFilter() {
	s.forceView = true
}

func (s *StatusItem) Pinned() bool {
	return s.pinned
}

func (s *StatusItem) Refetch(ac *AccountClient) bool {
	ns, err := ac.GetStatus(s.item.ID)
	if err != nil {
		return false
	}
	nsi := NewStatusItemID(ns, s.pinned, s.id)
	*s = *nsi.(*StatusItem)
	return true
}

func NewStatusHistoryItem(item *mastodon.StatusHistory) (sitem Item) {
	return &StatusHistoryItem{id: newID(), item: item, showSpoiler: false}
}

type StatusHistoryItem struct {
	id          uint
	item        *mastodon.StatusHistory
	showSpoiler bool
}

func (s *StatusHistoryItem) ID() uint {
	return s.id
}

func (s *StatusHistoryItem) Type() MastodonType {
	return StatusHistoryType
}

func (s *StatusHistoryItem) ToggleCW() {
	s.showSpoiler = !s.showSpoiler
}

func (s *StatusHistoryItem) ShowCW() bool {
	return s.showSpoiler
}

func (s *StatusHistoryItem) Raw() interface{} {
	return s.item
}

func (s *StatusHistoryItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	status := mastodon.Status{
		Content:          s.item.Content,
		SpoilerText:      s.item.SpoilerText,
		Account:          s.item.Account,
		Sensitive:        s.item.Sensitive,
		CreatedAt:        s.item.CreatedAt,
		Emojis:           s.item.Emojis,
		MediaAttachments: s.item.MediaAttachments,
	}
	return getUrlsStatus(&status)
}

func (s *StatusHistoryItem) Filtered(config.FeedType) (bool, string, string, bool) {
	return false, "", "", true
}

func (t *StatusHistoryItem) ForceViewFilter() {}

func (s *StatusHistoryItem) Pinned() bool {
	return false
}

func (s *StatusHistoryItem) Refetch(ac *AccountClient) bool {
	return false
}

func NewUserItem(item *User, profile bool) Item {
	return &UserItem{id: newID(), item: item, profile: profile}
}

type UserItem struct {
	id      uint
	item    *User
	profile bool
}

func (u *UserItem) ID() uint {
	return u.id
}

func (u *UserItem) Type() MastodonType {
	if u.profile {
		return ProfileType
	}
	return UserType
}

func (u *UserItem) ToggleCW() {
}

func (u *UserItem) ShowCW() bool {
	return false
}

func (u *UserItem) Raw() interface{} {
	return u.item
}

func (u *UserItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return getUrlsUser(u.item.Data)
}

func (u *UserItem) Filtered(config.FeedType) (bool, string, string, bool) {
	return false, "", "", true
}

func (u *UserItem) ForceViewFilter() {}

func (u *UserItem) Pinned() bool {
	return false
}

func (u *UserItem) Refetch(ac *AccountClient) bool {
	return false
}

func NewNotificationItem(item *mastodon.Notification, user *User) (nitem Item) {
	status := NewStatusItem(item.Status, false)
	nitem = &NotificationItem{
		id:          newID(),
		item:        item,
		showSpoiler: false,
		user:        NewUserItem(user, false),
		status:      status,
	}

	return nitem
}

type NotificationItem struct {
	id          uint
	item        *mastodon.Notification
	showSpoiler bool
	status      Item
	user        Item
}

type NotificationData struct {
	Item   *mastodon.Notification
	Status Item
	User   Item
}

func (n *NotificationItem) ID() uint {
	return n.id
}

func (n *NotificationItem) Type() MastodonType {
	return NotificationType
}

func (n *NotificationItem) ToggleCW() {
	n.showSpoiler = !n.showSpoiler
}

func (n *NotificationItem) ShowCW() bool {
	return n.showSpoiler
}

func (n *NotificationItem) Raw() interface{} {
	return &NotificationData{
		Item:   n.item,
		Status: n.status,
		User:   n.user,
	}
}

func (n *NotificationItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	nd := n.Raw().(*NotificationData)
	switch n.item.Type {
	case "favourite":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "reblog":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "mention":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "status":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "poll":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "update":
		return getUrlsStatus(nd.Status.Raw().(*mastodon.Status))
	case "follow":
		return getUrlsUser(nd.User.Raw().(*User).Data)
	case "follow_request":
		return getUrlsUser(nd.User.Raw().(*User).Data)
	default:
		return []util.URL{}, []mastodon.Mention{}, []mastodon.Tag{}, 0
	}
}

func (n *NotificationItem) Filtered(config.FeedType) (bool, string, string, bool) {
	return false, "", "", true
}

func (n *NotificationItem) ForceViewFilter() {}

func (n *NotificationItem) Pinned() bool {
	return false
}

func (n *NotificationItem) Refetch(ac *AccountClient) bool {
	return false
}

func NewListsItem(item *mastodon.List) Item {
	return &ListItem{id: newID(), item: item, showSpoiler: true}
}

type ListItem struct {
	id          uint
	item        *mastodon.List
	showSpoiler bool
}

func (s *ListItem) ID() uint {
	return s.id
}

func (s *ListItem) Type() MastodonType {
	return ListsType
}

func (s *ListItem) ToggleCW() {
}

func (s *ListItem) ShowCW() bool {
	return true
}

func (s *ListItem) Raw() interface{} {
	return s.item
}

func (s *ListItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return nil, nil, nil, 0
}

func (s *ListItem) Filtered(config.FeedType) (bool, string, string, bool) {
	return false, "", "", true
}

func (l *ListItem) ForceViewFilter() {}

func (n *ListItem) Pinned() bool {
	return false
}

func (l *ListItem) Refetch(ac *AccountClient) bool {
	return false
}

func NewTagItem(item *mastodon.Tag) Item {
	return &TagItem{id: newID(), item: item, showSpoiler: true}
}

type TagItem struct {
	id          uint
	item        *mastodon.Tag
	showSpoiler bool
}

func (t *TagItem) ID() uint {
	return t.id
}

func (t *TagItem) Type() MastodonType {
	return TagType
}

func (t *TagItem) ToggleCW() {
}

func (t *TagItem) ShowCW() bool {
	return true
}

func (t *TagItem) Raw() interface{} {
	return t.item
}

func (t *TagItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return nil, nil, nil, 0
}

func (t *TagItem) Filtered(config.FeedType) (bool, string, string, bool) {
	return false, "", "", true
}

func (t *TagItem) ForceViewFilter() {}

func (t *TagItem) Pinned() bool {
	return false
}

func (t *TagItem) Refetch(ac *AccountClient) bool {
	return false
}
