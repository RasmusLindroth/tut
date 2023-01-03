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

	"github.com/RasmusLindroth/tut/util"
	"github.com/gdamore/tcell/v2"
	"github.com/gobwas/glob"
	"github.com/pelletier/go-toml/v2"
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
	Hidden      bool
	HideBoosts  bool
	HideReplies bool
}

type General struct {
	Confirmation        bool
	MouseSupport        bool
	DateTodayFormat     string
	DateFormat          string
	DateRelative        int
	MaxWidth            int
	QuoteReply          bool
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
	Timelines           []*Timeline
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
	Key      Key
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

func NilDefaultBool(x *bool, def *bool) bool {
	if x == nil {
		return *def
	}
	return *x
}

func NilDefaultString(x *string, def *string) string {
	if x == nil {
		return *def
	}
	return *x
}

func NilDefaultInt(x *int, def *int) int {
	if x == nil {
		return *def
	}
	return *x
}

func NilDefaultInt64(x *int64, def *int64) int64 {
	if x == nil {
		return *def
	}
	return *x
}

var keyMatch = regexp.MustCompile(`^(.*?)\[(.*?)\](.*?)$`)

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

func NewKeyT(hint string, hintAlt string, keys []string, special []string) (Key, error) {
	k := Key{}
	if len(hint) > 0 && len(hintAlt) > 0 {
		k.Hint = [][]string{newHint(hint), newHint(hintAlt)}
	} else if len(hint) > 0 {
		k.Hint = [][]string{newHint(hint), newHint(hintAlt)}
	}
	var runes []rune
	var keysTcell []tcell.Key
	for _, r := range keys {
		if len(r) > 1 {
			return k, fmt.Errorf("key %s can only be one char", r)
		}
		if len(r) == 0 {
			continue
		}
		runes = append(runes, rune(r[0]))
	}
	for _, s := range special {
		found := false
		var fk tcell.Key
		for tk, tv := range tcell.KeyNames {
			if tv == s {
				found = true
				fk = tk
				break
			}
		}
		if found {
			keysTcell = append(keysTcell, fk)
		} else {
			return k, fmt.Errorf("no key named %s", s)
		}
	}
	k.Runes = runes
	k.Keys = keysTcell

	return k, nil
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

func parseTheme(cfg StyleTOML, xrdbColors map[string]string) Style {
	var style Style
	def := ConfigDefault.Style
	s := NilDefaultString(cfg.Background, def.Background)
	style.Background = parseColor(s, "#27822", xrdbColors)

	s = NilDefaultString(cfg.Text, def.Text)
	style.Text = parseColor(s, "#f8f8f2", xrdbColors)

	s = NilDefaultString(cfg.Subtle, def.Subtle)
	style.Subtle = parseColor(s, "#808080", xrdbColors)

	s = NilDefaultString(cfg.WarningText, def.WarningText)
	style.WarningText = parseColor(s, "#f92672", xrdbColors)

	s = NilDefaultString(cfg.TextSpecial1, def.TextSpecial1)
	style.TextSpecial1 = parseColor(s, "#ae81ff", xrdbColors)

	s = NilDefaultString(cfg.TextSpecial2, def.TextSpecial2)
	style.TextSpecial2 = parseColor(s, "#a6e22e", xrdbColors)

	s = NilDefaultString(cfg.TopBarBackground, def.TopBarBackground)
	style.TopBarBackground = parseColor(s, "#f92672", xrdbColors)

	s = NilDefaultString(cfg.TopBarText, def.TopBarText)
	style.TopBarText = parseColor(s, "white", xrdbColors)

	s = NilDefaultString(cfg.StatusBarBackground, def.StatusBarBackground)
	style.StatusBarBackground = parseColor(s, "#f92672", xrdbColors)

	s = NilDefaultString(cfg.StatusBarText, def.StatusBarText)
	style.StatusBarText = parseColor(s, "white", xrdbColors)

	s = NilDefaultString(cfg.StatusBarViewBackground, def.StatusBarViewBackground)
	style.StatusBarViewBackground = parseColor(s, "#ae81ff", xrdbColors)

	s = NilDefaultString(cfg.StatusBarViewText, def.StatusBarViewText)
	style.StatusBarViewText = parseColor(s, "white", xrdbColors)

	s = NilDefaultString(cfg.ListSelectedBackground, def.ListSelectedBackground)
	style.ListSelectedBackground = parseColor(s, "#f92672", xrdbColors)

	s = NilDefaultString(cfg.ListSelectedText, def.ListSelectedText)
	style.ListSelectedText = parseColor(s, "white", xrdbColors)

	s = NilDefaultString(cfg.ListSelectedInactiveBackground, sp(""))
	if len(s) > 0 {
		style.ListSelectedInactiveBackground = parseColor(s, "#ae81ff", xrdbColors)
	} else {
		style.ListSelectedInactiveBackground = style.StatusBarViewBackground
	}
	s = NilDefaultString(cfg.ListSelectedInactiveText, def.ListSelectedInactiveText)
	if len(s) > 0 {
		style.ListSelectedInactiveText = parseColor(s, "#f8f8f2", xrdbColors)
	} else {
		style.ListSelectedInactiveText = style.StatusBarViewText
	}

	s = NilDefaultString(cfg.ControlsText, sp(""))
	if len(s) > 0 {
		style.ControlsText = parseColor(s, "#f8f8f2", xrdbColors)
	} else {
		style.ControlsText = style.Text
	}
	s = NilDefaultString(cfg.ControlsHighlight, sp(""))
	if len(s) > 0 {
		style.ControlsHighlight = parseColor(s, "#a6e22e", xrdbColors)
	} else {
		style.ControlsHighlight = style.TextSpecial2
	}

	s = NilDefaultString(cfg.AutocompleteBackground, sp(""))
	if len(s) > 0 {
		style.AutocompleteBackground = parseColor(s, "#272822", xrdbColors)
	} else {
		style.AutocompleteBackground = style.Background
	}
	s = NilDefaultString(cfg.AutocompleteText, sp(""))
	if len(s) > 0 {
		style.AutocompleteText = parseColor(s, "#f8f8f2", xrdbColors)
	} else {
		style.AutocompleteText = style.Text
	}
	s = NilDefaultString(cfg.AutocompleteSelectedBackground, sp(""))
	if len(s) > 0 {
		style.AutocompleteSelectedBackground = parseColor(s, "#ae81ff", xrdbColors)
	} else {
		style.AutocompleteSelectedBackground = style.StatusBarViewBackground
	}
	s = NilDefaultString(cfg.AutocompleteSelectedText, sp(""))
	if len(s) > 0 {
		style.AutocompleteSelectedText = parseColor(s, "#f8f8f2", xrdbColors)
	} else {
		style.AutocompleteSelectedText = style.StatusBarViewText
	}

	s = NilDefaultString(cfg.ButtonColorOne, sp(""))
	if len(s) > 0 {
		style.ButtonColorOne = parseColor(s, "#ae81ff", xrdbColors)
	} else {
		style.ButtonColorOne = style.StatusBarViewBackground
	}
	s = NilDefaultString(cfg.ButtonColorTwo, sp(""))
	if len(s) > 0 {
		style.ButtonColorTwo = parseColor(s, "#272822", xrdbColors)
	} else {
		style.ButtonColorTwo = style.Background
	}

	s = NilDefaultString(cfg.TimelineNameBackground, sp(""))
	if len(s) > 0 {
		style.TimelineNameBackground = parseColor(s, "#272822", xrdbColors)
	} else {
		style.TimelineNameBackground = style.Background
	}
	s = NilDefaultString(cfg.TimelineNameText, sp(""))
	if len(s) > 0 {
		style.TimelineNameText = parseColor(s, "gray", xrdbColors)
	} else {
		style.TimelineNameText = style.Subtle
	}

	s = NilDefaultString(cfg.CommandText, sp(""))
	if len(s) > 0 {
		style.CommandText = parseColor(s, "white", xrdbColors)
	} else {
		style.CommandText = style.StatusBarText
	}
	return style
}

func parseStyle(cfg StyleTOML, cnfPath string, cnfDir string) Style {
	var xrdbColors map[string]string
	xrdbMap, _ := GetXrdbColors()
	def := ConfigDefault.Style
	prefix := NilDefaultString(cfg.XrdbPrefix, def.XrdbPrefix)
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
	theme := NilDefaultString(cfg.Theme, def.Theme)
	if theme != "none" && theme != "" {
		bundled, local, err := getThemes(cnfPath, cnfDir)
		if err != nil {
			log.Fatalf("Couldn't load themes. Error: %s\n", err)
		}
		found := false
		isLocal := false
		for _, t := range local {
			if filepath.Base(t) == fmt.Sprintf("%s.toml", theme) {
				found = true
				isLocal = true
				break
			}
		}
		if !found {
			for _, t := range bundled {
				if filepath.Base(t) == fmt.Sprintf("%s.toml", theme) {
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
		style = parseTheme(tcfg, xrdbColors)

	} else {
		style = parseTheme(cfg, xrdbColors)
	}

	return style
}

func getViewer(v *ViewerTOML, def *ViewerTOML) (program, args string, terminal, single, reverse bool) {
	program = *def.Program
	args = *def.Args
	terminal = *def.Terminal
	single = *def.Single
	reverse = *def.Reverse
	if v == nil {
		return
	}
	if v.Program != nil {
		program = *v.Program
	}
	if v.Args != nil {
		args = *v.Args
	}
	if v.Terminal != nil {
		terminal = *v.Terminal
	}
	if v.Single != nil {
		single = *v.Single
	}
	if v.Reverse != nil {
		reverse = *v.Reverse
	}
	return
}

func parseMedia(cfg MediaTOML) Media {
	media := Media{}
	var program, args string
	var terminal, single, reverse bool

	program, args, terminal, single, reverse = getViewer(cfg.Image, ConfigDefault.Media.Image)
	media.ImageViewer = program
	media.ImageArgs = strings.Fields(args)
	media.ImageTerminal = terminal
	media.ImageSingle = single
	media.ImageReverse = reverse

	program, args, terminal, single, reverse = getViewer(cfg.Video, ConfigDefault.Media.Video)
	media.VideoViewer = program
	media.VideoArgs = strings.Fields(args)
	media.VideoTerminal = terminal
	media.VideoSingle = single
	media.VideoReverse = reverse

	program, args, terminal, single, reverse = getViewer(cfg.Audio, ConfigDefault.Media.Audio)
	media.AudioViewer = program
	media.AudioArgs = strings.Fields(args)
	media.AudioTerminal = terminal
	media.AudioSingle = single
	media.AudioReverse = reverse

	program, args, terminal, _, _ = getViewer(cfg.Link, ConfigDefault.Media.Link)
	media.LinkViewer = program
	media.LinkArgs = strings.Fields(args)
	media.LinkTerminal = terminal

	return media
}

func parseGeneral(cfg GeneralTOML) General {
	general := General{}

	def := ConfigDefault.General
	general.Confirmation = NilDefaultBool(cfg.Confirmation, def.Confirmation)
	general.Confirmation = NilDefaultBool(cfg.MouseSupport, def.MouseSupport)

	dateFormat := NilDefaultString(cfg.DateFormat, def.DateFormat)
	if dateFormat == "" {
		dateFormat = "2006-01-02 15:04"
	}
	general.DateFormat = dateFormat

	dateTodayFormat := NilDefaultString(cfg.DateTodayFormat, def.DateTodayFormat)
	if dateTodayFormat == "" {
		dateTodayFormat = "15:04"
	}
	general.DateTodayFormat = dateTodayFormat

	general.DateRelative = NilDefaultInt(cfg.DateRelative, def.DateRelative)

	general.QuoteReply = NilDefaultBool(cfg.QuoteReply, def.QuoteReply)
	general.MaxWidth = NilDefaultInt(cfg.MaxWidth, def.MaxWidth)
	general.ShortHints = NilDefaultBool(cfg.ShortHints, def.ShortHints)
	general.ShowFilterPhrase = NilDefaultBool(cfg.ShowFilterPhrase, def.ShowFilterPhrase)
	general.ShowIcons = NilDefaultBool(cfg.ShowIcons, def.ShowIcons)
	general.ShowHelp = NilDefaultBool(cfg.ShowHelp, def.ShowHelp)
	general.RedrawUI = NilDefaultBool(cfg.RedrawUI, def.RedrawUI)
	general.StickToTop = NilDefaultBool(cfg.StickToTop, def.StickToTop)
	general.ShowBoostedUser = NilDefaultBool(cfg.ShowBoostedUser, def.ShowBoostedUser)

	lp := NilDefaultString(cfg.ListPlacement, def.ListPlacement)
	switch lp {
	case "left":
		general.ListPlacement = ListPlacementLeft
	case "right":
		general.ListPlacement = ListPlacementRight
	case "top":
		general.ListPlacement = ListPlacementTop
	case "bottom":
		general.ListPlacement = ListPlacementBottom
	default:
		general.ListPlacement = ListPlacementLeft
	}
	ls := NilDefaultString(cfg.ListSplit, def.ListSplit)
	switch ls {
	case "row":
		general.ListSplit = ListRow
	case "column":
		general.ListSplit = ListColumn
	}

	listProp := NilDefaultInt(cfg.ListProportion, def.ListProportion)
	if listProp < 1 {
		listProp = 1
	}
	contentProp := NilDefaultInt(cfg.ContentProportion, def.ContentProportion)
	if contentProp < 1 {
		contentProp = 1
	}
	general.ListProportion = listProp
	general.ContentProportion = contentProp

	leaderString := NilDefaultString(cfg.LeaderKey, def.LeaderKey)
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
	general.LeaderTimeout = NilDefaultInt64(cfg.LeaderTimeout, def.LeaderTimeout)
	if general.LeaderKey != rune(0) {
		var las []LeaderAction
		if cfg.LeaderActions != nil {
			lactions := *cfg.LeaderActions
			for _, l := range lactions {
				la := LeaderAction{}
				ltype := NilDefaultString(l.Type, sp(""))
				ldata := NilDefaultString(l.Data, sp(""))
				lshortcut := NilDefaultString(l.Shortcut, sp(""))
				switch ltype {
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
					la.Subaction = ldata
				case "tags":
					la.Command = LeaderTags
				case "list-placement":
					la.Command = LeaderListPlacement
					la.Subaction = ldata
				case "list-split":
					la.Command = LeaderListSplit
					la.Subaction = ldata
				case "proportions":
					la.Command = LeaderProportions
					la.Subaction = ldata
				case "window":
					la.Command = LeaderWindow
					la.Subaction = ldata
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
				case "newer":
					la.Command = LeaderLoadNewer
				default:
					fmt.Printf("leader-action %s is invalid\n", ltype)
					os.Exit(1)
				}
				la.Shortcut = lshortcut
				las = append(las, la)
			}
		}
		general.LeaderActions = las
	}

	var tls []*Timeline
	timelines := cfg.Timelines
	if cfg.Timelines != nil {
		for _, l := range *timelines {
			tl := Timeline{}
			if l.Type == nil {
				fmt.Println("timelines must have a type")
				os.Exit(1)
			}
			switch *l.Type {
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
				tl.Subaction = NilDefaultString(l.Data, sp(""))
			default:
				fmt.Printf("timeline %s is invalid\n", *l.Type)
				os.Exit(1)
			}
			tl.Name = NilDefaultString(l.Name, sp(""))
			tl.HideBoosts = NilDefaultBool(l.HideBoosts, bf)
			tl.HideReplies = NilDefaultBool(l.HideReplies, bf)
			if l.Keys != nil {
				var keys []string
				var special []string
				if l.Keys != nil {
					keys = *l.Keys
				}
				if l.SpecialKeys != nil {
					special = *l.SpecialKeys
				}
				var err error
				tl.Key, err = NewKeyT("", "", keys, special)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			tls = append(tls, &tl)
		}
	}
	if len(tls) == 0 {
		tls = append(tls,
			&Timeline{
				FeedType: TimelineHome,
				Name:     "Home",
			},
		)
		tls = append(tls,
			&Timeline{
				FeedType: Notifications,
				Name:     "[N]otifications",
				Key:      inputStrOrErr([]string{"", "'n'", "'N'"}, false),
			},
		)
	}
	general.Timelines = tls

	general.TerminalTitle = NilDefaultInt(cfg.TerminalTitle, def.TerminalTitle)
	/*
		0 = No terminal title
		1 = Show title in terminal and top bar
		2 = Only show terminal title, and no top bar
	*/
	if general.TerminalTitle < 0 || general.TerminalTitle > 2 {
		general.TerminalTitle = 0
	}

	nths := []NotificationToHide{}
	nth := cfg.NotificationsToHide
	if nth != nil {
		for _, n := range *nth {
			switch n {
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
				log.Fatalf("%s in notifications-to-hide is invalid\n", n)
				os.Exit(1)
			}
			general.NotificationsToHide = nths
		}
	}
	return general
}

func parseOpenPattern(cfg OpenPatternTOML) OpenPattern {
	om := OpenPattern{}
	if cfg.Patterns == nil {
		return om
	}
	for _, p := range *cfg.Patterns {
		pattern := Pattern{
			Program:  NilDefaultString(p.Program, sp("")),
			Terminal: NilDefaultBool(p.Terminal, bf),
		}
		if p.Args != nil {
			pattern.Args = strings.Fields(*p.Args)
		}
		pg := NilDefaultString(p.Matching, sp(""))
		compiled, err := glob.Compile(pg)
		if err != nil {
			panic(fmt.Sprintf("Couldn't compile pattern %s in config. Error: %v", pg, err))
		}
		pattern.Compiled = compiled
		om.Patterns = append(om.Patterns, pattern)
	}

	return om
}

func parseCustom(cfg OpenCustomTOML) OpenCustom {
	oc := OpenCustom{}
	if cfg.Programs == nil {
		return oc
	}
	for _, x := range *cfg.Programs {
		keys, special := []string{}, []string{}
		if x.Keys != nil {
			keys = *x.Keys
		}
		if x.SpecialKeys != nil {
			special = *x.SpecialKeys
		}
		key, err := NewKeyT(
			NilDefaultString(x.Hint, sp("")),
			"", keys, special,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		use := NilDefaultString(x.Program, sp(""))
		terminal := NilDefaultBool(x.Terminal, bf)
		if use == "" {
			continue
		}
		args := strings.Fields(NilDefaultString(x.Args, sp("")))
		c := Custom{
			Program:  use,
			Args:     args,
			Terminal: terminal,
			Key:      key,
		}
		oc.OpenCustoms = append(oc.OpenCustoms, c)
	}
	return oc
}

func parseNotifications(cfg NotificationsTOML) Notification {
	nc := Notification{}
	def := ConfigDefault.NotificationConfig
	nc.NotificationFollower = NilDefaultBool(cfg.Followers, def.Followers)
	nc.NotificationFavorite = NilDefaultBool(cfg.Favorite, def.Favorite)
	nc.NotificationMention = NilDefaultBool(cfg.Mention, def.Mention)
	nc.NotificationUpdate = NilDefaultBool(cfg.Update, def.Update)
	nc.NotificationBoost = NilDefaultBool(cfg.Boost, def.Followers)
	nc.NotificationPoll = NilDefaultBool(cfg.Poll, def.Poll)
	nc.NotificationPost = NilDefaultBool(cfg.Posts, def.Posts)
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

func inputStrOrErr(vals []string, double bool) Key {
	k, err := NewKey(vals, double)
	if err != nil {
		fmt.Printf("error parsing config. Error: %v\n", err)
		os.Exit(1)
	}
	return k
}

func inputOrDef(keyName string, user *KeyHintTOML, def *KeyHintTOML, double bool) Key {
	values := *def
	if user != nil {
		values = *user
	}
	keys, special := []string{}, []string{}
	if values.Keys != nil {
		keys = *values.Keys
	}
	if values.SpecialKeys != nil {
		special = *values.SpecialKeys
	}
	key, err := NewKeyT(
		NilDefaultString(values.Hint, sp("")),
		NilDefaultString(values.HintAlt, sp("")),
		keys, special,
	)
	if err != nil {
		fmt.Printf("error parsing config for key %s. Error: %v\n", keyName, err)
		os.Exit(1)
	}
	return key
}

func parseInput(cfg InputTOML) Input {
	def := ConfigDefault.Input
	ic := Input{}
	ic.GlobalDown = inputOrDef("global-down", cfg.GlobalDown, def.GlobalDown, false)
	ic.GlobalUp = inputOrDef("global-up", cfg.GlobalUp, def.GlobalUp, false)
	ic.GlobalEnter = inputOrDef("global-enter", cfg.GlobalEnter, def.GlobalEnter, false)
	ic.GlobalBack = inputOrDef("global-back", cfg.GlobalBack, def.GlobalBack, false)
	ic.GlobalExit = inputOrDef("global-exit", cfg.GlobalExit, def.GlobalExit, false)

	ic.MainHome = inputOrDef("main-home", cfg.MainHome, def.MainHome, false)
	ic.MainEnd = inputOrDef("main-end", cfg.MainEnd, def.MainEnd, false)
	ic.MainPrevFeed = inputOrDef("main-prev-feed", cfg.MainPrevFeed, def.MainPrevFeed, false)
	ic.MainNextFeed = inputOrDef("main-next-feed", cfg.MainNextFeed, def.MainNextFeed, false)
	ic.MainNextWindow = inputOrDef("main-next-window", cfg.MainNextWindow, def.MainNextWindow, false)
	ic.MainPrevWindow = inputOrDef("main-prev-window", cfg.MainPrevWindow, def.MainPrevWindow, false)
	ic.MainCompose = inputOrDef("main-compose", cfg.MainCompose, def.MainCompose, false)

	ic.StatusAvatar = inputOrDef("status-avatar", cfg.StatusAvatar, def.StatusAvatar, false)
	ic.StatusBoost = inputOrDef("status-boost", cfg.StatusBoost, def.StatusBoost, true)
	ic.StatusDelete = inputOrDef("status-delete", cfg.StatusDelete, def.StatusDelete, false)
	ic.StatusEdit = inputOrDef("status-edit", cfg.StatusEdit, def.StatusEdit, false)
	ic.StatusFavorite = inputOrDef("status-favorite", cfg.StatusFavorite, def.StatusFavorite, true)
	ic.StatusMedia = inputOrDef("status-media", cfg.StatusMedia, def.StatusMedia, false)
	ic.StatusLinks = inputOrDef("status-links", cfg.StatusLinks, def.StatusLinks, false)
	ic.StatusPoll = inputOrDef("status-poll", cfg.StatusPoll, def.StatusPoll, false)
	ic.StatusReply = inputOrDef("status-reply", cfg.StatusReply, def.StatusReply, false)
	ic.StatusBookmark = inputOrDef("status-bookmark", cfg.StatusBookmark, def.StatusBookmark, true)
	ic.StatusThread = inputOrDef("status-thread", cfg.StatusThread, def.StatusThread, false)
	ic.StatusUser = inputOrDef("status-user", cfg.StatusUser, def.StatusUser, false)
	ic.StatusViewFocus = inputOrDef("status-view-focus", cfg.StatusViewFocus, def.StatusViewFocus, false)
	ic.StatusYank = inputOrDef("status-yank", cfg.StatusYank, def.StatusYank, false)
	ic.StatusToggleCW = inputOrDef("status-toggle-cw", cfg.StatusToggleCW, def.StatusToggleCW, false)

	ic.UserAvatar = inputOrDef("user-avatar", cfg.UserAvatar, def.UserAvatar, false)
	ic.UserBlock = inputOrDef("user-block", cfg.UserBlock, def.UserBlock, true)
	ic.UserFollow = inputOrDef("user-follow", cfg.UserFollow, def.UserFollow, true)
	ic.UserFollowRequestDecide = inputOrDef("user-follow-request-decide", cfg.UserFollowRequestDecide, def.UserFollowRequestDecide, true)
	ic.UserMute = inputOrDef("user-mute", cfg.UserMute, def.UserMute, true)
	ic.UserLinks = inputOrDef("user-links", cfg.UserLinks, def.UserLinks, false)
	ic.UserUser = inputOrDef("user-user", cfg.UserUser, def.UserUser, false)
	ic.UserViewFocus = inputOrDef("user-view-focus", cfg.UserViewFocus, def.UserViewFocus, false)
	ic.UserYank = inputOrDef("user-yank", cfg.UserYank, def.UserYank, false)

	ic.ListOpenFeed = inputOrDef("list-open-feed", cfg.ListOpenFeed, def.ListOpenFeed, false)
	ic.ListUserList = inputOrDef("list-user-list", cfg.ListUserList, def.ListUserList, false)
	ic.ListUserAdd = inputOrDef("list-user-add", cfg.ListUserAdd, def.ListUserAdd, false)
	ic.ListUserDelete = inputOrDef("list-user-delete", cfg.ListUserDelete, def.ListUserDelete, false)

	ic.TagOpenFeed = inputOrDef("tag-open-feed", cfg.TagOpenFeed, def.TagOpenFeed, false)
	ic.TagFollow = inputOrDef("tag-follow", cfg.TagFollow, def.TagFollow, true)
	ic.LinkOpen = inputOrDef("link-open", cfg.LinkOpen, def.LinkOpen, false)
	ic.LinkYank = inputOrDef("link-yank", cfg.LinkYank, def.LinkYank, false)

	ic.ComposeEditCW = inputOrDef("compose-edit-cw", cfg.ComposeEditCW, def.ComposeEditCW, false)
	ic.ComposeEditText = inputOrDef("compose-edit-text", cfg.ComposeEditText, def.ComposeEditText, false)
	ic.ComposeIncludeQuote = inputOrDef("compose-include-quote", cfg.ComposeIncludeQuote, def.ComposeIncludeQuote, false)
	ic.ComposeMediaFocus = inputOrDef("compose-media-focus", cfg.ComposeMediaFocus, def.ComposeMediaFocus, false)
	ic.ComposePost = inputOrDef("compose-post", cfg.ComposePost, def.ComposePost, false)
	ic.ComposeToggleContentWarning = inputOrDef("compose-toggle-content-warning", cfg.ComposeToggleContentWarning, def.ComposeToggleContentWarning, false)
	ic.ComposeVisibility = inputOrDef("compose-visibility", cfg.ComposeVisibility, def.ComposeVisibility, false)
	ic.ComposeLanguage = inputOrDef("compose-language", cfg.ComposeLanguage, def.ComposeLanguage, false)
	ic.ComposePoll = inputOrDef("compose-poll", cfg.ComposePoll, def.ComposePoll, false)

	ic.MediaDelete = inputOrDef("media-delete", cfg.MediaDelete, def.MediaDelete, false)
	ic.MediaEditDesc = inputOrDef("media-edit-desc", cfg.MediaEditDesc, def.MediaEditDesc, false)
	ic.MediaAdd = inputOrDef("media-add", cfg.MediaAdd, def.MediaAdd, false)

	ic.VoteVote = inputOrDef("vote-vote", cfg.VoteVote, def.VoteVote, false)
	ic.VoteSelect = inputOrDef("vote-select", cfg.VoteSelect, def.VoteSelect, false)

	ic.PollAdd = inputOrDef("poll-add", cfg.PollAdd, def.PollAdd, false)
	ic.PollEdit = inputOrDef("poll-edit", cfg.PollEdit, def.PollEdit, false)
	ic.PollDelete = inputOrDef("poll-delete", cfg.PollDelete, def.PollDelete, false)
	ic.PollMultiToggle = inputOrDef("poll-multi-toggle", cfg.PollMultiToggle, def.PollMultiToggle, false)
	ic.PollExpiration = inputOrDef("poll-expiration", cfg.PollExpiration, def.PollExpiration, false)

	ic.PreferenceName = inputOrDef("preference-name", cfg.PreferenceName, def.PreferenceName, false)
	ic.PreferenceVisibility = inputOrDef("preference-visibility", cfg.PreferenceVisibility, def.PreferenceVisibility, false)
	ic.PreferenceBio = inputOrDef("preference-bio", cfg.PreferenceBio, def.PreferenceBio, false)
	ic.PreferenceSave = inputOrDef("preference-save", cfg.PreferenceSave, def.PreferenceSave, false)
	ic.PreferenceFields = inputOrDef("preference-fields", cfg.PreferenceFields, def.PreferenceFields, false)
	ic.PreferenceFieldsAdd = inputOrDef("preference-fields-add", cfg.PreferenceFieldsAdd, def.PreferenceFieldsAdd, false)
	ic.PreferenceFieldsEdit = inputOrDef("preference-fields-edit", cfg.PreferenceFieldsEdit, def.PreferenceFieldsEdit, false)
	ic.PreferenceFieldsDelete = inputOrDef("preference-fields-delete", cfg.PreferenceFieldsDelete, def.PreferenceFieldsDelete, false)
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
	f, err := os.Open(strings.TrimSuffix(filepath, "ini") + "toml")
	if err != nil {
		log.Fatalln(err)
	}
	var a ConfigTOML
	toml.NewDecoder(f).Decode(&a)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
	conf.General = parseGeneral(a.General)
	conf.Media = parseMedia(a.Media)
	conf.Style = parseStyle(a.Style, cnfPath, cnfDir)
	conf.OpenPattern = parseOpenPattern(a.OpenPattern)
	conf.OpenCustom = parseCustom(a.OpenCustom)
	conf.NotificationConfig = parseNotifications(a.NotificationConfig)
	conf.Templates = parseTemplates(cfg, cnfPath, cnfDir)
	conf.Input = parseInput(a.Input)

	return conf, nil
}

func createConfigDir() error {
	cd, err := util.GetConfigDir()
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
	cd, err := util.GetConfigDir()
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
		cd, err := util.GetConfigDir()
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

func getTheme(fname string, isLocal bool, cnfDir string) (StyleTOML, error) {
	var f io.Reader
	var err error
	if isLocal {
		var dir string
		if cnfDir != "" {
			dir = filepath.Join(cnfDir, "themes")
		} else {
			cd, err := util.GetConfigDir()
			if err != nil {
				log.Fatalf("couldn't find config dir. Err %v", err)
			}
			dir = filepath.Join(cd, "/tut/themes")
		}
		f, err = os.Open(
			filepath.Join(dir, fmt.Sprintf("%s.toml", strings.TrimSpace(fname))),
		)
	} else {
		f, err = themesFS.Open(fmt.Sprintf("themes/%s.toml", strings.TrimSpace(fname)))
	}
	if err != nil {
		return StyleTOML{}, err
	}
	var style StyleTOML
	toml.NewDecoder(f).Decode(&style)
	if err != nil {
		return style, err
	}
	switch x := f.(type) {
	case *os.File:
		x.Close()
	}
	return style, nil
}
