package api

import (
	"sync"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/util"
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
	ToggleSpoiler()
	ShowSpoiler() bool
	Raw() interface{}
	URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int)
}

func NewStatusItem(item *mastodon.Status) Item {
	return &StatusItem{id: newID(), item: item, showSpoiler: false}
}

type StatusItem struct {
	id          uint
	item        *mastodon.Status
	showSpoiler bool
}

func (s *StatusItem) ID() uint {
	return s.id
}

func (s *StatusItem) Type() MastodonType {
	return StatusType
}

func (s *StatusItem) ToggleSpoiler() {
	s.showSpoiler = !s.showSpoiler
}

func (s *StatusItem) ShowSpoiler() bool {
	return s.showSpoiler
}

func (s *StatusItem) Raw() interface{} {
	return s.item
}

func (s *StatusItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	status := s.item
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

func (u *UserItem) ToggleSpoiler() {
}

func (u *UserItem) ShowSpoiler() bool {
	return false
}

func (u *UserItem) Raw() interface{} {
	return u.item
}

func (u *UserItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return []util.URL{}, []mastodon.Mention{}, []mastodon.Tag{}, 0
}

func NewNotificationItem(item *mastodon.Notification, user *User) Item {
	n := &NotificationItem{
		id:          newID(),
		item:        item,
		showSpoiler: false,
		user:        NewUserItem(user, false),
		status:      NewStatusItem(item.Status),
	}

	return n
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

func (n *NotificationItem) ToggleSpoiler() {
	n.showSpoiler = !n.showSpoiler
}

func (n *NotificationItem) ShowSpoiler() bool {
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
	return nil, nil, nil, 0
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

func (s *ListItem) ToggleSpoiler() {
}

func (s *ListItem) ShowSpoiler() bool {
	return true
}

func (s *ListItem) Raw() interface{} {
	return s.item
}

func (s *ListItem) URLs() ([]util.URL, []mastodon.Mention, []mastodon.Tag, int) {
	return nil, nil, nil, 0
}
