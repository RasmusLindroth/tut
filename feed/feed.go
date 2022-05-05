package feed

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
)

type apiFunc func(pg *mastodon.Pagination) ([]api.Item, error)
type apiEmptyFunc func() ([]api.Item, error)
type apiIDFunc func(pg *mastodon.Pagination, id mastodon.ID) ([]api.Item, error)
type apiSearchFunc func(search string) ([]api.Item, error)
type apiSearchPGFunc func(pg *mastodon.Pagination, search string) ([]api.Item, error)
type apiThreadFunc func(status *mastodon.Status) ([]api.Item, int, error)

type FeedType uint

const (
	Favorites FeedType = iota
	Favorited
	Boosts
	Followers
	Following
	Blocking
	Muting
	InvalidFeed
	Notification
	Tag
	Thread
	TimelineFederated
	TimelineHome
	TimelineLocal
	Conversations
	User
	UserList
	Lists
	List
)

type LoadingLock struct {
	mux  sync.Mutex
	last time.Time
}

type Feed struct {
	accountClient *api.AccountClient
	feedType      FeedType
	items         []api.Item
	itemsMux      sync.RWMutex
	loadingNewer  *LoadingLock
	loadingOlder  *LoadingLock
	loadNewer     func()
	loadOlder     func()
	Update        chan bool
	apiData       *api.RequestData
	apiDataMux    sync.Mutex
	stream        *api.Receiver
	name          string
	close         func()
}

func (f *Feed) Type() FeedType {
	return f.feedType
}

func (f *Feed) List() []api.Item {
	f.itemsMux.RLock()
	defer f.itemsMux.RUnlock()
	return f.items
}

func (f *Feed) Item(index int) (api.Item, error) {
	f.itemsMux.RLock()
	defer f.itemsMux.RUnlock()
	if index < 0 || index >= len(f.items) {
		return nil, errors.New("item out of range")
	}
	return f.items[index], nil
}

func (f *Feed) Updated() {
	if len(f.Update) > 0 {
		return
	}
	f.Update <- true
}

func (f *Feed) LoadNewer() {
	if f.loadNewer == nil {
		return
	}
	lock := f.loadingNewer.mux.TryLock()
	if !lock {
		return
	}
	if time.Since(f.loadingNewer.last) < (500 * time.Millisecond) {
		f.loadingNewer.mux.Unlock()
		return
	}
	f.loadNewer()
	f.Updated()
	f.loadingNewer.last = time.Now()
	f.loadingNewer.mux.Unlock()
}

func (f *Feed) LoadOlder() {
	if f.loadOlder == nil {
		return
	}
	lock := f.loadingOlder.mux.TryLock()
	if !lock {
		return
	}
	if time.Since(f.loadingOlder.last) < (500 * time.Microsecond) {
		f.loadingOlder.mux.Unlock()
		return
	}
	f.loadOlder()
	f.Updated()
	f.loadingOlder.last = time.Now()
	f.loadingOlder.mux.Unlock()
}

func (f *Feed) HasStream() bool {
	return f.stream != nil
}

func (f *Feed) Close() {
	if f.close != nil {
		f.close()
	}
}

func (f *Feed) Name() string {
	return f.name
}

func (f *Feed) singleNewerSearch(fn apiSearchFunc, search string) {
	items, err := fn(search)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) singleThread(fn apiThreadFunc, status *mastodon.Status) {
	items, _, err := fn(status)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalNewer(fn apiFunc) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[0].Raw().(type) {
		case *mastodon.Status:
			pg.MinID = item.ID
		case *api.NotificationData:
			pg.MinID = item.Item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalOlder(fn apiFunc) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[len(f.items)-1].Raw().(type) {
		case *mastodon.Status:
			pg.MaxID = item.ID
		case *api.NotificationData:
			pg.MaxID = item.Item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) newerSearchPG(fn apiSearchPGFunc, search string) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[0].Raw().(type) {
		case *mastodon.Status:
			pg.MinID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, search)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) olderSearchPG(fn apiSearchPGFunc, search string) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[len(f.items)-1].Raw().(type) {
		case *mastodon.Status:
			pg.MaxID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, search)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalNewerUser(fn apiIDFunc, id mastodon.ID) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 1 {
		switch item := f.items[1].Raw().(type) {
		case *mastodon.Status:
			pg.MinID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, id)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		newItems := []api.Item{f.items[0]}
		newItems = append(newItems, items...)
		if len(f.items) > 1 {
			newItems = append(newItems, f.items[1:]...)
		}
		f.items = newItems
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalOlderUser(fn apiIDFunc, id mastodon.ID) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 1 {
		switch item := f.items[len(f.items)-1].Raw().(type) {
		case *mastodon.Status:
			pg.MaxID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, id)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalNewerID(fn apiIDFunc, id mastodon.ID) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[0].Raw().(type) {
		case *mastodon.Status:
			pg.MinID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, id)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalOlderID(fn apiIDFunc, id mastodon.ID) {
	f.itemsMux.RLock()
	pg := mastodon.Pagination{}
	if len(f.items) > 0 {
		switch item := f.items[len(f.items)-1].Raw().(type) {
		case *mastodon.Status:
			pg.MaxID = item.ID
		}
	}
	f.itemsMux.RUnlock()
	items, err := fn(&pg, id)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalEmpty(fn apiEmptyFunc) {
	items, err := fn()
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkNewer(fn apiFunc) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MinID = f.apiData.MinID
	maxTmp := f.apiData.MaxID

	items, err := fn(pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MinID = pg.MinID
	if pg.MaxID == "" {
		f.apiData.MaxID = maxTmp
	} else {
		f.apiData.MaxID = pg.MaxID
	}
	f.apiDataMux.Unlock()
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkOlder(fn apiFunc) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MaxID = f.apiData.MaxID
	if pg.MaxID == "" {
		f.apiDataMux.Unlock()
		return
	}

	items, err := fn(pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MaxID = pg.MaxID
	f.apiDataMux.Unlock()

	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkNewerID(fn apiIDFunc, id mastodon.ID) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MinID = f.apiData.MinID
	maxTmp := f.apiData.MaxID

	items, err := fn(pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MinID = pg.MinID
	if pg.MaxID == "" {
		f.apiData.MaxID = maxTmp
	} else {
		f.apiData.MaxID = pg.MaxID
	}
	f.apiDataMux.Unlock()
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkOlderID(fn apiIDFunc, id mastodon.ID) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MaxID = f.apiData.MaxID
	if pg.MaxID == "" {
		f.apiDataMux.Unlock()
		return
	}

	items, err := fn(pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MaxID = pg.MaxID
	f.apiDataMux.Unlock()

	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.Updated()
	}
	f.itemsMux.Unlock()
}

func (f *Feed) startStream(rec *api.Receiver, err error) {
	if err != nil {
		log.Fatalln("Couldn't open stream")
	}
	f.stream = rec
	go func() {
		for e := range f.stream.Ch {
			switch t := e.(type) {
			case *mastodon.UpdateEvent:
				s := api.NewStatusItem(t.Status)
				f.itemsMux.Lock()
				f.items = append([]api.Item{s}, f.items...)
				f.Updated()
				f.itemsMux.Unlock()
			}
		}
	}()
}

func (f *Feed) startStreamNotification(rec *api.Receiver, err error) {
	if err != nil {
		log.Fatalln("Couldn't open stream")
	}
	f.stream = rec
	go func() {
		for e := range f.stream.Ch {
			switch t := e.(type) {
			case *mastodon.NotificationEvent:
				rel, err := f.accountClient.Client.GetAccountRelationships(context.Background(), []string{string(t.Notification.Account.ID)})
				if err != nil {
					continue
				}
				if len(rel) == 0 {
					log.Fatalln(t.Notification.Account.Acct)
					continue
				}
				s := api.NewNotificationItem(t.Notification,
					&api.User{
						Data:     &t.Notification.Account,
						Relation: rel[0],
					})
				f.itemsMux.Lock()
				f.items = append([]api.Item{s}, f.items...)
				f.Updated()
				f.itemsMux.Unlock()
			}
		}
	}()
}

func NewTimelineHome(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      TimelineHome,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalNewer(feed.accountClient.GetTimeline) }
	feed.loadOlder = func() { feed.normalOlder(feed.accountClient.GetTimeline) }
	feed.startStream(feed.accountClient.NewHomeStream())
	feed.close = func() { feed.accountClient.RemoveHomeReceiver(feed.stream) }

	return feed
}

func NewTimelineFederated(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      TimelineFederated,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalNewer(feed.accountClient.GetTimelineFederated) }
	feed.loadOlder = func() { feed.normalOlder(feed.accountClient.GetTimelineFederated) }
	feed.startStream(feed.accountClient.NewFederatedStream())
	feed.close = func() { feed.accountClient.RemoveFederatedReceiver(feed.stream) }

	return feed
}

func NewTimelineLocal(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      TimelineLocal,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalNewer(feed.accountClient.GetTimelineLocal) }
	feed.loadOlder = func() { feed.normalOlder(feed.accountClient.GetTimelineLocal) }
	feed.startStream(feed.accountClient.NewLocalStream())
	feed.close = func() { feed.accountClient.RemoveFederatedReceiver(feed.stream) }

	return feed
}

func NewConversations(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Conversations,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalNewer(feed.accountClient.GetConversations) }
	feed.loadOlder = func() { feed.normalOlder(feed.accountClient.GetConversations) }
	feed.startStream(feed.accountClient.NewDirectStream())
	feed.close = func() { feed.accountClient.RemoveFederatedReceiver(feed.stream) }

	return feed
}

func NewNotifications(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Notification,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalNewer(feed.accountClient.GetNotifications) }
	feed.loadOlder = func() { feed.normalOlder(feed.accountClient.GetNotifications) }
	feed.startStreamNotification(feed.accountClient.NewHomeStream())
	feed.close = func() { feed.accountClient.RemoveHomeReceiver(feed.stream) }

	return feed
}

func NewFavorites(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Favorited,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.linkNewer(feed.accountClient.GetFavorites) }
	feed.loadOlder = func() { feed.linkOlder(feed.accountClient.GetFavorites) }

	return feed
}

func NewBookmarks(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Favorites,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.linkNewer(feed.accountClient.GetBookmarks) }
	feed.loadOlder = func() { feed.linkOlder(feed.accountClient.GetBookmarks) }

	return feed
}

func NewUserSearch(ac *api.AccountClient, search string) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      UserList,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		name:          search,
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.singleNewerSearch(feed.accountClient.GetUsers, search) }

	return feed
}

func NewUserProfile(ac *api.AccountClient, user *api.User) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      User,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}
	feed.items = append(feed.items, api.NewUserItem(user, true))
	feed.loadNewer = func() { feed.normalNewerUser(feed.accountClient.GetUser, user.Data.ID) }
	feed.loadOlder = func() { feed.normalOlderUser(feed.accountClient.GetUser, user.Data.ID) }

	return feed
}

func NewThread(ac *api.AccountClient, status *mastodon.Status) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Thread,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.singleThread(feed.accountClient.GetThread, status) }

	return feed
}

func NewTag(ac *api.AccountClient, search string) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Tag,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		name:          search,
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.newerSearchPG(feed.accountClient.GetTag, search) }
	feed.loadOlder = func() { feed.olderSearchPG(feed.accountClient.GetTag, search) }
	feed.startStream(feed.accountClient.NewTagStream(search))
	feed.close = func() { feed.accountClient.RemoveTagReceiver(feed.stream, search) }

	return feed
}

func NewListList(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Lists,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.normalEmpty(feed.accountClient.GetLists) }

	return feed
}

func NewList(ac *api.AccountClient, list *mastodon.List) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      List,
		items:         make([]api.Item, 0),
		loadOlder:     func() {},
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		name:          list.Title,
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}
	feed.loadNewer = func() { feed.normalNewerID(feed.accountClient.GetListStatuses, list.ID) }
	feed.loadOlder = func() { feed.normalOlderID(feed.accountClient.GetListStatuses, list.ID) }
	feed.startStream(feed.accountClient.NewListStream(list.ID))
	feed.close = func() { feed.accountClient.RemoveListReceiver(feed.stream, list.ID) }

	return feed
}

func NewFavoritesStatus(ac *api.AccountClient, id mastodon.ID) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Favorites,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadOlder:     func() {},
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewerID(feed.accountClient.GetFavoritesStatus, id)
		}
		once = false
	}

	return feed
}

func NewBoosts(ac *api.AccountClient, id mastodon.ID) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Boosts,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadOlder:     func() {},
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewerID(feed.accountClient.GetBoostsStatus, id)
		}
		once = false
	}

	return feed
}

func NewFollowers(ac *api.AccountClient, id mastodon.ID) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Followers,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewerID(feed.accountClient.GetFollowers, id)
		}
		once = false
	}
	feed.loadOlder = func() { feed.linkOlderID(feed.accountClient.GetFollowers, id) }

	return feed
}

func NewFollowing(ac *api.AccountClient, id mastodon.ID) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Following,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewerID(feed.accountClient.GetFollowing, id)
		}
		once = false
	}
	feed.loadOlder = func() { feed.linkOlderID(feed.accountClient.GetFollowing, id) }

	return feed
}

func NewBlocking(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Blocking,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	feed.loadNewer = func() { feed.linkNewer(feed.accountClient.GetBlocking) }
	feed.loadOlder = func() { feed.linkOlder(feed.accountClient.GetBlocking) }

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewer(feed.accountClient.GetBlocking)
		}
		once = false
	}
	feed.loadOlder = func() { feed.linkOlder(feed.accountClient.GetBlocking) }

	return feed
}

func NewMuting(ac *api.AccountClient) *Feed {
	feed := &Feed{
		accountClient: ac,
		feedType:      Muting,
		items:         make([]api.Item, 0),
		apiData:       &api.RequestData{},
		Update:        make(chan bool, 1),
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
	}

	once := true
	feed.loadNewer = func() {
		if once {
			feed.linkNewer(feed.accountClient.GetMuting)
		}
		once = false
	}
	feed.loadOlder = func() { feed.linkOlder(feed.accountClient.GetMuting) }

	return feed
}
