package main

import (
	"os"

	"github.com/gdamore/tcell"
	"github.com/kyoh86/xdg"
)

type Config struct {
	General GeneralConfig
	Style   StyleConfig
	Media   MediaConfig
}

type GeneralConfig struct {
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
	ImageSingle bool
	VideoViewer string
	VideoSingle bool
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
