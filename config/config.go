package config

import (
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/RasmusLindroth/tut/feed"
	"github.com/gdamore/tcell/v2"
	"github.com/gobwas/glob"
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
	LeaderCompose
	LeaderBlocking
	LeaderBookmarks
	LeaderSaved
	LeaderFavorited
	LeaderBoosts
	LeaderFavorites
	LeaderFollowing
	LeaderFollowers
	LeaderMuting
	LeaderProfile
	LeaderNotifications
	LeaderLists
	LeaderTag
	LeaderUser
)

type General struct {
	Confirmation         bool
	DateTodayFormat      string
	DateFormat           string
	DateRelative         int
	MaxWidth             int
	StartTimeline        feed.FeedType
	NotificationFeed     bool
	QuoteReply           bool
	CharLimit            int
	ShortHints           bool
	ListPlacement        ListPlacement
	ListSplit            ListSplit
	HideNotificationText bool
	ListProportion       int
	ContentProportion    int
	ShowIcons            bool
	ShowHelp             bool
	RedrawUI             bool
	LeaderKey            rune
	LeaderTimeout        int64
	LeaderActions        []LeaderAction
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
	NotificationBoost
	NotificationPoll
	NotificationPost
)

type Notification struct {
	NotificationFollower bool
	NotificationFavorite bool
	NotificationMention  bool
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

	MainHome              Key
	MainEnd               Key
	MainPrevFeed          Key
	MainNextFeed          Key
	MainNotificationFocus Key
	MainCompose           Key

	StatusAvatar        Key
	StatusBoost         Key
	StatusDelete        Key
	StatusFavorite      Key
	StatusMedia         Key
	StatusLinks         Key
	StatusPoll          Key
	StatusReply         Key
	StatusBookmark      Key
	StatusThread        Key
	StatusUser          Key
	StatusViewFocus     Key
	StatusYank          Key
	StatusToggleSpoiler Key

	UserAvatar    Key
	UserBlock     Key
	UserFollow    Key
	UserMute      Key
	UserLinks     Key
	UserUser      Key
	UserViewFocus Key
	UserYank      Key

	ListOpenFeed Key

	LinkOpen Key
	LinkYank Key

	ComposeEditSpoiler          Key
	ComposeEditText             Key
	ComposeIncludeQuote         Key
	ComposeMediaFocus           Key
	ComposePost                 Key
	ComposeToggleContentWarning Key
	ComposeVisibility           Key

	MediaDelete   Key
	MediaEditDesc Key
	MediaAdd      Key

	VoteVote   Key
	VoteSelect Key
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

func parseStyle(cfg *ini.File) Style {
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
		themes, err := getThemes()
		if err != nil {
			log.Fatalf("Couldn't load themes. Error: %s\n", err)
		}
		found := false
		for _, t := range themes {
			if filepath.Base(t) == fmt.Sprintf("%s.ini", theme) {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Couldn't find theme %s\n", theme)
		}
		tcfg, err := getTheme(theme)
		if err != nil {
			log.Fatalf("Couldn't load theme. Error: %s\n", err)
		}
		bg := tcfg.Section("").Key("background").String()
		style.Background = parseColor(bg, "default", xrdbColors)

		text := tcfg.Section("").Key("text").String()
		style.Text = tcell.GetColor(text)

		subtle := tcfg.Section("").Key("subtle").String()
		style.Subtle = tcell.GetColor(subtle)

		warningText := tcfg.Section("").Key("warning-text").String()
		style.WarningText = tcell.GetColor(warningText)

		textSpecial1 := tcfg.Section("").Key("text-special-one").String()
		style.TextSpecial1 = tcell.GetColor(textSpecial1)

		textSpecial2 := tcfg.Section("").Key("text-special-two").String()
		style.TextSpecial2 = tcell.GetColor(textSpecial2)

		topBarBackround := tcfg.Section("").Key("top-bar-background").String()
		style.TopBarBackground = tcell.GetColor(topBarBackround)

		topBarText := tcfg.Section("").Key("top-bar-text").String()
		style.TopBarText = tcell.GetColor(topBarText)

		statusBarBackround := tcfg.Section("").Key("status-bar-background").String()
		style.StatusBarBackground = tcell.GetColor(statusBarBackround)

		statusBarText := tcfg.Section("").Key("status-bar-text").String()
		style.StatusBarText = tcell.GetColor(statusBarText)

		statusBarViewBackround := tcfg.Section("").Key("status-bar-view-background").String()
		style.StatusBarViewBackground = tcell.GetColor(statusBarViewBackround)

		statusBarViewText := tcfg.Section("").Key("status-bar-view-text").String()
		style.StatusBarViewText = tcell.GetColor(statusBarViewText)

		listSelectedBackground := tcfg.Section("").Key("list-selected-background").String()
		style.ListSelectedBackground = tcell.GetColor(listSelectedBackground)

		listSelectedText := tcfg.Section("").Key("list-selected-text").String()
		style.ListSelectedText = tcell.GetColor(listSelectedText)
	} else {
		bg := cfg.Section("style").Key("background").String()
		style.Background = parseColor(bg, "default", xrdbColors)

		text := cfg.Section("style").Key("text").String()
		style.Text = parseColor(text, "white", xrdbColors)

		subtle := cfg.Section("style").Key("subtle").String()
		style.Subtle = parseColor(subtle, "gray", xrdbColors)

		warningText := cfg.Section("style").Key("warning-text").String()
		style.WarningText = parseColor(warningText, "#f92672", xrdbColors)

		textSpecial1 := cfg.Section("style").Key("text-special-one").String()
		style.TextSpecial1 = parseColor(textSpecial1, "#ae81ff", xrdbColors)

		textSpecial2 := cfg.Section("style").Key("text-special-two").String()
		style.TextSpecial2 = parseColor(textSpecial2, "#a6e22e", xrdbColors)

		topBarBackround := cfg.Section("style").Key("top-bar-background").String()
		style.TopBarBackground = parseColor(topBarBackround, "#f92672", xrdbColors)

		topBarText := cfg.Section("style").Key("top-bar-text").String()
		style.TopBarText = parseColor(topBarText, "white", xrdbColors)

		statusBarBackround := cfg.Section("style").Key("status-bar-background").String()
		style.StatusBarBackground = parseColor(statusBarBackround, "#f92672", xrdbColors)

		statusBarText := cfg.Section("style").Key("status-bar-text").String()
		style.StatusBarText = parseColor(statusBarText, "white", xrdbColors)

		statusBarViewBackround := cfg.Section("style").Key("status-bar-view-background").String()
		style.StatusBarViewBackground = parseColor(statusBarViewBackround, "#ae81ff", xrdbColors)

		statusBarViewText := cfg.Section("style").Key("status-bar-view-text").String()
		style.StatusBarViewText = parseColor(statusBarViewText, "white", xrdbColors)

		listSelectedBackground := cfg.Section("style").Key("list-selected-background").String()
		style.ListSelectedBackground = parseColor(listSelectedBackground, "#f92672", xrdbColors)

		listSelectedText := cfg.Section("style").Key("list-selected-text").String()
		style.ListSelectedText = parseColor(listSelectedText, "white", xrdbColors)
	}

	return style
}

func parseGeneral(cfg *ini.File) General {
	general := General{}

	general.Confirmation = cfg.Section("general").Key("confirmation").MustBool(true)
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

	tl := cfg.Section("general").Key("timeline").In("home", []string{"home", "direct", "local", "federated"})
	switch tl {
	case "direct":
		general.StartTimeline = feed.Conversations
	case "local":
		general.StartTimeline = feed.TimelineLocal
	case "federated":
		general.StartTimeline = feed.TimelineFederated
	default:
		general.StartTimeline = feed.TimelineHome
	}

	general.NotificationFeed = cfg.Section("general").Key("notification-feed").MustBool(true)
	general.QuoteReply = cfg.Section("general").Key("quote-reply").MustBool(false)
	general.CharLimit = cfg.Section("general").Key("char-limit").MustInt(500)
	general.MaxWidth = cfg.Section("general").Key("max-width").MustInt(0)
	general.ShortHints = cfg.Section("general").Key("short-hints").MustBool(false)
	general.HideNotificationText = cfg.Section("general").Key("hide-notification-text").MustBool(false)
	general.ShowIcons = cfg.Section("general").Key("show-icons").MustBool(true)
	general.ShowHelp = cfg.Section("general").Key("show-help").MustBool(true)
	general.RedrawUI = cfg.Section("general").Key("redraw-ui").MustBool(true)

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
			if len(parts) != 2 {
				fmt.Printf("leader-action must consist of two parts seperated by a comma. Your value is: %s\n", strings.Join(parts, ","))
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
			case "compose":
				la.Command = LeaderCompose
			case "blocking":
				la.Command = LeaderBlocking
			case "bookmarks":
				la.Command = LeaderBookmarks
			case "saved":
				la.Command = LeaderSaved
			case "favorited":
				la.Command = LeaderFavorited
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
			case "profile":
				la.Command = LeaderProfile
			case "notifications":
				la.Command = LeaderNotifications
			case "lists":
				la.Command = LeaderLists
			case "tag":
				la.Command = LeaderTag
				la.Subaction = subaction
			default:
				fmt.Printf("leader-action %s is invalid\n", parts[0])
				os.Exit(1)
			}
			la.Shortcut = parts[1]
			las = append(las, la)
		}
		general.LeaderActions = las
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
	nc.NotificationBoost = cfg.Section("desktop-notification").Key("boost").MustBool(false)
	nc.NotificationPoll = cfg.Section("desktop-notification").Key("poll").MustBool(false)
	nc.NotificationPost = cfg.Section("desktop-notification").Key("posts").MustBool(false)
	return nc
}

func parseTemplates(cfg *ini.File) Templates {
	var tootTmpl *template.Template
	tootTmplPath, exists, err := checkConfig("toot.tmpl")
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
	userTmplPath, exists, err := checkConfig("user.tmpl")
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

		MainHome:              inputStrOrErr([]string{"\"\"", "'g'", "\"Home\""}, false),
		MainEnd:               inputStrOrErr([]string{"\"\"", "'G'", "\"End\""}, false),
		MainPrevFeed:          inputStrOrErr([]string{"\"\"", "'h'", "'H'", "\"Left\""}, false),
		MainNextFeed:          inputStrOrErr([]string{"\"\"", "'l'", "'L'", "\"Right\""}, false),
		MainNotificationFocus: inputStrOrErr([]string{"\"[N]otifications\"", "'n'", "'N'"}, false),
		MainCompose:           inputStrOrErr([]string{"\"\"", "'c'", "'C'"}, false),

		StatusAvatar:        inputStrOrErr([]string{"\"[A]vatar\"", "'a'", "'A'"}, false),
		StatusBoost:         inputStrOrErr([]string{"\"[B]oost\"", "\"Un[B]oost\"", "'b'", "'B'"}, true),
		StatusDelete:        inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
		StatusFavorite:      inputStrOrErr([]string{"\"[F]avorite\"", "\"Un[F]avorite\"", "'f'", "'F'"}, true),
		StatusMedia:         inputStrOrErr([]string{"\"[M]edia\"", "'m'", "'M'"}, false),
		StatusLinks:         inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		StatusPoll:          inputStrOrErr([]string{"\"[P]oll\"", "'p'", "'P'"}, false),
		StatusReply:         inputStrOrErr([]string{"\"[R]eply\"", "'r'", "'R'"}, false),
		StatusBookmark:      inputStrOrErr([]string{"\"[S]ave\"", "\"Un[S]ave\"", "'s'", "'S'"}, true),
		StatusThread:        inputStrOrErr([]string{"\"[T]hread\"", "'t'", "'T'"}, false),
		StatusUser:          inputStrOrErr([]string{"\"[U]ser\"", "'u'", "'U'"}, false),
		StatusViewFocus:     inputStrOrErr([]string{"\"[V]iew\"", "'v'", "'V'"}, false),
		StatusYank:          inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),
		StatusToggleSpoiler: inputStrOrErr([]string{"\"Press [Z] to toggle spoiler\"", "'z'", "'Z'"}, false),

		UserAvatar:    inputStrOrErr([]string{"\"[A]vatar\"", "'a'", "'A'"}, false),
		UserBlock:     inputStrOrErr([]string{"\"[B]lock\"", "\"Un[B]lock\"", "'b'", "'B'"}, true),
		UserFollow:    inputStrOrErr([]string{"\"[F]ollow\"", "\"Un[F]ollow\"", "'f'", "'F'"}, true),
		UserMute:      inputStrOrErr([]string{"\"[M]ute\"", "\"Un[M]ute\"", "'m'", "'M'"}, true),
		UserLinks:     inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		UserUser:      inputStrOrErr([]string{"\"[U]ser\"", "'u'", "'U'"}, false),
		UserViewFocus: inputStrOrErr([]string{"\"[V]iew\"", "'v'", "'V'"}, false),
		UserYank:      inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),

		ListOpenFeed: inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),

		LinkOpen: inputStrOrErr([]string{"\"[O]pen\"", "'o'", "'O'"}, false),
		LinkYank: inputStrOrErr([]string{"\"[Y]ank\"", "'y'", "'Y'"}, false),

		ComposeEditSpoiler:          inputStrOrErr([]string{"\"[C]W Text\"", "'c'", "'C'"}, false),
		ComposeEditText:             inputStrOrErr([]string{"\"[E]dit text\"", "'e'", "'E'"}, false),
		ComposeIncludeQuote:         inputStrOrErr([]string{"\"[I]nclude quote\"", "'i'", "'I'"}, false),
		ComposeMediaFocus:           inputStrOrErr([]string{"\"[M]edia\"", "'m'", "'M'"}, false),
		ComposePost:                 inputStrOrErr([]string{"\"[P]ost\"", "'p'", "'P'"}, false),
		ComposeToggleContentWarning: inputStrOrErr([]string{"\"[T]oggle CW\"", "'t'", "'T'"}, false),
		ComposeVisibility:           inputStrOrErr([]string{"\"[V]isibility\"", "'v'", "'V'"}, false),

		MediaDelete:   inputStrOrErr([]string{"\"[D]elete\"", "'d'", "'D'"}, false),
		MediaEditDesc: inputStrOrErr([]string{"\"[E]dit desc\"", "'e'", "'E'"}, false),
		MediaAdd:      inputStrOrErr([]string{"\"[A]dd\"", "'a'", "'A'"}, false),

		VoteVote:   inputStrOrErr([]string{"\"[V]ote\"", "'v'", "'V'"}, false),
		VoteSelect: inputStrOrErr([]string{"\"[Enter] to select\"", "' '", "\"Enter\""}, false),
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
	ic.MainNotificationFocus = inputOrErr(cfg, "main-notification-focus", false, ic.MainNotificationFocus)
	ic.MainCompose = inputOrErr(cfg, "main-compose", false, ic.MainCompose)

	ic.StatusAvatar = inputOrErr(cfg, "status-avatar", false, ic.StatusAvatar)
	ic.StatusBoost = inputOrErr(cfg, "status-boost", true, ic.StatusBoost)
	ic.StatusDelete = inputOrErr(cfg, "status-delete", false, ic.StatusDelete)
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
	ic.StatusToggleSpoiler = inputOrErr(cfg, "status-toggle-spoiler", false, ic.StatusToggleSpoiler)

	ic.UserAvatar = inputOrErr(cfg, "user-avatar", false, ic.UserAvatar)
	ic.UserBlock = inputOrErr(cfg, "user-block", true, ic.UserBlock)
	ic.UserFollow = inputOrErr(cfg, "user-follow", true, ic.UserFollow)
	ic.UserMute = inputOrErr(cfg, "user-mute", true, ic.UserMute)
	ic.UserLinks = inputOrErr(cfg, "user-links", false, ic.UserLinks)
	ic.UserUser = inputOrErr(cfg, "user-user", false, ic.UserUser)
	ic.UserViewFocus = inputOrErr(cfg, "user-view-focus", false, ic.UserViewFocus)
	ic.UserYank = inputOrErr(cfg, "user-yank", false, ic.UserYank)

	ic.ListOpenFeed = inputOrErr(cfg, "list-open-feed", false, ic.ListOpenFeed)

	ic.LinkOpen = inputOrErr(cfg, "link-open", false, ic.LinkOpen)
	ic.LinkYank = inputOrErr(cfg, "link-yank", false, ic.LinkYank)

	ic.ComposeEditSpoiler = inputOrErr(cfg, "compose-edit-spoiler", false, ic.ComposeEditSpoiler)
	ic.ComposeEditText = inputOrErr(cfg, "compose-edit-text", false, ic.ComposeEditText)
	ic.ComposeIncludeQuote = inputOrErr(cfg, "compose-include-quote", false, ic.ComposeIncludeQuote)
	ic.ComposeMediaFocus = inputOrErr(cfg, "compose-media-focus", false, ic.ComposeMediaFocus)
	ic.ComposePost = inputOrErr(cfg, "compose-post", false, ic.ComposePost)
	ic.ComposeToggleContentWarning = inputOrErr(cfg, "compose-toggle-content-warning", false, ic.ComposeToggleContentWarning)
	ic.ComposeVisibility = inputOrErr(cfg, "compose-visibility", false, ic.ComposeVisibility)

	ic.MediaDelete = inputOrErr(cfg, "media-delete", false, ic.MediaDelete)
	ic.MediaEditDesc = inputOrErr(cfg, "media-edit-desc", false, ic.MediaEditDesc)
	ic.MediaAdd = inputOrErr(cfg, "media-add", false, ic.MediaAdd)

	ic.VoteVote = inputOrErr(cfg, "vote-vote", false, ic.VoteVote)
	ic.VoteSelect = inputOrErr(cfg, "vote-select", false, ic.VoteSelect)
	return ic
}

func parseConfig(filepath string) (Config, error) {
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
	conf.Style = parseStyle(cfg)
	conf.OpenPattern = parseOpenPattern(cfg)
	conf.OpenCustom = parseCustom(cfg)
	conf.NotificationConfig = parseNotifications(cfg)
	conf.Templates = parseTemplates(cfg)
	conf.Input = parseInput(cfg)

	return conf, nil
}

func createConfigDir() error {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't find $HOME. Err %v", err)
	}
	path := cd + "/tut"
	return os.MkdirAll(path, os.ModePerm)
}

func checkConfig(filename string) (path string, exists bool, err error) {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't find $HOME. Err %v", err)
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

func getThemes() ([]string, error) {
	entries, err := themesFS.ReadDir("themes")
	files := []string{}
	if err != nil {
		return []string{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fp := filepath.Join("themes/", entry.Name())
		files = append(files, fp)
	}
	return files, nil
}

func getTheme(fname string) (*ini.File, error) {
	f, err := themesFS.Open(fmt.Sprintf("themes/%s.ini", strings.TrimSpace(fname)))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
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
