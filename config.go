package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

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

type Config struct {
	General            GeneralConfig
	Style              StyleConfig
	Media              MediaConfig
	OpenPattern        OpenPatternConfig
	OpenCustom         OpenCustomConfig
	NotificationConfig NotificationConfig
	Templates          TemplatesConfig
}

type GeneralConfig struct {
	AutoLoadNewer        bool
	AutoLoadSeconds      int
	DateTodayFormat      string
	DateFormat           string
	DateRelative         int
	MaxWidth             int
	StartTimeline        TimelineType
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
}

type StyleConfig struct {
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

type MediaConfig struct {
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

type OpenPatternConfig struct {
	Patterns []Pattern
}

type OpenCustom struct {
	Index    int
	Name     string
	Program  string
	Args     []string
	Terminal bool
}
type OpenCustomConfig struct {
	OpenCustoms []OpenCustom
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
	NotificationFollower = iota
	NotificationFavorite
	NotificationMention
	NotificationBoost
	NotificationPoll
	NotificationPost
)

type NotificationConfig struct {
	NotificationFollower bool
	NotificationFavorite bool
	NotificationMention  bool
	NotificationBoost    bool
	NotificationPoll     bool
	NotificationPost     bool
}

type TemplatesConfig struct {
	TootTemplate *template.Template
	UserTemplate *template.Template
	HelpTemplate *template.Template
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

func parseStyle(cfg *ini.File) StyleConfig {
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

	style := StyleConfig{}

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

	return style
}

func parseGeneral(cfg *ini.File) GeneralConfig {
	general := GeneralConfig{}

	general.AutoLoadNewer = cfg.Section("media").Key("auto-load-newer").MustBool(true)
	autoLoadSeconds, err := cfg.Section("general").Key("auto-load-seconds").Int()
	if err != nil {
		autoLoadSeconds = 60
	}
	general.AutoLoadSeconds = autoLoadSeconds

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
	case "home":
		general.StartTimeline = TimelineHome
	case "direct":
		general.StartTimeline = TimelineDirect
	case "local":
		general.StartTimeline = TimelineLocal
	case "federated":
		general.StartTimeline = TimelineFederated
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

	return general
}

func parseMedia(cfg *ini.File) MediaConfig {
	media := MediaConfig{}
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

func ParseOpenPattern(cfg *ini.File) OpenPatternConfig {
	om := OpenPatternConfig{}

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

func ParseCustom(cfg *ini.File) OpenCustomConfig {
	oc := OpenCustomConfig{}

	for i := 1; i < 6; i++ {
		name := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-name", i)).MustString("")
		use := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-use", i)).MustString("")
		terminal := cfg.Section("open-custom").Key(fmt.Sprintf("c%d-terminal", i)).MustBool(false)
		if use == "" {
			continue
		}
		comp := strings.Fields(use)
		c := OpenCustom{}
		c.Index = i
		c.Name = name
		c.Program = comp[0]
		c.Args = comp[1:]
		c.Terminal = terminal
		oc.OpenCustoms = append(oc.OpenCustoms, c)
	}
	return oc
}

func ParseNotifications(cfg *ini.File) NotificationConfig {
	nc := NotificationConfig{}
	nc.NotificationFollower = cfg.Section("desktop-notification").Key("followers").MustBool(false)
	nc.NotificationFavorite = cfg.Section("desktop-notification").Key("favorite").MustBool(false)
	nc.NotificationMention = cfg.Section("desktop-notification").Key("mention").MustBool(false)
	nc.NotificationBoost = cfg.Section("desktop-notification").Key("boost").MustBool(false)
	nc.NotificationPoll = cfg.Section("desktop-notification").Key("poll").MustBool(false)
	nc.NotificationPost = cfg.Section("desktop-notification").Key("posts").MustBool(false)
	return nc
}

func ParseTemplates(cfg *ini.File) TemplatesConfig {
	var tootTmpl *template.Template
	tootTmplPath, exists, err := CheckConfig("toot.tmpl")
	if err != nil {
		log.Fatalln(
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
	userTmplPath, exists, err := CheckConfig("user.tmpl")
	if err != nil {
		log.Fatalln(
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
	return TemplatesConfig{
		TootTemplate: tootTmpl,
		UserTemplate: userTmpl,
		HelpTemplate: helpTmpl,
	}
}

func ParseConfig(filepath string) (Config, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}, filepath)
	conf := Config{}
	if err != nil {
		return conf, err
	}
	conf.General = parseGeneral(cfg)
	conf.Media = parseMedia(cfg)
	conf.Style = parseStyle(cfg)
	conf.OpenPattern = ParseOpenPattern(cfg)
	conf.OpenCustom = ParseCustom(cfg)
	conf.NotificationConfig = ParseNotifications(cfg)
	conf.Templates = ParseTemplates(cfg)

	return conf, nil
}

func CreateConfigDir() error {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't find $HOME. Err %v", err)
	}
	path := cd + "/tut"
	return os.MkdirAll(path, os.ModePerm)
}

func CheckConfig(filename string) (path string, exists bool, err error) {
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
