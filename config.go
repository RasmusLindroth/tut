package main

import (
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/kyoh86/xdg"
	"gopkg.in/ini.v1"
)

type Config struct {
	General GeneralConfig
	Style   StyleConfig
	Media   MediaConfig
}

type GeneralConfig struct {
	AutoLoadNewer   bool
	AutoLoadSeconds int
	DateTodayFormat string
	DateFormat      string
	StartTimeline   TimelineType
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

	ListSelectedBackground tcell.Color
	ListSelectedText       tcell.Color
}

type MediaConfig struct {
	ImageViewer string
	ImageArgs   []string
	ImageSingle bool
	VideoViewer string
	VideoArgs   []string
	VideoSingle bool
	AudioViewer string
	AudioArgs   []string
	AudioSingle bool
	LinkViewer  string
	LinkArgs    []string
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

	videoViewerComponents := strings.Fields(cfg.Section("media").Key("video-viewer").String())
	if len(videoViewerComponents) == 0 {
		media.VideoViewer = "xdg-open"
		media.VideoArgs = []string{}
	} else {
		media.VideoViewer = videoViewerComponents[0]
		media.VideoArgs = videoViewerComponents[1:]
	}
	media.VideoSingle = cfg.Section("media").Key("video-single").MustBool(true)

	audioViewerComponents := strings.Fields(cfg.Section("media").Key("audio-viewer").String())
	if len(audioViewerComponents) == 0 {
		media.AudioViewer = "xdg-open"
		media.AudioArgs = []string{}
	} else {
		media.AudioViewer = audioViewerComponents[0]
		media.AudioArgs = audioViewerComponents[1:]
	}
	media.AudioSingle = cfg.Section("media").Key("audio-single").MustBool(true)

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

# The timeline that opens up when you start tut
# Valid values: home, direct, local, federated
# default=home
timeline=home

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

# Your video viewer
# default=xdg-open
video-viewer=xdg-open

# If videos should open one by one. See above comment about image-single
# default=true
video-single=true

# Your audio viewer
# default=xdg-open
audio-viewer=xdg-open

# If audio files should open one by one. See above comment about image-single
# default=true
audio-single=true

# Your web browser
# default=xdg-open
link-viewer=xdg-open

[style]
# All styles can be represented in their HEX value like #ffffff or
# with their name, so in this case white.
# The only special value is "default" which equals to transparent,
# so it will be the same color as your terminal. But this can lead
# to some artifacts left from a previous paint

# You can also use xrdb colors like this xrdb:color1
# The program will use colors prefixed with an * first then look
# for URxvt or XTerm if it can't find any color prefixed with an asterix.
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
# default=xrdb:color0
status-bar-background=xrdb:color0

# The color of the text in the bar at the bottom
# default=xrdb:foreground
status-bar-text=xrdb:foreground

# Background of selected list items
# default=xrdb:color5
list-selected-background=xrdb:color5

# The text color of selected list items
# default=xrdb:background
list-selected-text=xrdb:background
`
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(conf)
	if err != nil {
		return err
	}
	return nil
}
