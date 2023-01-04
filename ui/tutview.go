package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	SubFocus      SubFocusAt
	Leader        *Leader
	Shared        *Shared
	View          *tview.Pages

	LoginView      *LoginView
	MainView       *MainView
	LinkView       *LinkView
	ComposeView    *ComposeView
	VoteView       *VoteView
	PollView       *PollView
	PreferenceView *PreferenceView
	HelpView       *HelpView
	EditorView     *EditorView
	ModalView      *ModalView
}

func (tv *TutView) CleanExit(code int) {
	os.Exit(code)
}

func NewLeader(tv *TutView) *Leader {
	return &Leader{
		tv: tv,
	}
}

type Leader struct {
	tv        *TutView
	timeStart time.Time
	content   string
}

func (l *Leader) IsActive() bool {
	td := time.Duration(l.tv.tut.Config.General.LeaderTimeout)
	return time.Since(l.timeStart) < td*time.Millisecond
}

func (l *Leader) Reset() {
	l.timeStart = time.Now()
	l.content = ""
}

func (l *Leader) ResetInactive() {
	l.timeStart = time.Now().Add(-1 * time.Hour)
	l.content = ""
}

func (l *Leader) AddRune(r rune) {
	l.content += string(r)
}

func (l *Leader) Content() string {
	return l.content
}

func NewTutView(t *Tut, accs *auth.AccountData, selectedUser string) *TutView {
	tv := &TutView{
		tut:  t,
		View: tview.NewPages(),
	}
	tv.Leader = NewLeader(tv)
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
				break
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
	if tv.tut.Config.General.ShowHelp {
		tv.Shared.Bottom.Cmd.ShowMsg("Press ? or :help to learn how tut functions")
	}
	client := mastodon.NewClient(conf)
	me, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		fmt.Printf("Couldn't login. Error %s\n", err)
		tv.tut.App.Stop()
		tv.CleanExit(1)
	}
	ac := &api.AccountClient{
		Me:       me,
		Client:   client,
		Streams:  make(map[string]*api.Stream),
		WSClient: client.NewWSClient(),
	}
	inst, err := ac.Client.GetInstanceV2(context.Background())
	if err != nil {
		inst, err := ac.Client.GetInstance(context.Background())
		if err != nil {
			fmt.Printf("Couldn't get instance. Error %s\n", err)
			tv.tut.App.Stop()
			tv.CleanExit(1)
		}
		ac.InstanceOld = inst
	} else {
		ac.Instance = inst
	}
	tv.tut.Client = ac

	update := make(chan bool, 1)
	tv.SubFocus = ListFocus
	tv.LinkView = NewLinkView(tv)
	tv.Timeline = NewTimeline(tv, update)
	tv.MainView = NewMainView(tv, update)
	tv.ComposeView = NewComposeView(tv)
	tv.VoteView = NewVoteView(tv)
	tv.PollView = NewPollView(tv)
	tv.PreferenceView = NewPreferenceView(tv)
	tv.HelpView = NewHelpView(tv)
	tv.EditorView = NewEditorView(tv)
	tv.ModalView = NewModalView(tv)

	tv.View.AddPage("main", tv.MainView.View, true, false)
	tv.View.AddPage("link", tv.LinkView.View, true, false)
	tv.View.AddPage("compose", tv.ComposeView.View, true, false)
	tv.View.AddPage("vote", tv.VoteView.View, true, false)
	tv.View.AddPage("help", tv.HelpView.View, true, false)
	tv.View.AddPage("editor", tv.EditorView.View, true, false)
	tv.View.AddPage("poll", tv.PollView.View, true, false)
	tv.View.AddPage("preference", tv.PreferenceView.View, true, false)
	tv.View.AddPage("modal", tv.ModalView.View, true, false)
	tv.SetPage(MainFocus)
}

func (tv *TutView) FocusFeed(index int, ct *config.Timeline) {
	if index < 0 || index >= len(tv.Timeline.Feeds) {
		return
	}
	tv.Timeline.FeedFocusIndex = index
	for i := 0; i < len(tv.Timeline.Feeds); i++ {
		if i == index {
			for _, f := range tv.Timeline.Feeds[i].Feeds {
				f.ListInFocus()
			}
		} else {
			for _, f := range tv.Timeline.Feeds[i].Feeds {
				f.ListOutFocus()
			}
		}
	}
	for i, tl := range tv.Timeline.Feeds[index].Feeds {
		if ct == tl.Timeline {
			tv.Timeline.Feeds[index].FeedIndex = i
			break
		}
	}
	tv.Shared.Top.SetText(tv.Timeline.GetTitle())
	tv.Timeline.update <- true
}

func (tv *TutView) NextFeed() {
	index := tv.Timeline.FeedFocusIndex + 1
	if index >= len(tv.Timeline.Feeds) {
		index = 0
	}
	tv.FocusFeed(index, nil)
}

func (tv *TutView) PrevFeed() {
	index := tv.Timeline.FeedFocusIndex - 1
	if index < 0 {
		index = len(tv.Timeline.Feeds) - 1
	}
	tv.FocusFeed(index, nil)
}
