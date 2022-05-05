package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/auth"
	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type TimelineFocusAt uint

const (
	FeedFocus TimelineFocusAt = iota
	NotificationFocus
)

type SubFocusAt uint

const (
	ListFocus SubFocusAt = iota
	ContentFocus
)

type Tut struct {
	Client *api.AccountClient
	App    *tview.Application
	Config *config.Config
}

type TutView struct {
	tut           *Tut
	Timeline      *Timeline
	PageFocus     PageFocusAt
	PrevPageFocus PageFocusAt
	TimelineFocus TimelineFocusAt
	SubFocus      SubFocusAt
	Shared        *Shared
	View          *tview.Pages

	LoginView   *LoginView
	MainView    *MainView
	LinkView    *LinkView
	ComposeView *ComposeView
	VoteView    *VoteView
	HelpView    *HelpView
	ModalView   *ModalView

	FileList []string
}

func NewTutView(t *Tut, accs *auth.AccountData, selectedUser string) *TutView {
	tv := &TutView{
		tut:      t,
		View:     tview.NewPages(),
		FileList: []string{},
	}
	tv.Shared = NewShared(tv)
	if selectedUser != "" {
		useHost := false
		found := false
		if strings.Contains(selectedUser, "@") {
			useHost = true
		}
		for _, acc := range accs.Accounts {
			accName := acc.Name
			if useHost {
				host := strings.TrimPrefix(acc.Server, "https://")
				host = strings.TrimPrefix(host, "http://")
				accName += "@" + host
			}
			if accName == selectedUser {
				tv.loggedIn(acc)
				found = true
			}
		}
		if !found {
			log.Fatalf("Couldn't find a user named %s. Try again", selectedUser)
		}
	} else if len(accs.Accounts) > 1 {
		tv.LoginView = NewLoginView(tv, accs)
		tv.View.AddPage("login", tv.LoginView.View, true, true)
		tv.SetPage(LoginFocus)
	} else {
		tv.loggedIn(accs.Accounts[0])
	}
	return tv
}

func (tv *TutView) loggedIn(acc auth.Account) {
	conf := &mastodon.Config{
		Server:       acc.Server,
		ClientID:     acc.ClientID,
		ClientSecret: acc.ClientSecret,
		AccessToken:  acc.AccessToken,
	}
	client := mastodon.NewClient(conf)
	me, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		fmt.Printf("Couldn't login. Error %s\n", err)
		os.Exit(1)
	}
	ac := &api.AccountClient{
		Me:      me,
		Client:  client,
		Streams: make(map[string]*api.Stream),
	}
	tv.tut.Client = ac

	update := make(chan bool, 1)
	tv.TimelineFocus = FeedFocus
	tv.SubFocus = ListFocus
	tv.LinkView = NewLinkView(tv)
	tv.Timeline = NewTimeline(tv, update)
	tv.MainView = NewMainView(tv, update)
	tv.ComposeView = NewComposeView(tv)
	tv.VoteView = NewVoteView(tv)
	tv.HelpView = NewHelpView(tv)
	tv.ModalView = NewModalView(tv)

	tv.View.AddPage("main", tv.MainView.View, true, false)
	tv.View.AddPage("link", tv.LinkView.View, true, false)
	tv.View.AddPage("compose", tv.ComposeView.View, true, false)
	tv.View.AddPage("vote", tv.VoteView.View, true, false)
	tv.View.AddPage("help", tv.HelpView.View, true, false)
	tv.View.AddPage("modal", tv.ModalView.View, true, false)
	tv.SetPage(MainFocus)
}

func (tv *TutView) FocusNotification() {
	tv.TimelineFocus = NotificationFocus
	for _, f := range tv.Timeline.Feeds {
		f.ListOutFocus()
	}
	tv.Timeline.Notifications.ListInFocus()
	tv.Timeline.update <- true
}

func (tv *TutView) FocusFeed() {
	tv.TimelineFocus = FeedFocus
	for _, f := range tv.Timeline.Feeds {
		f.ListInFocus()
	}
	tv.Timeline.Notifications.ListOutFocus()
	tv.Timeline.update <- true
}
