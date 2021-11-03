package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-mastodon"
	"github.com/rivo/tview"
)

type FeedType uint

const (
	TimelineFeedType FeedType = iota
	ThreadFeedType
	UserFeedType
	UserListFeedType
	NotificationFeedType
	TagFeedType
	ListFeedType
)

type Feed interface {
	GetFeedList() <-chan ListItem
	LoadNewer() int
	LoadOlder() int
	DrawList()
	DrawToot()
	DrawSpoiler()
	RedrawControls()
	GetCurrentUser() *mastodon.Account
	GetCurrentStatus() *mastodon.Status
	FeedType() FeedType
	GetSavedIndex() int
	GetDesc() string
	Input(event *tcell.EventKey)
}

type Poll struct {
	ID         string
	ExpiresAt  time.Time
	Expired    bool
	Multiple   bool
	VotesCount int64
	Options    []PollOption
	Voted      bool
}

type PollOption struct {
	Title      string
	VotesCount int64
	Percent    string
}

type Media struct {
	Type        string
	Description string
	URL         string
}

type Card struct {
	Type        string
	Title       string
	Description string
	URL         string
}

type Toot struct {
	Visibility         string
	Boosted            bool
	BoostedDisplayName string
	BoostedAcct        string
	Bookmarked         bool
	AccountDisplayName string
	Account            string
	Spoiler            bool
	SpoilerText        string
	ShowSpoiler        bool
	ContentText        string
	Width              int
	HasExtra           bool
	Poll               Poll
	Media              []Media
	Card               Card
	Replies            int
	Boosts             int
	Favorites          int
	Controls           string
}

type DisplayTootData struct {
	Toot  Toot
	Style StyleConfig
}

func showTootOptions(app *App, status *mastodon.Status, showSensitive bool) (string, string) {
	var strippedContent string
	var strippedSpoiler string
	var urls []URL
	var u []URL

	strippedContent, urls = cleanTootHTML(status.Content)
	strippedContent = tview.Escape(strippedContent)

	toot := Toot{
		Width:              app.UI.StatusView.GetTextWidth(),
		ContentText:        strippedContent,
		Boosted:            status.Reblog != nil,
		BoostedDisplayName: tview.Escape(status.Account.DisplayName),
		BoostedAcct:        tview.Escape(status.Account.Acct),
		ShowSpoiler:        showSensitive,
	}

	if status.Reblog != nil {
		status = status.Reblog
	}

	toot.AccountDisplayName = tview.Escape(status.Account.DisplayName)
	toot.Account = tview.Escape(status.Account.Acct)
	toot.Bookmarked = status.Bookmarked.(bool)
	toot.Visibility = status.Visibility
	toot.Spoiler = status.Sensitive

	if status.Poll != nil {
		p := *status.Poll
		toot.Poll = Poll{
			ID:         string(p.ID),
			ExpiresAt:  p.ExpiresAt,
			Expired:    p.Expired,
			Multiple:   p.Multiple,
			VotesCount: p.VotesCount,
			Voted:      p.Voted,
			Options:    []PollOption{},
		}
		for _, item := range p.Options {
			percent := 0.0
			if p.VotesCount > 0 {
				percent = float64(item.VotesCount) / float64(p.VotesCount) * 100
			}

			o := PollOption{
				Title:      tview.Escape(item.Title),
				VotesCount: item.VotesCount,
				Percent:    fmt.Sprintf("%.2f", percent),
			}
			toot.Poll.Options = append(toot.Poll.Options, o)
		}

	} else {
		toot.Poll = Poll{}
	}

	if status.Sensitive {
		strippedSpoiler, u = cleanTootHTML(status.SpoilerText)
		strippedSpoiler = tview.Escape(strippedSpoiler)
		urls = append(urls, u...)
	}

	toot.SpoilerText = strippedSpoiler
	app.UI.LinkOverlay.SetLinks(urls, status)

	media := []Media{}
	for _, att := range status.MediaAttachments {
		m := Media{
			Type:        att.Type,
			Description: tview.Escape(att.Description),
			URL:         att.URL,
		}
		media = append(media, m)
	}
	toot.Media = media

	if status.Card != nil {
		toot.Card = Card{
			Type:        status.Card.Type,
			Title:       tview.Escape(strings.TrimSpace(status.Card.Title)),
			Description: tview.Escape(strings.TrimSpace(status.Card.Description)),
			URL:         status.Card.URL,
		}
	} else {
		toot.Card = Card{}
	}

	toot.HasExtra = len(status.MediaAttachments) > 0 || status.Card != nil || status.Poll != nil
	toot.Replies = int(status.RepliesCount)
	toot.Boosts = int(status.ReblogsCount)
	toot.Favorites = int(status.FavouritesCount)

	app.UI.StatusView.ScrollToBeginning()

	var info []string
	if status.Favourited == true {
		info = append(info, ColorKey(app.Config, "Un", "F", "avorite"))
	} else {
		info = append(info, ColorKey(app.Config, "", "F", "avorite"))
	}
	if status.Reblogged == true {
		info = append(info, ColorKey(app.Config, "Un", "B", "oost"))
	} else {
		info = append(info, ColorKey(app.Config, "", "B", "oost"))
	}
	info = append(info, ColorKey(app.Config, "", "T", "hread"))
	info = append(info, ColorKey(app.Config, "", "R", "eply"))
	info = append(info, ColorKey(app.Config, "", "V", "iew"))
	info = append(info, ColorKey(app.Config, "", "U", "ser"))
	if len(status.MediaAttachments) > 0 {
		info = append(info, ColorKey(app.Config, "", "M", "edia"))
	}
	if len(urls)+len(status.Mentions)+len(status.Tags) > 0 {
		info = append(info, ColorKey(app.Config, "", "O", "pen"))
	}
	info = append(info, ColorKey(app.Config, "", "A", "vatar"))
	if status.Account.ID == app.Me.ID {
		info = append(info, ColorKey(app.Config, "", "D", "elete"))
	}

	if status.Bookmarked == false {
		info = append(info, ColorKey(app.Config, "", "S", "ave"))
	} else {
		info = append(info, ColorKey(app.Config, "Un", "S", "ave"))
	}
	info = append(info, ColorKey(app.Config, "", "Y", "ank"))

	controls := strings.Join(info, " ")

	td := DisplayTootData{
		Toot:  toot,
		Style: app.Config.Style,
	}
	var output bytes.Buffer
	err := app.Config.Templates.TootTemplate.ExecuteTemplate(&output, "toot.tmpl", td)
	if err != nil {
		panic(err)
	}

	return output.String(), controls
}

type User struct {
	Username       string
	Account        string
	DisplayName    string
	Locked         bool
	CreatedAt      time.Time
	FollowersCount int64
	FollowingCount int64
	StatusCount    int64
	Note           string
	URL            string
	Avatar         string
	AvatarStatic   string
	Header         string
	HeaderStatic   string
	Fields         []Field
	Bot            bool
	//Emojis         []Emoji
	//Moved *Account `json:"moved"`
}

type Field struct {
	Name       string
	Value      string
	VerifiedAt time.Time
}

type DisplayUserData struct {
	User  User
	Style StyleConfig
}

func showUser(app *App, user *mastodon.Account, relation *mastodon.Relationship, showUserControl bool) (string, string) {
	u := User{
		Username:       tview.Escape(user.Username),
		Account:        tview.Escape(user.Acct),
		DisplayName:    tview.Escape(user.DisplayName),
		Locked:         user.Locked,
		CreatedAt:      user.CreatedAt,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		StatusCount:    user.StatusesCount,
		URL:            user.URL,
		Avatar:         user.Avatar,
		AvatarStatic:   user.AvatarStatic,
		Header:         user.Header,
		HeaderStatic:   user.HeaderStatic,
	}

	var controls string

	var urls []URL
	fields := []Field{}
	u.Note, urls = cleanTootHTML(user.Note)
	for _, f := range user.Fields {
		value, fu := cleanTootHTML(f.Value)
		fields = append(fields, Field{
			Name:       tview.Escape(f.Name),
			Value:      tview.Escape(value),
			VerifiedAt: f.VerifiedAt,
		})
		urls = append(urls, fu...)
	}
	u.Fields = fields

	app.UI.LinkOverlay.SetLinks(urls, nil)

	var controlItems []string
	if app.Me.ID != user.ID {
		if relation.Following {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "F", "ollow"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "F", "ollow"))
		}
		if relation.Blocking {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "B", "lock"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "B", "lock"))
		}
		if relation.Muting {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "M", "ute"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "M", "ute"))
		}
		if len(urls) > 0 {
			controlItems = append(controlItems, ColorKey(app.Config, "", "O", "pen"))
		}
	}
	if showUserControl {
		controlItems = append(controlItems, ColorKey(app.Config, "", "U", "ser"))
	}
	controlItems = append(controlItems, ColorKey(app.Config, "", "A", "vatar"))
	controlItems = append(controlItems, ColorKey(app.Config, "", "Y", "ank"))
	controls = strings.Join(controlItems, " ")

	ud := DisplayUserData{
		User:  u,
		Style: app.Config.Style,
	}
	var output bytes.Buffer
	err := app.Config.Templates.UserTemplate.ExecuteTemplate(&output, "user.tmpl", ud)
	if err != nil {
		panic(err)
	}

	return output.String(), controls
}

func showUserBak(app *App, user *mastodon.Account, relation *mastodon.Relationship, showUserControl bool) (string, string) {
	var text string
	var controls string

	n := ColorMark(app.Config.Style.Text)
	s1 := ColorMark(app.Config.Style.TextSpecial1)
	s2 := ColorMark(app.Config.Style.TextSpecial2)

	if user.DisplayName != "" {
		text = fmt.Sprintf(s2+"%s\n", user.DisplayName)
	}
	text += fmt.Sprintf(s1+"%s\n\n", user.Acct)

	text += fmt.Sprintf("%sToots %s%d %sFollowers %s%d %sFollowing %s%d\n\n",
		n, s2, user.StatusesCount, n, s2, user.FollowersCount, n, s2, user.FollowingCount)

	note, urls := cleanTootHTML(user.Note)
	text += n + note + "\n\n"

	for _, f := range user.Fields {
		value, fu := cleanTootHTML(f.Value)
		text += fmt.Sprintf("%s%s: %s%s\n", s2, f.Name, n, value)
		urls = append(urls, fu...)
	}

	app.UI.LinkOverlay.SetLinks(urls, nil)

	var controlItems []string
	if app.Me.ID != user.ID {
		if relation.Following {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "F", "ollow"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "F", "ollow"))
		}
		if relation.Blocking {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "B", "lock"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "B", "lock"))
		}
		if relation.Muting {
			controlItems = append(controlItems, ColorKey(app.Config, "Un", "M", "ute"))
		} else {
			controlItems = append(controlItems, ColorKey(app.Config, "", "M", "ute"))
		}
		if len(urls) > 0 {
			controlItems = append(controlItems, ColorKey(app.Config, "", "O", "pen"))
		}
	}
	if showUserControl {
		controlItems = append(controlItems, ColorKey(app.Config, "", "U", "ser"))
	}
	controlItems = append(controlItems, ColorKey(app.Config, "", "A", "vatar"))
	controlItems = append(controlItems, ColorKey(app.Config, "", "Y", "ank"))
	controls = strings.Join(controlItems, " ")

	return text, controls
}

func drawStatusList(statuses []*mastodon.Status, longFormat, shortFormat string, relativeDate int) <-chan ListItem {
	ch := make(chan ListItem)
	go func() {
		today := time.Now()
		for _, s := range statuses {
			sLocal := s.CreatedAt.Local()
			dateOutput := OutputDate(sLocal, today, longFormat, shortFormat, relativeDate)

			content := fmt.Sprintf("%s %s", dateOutput, s.Account.Acct)
			iconText := ""
			rs := s
			if s.Reblog != nil {
				rs = s.Reblog
			}
			if rs.RepliesCount > 0 || (rs.InReplyToID != nil && rs.InReplyToID != "") {
				iconText = " ⤶ "
			}
			ch <- ListItem{Text: content, Icons: iconText}
		}
		close(ch)
	}()
	return ch
}

func openAvatar(app *App, user mastodon.Account) {
	f, err := downloadFile(user.AvatarStatic)
	if err != nil {
		app.UI.CmdBar.ShowError("Couldn't open avatar")
		return
	}
	openMediaType(app.UI.Root, app.Config.Media, []string{f}, "image")
}

type ControlItem uint

const (
	ControlAvatar ControlItem = 1 << iota
	ControlBlock
	ControlBoost
	ControlCompose
	ControlDelete
	ControlEnter
	ControlFavorite
	ControlFollow
	ControlList
	ControlMedia
	ControlMute
	ControlOpen
	ControlReply
	ControlThread
	ControlUser
	ControlSpoiler
	ControlBookmark
	ControlYankStatus
	ControlYankUser
)

func inputOptions(options []ControlItem) ControlItem {
	var controls ControlItem
	for _, o := range options {
		controls = controls | o
	}
	return controls
}

func inputSimple(app *App, event *tcell.EventKey, controls ControlItem,
	user mastodon.Account, status *mastodon.Status, originalStatus *mastodon.Status, relation *mastodon.Relationship, feed Feed, listInfo *ListInfo) (updated bool,
	redrawControls bool, redrawToot bool, newStatus *mastodon.Status, newRelation *mastodon.Relationship) {

	newStatus = status
	newRelation = relation
	var err error

	if event.Key() == tcell.KeyEnter {
		if controls&ControlEnter == 0 {
			return
		}
		if controls&ControlUser != 0 {
			app.UI.StatusView.AddFeed(
				NewUserFeed(app, user),
			)
		}
		if controls&ControlList != 0 {
			app.UI.StatusView.AddFeed(
				NewTimelineFeed(app, TimelineList, listInfo),
			)
		}
	}

	if event.Key() != tcell.KeyRune {
		return
	}

	switch event.Rune() {
	case 'a', 'A':
		if controls&ControlAvatar != 0 {
			openAvatar(app, user)
		}
	case 'b', 'B':
		if controls&ControlBoost != 0 {
			newStatus, err = app.API.BoostToggle(status)
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't boost toot. Error: %v\n", err))
				return
			}
			updated = true
			redrawControls = true
		}
		if controls&ControlBlock != 0 {
			if relation.Blocking {
				newRelation, err = app.API.UnblockUser(user)
			} else {
				newRelation, err = app.API.BlockUser(user)
			}
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't block/unblock user. Error: %v\n", err))
				return
			}
			updated = true
			redrawToot = true
			redrawControls = true
		}
	case 'd', 'D':
		if controls&ControlDelete != 0 {
			err = app.API.DeleteStatus(status)
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't delete toot. Error: %v\n", err))
			} else {
				status.Card = nil
				status.Sensitive = false
				status.SpoilerText = ""
				status.Favourited = false
				status.MediaAttachments = nil
				status.Reblogged = false
				status.Content = "Deleted"
				updated = true
				redrawToot = true
			}
		}
	case 'f', 'F':
		if controls&ControlFavorite != 0 {
			newStatus, err = app.API.FavoriteToogle(status)
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't toggle favorite on toot. Error: %v\n", err))
				return
			}
			updated = true
			redrawControls = true
		}
		if controls&ControlFollow != 0 {
			if relation.Following {
				newRelation, err = app.API.UnfollowUser(user)
			} else {
				newRelation, err = app.API.FollowUser(user)
			}
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't follow/unfollow user. Error: %v\n", err))
				return
			}
			updated = true
			redrawToot = true
			redrawControls = true
		}
	case 'c', 'C':
		if controls&ControlCompose != 0 {
			app.UI.NewToot()
		}
	case 'm', 'M':
		if controls&ControlMedia != 0 {
			app.UI.OpenMedia(status)
		}
		if controls&ControlMute != 0 {
			if relation.Muting {
				newRelation, err = app.API.UnmuteUser(user)
			} else {
				newRelation, err = app.API.MuteUser(user)
			}
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't mute/unmute user. Error: %v\n", err))
				return
			}
			updated = true
			redrawToot = true
			redrawControls = true
		}
	case 'o', 'O':
		if controls&ControlOpen != 0 {
			app.UI.ShowLinks()
		}
		if controls&ControlList != 0 {
			app.UI.StatusView.AddFeed(
				NewTimelineFeed(app, TimelineList, listInfo),
			)
		}
	case 'r', 'R':
		if controls&ControlReply != 0 {
			app.UI.Reply(status)
		}
	case 's', 'S':
		if controls&ControlBookmark != 0 {
			tmpStatus, err := app.API.BookmarkToogle(status)
			newStatus = originalStatus
			if newStatus.Reblog != nil {
				newStatus.Reblog.Bookmarked = tmpStatus.Bookmarked
			} else {
				newStatus.Bookmarked = tmpStatus.Bookmarked
			}
			if err != nil {
				app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't toggle bookmark on toot. Error: %v\n", err))
				return
			}
			updated = true
			redrawControls = true
			redrawToot = true
		}
	case 't', 'T':
		if controls&ControlThread != 0 {
			app.UI.StatusView.AddFeed(
				NewThreadFeed(app, status),
			)
		}
	case 'u', 'U':
		if controls&ControlUser != 0 {
			app.UI.StatusView.AddFeed(
				NewUserFeed(app, user),
			)
		}
	case 'y', 'Y':
		if controls&ControlYankStatus != 0 {
			copyToClipboard(status.URL)
		}
		if controls&ControlYankUser != 0 {
			copyToClipboard(user.URL)
		}
	case 'z', 'Z':
		if controls&ControlSpoiler != 0 {
			feed.DrawSpoiler()
		}
	}
	return
}

func userFromStatus(s *mastodon.Status) *mastodon.Account {
	if s == nil {
		return nil
	}
	if s.Reblog != nil {
		s = s.Reblog
	}
	return &s.Account
}

func NewTimelineFeed(app *App, tl TimelineType, listInfo *ListInfo) *TimelineFeed {
	t := &TimelineFeed{
		app:          app,
		timelineType: tl,
		linkPrev:     "",
		linkNext:     "",
		listInfo:     listInfo,
	}
	var err error
	t.statuses, err = t.app.API.GetStatuses(t)
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load timeline toots. Error: %v\n", err))
	}
	return t
}

type TimelineFeed struct {
	app          *App
	timelineType TimelineType
	statuses     []*mastodon.Status
	linkPrev     mastodon.ID //Only bm and fav
	linkNext     mastodon.ID //Only bm and fav
	listInfo     *ListInfo   //only lists
	index        int
	showSpoiler  bool
}

func (t *TimelineFeed) FeedType() FeedType {
	return TimelineFeedType
}

func (t *TimelineFeed) GetDesc() string {
	switch t.timelineType {
	case TimelineHome:
		return "Timeline home"
	case TimelineDirect:
		return "Timeline direct"
	case TimelineLocal:
		return "Timeline local"
	case TimelineFederated:
		return "Timeline federated"
	case TimelineBookmarked:
		return "Bookmarks"
	case TimelineFavorited:
		return "Favorited"
	case TimelineList:
		return fmt.Sprintf("List: %s", t.listInfo.name)
	}
	return "Timeline"
}

func (t *TimelineFeed) GetCurrentStatus() *mastodon.Status {
	index := t.app.UI.StatusView.GetCurrentItem()
	if index >= len(t.statuses) {
		return nil
	}
	return t.statuses[index]
}

func (t *TimelineFeed) GetCurrentUser() *mastodon.Account {
	return userFromStatus(t.GetCurrentStatus())
}

func (t *TimelineFeed) GetFeedList() <-chan ListItem {
	return drawStatusList(t.statuses, t.app.Config.General.DateFormat, t.app.Config.General.DateTodayFormat, t.app.Config.General.DateRelative)
}

func (t *TimelineFeed) LoadNewer() int {
	var statuses []*mastodon.Status
	var err error
	if len(t.statuses) == 0 {
		statuses, err = t.app.API.GetStatuses(t)
	} else {
		statuses, err = t.app.API.GetStatusesNewer(t)
		newCount := len(statuses)
		if newCount > 0 {
			Notify(t.app.Config.NotificationConfig, NotificationPost,
				fmt.Sprintf("%d new toots", newCount), "")
		}
	}
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load new toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	old := t.statuses
	t.statuses = append(statuses, old...)
	return len(statuses)
}

func (t *TimelineFeed) LoadOlder() int {
	var statuses []*mastodon.Status
	var err error
	if len(t.statuses) == 0 {
		statuses, err = t.app.API.GetStatuses(t)
	} else {
		statuses, err = t.app.API.GetStatusesOlder(t)
	}
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load older toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	t.statuses = append(t.statuses, statuses...)
	return len(statuses)
}

func (t *TimelineFeed) DrawList() {
	t.app.UI.StatusView.SetList(t.GetFeedList())
}

func (t *TimelineFeed) DrawSpoiler() {
	t.showSpoiler = true
	t.DrawToot()
}

func (t *TimelineFeed) DrawToot() {
	if len(t.statuses) == 0 {
		t.app.UI.StatusView.SetText("")
		t.app.UI.StatusView.SetControls("")
		return
	}
	t.index = t.app.UI.StatusView.GetCurrentItem()
	text, controls := showTootOptions(t.app, t.statuses[t.index], t.showSpoiler)
	t.showSpoiler = false
	t.app.UI.StatusView.SetText(text)
	t.app.UI.StatusView.SetControls(controls)
	t.app.UI.Root.Sync()
}

func (t *TimelineFeed) RedrawControls() {
	status := t.GetCurrentStatus()
	if status == nil {
		return
	}
	_, controls := showTootOptions(t.app, status, t.showSpoiler)
	t.app.UI.StatusView.SetControls(controls)
}

func (t *TimelineFeed) GetSavedIndex() int {
	return t.index
}

func (t *TimelineFeed) Input(event *tcell.EventKey) {
	status := t.GetCurrentStatus()
	originalStatus := status
	if status == nil {
		return
	}
	if status.Reblog != nil {
		status = status.Reblog
	}
	user := status.Account

	controls := []ControlItem{
		ControlAvatar, ControlThread, ControlUser, ControlSpoiler,
		ControlCompose, ControlOpen, ControlReply, ControlMedia,
		ControlFavorite, ControlBoost, ControlDelete, ControlBookmark,
		ControlYankStatus,
	}
	options := inputOptions(controls)

	updated, rc, rt, newS, _ := inputSimple(t.app, event, options, user, status, originalStatus, nil, t, nil)
	if updated {
		index := t.app.UI.StatusView.GetCurrentItem()
		t.statuses[index] = newS
	}
	if rc {
		t.RedrawControls()
	}
	if rt {
		t.DrawToot()
	}
}

func NewThreadFeed(app *App, s *mastodon.Status) *ThreadFeed {
	t := &ThreadFeed{
		app: app,
	}
	statuses, index, err := t.app.API.GetThread(s)
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load thread. Error: %v\n", err))
	}
	t.statuses = statuses
	t.status = s
	t.index = index
	return t
}

type ThreadFeed struct {
	app         *App
	statuses    []*mastodon.Status
	status      *mastodon.Status
	index       int
	showSpoiler bool
}

func (t *ThreadFeed) FeedType() FeedType {
	return ThreadFeedType
}

func (t *ThreadFeed) GetDesc() string {
	return "Thread"
}

func (t *ThreadFeed) GetCurrentStatus() *mastodon.Status {
	index := t.app.UI.StatusView.GetCurrentItem()
	if index >= len(t.statuses) {
		return nil
	}
	return t.statuses[t.app.UI.StatusView.GetCurrentItem()]
}

func (t *ThreadFeed) GetCurrentUser() *mastodon.Account {
	return userFromStatus(t.GetCurrentStatus())
}

func (t *ThreadFeed) GetFeedList() <-chan ListItem {
	return drawStatusList(t.statuses, t.app.Config.General.DateFormat, t.app.Config.General.DateTodayFormat, t.app.Config.General.DateRelative)
}

func (t *ThreadFeed) LoadNewer() int {
	return 0
}

func (t *ThreadFeed) LoadOlder() int {
	return 0
}

func (t *ThreadFeed) DrawList() {
	t.app.UI.StatusView.SetList(t.GetFeedList())
}

func (t *ThreadFeed) DrawSpoiler() {
	t.showSpoiler = true
	t.DrawToot()
}

func (t *ThreadFeed) DrawToot() {
	status := t.GetCurrentStatus()
	if status == nil {
		t.app.UI.StatusView.SetText("")
		t.app.UI.StatusView.SetControls("")
		return
	}
	t.index = t.app.UI.StatusView.GetCurrentItem()
	text, controls := showTootOptions(t.app, status, t.showSpoiler)
	t.showSpoiler = false
	t.app.UI.StatusView.SetText(text)
	t.app.UI.StatusView.SetControls(controls)
	t.app.UI.Root.Sync()
}

func (t *ThreadFeed) RedrawControls() {
	status := t.GetCurrentStatus()
	if status == nil {
		t.app.UI.StatusView.SetText("")
		t.app.UI.StatusView.SetControls("")
		return
	}
	_, controls := showTootOptions(t.app, status, t.showSpoiler)
	t.app.UI.StatusView.SetControls(controls)
}

func (t *ThreadFeed) GetSavedIndex() int {
	return t.index
}

func (t *ThreadFeed) Input(event *tcell.EventKey) {
	status := t.GetCurrentStatus()
	originalStatus := status
	if status == nil {
		return
	}
	if status.Reblog != nil {
		status = status.Reblog
	}
	user := status.Account

	controls := []ControlItem{
		ControlAvatar, ControlUser, ControlSpoiler,
		ControlCompose, ControlOpen, ControlReply, ControlMedia,
		ControlFavorite, ControlBoost, ControlDelete, ControlBookmark,
		ControlYankStatus,
	}
	if status.ID != t.status.ID {
		controls = append(controls, ControlThread)
	}
	options := inputOptions(controls)

	updated, rc, rt, newS, _ := inputSimple(t.app, event, options, user, status, originalStatus, nil, t, nil)
	if updated {
		index := t.app.UI.StatusView.GetCurrentItem()
		t.statuses[index] = newS
	}
	if rc {
		t.RedrawControls()
	}
	if rt {
		t.DrawToot()
	}
}

func NewUserFeed(app *App, a mastodon.Account) *UserFeed {
	u := &UserFeed{
		app: app,
	}
	statuses, err := app.API.GetUserStatuses(a)
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load user toots. Error: %v\n", err))
		return u
	}
	u.statuses = statuses
	relation, err := app.API.UserRelation(a)
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load user data. Error: %v\n", err))
		return u
	}
	u.relation = relation
	u.user = a
	return u
}

type UserFeed struct {
	app         *App
	statuses    []*mastodon.Status
	user        mastodon.Account
	relation    *mastodon.Relationship
	index       int
	showSpoiler bool
}

func (u *UserFeed) FeedType() FeedType {
	return UserFeedType
}

func (u *UserFeed) GetDesc() string {
	return "User " + u.user.Acct
}

func (u *UserFeed) GetCurrentStatus() *mastodon.Status {
	index := u.app.UI.app.UI.StatusView.GetCurrentItem()
	if index > 0 && index-1 >= len(u.statuses) {
		return nil
	}
	return u.statuses[index-1]
}

func (u *UserFeed) GetCurrentUser() *mastodon.Account {
	return &u.user
}

func (u *UserFeed) GetFeedList() <-chan ListItem {
	ch := make(chan ListItem)
	go func() {
		ch <- ListItem{Text: "Profile", Icons: ""}
		for s := range drawStatusList(u.statuses, u.app.Config.General.DateFormat, u.app.Config.General.DateTodayFormat, u.app.Config.General.DateRelative) {
			ch <- s
		}
		close(ch)
	}()
	return ch
}

func (u *UserFeed) LoadNewer() int {
	var statuses []*mastodon.Status
	var err error
	if len(u.statuses) == 0 {
		statuses, err = u.app.API.GetUserStatuses(u.user)
	} else {
		statuses, err = u.app.API.GetUserStatusesNewer(u.user, u.statuses[0])
	}
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load new toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	old := u.statuses
	u.statuses = append(statuses, old...)
	return len(statuses)
}

func (u *UserFeed) LoadOlder() int {
	var statuses []*mastodon.Status
	var err error
	if len(u.statuses) == 0 {
		statuses, err = u.app.API.GetUserStatuses(u.user)
	} else {
		statuses, err = u.app.API.GetUserStatusesOlder(u.user, u.statuses[len(u.statuses)-1])
	}
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load older toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	u.statuses = append(u.statuses, statuses...)
	return len(statuses)
}

func (u *UserFeed) DrawList() {
	u.app.UI.StatusView.SetList(u.GetFeedList())
}

func (u *UserFeed) DrawSpoiler() {
	u.showSpoiler = true
	u.DrawToot()
}

func (u *UserFeed) DrawToot() {
	u.index = u.app.UI.StatusView.GetCurrentItem()

	var text string
	var controls string

	if u.index == 0 {
		text, controls = showUser(u.app, &u.user, u.relation, false)
	} else {
		status := u.GetCurrentStatus()
		if status == nil {
			text = ""
			controls = ""
		} else {
			text, controls = showTootOptions(u.app, status, u.showSpoiler)
		}
		u.showSpoiler = false
	}

	u.app.UI.StatusView.SetText(text)
	u.app.UI.StatusView.SetControls(controls)
	u.app.UI.Root.Sync()
}

func (u *UserFeed) RedrawControls() {
	var controls string
	status := u.GetCurrentStatus()
	if status == nil {
		controls = ""
	} else {
		_, controls = showTootOptions(u.app, status, u.showSpoiler)
	}
	u.app.UI.StatusView.SetControls(controls)
}

func (u *UserFeed) GetSavedIndex() int {
	return u.index
}

func (u *UserFeed) Input(event *tcell.EventKey) {
	index := u.GetSavedIndex()

	if index == 0 {
		controls := []ControlItem{
			ControlAvatar, ControlFollow, ControlBlock, ControlMute, ControlOpen,
			ControlYankUser,
		}
		options := inputOptions(controls)

		updated, _, _, _, newRel := inputSimple(u.app, event, options, u.user, nil, nil, u.relation, u, nil)
		if updated {
			u.relation = newRel
			u.DrawToot()
		}
		return
	}

	status := u.GetCurrentStatus()
	originalStatus := status
	if status == nil {
		return
	}
	if status.Reblog != nil {
		status = status.Reblog
	}
	user := status.Account

	controls := []ControlItem{
		ControlAvatar, ControlThread, ControlSpoiler, ControlCompose,
		ControlOpen, ControlReply, ControlMedia, ControlFavorite, ControlBoost,
		ControlDelete, ControlUser, ControlBookmark, ControlYankStatus,
	}
	options := inputOptions(controls)

	updated, rc, rt, newS, _ := inputSimple(u.app, event, options, user, status, originalStatus, nil, u, nil)
	if updated {
		index := u.app.UI.StatusView.GetCurrentItem()
		u.statuses[index-1] = newS
	}
	if rc {
		u.RedrawControls()
	}
	if rt {
		u.DrawToot()
	}
}

func NewNotificationFeed(app *App, docked bool) *NotificationsFeed {
	n := &NotificationsFeed{
		app:    app,
		docked: docked,
	}
	n.notifications, _ = n.app.API.GetNotifications()
	return n
}

type Notification struct {
	N *mastodon.Notification
	R *mastodon.Relationship
}

type NotificationsFeed struct {
	app           *App
	notifications []*Notification
	docked        bool
	index         int
	showSpoiler   bool
}

func (n *NotificationsFeed) FeedType() FeedType {
	return NotificationFeedType
}

func (n *NotificationsFeed) GetDesc() string {
	return "Notifications"
}

func (n *NotificationsFeed) GetCurrentNotification() *Notification {
	var index int
	if n.docked {
		index = n.app.UI.StatusView.notificationView.list.GetCurrentItem()
	} else {
		index = n.app.UI.StatusView.GetCurrentItem()
	}
	if index >= len(n.notifications) {
		return nil
	}
	return n.notifications[index]
}

func (n *NotificationsFeed) GetCurrentStatus() *mastodon.Status {
	notification := n.GetCurrentNotification()
	if notification.N == nil {
		return nil
	}
	return notification.N.Status
}

func (n *NotificationsFeed) GetCurrentUser() *mastodon.Account {
	notification := n.GetCurrentNotification()
	if notification.N == nil {
		return nil
	}
	return &notification.N.Account
}

func (n *NotificationsFeed) GetFeedList() <-chan ListItem {
	ch := make(chan ListItem)
	notifications := n.notifications
	go func() {
		today := time.Now()
		for _, item := range notifications {
			sLocal := item.N.CreatedAt.Local()
			long := n.app.Config.General.DateFormat
			short := n.app.Config.General.DateTodayFormat
			relative := n.app.Config.General.DateRelative

			dateOutput := OutputDate(sLocal, today, long, short, relative)

			iconText := ""
			switch item.N.Type {
			case "follow", "follow_request":
				iconText += " + "
			case "favourite":
				iconText = " ★ "
			case "reblog":
				iconText = " ♺ "
			case "mention":
				iconText = " ⤶ "
			case "poll":
				iconText = " = "
			}

			content := fmt.Sprintf("%s %s", dateOutput, item.N.Account.Acct)
			ch <- ListItem{Text: content, Icons: iconText}
		}
		close(ch)
	}()
	return ch
}

func (n *NotificationsFeed) LoadNewer() int {
	var notifications []*Notification
	var err error
	if len(n.notifications) == 0 {
		notifications, err = n.app.API.GetNotifications()
	} else {
		notifications, err = n.app.API.GetNotificationsNewer(n.notifications[0])
		for _, o := range notifications {
			switch o.N.Type {
			case "follow":
				Notify(n.app.Config.NotificationConfig, NotificationFollower,
					"New follower", fmt.Sprintf("%s follows you", o.N.Account.Username))
			case "favourite":
				Notify(n.app.Config.NotificationConfig, NotificationFavorite,
					"Favorited your toot", fmt.Sprintf("%s favorited your toot", o.N.Account.Username))
			case "reblog":
				Notify(n.app.Config.NotificationConfig, NotificationBoost,
					"Boosted your toot", fmt.Sprintf("%s boosted your toot", o.N.Account.Username))
			case "mention":
				Notify(n.app.Config.NotificationConfig, NotificationMention,
					"Mentioned in toot", fmt.Sprintf("%s mentioned you", o.N.Account.Username))
			case "poll":
				Notify(n.app.Config.NotificationConfig, NotificationPoll,
					"Poll has ended", "")
			default:
			}
		}
	}
	if err != nil {
		n.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load new toots. Error: %v\n", err))
		return 0
	}
	if len(notifications) == 0 {
		return 0
	}
	old := n.notifications
	n.notifications = append(notifications, old...)
	return len(notifications)
}

func (n *NotificationsFeed) LoadOlder() int {
	var notifications []*Notification
	var err error
	if len(n.notifications) == 0 {
		notifications, err = n.app.API.GetNotifications()
	} else {
		notifications, err = n.app.API.GetNotificationsOlder(n.notifications[len(n.notifications)-1])
	}
	if err != nil {
		n.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load older toots. Error: %v\n", err))
		return 0
	}
	if len(notifications) == 0 {
		return 0
	}
	n.notifications = append(n.notifications, notifications...)
	return len(notifications)
}

func (n *NotificationsFeed) DrawList() {
	if n.docked {
		n.app.UI.StatusView.notificationView.SetList(n.GetFeedList())
	} else {
		n.app.UI.StatusView.SetList(n.GetFeedList())
	}
}

func (n *NotificationsFeed) DrawSpoiler() {
	n.showSpoiler = true
	n.DrawToot()
}

func (n *NotificationsFeed) DrawToot() {
	if n.docked {
		n.index = n.app.UI.StatusView.notificationView.list.GetCurrentItem()
	} else {
		n.index = n.app.UI.StatusView.GetCurrentItem()
	}
	notification := n.GetCurrentNotification()
	if notification == nil {
		n.app.UI.StatusView.SetText("")
		n.app.UI.StatusView.SetControls("")
		return
	}
	var text string
	var controls string
	defer func() { n.showSpoiler = false }()

	switch notification.N.Type {
	case "follow":
		text = SublteText(n.app.Config.Style, FormatUsername(notification.N.Account)+" started following you\n\n")
		var t string
		t, controls = showUser(n.app, &notification.N.Account, notification.R, true)
		text += t
	case "favourite":
		pre := SublteText(n.app.Config.Style, FormatUsername(notification.N.Account)+" favorited your toot") + "\n\n"
		text, controls = showTootOptions(n.app, notification.N.Status, n.showSpoiler)
		text = pre + text
	case "reblog":
		pre := SublteText(n.app.Config.Style, FormatUsername(notification.N.Account)+" boosted your toot") + "\n\n"
		text, controls = showTootOptions(n.app, notification.N.Status, n.showSpoiler)
		text = pre + text
	case "mention":
		pre := SublteText(n.app.Config.Style, FormatUsername(notification.N.Account)+" mentioned you") + "\n\n"
		text, controls = showTootOptions(n.app, notification.N.Status, n.showSpoiler)
		text = pre + text
	case "poll":
		pre := SublteText(n.app.Config.Style, "A poll of yours or one you participated in has ended") + "\n\n"
		text, controls = showTootOptions(n.app, notification.N.Status, n.showSpoiler)
		text = pre + text
	case "follow_request":
		text = SublteText(n.app.Config.Style, FormatUsername(notification.N.Account)+" wants to follow you. This is currently not implemented, so use another app to accept or reject the request.\n\n")
	default:
	}

	n.app.UI.StatusView.SetText(text)
	n.app.UI.StatusView.SetControls(controls)
	n.app.UI.Root.Sync()
}

func (n *NotificationsFeed) RedrawControls() {
	notification := n.GetCurrentNotification()
	if notification == nil {
		n.app.UI.StatusView.SetControls("")
		return
	}
	switch notification.N.Type {
	case "favourite", "reblog", "mention", "poll":
		_, controls := showTootOptions(n.app, notification.N.Status, n.showSpoiler)
		n.app.UI.StatusView.SetControls(controls)
	case "follow":
		_, controls := showUser(n.app, &notification.N.Account, notification.R, true)
		n.app.UI.StatusView.SetControls(controls)
	}
}

func (n *NotificationsFeed) GetSavedIndex() int {
	return n.index
}

func (n *NotificationsFeed) Input(event *tcell.EventKey) {
	notification := n.GetCurrentNotification()
	if notification == nil {
		return
	}
	if notification.N.Type == "follow" {
		controls := []ControlItem{
			ControlUser, ControlFollow, ControlBlock,
			ControlMute, ControlAvatar, ControlOpen,
			ControlYankUser,
		}
		options := inputOptions(controls)
		_, rc, _, _, rel := inputSimple(n.app, event, options, notification.N.Account, nil, nil, notification.R, n, nil)
		if rc {
			var index int
			if n.docked {
				index = n.app.UI.StatusView.notificationView.list.GetCurrentItem()
			} else {
				index = n.app.UI.StatusView.GetCurrentItem()
			}
			n.notifications[index].R = rel
			n.RedrawControls()
		}
		return
	}

	if notification.N.Type == "follow_request" {
		return
	}
	status := notification.N.Status
	originalStatus := status
	if status.Reblog != nil {
		status = status.Reblog
	}

	controls := []ControlItem{
		ControlAvatar, ControlThread, ControlUser, ControlSpoiler,
		ControlCompose, ControlOpen, ControlReply, ControlMedia,
		ControlFavorite, ControlBoost, ControlDelete, ControlBookmark,
		ControlYankStatus,
	}
	options := inputOptions(controls)

	updated, rc, rt, newS, _ := inputSimple(n.app, event, options, notification.N.Account, status, originalStatus, nil, n, nil)
	if updated {
		var index int
		if n.docked {
			index = n.app.UI.StatusView.notificationView.list.GetCurrentItem()
		} else {
			index = n.app.UI.StatusView.GetCurrentItem()
		}
		n.notifications[index].N.Status = newS
	}
	if rc {
		n.RedrawControls()
	}
	if rt {
		n.DrawToot()
	}
}

func NewTagFeed(app *App, tag string) *TagFeed {
	t := &TagFeed{
		app: app,
		tag: tag,
	}
	t.statuses, _ = t.app.API.GetTags(tag)
	return t
}

type TagFeed struct {
	app         *App
	tag         string
	statuses    []*mastodon.Status
	index       int
	showSpoiler bool
}

func (t *TagFeed) FeedType() FeedType {
	return TagFeedType
}

func (t *TagFeed) GetDesc() string {
	return "Tag #" + t.tag
}

func (t *TagFeed) GetCurrentStatus() *mastodon.Status {
	index := t.app.UI.StatusView.GetCurrentItem()
	if index >= len(t.statuses) {
		return nil
	}
	return t.statuses[t.app.UI.StatusView.GetCurrentItem()]
}

func (t *TagFeed) GetCurrentUser() *mastodon.Account {
	return userFromStatus(t.GetCurrentStatus())
}

func (t *TagFeed) GetFeedList() <-chan ListItem {
	return drawStatusList(t.statuses, t.app.Config.General.DateFormat, t.app.Config.General.DateTodayFormat, t.app.Config.General.DateRelative)
}

func (t *TagFeed) LoadNewer() int {
	var statuses []*mastodon.Status
	var err error
	if len(t.statuses) == 0 {
		statuses, err = t.app.API.GetTags(t.tag)
	} else {
		statuses, err = t.app.API.GetTagsNewer(t.tag, t.statuses[0])
	}
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load new toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	old := t.statuses
	t.statuses = append(statuses, old...)
	return len(statuses)
}

func (t *TagFeed) LoadOlder() int {
	var statuses []*mastodon.Status
	var err error
	if len(t.statuses) == 0 {
		statuses, err = t.app.API.GetTags(t.tag)
	} else {
		statuses, err = t.app.API.GetTagsOlder(t.tag, t.statuses[len(t.statuses)-1])
	}
	if err != nil {
		t.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load older toots. Error: %v\n", err))
		return 0
	}
	if len(statuses) == 0 {
		return 0
	}
	t.statuses = append(t.statuses, statuses...)
	return len(statuses)
}

func (t *TagFeed) DrawList() {
	t.app.UI.StatusView.SetList(t.GetFeedList())
}

func (t *TagFeed) DrawSpoiler() {
	t.showSpoiler = true
	t.DrawToot()
}

func (t *TagFeed) DrawToot() {
	if len(t.statuses) == 0 {
		t.app.UI.StatusView.SetText("")
		t.app.UI.StatusView.SetControls("")
		return
	}
	t.index = t.app.UI.StatusView.GetCurrentItem()
	text, controls := showTootOptions(t.app, t.statuses[t.index], t.showSpoiler)
	t.showSpoiler = false
	t.app.UI.StatusView.SetText(text)
	t.app.UI.StatusView.SetControls(controls)
	t.app.UI.Root.Sync()
}

func (t *TagFeed) RedrawControls() {
	status := t.GetCurrentStatus()
	if status == nil {
		return
	}
	_, controls := showTootOptions(t.app, status, t.showSpoiler)
	t.app.UI.StatusView.SetControls(controls)
}

func (t *TagFeed) GetSavedIndex() int {
	return t.index
}

func (t *TagFeed) Input(event *tcell.EventKey) {
	status := t.GetCurrentStatus()
	originalStatus := status
	if status == nil {
		return
	}
	if status.Reblog != nil {
		status = status.Reblog
	}
	user := status.Account

	controls := []ControlItem{
		ControlAvatar, ControlThread, ControlUser, ControlSpoiler,
		ControlCompose, ControlOpen, ControlReply, ControlMedia,
		ControlFavorite, ControlBoost, ControlDelete, ControlBookmark,
		ControlYankStatus,
	}
	options := inputOptions(controls)

	updated, rc, rt, newS, _ := inputSimple(t.app, event, options, user, status, originalStatus, nil, t, nil)
	if updated {
		index := t.app.UI.StatusView.GetCurrentItem()
		t.statuses[index] = newS
	}
	if rc {
		t.RedrawControls()
	}
	if rt {
		t.DrawToot()
	}
}

func NewUserListFeed(app *App, t UserListType, s string) *UserListFeed {
	u := &UserListFeed{
		app:      app,
		listType: t,
		input:    s,
	}
	users, err := app.API.GetUserList(t, s)
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load users. Error: %v\n", err))
		return u
	}
	u.users = users
	return u
}

type UserListFeed struct {
	app      *App
	users    []*UserData
	index    int
	input    string
	listType UserListType
}

func (u *UserListFeed) FeedType() FeedType {
	return UserListFeedType
}

func (u *UserListFeed) GetDesc() string {
	var output string
	switch u.listType {
	case UserListSearch:
		output = "User search: " + u.input
	case UserListBoosts:
		output = "Boosts"
	case UserListFavorites:
		output = "Favorites"
	case UserListFollowers:
		output = "Followers"
	case UserListFollowing:
		output = "Following"
	case UserListBlocking:
		output = "Blocking"
	case UserListMuting:
		output = "Muting"
	}
	return output
}

func (u *UserListFeed) GetCurrentStatus() *mastodon.Status {
	return nil
}

func (u *UserListFeed) GetCurrentUser() *mastodon.Account {
	ud := u.GetCurrentUserData()
	if ud == nil {
		return nil
	}
	return ud.User
}

func (u *UserListFeed) GetCurrentUserData() *UserData {
	index := u.app.UI.app.UI.StatusView.GetCurrentItem()
	if len(u.users) == 0 || index > len(u.users)-1 {
		return nil
	}
	return u.users[index-1]
}

func (u *UserListFeed) GetFeedList() <-chan ListItem {
	ch := make(chan ListItem)
	users := u.users
	go func() {
		for _, user := range users {
			var username string
			if user.User.DisplayName == "" {
				username = user.User.Acct
			} else {
				username = fmt.Sprintf("%s (%s)", user.User.DisplayName, user.User.Acct)
			}
			ch <- ListItem{Text: username, Icons: ""}
		}
		close(ch)
	}()
	return ch
}

func (u *UserListFeed) LoadNewer() int {
	var users []*UserData
	var err error
	if len(u.users) == 0 {
		users, err = u.app.API.GetUserList(u.listType, u.input)
	} else {
		users, err = u.app.API.GetUserListNewer(u.listType, u.input, u.users[0].User)
	}
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load new users. Error: %v\n", err))
		return 0
	}
	if len(users) == 0 {
		return 0
	}
	old := u.users
	u.users = append(users, old...)
	return len(users)
}

func (u *UserListFeed) LoadOlder() int {
	var users []*UserData
	var err error
	if len(u.users) == 0 {
		users, err = u.app.API.GetUserList(u.listType, u.input)
	} else {
		users, err = u.app.API.GetUserListOlder(u.listType, u.input, u.users[len(u.users)-1].User)
	}
	if err != nil {
		u.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load more users. Error: %v\n", err))
		return 0
	}
	if len(users) == 0 {
		return 0
	}
	u.users = append(u.users, users...)
	return len(users)
}

func (u *UserListFeed) DrawList() {
	u.app.UI.StatusView.SetList(u.GetFeedList())
}

func (u *UserListFeed) RedrawControls() {
	//Does not implement
}

func (u *UserListFeed) DrawSpoiler() {
	//Does not implement
}

func (u *UserListFeed) DrawToot() {
	u.index = u.app.UI.StatusView.GetCurrentItem()
	index := u.index
	if index > len(u.users)-1 || len(u.users) == 0 {
		return
	}
	user := u.users[index]

	text, controls := showUser(u.app, user.User, user.Relationship, true)

	u.app.UI.StatusView.SetText(text)
	u.app.UI.StatusView.SetControls(controls)
	u.app.UI.Root.Sync()
}

func (u *UserListFeed) GetSavedIndex() int {
	return u.index
}

func (u *UserListFeed) Input(event *tcell.EventKey) {
	index := u.GetSavedIndex()
	if index > len(u.users)-1 || len(u.users) == 0 {
		return
	}
	user := u.users[index]

	controls := []ControlItem{
		ControlAvatar, ControlFollow, ControlBlock, ControlMute, ControlOpen,
		ControlUser, ControlEnter, ControlYankUser,
	}
	options := inputOptions(controls)

	updated, _, _, _, newRel := inputSimple(u.app, event, options, *user.User, nil, nil, user.Relationship, u, nil)
	if updated {
		u.users[index].Relationship = newRel
		u.DrawToot()
	}
}

func NewListFeed(app *App) *ListFeed {
	l := &ListFeed{
		app: app,
	}
	lists, err := app.API.GetLists()
	if err != nil {
		l.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load lists. Error: %v\n", err))
		return l
	}
	l.lists = lists
	return l
}

type ListInfo struct {
	name string
	id   mastodon.ID
}

type ListFeed struct {
	app   *App
	lists []*mastodon.List
	index int
}

func (l *ListFeed) FeedType() FeedType {
	return ListFeedType
}

func (l *ListFeed) GetDesc() string {
	return "Lists"
}

func (l *ListFeed) GetCurrentStatus() *mastodon.Status {
	return nil
}

func (l *ListFeed) GetCurrentUser() *mastodon.Account {
	return nil
}

func (l *ListFeed) GetFeedList() <-chan ListItem {
	ch := make(chan ListItem)
	go func() {
		for _, list := range l.lists {
			ch <- ListItem{Text: list.Title, Icons: ""}
		}
		close(ch)
	}()
	return ch
}

func (l *ListFeed) LoadNewer() int {
	return 0
}

func (l *ListFeed) LoadOlder() int {
	return 0
}

func (l *ListFeed) DrawList() {
	l.app.UI.StatusView.SetList(l.GetFeedList())
}

func (l *ListFeed) RedrawControls() {
	//Does not implement
}

func (l *ListFeed) DrawSpoiler() {
	//Does not implement
}

func (l *ListFeed) DrawToot() {
	l.index = l.app.UI.StatusView.GetCurrentItem()
	index := l.index
	if index > len(l.lists)-1 || len(l.lists) == 0 {
		return
	}
	list := l.lists[index]

	text := ColorKey(l.app.Config, "", "O", "pen")
	text += fmt.Sprintf(" list %s", list.Title)

	l.app.UI.StatusView.SetText(text)
	l.app.UI.StatusView.SetControls("")
	l.app.UI.Root.Sync()
}

func (l *ListFeed) GetSavedIndex() int {
	return l.index
}

func (l *ListFeed) Input(event *tcell.EventKey) {
	index := l.GetSavedIndex()
	if index > len(l.lists)-1 || len(l.lists) == 0 {
		return
	}
	list := l.lists[index]
	li := ListInfo{
		name: list.Title,
		id:   list.ID,
	}

	controls := []ControlItem{ControlEnter, ControlList}
	options := inputOptions(controls)

	inputSimple(l.app, event, options, mastodon.Account{}, nil, nil, nil, nil, &li)
}
