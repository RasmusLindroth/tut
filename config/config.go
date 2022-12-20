package config

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gdamore/tcell/v2"
	"github.com/gobwas/glob"
	"golang.org/x/exp/slices"
	"gopkg.in/ini.v1"
)

//go:embed toot.tmpl
var tootTemplate string

//go:embed user.tmpl
var userTemplate string

//go:embed help.tmpl
var helpTemplate string

//go:embed themes/*
var themesFS embed.FS

type Config struct {
	General            General
	Style              Style
	Media              Media
	OpenPattern        OpenPattern
	OpenCustom         OpenCustom
	NotificationConfig Notification
	Templates          Templates
	Input              Input
}

type LeaderAction struct {
	Command   LeaderCommand
	Subaction string
	Shortcut  string
}

type LeaderCommand uint

const (
	LeaderNone LeaderCommand = iota
	LeaderHome
	LeaderDirect
	LeaderLocal
	LeaderFederated
	LeaderSpecialAll
	LeaderSpecialBoosts
	LeaderSpecialReplies
	LeaderClearNotifications
	LeaderCompose
	LeaderEdit
	LeaderBlocking
	LeaderBookmarks
	LeaderSaved
	LeaderFavorited
	LeaderBoosts
	LeaderFavorites
	LeaderFollowing
	LeaderFollowers
	LeaderListPlacement
	LeaderListSplit
	LeaderMuting
	LeaderPreferences
	LeaderProfile
	LeaderProportions
	LeaderNotifications
	LeaderMentions
	LeaderLists
	LeaderRefetch
	LeaderTag
	LeaderTags
	LeaderStickToTop
	LeaderHistory
	LeaderUser
	LeaderLoadNewer
	LeaderWindow
	LeaderCloseWindow
	LeaderSwitch
	LeaderMoveWindowLeft
	LeaderMoveWindowRight
	LeaderMoveWindowHome
	LeaderMoveWindowEnd
)

type FeedType uint

const (
	Favorites FeedType = iota
	Favorited
	Boosts
	Followers
	Following
	FollowRequests
	Blocking
	Muting
	History
	InvalidFeed
	Notifications
	Saved
	Tag
	Tags
	Thread
	TimelineFederated
	TimelineHome
	TimelineHomeSpecial
	TimelineLocal
	Mentions
	Conversations
	User
	UserList
	Lists
	List
	ListUsersIn
	ListUsersAdd
)

type NotificationToHide string

const (
	HideMention       NotificationToHide = "mention"
	HideStatus        NotificationToHide = "status"
	HideBoost         NotificationToHide = "reblog"
	HideFollow        NotificationToHide = "follow"
	HideFollowRequest NotificationToHide = "follow_request"
	HideFavorite      NotificationToHide = "favourite"
	HidePoll          NotificationToHide = "poll"
	HideEdited        NotificationToHide = "update"
)

type Timeline struct {
	FeedType    FeedType
	Subaction   string
	Name        string
	Key         Key
	ShowBoosts  bool
	ShowReplies bool
}

type General struct {
	Confirmation        bool
	MouseSupport        bool
	DateTodayFormat     string
	DateFormat          string
	DateRelative        int
	MaxWidth            int
	NotificationFeed    bool
	QuoteReply          bool
	CharLimit           int
	ShortHints          bool
	ShowFilterPhrase    bool
	ListPlacement       ListPlacement
	ListSplit           ListSplit
	ListProportion      int
	ContentProportion   int
	TerminalTitle       int
	ShowIcons           bool
	ShowHelp            bool
	RedrawUI            bool
	LeaderKey           rune
	LeaderTimeout       int64
	LeaderActions       []LeaderAction
	TimelineName        bool
	Timelines           []Timeline
	StickToTop          bool
	NotificationsToHide []NotificationToHide
	ShowBoostedUser     bool
}

type Style struct {
	Theme string

	Background tcell.Color
	Text       tcell.Color

	Subtle      tcell.Color
	WarningText tcell.Color

	TextSpecial1 tcell.Color
	TextSpecial2 tcell.Color

	TopBarBackground tcell.Color
	TopBarText       tcell.Color

	StatusBarBackground tcell.Color
	StatusBarText       tcell.Color

	StatusBarViewBackground tcell.Color
	StatusBarViewText       tcell.Color

	ListSelectedBackground tcell.Color
	ListSelectedText       tcell.Color

	ListSelectedInactiveBackground tcell.Color
	ListSelectedInactiveText       tcell.Color

	ControlsText      tcell.Color
	ControlsHighlight tcell.Color

	AutocompleteBackground tcell.Color
	AutocompleteText       tcell.Color

	AutocompleteSelectedBackground tcell.Color
	AutocompleteSelectedText       tcell.Color

	ButtonColorOne tcell.Color
	ButtonColorTwo tcell.Color

	TimelineNameBackground tcell.Color
	TimelineNameText       tcell.Color

	IconColor tcell.Color

	CommandText tcell.Color
}

type Media struct {
	ImageViewer   string
	ImageArgs     []string
	ImageTerminal bool
	ImageSingle   bool
	ImageReverse  bool
	VideoViewer   string
	VideoArgs     []string
	VideoTerminal bool
	VideoSingle   bool
	VideoReverse  bool
	AudioViewer   string
	AudioArgs     []string
	AudioTerminal bool
	AudioSingle   bool
	AudioReverse  bool
	LinkViewer    string
	LinkArgs      []string
	LinkTerminal  bool
}

type Pattern struct {
	Pattern  string
	Open     string
	Compiled glob.Glob
	Program  string
	Args     []string
	Terminal bool
}

type OpenPattern struct {
	Patterns []Pattern
}

type Custom struct {
	Index    int
	Name     string
	Program  string
	Args     []string
	Terminal bool
}
type OpenCustom struct {
	OpenCustoms []Custom
}

type ListPlacement uint

const (
	ListPlacementTop ListPlacement = iota
	ListPlacementBottom
	ListPlacementLeft
	ListPlacementRight
)

type ListSplit uint

const (
	ListRow ListSplit = iota
	ListColumn
)

type NotificationType uint

const (
	NotificationFollower NotificationType = iota
	NotificationFavorite
	NotificationMention
	NotificationUpdate
	NotificationBoost
	NotificationPoll
	NotificationPost
)

type Notification struct {
	NotificationFollower bool
	NotificationFavorite bool
	NotificationMention  bool
	NotificationUpdate   bool
	NotificationBoost    bool
	NotificationPoll     bool
	NotificationPost     bool
}

type Templates struct {
	Toot *template.Template
	User *template.Template
	Help *template.Template
}

var keyMatch = regexp.MustCompile("^\"(.*?)\\[(.*?)\\](.*?)\"$")

func newHint(s string) []string {
	matches := keyMatch.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return []string{"", "", ""}
	}
	if len(matches[0]) != 4 {
		return []string{"", "", ""}
	}
	return []string{matches[0][1], matches[0][2], matches[0][3]}
}

func NewKey(s []string, double bool) (Key, error) {
	var k Key
	if len(s) < 2 {
		return k, errors.New("key must have a minimum length of 2")
	}
	var start int
	if double {
		start = 1
		k = Key{
			Hint: [][]string{newHint(s[0]), newHint(s[1])},
		}
	} else {
		start = 0
		k = Key{
			Hint: [][]string{newHint(s[0])},
		}
	}
	var runes []rune
	var keys []tcell.Key
	for _, v := range s[start+1:] {
		value := []rune(strings.TrimSpace(v))
		if len(value) < 3 {
			return k, errors.New("key value must have a minimum length of 3")
		}
		if value[0] == '\'' {
			if len(value) != 3 {
				return k, fmt.Errorf("rune %s must only contain one char", string(value))
			}
			runes = append(runes, value[1])
		} else if value[0] == '"' {
			if value[len(value)-1] != '"' {
				return k, fmt.Errorf("key %s must end with \"", string(value))
			}
			keyName := string(value[1 : len(value)-1])
			found := false
			var fk tcell.Key
			for tk, tv := range tcell.KeyNames {
				if tv == keyName {
					found = true
					fk = tk
					break
				}
			}
			if found {
				keys = append(keys, fk)
			} else {
				return k, fmt.Errorf("no key named %s", keyName)
			}
		} else {
			return k, fmt.Errorf("input %s is in the wrong format", string(value))
		}
	}
	k.Runes = runes
	k.Keys = keys

	return k, nil
}

type Key struct {
	Hint  [][]string
	Runes []rune
	Keys  []tcell.Key
}

func (k Key) Match(kb tcell.Key, rb rune) bool {
	for _, ka := range k.Keys {
		if ka == kb {
			return true
		}
	}
	for _, ra := range k.Runes {
		if ra == rb {
			return true
		}
	}
	return false
}

type Input struct {
	GlobalDown  Key
	GlobalUp    Key
	GlobalEnter Key
	GlobalBack  Key
	GlobalExit  Key

	MainHome       Key
	MainEnd        Key
	MainPrevFeed   Key
	MainNextFeed   Key
	MainPrevWindow Key
	MainNextWindow Key
	MainCompose    Key

	StatusAvatar       Key
	StatusBoost        Key
	StatusDelete       Key
	StatusEdit         Key
	StatusFavorite     Key
	StatusMedia        Key
	StatusLinks        Key
	StatusPoll         Key
	StatusReply        Key
	StatusBookmark     Key
	StatusThread       Key
	StatusUser         Key
	StatusViewFocus    Key
	StatusYank         Key
	StatusToggleCW     Key
	StatusShowFiltered Key

	UserAvatar              Key
	UserBlock               Key
	UserFollow              Key
	UserFollowRequestDecide Key
	UserMute                Key
	UserLinks               Key
	UserUser                Key
	UserViewFocus           Key
	UserYank                Key

	ListOpenFeed   Key
	ListUserList   Key
	ListUserAdd    Key
	ListUserDelete Key

	TagOpenFeed Key
	TagFollow   Key

	LinkOpen Key
	LinkYank Key

	ComposeEditCW               Key
	ComposeEditText             Key
	ComposeIncludeQuote         Key
	ComposeMediaFocus           Key
	ComposePost                 Key
	ComposeToggleContentWarning Key
	ComposeVisibility           Key
	ComposeLanguage             Key
	ComposePoll                 Key

	MediaDelete   Key
	MediaEditDesc Key
	MediaAdd      Key

	VoteVote   Key
	VoteSelect Key

	PollAdd         Key
	PollEdit        Key
	PollDelete      Key
	PollMultiToggle Key
	PollExpiration  Key

	PreferenceName         Key
	PreferenceVisibility   Key
	PreferenceBio          Key
	PreferenceSave         Key
	PreferenceFields       Key
	PreferenceFieldsAdd    Key
	PreferenceFieldsEdit   Key
	PreferenceFieldsDelete Key
}

func parseColor(input string, def string, xrdb map[string]string) tcell.Color {
	if input == "" {
		return tcell.GetColor(def)
	}

	if strings.HasPrefix(input, "xrdb:") {
		key := strings.TrimPrefix(input, "xrdb:")
		if c, ok := xrdb[key]; ok {
			return tcell.GetColor(c)
		} else {
			return tcell.GetColor(def)
		}
	}
	return tcell.GetColor(input)
}

func parseStyle(cfg *ini.File, cnfPath string, cnfDir string) Style {
	var xrdbColors map[string]string
	xrdbMap, _ := GetXrdbColors()
	prefix := cfg.Section("style").Key("xrdb-prefix").String()
	if prefix == "" {
		prefix = "guess"
	}

	if prefix == "guess" {
		if m, ok := xrdbMap["*"]; ok {
			xrdbColors = m
		} else if m, ok := xrdbMap["URxvt"]; ok {
			xrdbColors = m
		} else if m, ok := xrdbMap["XTerm"]; ok {
			xrdbColors = m
		}
	} else {
		if m, ok := xrdbMap[prefix]; ok {
			xrdbColors = m
		}
	}

	style := Style{}
	theme := cfg.Section("style").Key("theme").String()
	if theme != "none" && theme != "" {
		bundled, local, err := getThemes(cnfPath, cnfDir)
		if err != nil {
			log.Fatalf("Couldn't load themes. Error: %s\n", err)
		}
		found := false
		isLocal := false
		for _, t := range local {
			if filepath.Base(t) == fmt.Sprintf("%s.ini", theme) {
				found = true
				isLocal = true
				break
			}
		}
		if !found {
			for _, t := range bundled {
				if filepath.Base(t) == fmt.Sprintf("%s.ini", theme) {
					found = true
					break
				}
			}
		}
		if !found {
			log.Fatalf("Couldn't find theme %s\n", theme)
		}
		tcfg, err := getTheme(theme, isLocal, cnfDir)
		if err != nil {
			log.Fatalf("Couldn't load theme. Error: %s\n", err)
		}
		s := tcfg.Section("").Key("background").String()
		style.Background = parseColor(s, "default", xrdbColors)

		s = tcfg.Section("").Key("text").String()
		style.Text = tcell.GetColor(s)

		s = tcfg.Section("").Key("subtle").String()
		style.Subtle = tcell.GetColor(s)

		s = tcfg.Section("").Key("warning-text").String()
		style.WarningText = tcell.GetColor(s)

		s = tcfg.Section("").Key("text-special-one").String()
		style.TextSpecial1 = tcell.GetColor(s)

		s = tcfg.Section("").Key("text-special-two").String()
		style.TextSpecial2 = tcell.GetColor(s)

		s = tcfg.Section("").Key("top-bar-background").String()
		style.TopBarBackground = tcell.GetColor(s)

		s = tcfg.Section("").Key("top-bar-text").String()
		style.TopBarText = tcell.GetColor(s)

		s = tcfg.Section("").Key("status-bar-background").String()
		style.StatusBarBackground = tcell.GetColor(s)

		s = tcfg.Section("").Key("status-bar-text").String()
		style.StatusBarText = tcell.GetColor(s)

		s = tcfg.Section("").Key("status-bar-view-background").String()
		style.StatusBarViewBackground = tcell.GetColor(s)

		s = tcfg.Section("").Key("status-bar-view-text").String()
		style.StatusBarViewText = tcell.GetColor(s)

		s = tcfg.Section("").Key("list-selected-background").String()
		style.ListSelectedBackground = tcell.GetColor(s)

		s = tcfg.Section("").Key("list-selected-text").String()
		style.ListSelectedText = tcell.GetColor(s)

		s = tcfg.Section("").Key("list-selected-inactive-background").String()
		if len(s) > 0 {
			style.ListSelectedInactiveBackground = tcell.GetColor(s)
		} else {
			style.ListSelectedInactiveBackground = style.StatusBarViewBackground
		}
		s = tcfg.Section("").Key("list-selected-inactive-text").String()
		if len(s) > 0 {
			style.ListSelectedInactiveText = tcell.GetColor(s)
		} else {
			style.ListSelectedInactiveText = style.StatusBarViewText
		}

		s = tcfg.Section("").Key("controls-highlight").String()
		if len(s) > 0 {
			style.ControlsHighlight = tcell.GetColor(s)
		} else {
			style.ControlsHighlight = style.TextSpecial2
		}

		s = tcfg.Section("").Key("controls-text").String()
		if len(s) > 0 {
			style.ControlsText = tcell.GetColor(s)
		} else {
			style.ControlsText = style.Text
		}
		s = tcfg.Section("").Key("controls-highlight").String()
		if len(s) > 0 {
			style.ControlsHighlight = tcell.GetColor(s)
		} else {
			style.ControlsHighlight = style.TextSpecial2
		}

		s = tcfg.Section("").Key("autocomplete-background").String()
		if len(s) > 0 {
			style.AutocompleteBackground = tcell.GetColor(s)
		} else {
			style.AutocompleteBackground = style.Background
		}
		s = tcfg.Section("").Key("autocomplete-text").String()
		if len(s) > 0 {
			style.AutocompleteText = tcell.GetColor(s)
		} else {
			style.AutocompleteText = style.Text
		}
		s = tcfg.Section("").Key("autocomplete-selected-background").String()
		if len(s) > 0 {
			style.AutocompleteSelectedBackground = tcell.GetColor(s)
		} else {
			style.AutocompleteSelectedBackground = style.StatusBarViewBackground
		}
		s = tcfg.Section("").Key("autocomplete-selected-text").String()
		if len(s) > 0 {
			style.AutocompleteSelectedText = tcell.GetColor(s)
		} else {
			style.AutocompleteSelectedText = style.StatusBarViewText
		}

		s = tcfg.Section("").Key("button-color-one").String()
		if len(s) > 0 {
			style.ButtonColorOne = tcell.GetColor(s)
		} else {
			style.ButtonColorOne = style.StatusBarViewBackground
		}
		s = tcfg.Section("").Key("button-color-two").String()
		if len(s) > 0 {
			style.ButtonColorTwo = tcell.GetColor(s)
		} else {
			style.ButtonColorTwo = style.Background
		}

		s = tcfg.Section("").Key("timeline-name-background").String()
		if len(s) > 0 {
			style.TimelineNameBackground = tcell.GetColor(s)
		} else {
			style.TimelineNameBackground = style.Background
		}
		s = tcfg.Section("").Key("timeline-name-text").String()
		if len(s) > 0 {
			style.TimelineNameText = tcell.GetColor(s)
		} else {
			style.TimelineNameText = style.Subtle
		}
		s = tcfg.Section("").Key("command-text").String()
		if len(s) > 0 {
			style.CommandText = tcell.GetColor(s)
		} else {
			style.CommandText = style.StatusBarText
		}
	} else {
		s := cfg.Section("style").Key("background").String()
		style.Background = parseColor(s, "#27822", xrdbColors)

		s = cfg.Section("style").Key("text").String()
		style.Text = parseColor(s, "#f8f8f2", xrdbColors)

		s = cfg.Section("style").Key("subtle").String()
		style.Subtle = parseColor(s, "#808080", xrdbColors)

		s = cfg.Section("style").Key("warning-text").String()
		style.WarningText = parseColor(s, "#f92672", xrdbColors)

		s = cfg.Section("style").Key("text-special-one").String()
		style.TextSpecial1 = parseColor(s, "#ae81ff", xrdbColors)

		s = cfg.Section("style").Key("text-special-two").String()
		style.TextSpecial2 = parseColor(s, "#a6e22e", xrdbColors)

		s = cfg.Section("style").Key("top-bar-background").String()
		style.TopBarBackground = parseColor(s, "#f92672", xrdbColors)

		s = cfg.Section("style").Key("top-bar-text").String()
		style.TopBarText = parseColor(s, "white", xrdbColors)

		s = cfg.Section("style").Key("status-bar-background").String()
		style.StatusBarBackground = parseColor(s, "#f92672", xrdbColors)

		s = cfg.Section("style").Key("status-bar-text").String()
		style.StatusBarText = parseColor(s, "white", xrdbColors)

		s = cfg.Section("style").Key("status-bar-view-background").String()
		style.StatusBarViewBackground = parseColor(s, "#ae81ff", xrdbColors)

		s = cfg.Section("style").Key("status-bar-view-text").String()
		style.StatusBarViewText = parseColor(s, "white", xrdbColors)

		s = cfg.Section("style").Key("list-selected-background").String()
		style.ListSelectedBackground = parseColor(s, "#f92672", xrdbColors)

		s = cfg.Section("style").Key("list-selected-text").String()
		style.ListSelectedText = parseColor(s, "white", xrdbColors)

		s = cfg.Section("style").Key("list-selected-inactive-background").String()
		if len(s) > 0 {
			style.ListSelectedInactiveBackground = parseColor(s, "#ae81ff", xrdbColors)
		} else {
			style.ListSelectedInactiveBackground = style.StatusBarViewBackground
		}
		s = cfg.Section("style").Key("list-selected-inactive-text").String()
		if len(s) > 0 {
			style.ListSelectedInactiveText = parseColor(s, "#f8f8f2", xrdbColors)
		} else {
			style.ListSelectedInactiveText = style.StatusBarViewText
		}

		s = cfg.Section("style").Key("controls-text").String()
		if len(s) > 0 {
			style.ControlsText = parseColor(s, "#f8f8f2", xrdbColors)
		} else {
			style.ControlsText = style.Text
		}
		s = cfg.Section("style").Key("controls-highlight").String()
		if len(s) > 0 {
			style.ControlsHighlight = parseColor(s, "#a6e22e", xrdbColors)
		} else {
			style.ControlsHighlight = style.TextSpecial2
		}

		s = cfg.Section("style").Key("autocomplete-background").String()
		if len(s) > 0 {
			style.AutocompleteBackground = parseColor(s, "#272822", xrdbColors)
		} else {
			style.AutocompleteBackground = style.Background
		}
		s = cfg.Section("style").Key("autocomplete-text").String()
		if len(s) > 0 {
			style.AutocompleteText = parseColor(s, "#f8f8f2", xrdbColors)
		} else {
			style.AutocompleteText = style.Text
		}
		s = cfg.Section("style").Key("autocomplete-selected-background").String()
		if len(s) > 0 {
			style.AutocompleteSelectedBackground = parseColor(s, "#ae81ff", xrdbColors)
		} else {
			style.AutocompleteSelectedBackground = style.StatusBarViewBackground
		}
		s = cfg.Section("style").Key("autocomplete-selected-text").String()
		if len(s) > 0 {
			style.AutocompleteSelectedText = parseColor(s, "#f8f8f2", xrdbColors)
		} else {
			style.AutocompleteSelectedText = style.StatusBarViewText
		}

		s = cfg.Section("style").Key("button-color-one").String()
		if len(s) > 0 {
			style.ButtonColorOne = parseColor(s, "#ae81ff", xrdbColors)
		} else {
			style.ButtonColorOne = style.StatusBarViewBackground
		}
		s = cfg.Section("style").Key("button-color-two").String()
		if len(s) > 0 {
			style.ButtonColorTwo = parseColor(s, "#272822", xrdbColors)
		} else {
			style.ButtonColorTwo = style.Background
		}

		s = cfg.Section("style").Key("timeline-name-background").String()
		if len(s) > 0 {
			style.TimelineNameBackground = parseColor(s, "#272822", xrdbColors)
		} else {
			style.TimelineNameBackground = style.Background
		}
		s = cfg.Section("style").Key("timeline-name-text").String()
		if len(s) > 0 {
			style.TimelineNameText = parseColor(s, "gray", xrdbColors)
		} else {
			style.TimelineNameText = style.Subtle
		}

		s = cfg.Section("style").Key("command-text").String()
		if len(s) > 0 {
			style.CommandText = parseColor(s, "white", xrdbColors)
		} else {
			style.CommandText = style.StatusBarText
		}
	}

	return style
}

func parseGeneral(cfg *ini.File) General {
	general := General{}

	general.Confirmation = cfg.Section("general").Key("confirmation").MustBool(true)
	general.MouseSupport = cfg.Section("general").Key("mouse-support").MustBool(false)
	dateFormat := cfg.Section("general").Key("date-format").String()
	if dateFormat == "" {
		dateFormat = "2006-01-02 15:04"
	}
	general.DateFormat = dateFormat

	dateTodayFormat := cfg.Section("general").Key("date-today-format").String()
	if dateTodayFormat == "" {
		dateTodayFormat = "15:04"
	}
	general.DateTodayFormat = dateTodayFormat

	dateRelative, err := cfg.Section("general").Key("date-relative").Int()
	if err != nil {
		dateRelative = -1
	}
	general.DateRelative = dateRelative

	general.NotificationFeed = cfg.Section("general").Key("notification-feed").MustBool(true)
	general.QuoteReply = cfg.Section("general").Key("quote-reply").MustBool(false)
	general.CharLimit = cfg.Section("general").Key("char-limit").MustInt(500)
	general.MaxWidth = cfg.Section("general").Key("max-width").MustInt(0)
	general.ShortHints = cfg.Section("general").Key("short-hints").MustBool(false)
	general.ShowFilterPhrase = cfg.Section("general").Key("show-filter-phrase").MustBool(true)
	general.ShowIcons = cfg.Section("general").Key("show-icons").MustBool(true)
	general.ShowHelp = cfg.Section("general").Key("show-help").MustBool(true)
	general.RedrawUI = cfg.Section("general").Key("redraw-ui").MustBool(true)
	general.StickToTop = cfg.Section("general").Key("stick-to-top").MustBool(false)
	general.ShowBoostedUser = cfg.Section("general").Key("show-boosted-user").MustBool(false)

	lp := cfg.Section("general").Key("list-placement").In("left", []string{"left", "right", "top", "bottom"})
	switch lp {
	case "left":
		general.ListPlacement = ListPlacementLeft
	case "right":
		general.ListPlacement = ListPlacementRight
	case "top":
		general.ListPlacement = ListPlacementTop
	case "bottom":
		general.ListPlacement = ListPlacementBottom
	}
	ls := cfg.Section("general").Key("list-split").In("row", []string{"row", "column"})
	switch ls {
	case "row":
		general.ListSplit = ListRow
	case "column":
		general.ListSplit = ListColumn
	}

	listProp := cfg.Section("general").Key("list-proportion").MustInt(1)
	if listProp < 1 {
		listProp = 1
	}
	contentProp := cfg.Section("general").Key("content-proportion").MustInt(2)
	if contentProp < 1 {
		contentProp = 1
	}
	general.ListProportion = listProp
	general.ContentProportion = contentProp

	leaderString := cfg.Section("general").Key("leader-key").MustString("")
	leaderRunes := []rune(leaderString)
	if len(leaderRunes) > 1 {
		leaderRunes = []rune(strings.TrimSpace(leaderString))
	}
	if len(leaderRunes) > 1 {
		fmt.Println("error parsing leader-key. Error: leader-key can only be one char long")
		os.Exit(1)
	}
	if len(leaderRunes) == 1 {
		general.LeaderKey = leaderRunes[0]
	}
	if general.LeaderKey != rune(0) {
		general.LeaderTimeout = cfg.Section("general").Key("leader-timeout").MustInt64(1000)
		lactions := cfg.Section("general").Key("leader-action").ValueWithShadows()
		var las []LeaderAction
		for _, l := range lactions {
			parts := strings.Split(l, ",")
			if len(parts) < 2 {
				fmt.Printf("leader-action must consist of atleast two parts separated by a comma. Your value is: %s\n", strings.Join(parts, ","))
				os.Exit(1)
			}
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			cmd := parts[0]
			var subaction string
			if strings.Contains(parts[0], " ") {
				p := strings.Split(cmd, " ")
				cmd = p[0]
				subaction = strings.Join(p[1:], " ")
			}
			la := LeaderAction{}
			switch cmd {
			case "home":
				la.Command = LeaderHome
			case "direct":
				la.Command = LeaderDirect
			case "local":
				la.Command = LeaderLocal
			case "federated":
				la.Command = LeaderFederated
			case "special-all":
				la.Command = LeaderSpecialAll
			case "special-boosts":
				la.Command = LeaderSpecialBoosts
			case "special-replies":
				la.Command = LeaderSpecialReplies
			case "clear-notifications":
				la.Command = LeaderClearNotifications
			case "compose":
				la.Command = LeaderCompose
			case "edit":
				la.Command = LeaderEdit
			case "blocking":
				la.Command = LeaderBlocking
			case "bookmarks":
				la.Command = LeaderBookmarks
			case "saved":
				la.Command = LeaderSaved
			case "favorited":
				la.Command = LeaderFavorited
			case "history":
				la.Command = LeaderHistory
			case "boosts":
				la.Command = LeaderBoosts
			case "favorites":
				la.Command = LeaderFavorites
			case "following":
				la.Command = LeaderFollowing
			case "followers":
				la.Command = LeaderFollowers
			case "muting":
				la.Command = LeaderMuting
			case "preferences":
				la.Command = LeaderPreferences
			case "profile":
				la.Command = LeaderProfile
			case "notifications":
				la.Command = LeaderNotifications
			case "mentions":
				la.Command = LeaderMentions
			case "lists":
				la.Command = LeaderLists
			case "stick-to-top":
				la.Command = LeaderStickToTop
			case "refetch":
				la.Command = LeaderRefetch
			case "tag":
				la.Command = LeaderTag
				la.Subaction = subaction
			case "tags":
				la.Command = LeaderTags
			case "list-placement":
				la.Command = LeaderListPlacement
				la.Subaction = subaction
			case "list-split":
				la.Command = LeaderListSplit
				la.Subaction = subaction
			case "proportions":
				la.Command = LeaderProportions
				la.Subaction = subaction
			case "window":
				la.Command = LeaderWindow
				la.Subaction = subaction
			case "close-window":
				la.Command = LeaderCloseWindow
			case "move-window-left", "move-window-up":
				la.Command = LeaderMoveWindowLeft
			case "move-window-right", "move-window-down":
				la.Command = LeaderMoveWindowRight
			case "move-window-home":
				la.Command = LeaderMoveWindowHome
			case "move-window-end":
				la.Command = LeaderMoveWindowEnd
			case "switch":
				la.Command = LeaderSwitch
				sa := ""
				if len(parts) > 2 {
					sa = strings.Join(parts[2:], ",")
				}
				if len(sa) > 0 {
					la.Subaction = fmt.Sprintf("%s,%s", subaction, sa)
				} else {
					la.Subaction = subaction
				}
			case "newer":
				la.Command = LeaderLoadNewer
			default:
				fmt.Printf("leader-action %s is invalid\n", parts[0])
				os.Exit(1)
			}
			la.Shortcut = parts[1]
			las = append(las, la)
		}
		general.LeaderActions = las
	}

	general.TimelineName = cfg.Section("general").Key("timeline-show-name").MustBool(true)
	var tls []Timeline
	timelines := cfg.Section("general").Key("timelines").ValueWithShadows()
	for _, l := range timelines {
		parts := strings.Split(l, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		if len(parts) == 0 {
			fmt.Printf("timelines must consist of atleast one part seperated by a comma. Your value is: %s\n", strings.Join(parts, ","))
			os.Exit(1)
		}
		if len(parts) == 1 {
			parts = append(parts, "")
		}
		cmd := parts[0]
		var subaction string
		if strings.Contains(parts[0], " ") {
			p := strings.Split(cmd, " ")
			cmd = p[0]
			subaction = strings.Join(p[1:], " ")
		}
		tl := Timeline{}
		switch cmd {
		case "home":
			tl.FeedType = TimelineHome
		case "special":
			tl.FeedType = TimelineHomeSpecial
		case "direct":
			tl.FeedType = Conversations
		case "local":
			tl.FeedType = TimelineLocal
		case "federated":
			tl.FeedType = TimelineFederated
		case "bookmarks":
			tl.FeedType = Saved
		case "saved":
			tl.FeedType = Saved
		case "favorited":
			tl.FeedType = Favorited
		case "notifications":
			tl.FeedType = Notifications
		case "mentions":
			tl.FeedType = Mentions
		case "lists":
			tl.FeedType = Lists
		case "tag":
			tl.FeedType = Tag
			tl.Subaction = subaction
		default:
			fmt.Printf("timeline %s is invalid\n", parts[0])
			os.Exit(1)
		}
		tfStr := []string{"true", "false"}
		tl.Name = parts[1]
		if slices.Contains(tfStr, tl.Name) || (strings.HasPrefix(parts[1], "\"") && strings.HasSuffix(parts[1], "\"")) ||
			(strings.HasPrefix(parts[1], "'") && strings.HasSuffix(parts[1], "'")) {
			tl.Name = ""
		}
		tfs := []bool{true, true}
		stop := len(parts)
		if len(parts) > 1 {
			if len(parts) > 2 && slices.Contains(tfStr, parts[len(parts)-2]) &&
				slices.Contains(tfStr, parts[len(parts)-1]) &&
				len(parts)-2 > 0 {
				tfs[0] = parts[len(parts)-2] == "true"
				tfs[1] = parts[len(parts)-1] == "true"
				stop = len(parts) - 2
			} else if slices.Contains(tfStr, parts[len(parts)-1]) &&
				len(parts)-1 > 0 {
				tfs[0] = parts[len(parts)-1] == "true"
				stop = len(parts) - 1
			}
			if stop > 2 {
				vals := []string{""}
				start := 2
				if tl.Name == "" {
					start = 1
				}
				vals = append(vals, parts[start:stop]...)
				tl.Key = inputStrOrErr(vals, false)
			}
		}
		tl.ShowBoosts = tfs[0]
		tl.ShowReplies = tfs[1]
		tls = append(tls, tl)
	}
	if len(tls) == 0 {
		tls = append(tls,
			Timeline{
				FeedType:    TimelineHome,
				Name:        "",
				ShowBoosts:  true,
				ShowReplies: true,
			},
		)
		tls = append(tls,
			Timeline{
				FeedType:    Notifications,
				Name:        "[N]otifications",
				Key:         inputStrOrErr([]string{"", "'n'", "'N'"}, false),
				ShowBoosts:  true,
				ShowReplies: true,
			},
		)
	}
	general.Timelines = tls

	general.TerminalTitle = cfg.Section("general").Key("terminal-title").MustInt(0)
	/*
		0 = No terminal title
		1 = Show title in terminal and top bar
		2 = Only show terminal title, and no top bar
	*/
	if general.TerminalTitle < 0 || general.TerminalTitle > 2 {
		general.TerminalTitle = 0
	}

	nths := []NotificationToHide{}
	nth := cfg.Section("general").Key("notifications-to-hide").MustString("")
	parts := strings.Split(nth, ",")
	for _, p := range parts {
		s := strings.TrimSpace(p)
		switch s {
		case "mention":
			nths = append(nths, HideMention)
		case "status":
			nths = append(nths, HideStatus)
		case "boost":
			nths = append(nths, HideBoost)
		case "follow":
			nths = append(nths, HideFollow)
		case "follow_request":
			nths = append(nths, HideFollowRequest)
		case "favorite":
			nths = append(nths, HideFavorite)
		case "poll":
			nths = append(nths, HidePoll)
		case "edit":
			nths = append(nths, HideEdited)
		default:
			if len(s) > 0 {
				log.Fatalf("%s in notifications-to-hide is invalid\n", s)
				os.Exit(1)
			}
		}
		general.NotificationsToHide = nths
	}

	return general
}

func parseMedia(cfg *ini.File) Media {
	media := Media{}
	imageViewerComponents := strings.Fields(cfg.Section("media").Key("image-viewer").String())
	if len(imageViewerComponents) == 0 {
		media.ImageViewer = "xdg-open"
		media.ImageArgs = []string{}
	} else {
		media.ImageViewer = imageViewerComponents[0]
		media.ImageArgs = imageViewerComponents[1:]
	}
	media.ImageTerminal = cfg.Section("media").Key("image-terminal").MustBool(false)
	media.ImageSingle = cfg.Section("media").Key("image-single").MustBool(true)
	media.ImageReverse = cfg.Section("media").Key("image-reverse").MustBool(false)

	videoViewerComponents := strings.Fields(cfg.Section("media").Key("video-viewer").String())
	if len(videoViewerComponents) == 0 {
		media.VideoViewer = "xdg-open"
		media.VideoArgs = []string{}
	} else {
		media.VideoViewer = videoViewerComponents[0]
		media.VideoArgs = videoViewerComponents[1:]
	}
	media.VideoTerminal = cfg.Section("media").Key("video-terminal").MustBool(false)
	media.VideoSingle = cfg.Section("media").Key("video-single").MustBool(true)
	media.VideoReverse = cfg.Section("media").Key("video-reverse").MustBool(false)

	audioViewerComponents := strings.Fields(cfg.Section("media").Key("audio-viewer").String())
	if len(audioViewerComponents) == 0 {
		media.AudioViewer = "xdg-open"
		media.AudioArgs = []string{}
	} else {
		media.AudioViewer = audioViewerComponents[0]
		media.AudioArgs = audioViewerComponents[1:]
	}
	media.AudioTerminal = cfg.Section("media").Key("audio-terminal").MustBool(false)
	media.AudioSingle = cfg.Section("media").Key("audio-single").MustBool(true)
	media.AudioReverse = cfg.Section("media").Key("audio-reverse").MustBool(false)

	linkViewerComponents := strings.Fields(cfg.Section("media").Key("link-viewer").String())
	if len(linkViewerComponents) == 0 {
		media.LinkViewer = "xdg-open"
		media.LinkArgs = []string{}
	} else {
		media.LinkViewer = linkViewerComponents[0]
		media.LinkArgs = linkViewerComponents[1:]
	}
	media.LinkTerminal = cfg.Section("media").Key("link-terminal").MustBool(false)

	return media
}

func parseOpenPattern(cfg *ini.File) OpenPattern {
	om := OpenPattern{}

	keys := cfg.Section("open-pattern").KeyStrings()
	pairs := make(map[string]Pattern)
	for _, s := range keys {
		parts := strings.Split(s, "-")
		if len(parts) < 2 {
			panic(fmt.Sprintf("Invalid key %s in config. Must end in -pattern, -use or -terminal", s))
		}
		last := parts[len(parts)-1]
		if last != "pattern" && last != "use" && last != "terminal" {
			panic(fmt.Sprintf("Invalid key %s in config. Must end in -pattern, -use or -terminal", s))
		}

		name := strings.Join(parts[:len(parts)-1], "-")
		if _, ok := pairs[name]; !ok {
			pairs[name] = Pattern{}
		}
		if last == "pattern" {
			tmp := pairs[name]
			tmp.Pattern = cfg.Section("open-pattern").Key(s).MustString("")
			pairs[name] = tmp
		}
		if last == "use" {
			tmp := pairs[name]
			tmp.Open = cfg.Section("open-pattern").Key(s).MustString("")
			pairs[name] = tmp
		}
		if last == "terminal" {
			tmp := pairs[name]
			tmp.Terminal = cfg.Section("open-pattern").Key(s).MustBool(false)
			pairs[name] = tmp
		}
	}

	for key := range pairs {
		if pairs[key].Pattern == "" {
			panic(fmt.Sprintf("Invalid value for key %s in config. Can't be empty", key+"-pattern"))
		}
		if pairs[key].Open == "" {
			panic(fmt.Sprintf("Invalid value for key %s in config. Can't be empty", key+"-use"))
		}

		compiled, err := glob.Compile(pairs[key].Pattern)
		if err != nil {
			panic(fmt.Sprintf("Couldn't compile pattern for key %s in config. Error: %v", key+"-pattern", err))
		}
		tmp := pairs[key]
		tmp.Compiled = compiled
		comp := strings.Fields(tmp.Open)
		tmp.Program = comp[0]
		tmp.Args = comp[1:]
		om.Patterns = append(om.Patterns, tmp)
	}

	return om
}

func parseCustom(cfg *ini.File) OpenCustom {
	oc := OpenCustom{}

	for i := 1; i < 6; i++ {
		name := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-name", i)).MustString("")
		use := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-use", i)).MustString("")
		terminal := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-terminal", i)).MustBool(false)
		if use == "" {
			continue
		}
		comp := strings.Fields(use)
		c := Custom{}
		c.Index = i
		c.Name = name
		c.Program = comp[0]
		c.Args = comp[1:]
		c.Terminal = terminal
		oc.OpenCustoms = append(oc.OpenCustoms, c)
	}
	return oc
}

func parseNotifications(cfg *ini.File) Notification {
	nc := Notification{}
	nc.NotificationFollower = cfg.Section("desktop-notification").Key("followers").MustBool(false)
	nc.NotificationFavorite = cfg.Section("desktop-notification").Key("favorite").MustBool(false)
	nc.NotificationMention = cfg.Section("desktop-notification").Key("mention").MustBool(false)
	nc.NotificationUpdate = cfg.Section("desktop-notification").Key("update").MustBool(false)
	nc.NotificationBoost = cfg.Section("desktop-notification").Key("boost").MustBool(false)
	nc.NotificationPoll = cfg.Section("desktop-notification").Key("poll").MustBool(false)
	nc.NotificationPost = cfg.Section("desktop-notification").Key("posts").MustBool(false)
	return nc
}

func parseTemplates(cfg *ini.File, cnfPath string, cnfDir string) Templates {
	var tootTmpl *template.Template
	tootTmplPath, exists, err := checkConfig("toot.tmpl", cnfPath, cnfDir)
	if err != nil {
		log.Fatalf(
			fmt.Sprintf("Couldn't access toot.tmpl. Error: %v", err),
		)
	}
	if exists {
		tootTmpl, err = template.New("toot.tmpl").Funcs(template.FuncMap{
			"Color": ColorMark,
			"Flags": TextFlags,
		}).ParseFiles(tootTmplPath)
	}
	if !exists || err != nil {
		tootTmpl, err = template.New("toot.tmpl").Funcs(template.FuncMap{
			"Color": ColorMark,
			"Flags": TextFlags,
		}).Parse(tootTemplate)
	}
	if err != nil {
		log.Fatalf("Couldn't parse toot.tmpl. Error: %v", err)
	}
	var userTmpl *template.Template
	userTmplPath, exists, err := checkConfig("user.tmpl", cnfPath, cnfDir)
	if err != nil {
		log.Fatalf(
			fmt.Sprintf("Couldn't access user.tmpl. Error: %v", err),
		)
	}
	if exists {
		userTmpl, err = template.New("user.tmpl").Funcs(template.FuncMap{
			"Color": ColorMark,
			"Flags": TextFlags,
		}).ParseFiles(userTmplPath)
	}
	if !exists || err != nil {
		userTmpl, err = template.New("user.tmpl").Funcs(template.FuncMap{
			"Color": ColorMark,
			"Flags": TextFlags,
		}).Parse(userTemplate)
	}
	if err != nil {
		log.Fatalf("Couldn't parse user.tmpl. Error: %v", err)
	}
	var helpTmpl *template.Template
	helpTmpl, err = template.New("help.tmpl").Funcs(template.FuncMap{
		"Color": ColorMark,
		"Flags": TextFlags,
	}).Parse(helpTemplate)
	if err != nil {
		log.Fatalf("Couldn't parse help.tmpl. Error: %v", err)
	}
	return Templates{
		Toot: tootTmpl,
		User: userTmpl,
		Help: helpTmpl,
	}
}

func inputOrErr(cfg *ini.File, key string, double bool, def Key) Key {
	if !cfg.Section("input").HasKey(key) {
		return def
	}
	vals := cfg.Section("input").Key(key).Strings(",")
	k, err := NewKey(vals, double)
	if err != nil {
		fmt.Printf("error parsing config for key %s. Error: %v\n", key, err)
		os.Exit(1)
	}
	return k
}
func inputStrOrErr(vals []string, double bool) Key {
	k, err := NewKey(vals, double)
	if err != nil {
		fmt.Printf("error parsing config. Error: %v\n", err)
		os.Exit(1)
	}
	return k
}

func parseInput(cfg *ini.File) Input {
	ic := Input{
		GlobalDown:  inputStrOrErr([]string{"\"\"", "'j'", "'J'", "\"Down\""}, false),
		GlobalUp:    inputStrOrErr([]string{"\"\"", "'k'", "'k'", "\"Up\""}, false),
		GlobalEnter: inputStrOrErr([]string{"\"\"", "\"Enter\""}, false),
		GlobalBack:  inputStrOrErr([]string{"\"[Esc]\"", "\"Esc\""}, false),
		GlobalExit:  inputStrOrErr([]string{"\"[Q]uit\"", "'q'", "'Q'"}, false),

		MainHome:       inputStrOrErr([]string{"\"\"", "'g'", "\"Home\""}, false),
		MainEnd:        inputStrOrErr([]string{"\"\"", "'G'", "\"End\""}, false),
		MainPrevFeed:   inputStrOrErr([]string{"\"\"", "'h'", "'H'", "\"Left\""}, false),
		MainNextFeed:   inputStrOrErr([]string{"\"\"", "'l'", "'L'", "\"Right\""}, false),
		MainPrevWindow: inputStrOrErr([]string{"\"\"", "\"Backtab\""}, false),
		MainNextWindow: inputStrOrErr([]string{"\"\"", "\"Tab\""}, false),
		MainCompose:    inputStrOrErr([]string{"\"\"", "'c'", "'C'"}, false),

		StatusAvatar:       inputStrOrErr([]string{"\"[A]vatar\"", "'a'", "'A'"}, false),
		StatusBoost:        inputStrOrErr([]string{"\"[B]oost\"", "\"Un[B]oost\"", "'b'", "'B'"}, true),
		StatusDelete:       inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
		StatusEdit:         inputStrOrErr([]string{"\"[E]dit\"", "'e'", "'E'"}, false),
		StatusFavorite:     inputStrOrErr([]string{"\"[F]avorite\"", "\"Un[F]avorite\"", "'f'", "'F'"}, true),
		StatusMedia:        inputStrOrErr([]string{"\"[M]edia\"", "'m'", "'M'"}, false),
		StatusLinks:        inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		StatusPoll:         inputStrOrErr([]string{"\"[P]oll\"", "'p'", "'P'"}, false),
		StatusReply:        inputStrOrErr([]string{"\"[R]eply\"", "'r'", "'R'"}, false),
		StatusBookmark:     inputStrOrErr([]string{"\"[S]ave\"", "\"Un[S]ave\"", "'s'", "'S'"}, true),
		StatusThread:       inputStrOrErr([]string{"\"[T]hread\"", "'t'", "'T'"}, false),
		StatusUser:         inputStrOrErr([]string{"\"[U]ser\"", "'u'", "'U'"}, false),
		StatusViewFocus:    inputStrOrErr([]string{"\"[V]iew\"", "'v'", "'V'"}, false),
		StatusYank:         inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),
		StatusToggleCW:     inputStrOrErr([]string{"\"Press [Z] to toggle CW\"", "'z'", "'Z'"}, false),
		StatusShowFiltered: inputStrOrErr([]string{"\"Press [Z] to view filtered toot\"", "'z'", "'Z'"}, false),

		UserAvatar:              inputStrOrErr([]string{"\"[A]vatar\"", "'a'", "'A'"}, false),
		UserBlock:               inputStrOrErr([]string{"\"[B]lock\"", "\"Un[B]lock\"", "'b'", "'B'"}, true),
		UserFollow:              inputStrOrErr([]string{"\"[F]ollow\"", "\"Un[F]ollow\"", "'f'", "'F'"}, true),
		UserFollowRequestDecide: inputStrOrErr([]string{"\"Follow [R]equest\"", "\"Follow [R]equest\"", "'r'", "'R'"}, true),
		UserMute:                inputStrOrErr([]string{"\"[M]ute\"", "\"Un[M]ute\"", "'m'", "'M'"}, true),
		UserLinks:               inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		UserUser:                inputStrOrErr([]string{"\"[U]ser\"", "'u'", "'U'"}, false),
		UserViewFocus:           inputStrOrErr([]string{"\"[V]iew\"", "'v'", "'V'"}, false),
		UserYank:                inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),

		ListOpenFeed:   inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		ListUserList:   inputStrOrErr([]string{"\"[U]sers\"", "'u'", "'U'"}, false),
		ListUserAdd:    inputStrOrErr([]string{"\"[A]dd\"", "'a'", "'A'"}, false),
		ListUserDelete: inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),

		TagOpenFeed: inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		TagFollow:   inputStrOrErr([]string{"\"[F]ollow\"", "\"Un[F]ollow\"", "'f'", "'F'"}, true),

		LinkOpen: inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		LinkYank: inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),

		ComposeEditCW:               inputStrOrErr([]string{"\"[C]W Text\"", "'c'", "'C'"}, false),
		ComposeEditText:             inputStrOrErr([]string{"\"[E]dit text\"", "'e'", "'E'"}, false),
		ComposeIncludeQuote:         inputStrOrErr([]string{"\"[I]nclude quote\"", "'i'", "'I'"}, false),
		ComposeMediaFocus:           inputStrOrErr([]string{"\"[M]edia\"", "'m'", "'M'"}, false),
		ComposePost:                 inputStrOrErr([]string{"\"[P]ost\"", "'p'", "'P'"}, false),
		ComposeToggleContentWarning: inputStrOrErr([]string{"\"[T]oggle CW\"", "'t'", "'T'"}, false),
		ComposeVisibility:           inputStrOrErr([]string{"\"[V]isibility\"", "'v'", "'V'"}, false),
		ComposeLanguage:             inputStrOrErr([]string{"\"[L]ang\"", "'l'", "'L'"}, false),
		ComposePoll:                 inputStrOrErr([]string{"\"P[O]ll\"", "'o'", "'O'"}, false),

		MediaDelete:   inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
		MediaEditDesc: inputStrOrErr([]string{"\"[E]dit desc\"", "'e'", "'E'"}, false),
		MediaAdd:      inputStrOrErr([]string{"\"[A]dd\"", "'a'", "'A'"}, false),

		VoteVote:   inputStrOrErr([]string{"\"[V]ote\"", "'v'", "'V'"}, false),
		VoteSelect: inputStrOrErr([]string{"\"[Enter] to select\"", "' '", "\"Enter\""}, false),

		PollAdd:         inputStrOrErr([]string{"\"[A]dd\"", "'a'", "'A'"}, false),
		PollEdit:        inputStrOrErr([]string{"\"[E]dit\"", "'e'", "'E'"}, false),
		PollDelete:      inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
		PollMultiToggle: inputStrOrErr([]string{"\"Toggle [M]ultiple\"", "'m'", "'M'"}, false),
		PollExpiration:  inputStrOrErr([]string{"\"E[X]pires\"", "'x'", "'X'"}, false),

		PreferenceName:         inputStrOrErr([]string{"\"[N]ame\"", "'n'", "'N'"}, false),
		PreferenceBio:          inputStrOrErr([]string{"\"[B]io\"", "'b'", "'B'"}, false),
		PreferenceVisibility:   inputStrOrErr([]string{"\"[V]isibility\"", "'v'", "'V'"}, false),
		PreferenceSave:         inputStrOrErr([]string{"\"[S]ave\"", "'s'", "'S'"}, false),
		PreferenceFields:       inputStrOrErr([]string{"\"[F]ields\"", "'f'", "'F'"}, false),
		PreferenceFieldsAdd:    inputStrOrErr([]string{"\"[A]dd\"", "'a'", "'A'"}, false),
		PreferenceFieldsEdit:   inputStrOrErr([]string{"\"[E]dit\"", "'e'", "'E'"}, false),
		PreferenceFieldsDelete: inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
	}
	ic.GlobalDown = inputOrErr(cfg, "global-down", false, ic.GlobalDown)
	ic.GlobalUp = inputOrErr(cfg, "global-up", false, ic.GlobalUp)
	ic.GlobalEnter = inputOrErr(cfg, "global-enter", false, ic.GlobalEnter)
	ic.GlobalBack = inputOrErr(cfg, "global-back", false, ic.GlobalBack)
	ic.GlobalExit = inputOrErr(cfg, "global-exit", false, ic.GlobalExit)

	ic.MainHome = inputOrErr(cfg, "main-home", false, ic.MainHome)
	ic.MainEnd = inputOrErr(cfg, "main-end", false, ic.MainEnd)
	ic.MainPrevFeed = inputOrErr(cfg, "main-prev-feed", false, ic.MainPrevFeed)
	ic.MainNextFeed = inputOrErr(cfg, "main-next-feed", false, ic.MainNextFeed)
	ic.MainCompose = inputOrErr(cfg, "main-compose", false, ic.MainCompose)

	ic.StatusAvatar = inputOrErr(cfg, "status-avatar", false, ic.StatusAvatar)
	ic.StatusBoost = inputOrErr(cfg, "status-boost", true, ic.StatusBoost)
	ic.StatusDelete = inputOrErr(cfg, "status-delete", false, ic.StatusDelete)
	ic.StatusEdit = inputOrErr(cfg, "status-edit", false, ic.StatusEdit)
	ic.StatusFavorite = inputOrErr(cfg, "status-favorite", true, ic.StatusFavorite)
	ic.StatusMedia = inputOrErr(cfg, "status-media", false, ic.StatusMedia)
	ic.StatusLinks = inputOrErr(cfg, "status-links", false, ic.StatusLinks)
	ic.StatusPoll = inputOrErr(cfg, "status-poll", false, ic.StatusPoll)
	ic.StatusReply = inputOrErr(cfg, "status-reply", false, ic.StatusReply)
	ic.StatusBookmark = inputOrErr(cfg, "status-bookmark", true, ic.StatusBookmark)
	ic.StatusThread = inputOrErr(cfg, "status-thread", false, ic.StatusThread)
	ic.StatusUser = inputOrErr(cfg, "status-user", false, ic.StatusUser)
	ic.StatusViewFocus = inputOrErr(cfg, "status-view-focus", false, ic.StatusViewFocus)
	ic.StatusYank = inputOrErr(cfg, "status-yank", false, ic.StatusYank)
	ic.StatusToggleCW = inputOrErr(cfg, "status-toggle-spoiler", false, ic.StatusToggleCW)
	ts := cfg.Section("input").Key("status-toggle-spoiler").MustString("")
	if ts != "" {
		ic.StatusToggleCW = inputOrErr(cfg, "status-toggle-spoiler", false, ic.StatusToggleCW)
	} else {
		ic.StatusToggleCW = inputOrErr(cfg, "status-toggle-cw", false, ic.StatusToggleCW)
	}

	ic.UserAvatar = inputOrErr(cfg, "user-avatar", false, ic.UserAvatar)
	ic.UserBlock = inputOrErr(cfg, "user-block", true, ic.UserBlock)
	ic.UserFollow = inputOrErr(cfg, "user-follow", true, ic.UserFollow)
	ic.UserFollowRequestDecide = inputOrErr(cfg, "user-follow-request-decide", true, ic.UserFollowRequestDecide)
	ic.UserMute = inputOrErr(cfg, "user-mute", true, ic.UserMute)
	ic.UserLinks = inputOrErr(cfg, "user-links", false, ic.UserLinks)
	ic.UserUser = inputOrErr(cfg, "user-user", false, ic.UserUser)
	ic.UserViewFocus = inputOrErr(cfg, "user-view-focus", false, ic.UserViewFocus)
	ic.UserYank = inputOrErr(cfg, "user-yank", false, ic.UserYank)

	ic.ListOpenFeed = inputOrErr(cfg, "list-open-feed", false, ic.ListOpenFeed)
	ic.ListUserList = inputOrErr(cfg, "list-user-list", false, ic.ListUserList)
	ic.ListUserAdd = inputOrErr(cfg, "list-user-add", false, ic.ListUserAdd)
	ic.ListUserDelete = inputOrErr(cfg, "list-user-delete", false, ic.ListUserDelete)

	ic.TagOpenFeed = inputOrErr(cfg, "tag-open-feed", false, ic.TagOpenFeed)
	ic.TagFollow = inputOrErr(cfg, "tag-follow", true, ic.TagFollow)

	ic.LinkOpen = inputOrErr(cfg, "link-open", false, ic.LinkOpen)
	ic.LinkYank = inputOrErr(cfg, "link-yank", false, ic.LinkYank)

	es := cfg.Section("input").Key("compose-edit-spoiler").MustString("")
	if es != "" {
		ic.ComposeEditCW = inputOrErr(cfg, "compose-edit-spoiler", false, ic.ComposeEditCW)
	} else {
		ic.ComposeEditCW = inputOrErr(cfg, "compose-edit-cw", false, ic.ComposeEditCW)
	}
	ic.ComposeEditText = inputOrErr(cfg, "compose-edit-text", false, ic.ComposeEditText)
	ic.ComposeIncludeQuote = inputOrErr(cfg, "compose-include-quote", false, ic.ComposeIncludeQuote)
	ic.ComposeMediaFocus = inputOrErr(cfg, "compose-media-focus", false, ic.ComposeMediaFocus)
	ic.ComposePost = inputOrErr(cfg, "compose-post", false, ic.ComposePost)
	ic.ComposeToggleContentWarning = inputOrErr(cfg, "compose-toggle-content-warning", false, ic.ComposeToggleContentWarning)
	ic.ComposeVisibility = inputOrErr(cfg, "compose-visibility", false, ic.ComposeVisibility)
	ic.ComposeLanguage = inputOrErr(cfg, "compose-language", false, ic.ComposeLanguage)
	ic.ComposePoll = inputOrErr(cfg, "compose-poll", false, ic.ComposePoll)

	ic.MediaDelete = inputOrErr(cfg, "media-delete", false, ic.MediaDelete)
	ic.MediaEditDesc = inputOrErr(cfg, "media-edit-desc", false, ic.MediaEditDesc)
	ic.MediaAdd = inputOrErr(cfg, "media-add", false, ic.MediaAdd)

	ic.VoteVote = inputOrErr(cfg, "vote-vote", false, ic.VoteVote)
	ic.VoteSelect = inputOrErr(cfg, "vote-select", false, ic.VoteSelect)

	ic.PollAdd = inputOrErr(cfg, "poll-add", false, ic.PollAdd)
	ic.PollEdit = inputOrErr(cfg, "poll-edit", false, ic.PollEdit)
	ic.PollDelete = inputOrErr(cfg, "poll-delete", false, ic.PollDelete)
	ic.PollMultiToggle = inputOrErr(cfg, "poll-multi-toggle", false, ic.PollMultiToggle)
	ic.PollExpiration = inputOrErr(cfg, "poll-expiration", false, ic.PollExpiration)

	ic.PreferenceName = inputOrErr(cfg, "preference-name", false, ic.PreferenceName)
	ic.PreferenceVisibility = inputOrErr(cfg, "preference-visibility", false, ic.PreferenceVisibility)
	ic.PreferenceBio = inputOrErr(cfg, "preference-bio", false, ic.PreferenceBio)
	ic.PreferenceSave = inputOrErr(cfg, "preference-save", false, ic.PreferenceSave)
	ic.PreferenceFields = inputOrErr(cfg, "preference-fields", false, ic.PreferenceFields)
	ic.PreferenceFieldsAdd = inputOrErr(cfg, "preference-fields-add", false, ic.PreferenceFieldsAdd)
	ic.PreferenceFieldsEdit = inputOrErr(cfg, "preference-fields-edit", false, ic.PreferenceFieldsEdit)
	ic.PreferenceFieldsDelete = inputOrErr(cfg, "preference-fields-delete", false, ic.PreferenceFieldsDelete)
	return ic
}

func parseConfig(filepath string, cnfPath string, cnfDir string) (Config, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
		AllowShadows:             true,
	}, filepath)
	conf := Config{}
	if err != nil {
		return conf, err
	}
	conf.General = parseGeneral(cfg)
	conf.Media = parseMedia(cfg)
	conf.Style = parseStyle(cfg, cnfPath, cnfDir)
	conf.OpenPattern = parseOpenPattern(cfg)
	conf.OpenCustom = parseCustom(cfg)
	conf.NotificationConfig = parseNotifications(cfg)
	conf.Templates = parseTemplates(cfg, cnfPath, cnfDir)
	conf.Input = parseInput(cfg)

	return conf, nil
}

func createConfigDir() error {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't find config dir. Err %v", err)
	}
	path := cd + "/tut"
	return os.MkdirAll(path, os.ModePerm)
}

func checkConfig(filename string, cnfPath string, cnfDir string) (path string, exists bool, err error) {
	if cnfPath != "" && filename == "config.ini" {
		_, err = os.Stat(cnfPath)
		if os.IsNotExist(err) {
			return cnfPath, false, nil
		} else if err != nil {
			return cnfPath, true, err
		}
		return cnfPath, true, err
	}
	if cnfDir != "" {
		p := filepath.Join(cnfDir, filename)
		if os.IsNotExist(err) {
			return p, false, nil
		} else if err != nil {
			return p, true, err
		}
		return p, true, err
	}
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't find config dir. Err %v", err)
	}
	dir := cd + "/tut/"
	path = dir + filename
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, false, nil
	} else if err != nil {
		return path, true, err
	}
	return path, true, err
}

func CreateDefaultConfig(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(conftext)
	if err != nil {
		return err
	}
	return nil
}

func getThemes(cnfPath string, cnfDir string) (bundled []string, local []string, err error) {
	entries, err := themesFS.ReadDir("themes")
	if err != nil {
		return bundled, local, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fp := filepath.Join("themes/", entry.Name())
		bundled = append(bundled, fp)
	}
	_, exists, err := checkConfig("themes", cnfPath, cnfDir)
	if err != nil {
		return bundled, local, err
	}
	if !exists {
		return bundled, local, err
	}
	var dir string
	if cnfDir != "" {
		dir = filepath.Join(cnfDir, "themes")
	} else {
		cd, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("couldn't find config dir. Err %v", err)
		}
		dir = filepath.Join(cd, "/tut/themes")
	}
	entries, err = os.ReadDir(dir)
	if err != nil {
		return bundled, local, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fp := filepath.Join(dir, entry.Name())
		local = append(local, fp)
	}
	return bundled, local, nil
}

func getTheme(fname string, isLocal bool, cnfDir string) (*ini.File, error) {
	var f io.Reader
	var err error
	if isLocal {
		var dir string
		if cnfDir != "" {
			dir = filepath.Join(cnfDir, "themes")
		} else {
			cd, err := os.UserConfigDir()
			if err != nil {
				log.Fatalf("couldn't find config dir. Err %v", err)
			}
			dir = filepath.Join(cd, "/tut/themes")
		}
		f, err = os.Open(
			filepath.Join(dir, fmt.Sprintf("%s.ini", strings.TrimSpace(fname))),
		)
	} else {
		f, err = themesFS.Open(fmt.Sprintf("themes/%s.ini", strings.TrimSpace(fname)))
	}
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}, content)
	if err != nil {
		return nil, err
	}
	keys := []string{
		"background",
		"text",
		"subtle",
		"warning-text",
		"text-special-one",
		"text-special-two",
		"top-bar-background",
		"top-bar-text",
		"status-bar-background",
		"status-bar-text",
		"status-bar-view-background",
		"status-bar-view-text",
		"list-selected-background",
		"list-selected-text",
	}
	for _, k := range keys {
		if !cfg.Section("").HasKey(k) {
			return nil, fmt.Errorf("theme %s is missing %s", fname, k)
		}
	}
	return cfg, nil
}
