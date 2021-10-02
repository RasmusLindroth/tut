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
	"github.com/kyoh86/xdg"
	"gopkg.in/ini.v1"
)

//go:embed toot.tmpl
var tootTemplate string

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
	ImageViewer  string
	ImageArgs    []string
	ImageSingle  bool
	ImageReverse bool
	VideoViewer  string
	VideoArgs    []string
	VideoSingle  bool
	VideoReverse bool
	AudioViewer  string
	AudioArgs    []string
	AudioSingle  bool
	AudioReverse bool
	LinkViewer   string
	LinkArgs     []string
}

type Pattern struct {
	Pattern  string
	Open     string
	Compiled glob.Glob
	Program  string
	Args     []string
}

type OpenPatternConfig struct {
	Patterns []Pattern
}

type OpenCustom struct {
	Index   int
	Name    string
	Program string
	Args    []string
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

	return media
}

func ParseOpenPattern(cfg *ini.File) OpenPatternConfig {
	om := OpenPatternConfig{}

	keys := cfg.Section("open-pattern").KeyStrings()
	pairs := make(map[string]Pattern)
	for _, s := range keys {
		parts := strings.Split(s, "-")
		if len(parts) < 2 {
			panic(fmt.Sprintf("Invalid key %s in config. Must end in -pattern or -use", s))
		}
		last := parts[len(parts)-1]
		if last != "pattern" && last != "use" {
			panic(fmt.Sprintf("Invalid key %s in config. Must end in -pattern or -use", s))
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
		if use == "" {
			continue
		}
		comp := strings.Fields(use)
		c := OpenCustom{}
		c.Index = i
		c.Name = name
		c.Program = comp[0]
		c.Args = comp[1:]
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
	return TemplatesConfig{
		TootTemplate: tootTmpl,
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
	path := xdg.ConfigHome() + "/tut"
	return os.MkdirAll(path, os.ModePerm)
}

func CheckConfig(filename string) (path string, exists bool, err error) {
	dir := xdg.ConfigHome() + "/tut/"
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
	conf := `#
# Configuration file for tut

[general]
# If the program should check for new toots without user interaction.
# If you don't enable this the program will only look for new toots when 
# you reach the bottom or top of your feed. With this enabled it will check
# for new toots every x second.
# default=true
auto-load-newer=true

# How many seconds between each pulling of new toots if you have enabled
# auto-load-newer.
# default=60
auto-load-seconds=60

# The date format to be used
# See https://godoc.org/time#Time.Format
# default=2006-01-02 15:04
date-format=2006-01-02 15:04

# Format for dates the same day
# default=15:04
date-today-format=15:04

# This displays relative dates instead
# for statuses that are one day or older
# the output is 1y2m1d (1 year 2 months and 1 day) 
#
# The value is an integear
# -1     = don't use relative dates
#  0     = always use relative dates, except for dates < 1 day
#  1 - ∞ = number of days to use relative dates
#
# Example: date-relative=28 will display a relative date
# for toots that are between 1-28 days old. Otherwhise it
# will use the short or long format
# 
# default=-1
date-relative=-1

# The timeline that opens up when you start tut
# Valid values: home, direct, local, federated
# default=home
timeline=home

# The max width of text before it wraps when displaying toots
# 0 = no restriction
# default=0
max-width=0

# If you want to display a list of notifications
# under your timeline feed
# default=true
notification-feed=true

# Where do you want the list of toots to be placed
# Valid values: left, right, top, bottom
# default=left
list-placement=left

# If you have notification-feed set to true you can
# display it under the main list of toots (row)
# or place it to the right of the main list of toots (column)
# default=row
list-split=row

# Hide notification text above list in column split
# default=false
hide-notification-text=false

# You can change the proportions of the list view
# in relation to the content view
# list-proportion=1 and content-proportoin=3
# will result in the content taking up 3 times more space
# Must be n > 0
# defaults:
# 	list-proportion=1
# 	content-proportion=2
list-proportion=1
content-proportion=2

# If you always want to quote original message when replying
# default=false
quote-reply=false

# If you're on an instance with a custom character limit you can set it here 
# default=500
char-limit=500

# If you want to show icons in the list of toots
# default=true
show-icons=true

# If you've learnt all the shortcut keys you can remove the help text and 
# only show the key in tui. So it gets less cluttered.
# default=false
short-hints=false

[media]
# Your image viewer
# default=xdg-open
image-viewer=xdg-open

# If image should open one by one e.g. "imv image.png" multiple times
# If set to false all images will open at the same time like this 
# "imv image1.png image2.png image3.png".
# Not all image viewers support this, so try it first.
# default=true
image-single=true

# If you want to open the images in reverse order. In some image viewers 
# this will display the images in the "right" order.
# default=false
image-reverse=false

# Your video viewer
# default=xdg-open
video-viewer=xdg-open

# If videos should open one by one. See above comment about image-single
# default=true
video-single=true

# If you want to open the videos in reverse order. In some video apps 
# this will play the files in the "right" order.
# default=false
video-reverse=false

# Your audio viewer
# default=xdg-open
audio-viewer=xdg-open

# If you want to play the audio files in reverse order. In some audio apps 
# this will play the files in the "right" order.
# default=false
audio-reverse=false

# If audio files should open one by one. See above comment about image-single
# default=true
audio-single=true

# Your web browser
# default=xdg-open
link-viewer=xdg-open

[open-custom]
# This sections allows you to set up to five custom programs to upen URLs with.
# If the url points to an image, you can set c1-name to img and c1-use to imv.
# The name will show up in the UI, so keep it short so all five fits.
#
# c1-name=img
# c1-use=imv
# 
# c2-name=
# c2-use=
# 
# c3-name=
# c3-use=
# 
# c4-name=
# c4-use=
# 
# c5-name=
# c5-use=

[open-pattern]
# Here you can set your own glob patterns for opening matching URLs in the
# program you want them to open up in.
# You could for example open Youtube videos in your video player instead of
# your default browser.
#
# You must name the keys foo-pattern and foo-use, where use is the program 
# that will open up the URL. To see the syntax for glob pattern you can follow
# this URL https://github.com/gobwas/glob#syntax
#
# Example for youtube.com and youtu.be to open up in mpv instead of the browser
#
# y1-pattern=*youtube.com/watch*
# y1-use=mpv
# y2-pattern=*youtu.be/*
# y2-use=mpv

[desktop-notification]
# Under this section you can turn on desktop notifications

# Notification when someone follows you
# default=false
followers=false

# Notification when someone favorites one of your toots
# default=false
favorite=false

# Notification when someone mentions you
# default=false
mention=false

# Notification when someone boosts one of your toots
# default=false
boost=false

# Notification of poll results
# default=false
poll=false

# New posts in current timeline
# default=false
posts=false

[style]
# All styles can be represented in their HEX value like #ffffff or
# with their name, so in this case white.
# The only special value is "default" which equals to transparent,
# so it will be the same color as your terminal. But this can lead
# to some artifacts left from a previous paint

# You can also use xrdb colors like this xrdb:color1
# The program will use colors prefixed with an * first then look
# for URxvt or XTerm if it can't find any color prefixed with an asterik.
# If you don't want tut to guess the prefix you can set the prefix yourself.
# If the xrdb color can't be found a preset color will be used.

# The xrdb prefix used for colors in .Xresources
# default=guess
xrdb-prefix=guess

# The background color used on most elements
# default=xrdb:background
background=xrdb:background

# The text color used on most of the text
# default=xrdb:foreground
text=xrdb:foreground

# The color to display sublte elements or subtle text. Like lines and help text
# default=xrdb:color14
subtle=xrdb:color14

# The color for errors or warnings
# default=xrdb:color1
warning-text=xrdb:color1

# This color is used to display username
# default=xrdb:color5
text-special-one=xrdb:color5

# This color is used to display username and keys
# default=xrdb:color2
text-special-two=xrdb:color2

# The color of the bar at the top
# default=xrdb:color5
top-bar-background=xrdb:color5

# The color of the text in the bar at the top
# default=xrdb:background
top-bar-text=xrdb:background

# The color of the bar at the bottom
# default=xrdb:color5
status-bar-background=xrdb:color5

# The color of the text in the bar at the bottom
# default=xrdb:foreground
status-bar-text=xrdb:foreground

# The color of the bar at the bottom in view mode
# default=xrdb:color4
status-bar-view-background=xrdb:color4

# The color of the text in the bar at the bottom in view mode
# default=xrdb:foreground
status-bar-view-text=xrdb:foreground

# Background of selected list items
# default=xrdb:color5
list-selected-background=xrdb:color5

# The text color of selected list items
# default=xrdb:background
list-selected-text=xrdb:background
`
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(conf)
	if err != nil {
		return err
	}
	return nil
}
