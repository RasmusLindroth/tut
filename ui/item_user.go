package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

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
	Style config.Style
}

func drawUser(tut *Tut, data *api.User, main *tview.TextView, controls *tview.TextView, additional string) {
	user := data.Data
	relation := data.Relation
	showUserControl := true
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
		Bot:            user.Bot,
	}

	var controlsS string

	var urls []util.URL
	fields := []Field{}
	u.Note, urls = util.CleanHTML(user.Note)
	for _, f := range user.Fields {
		value, fu := util.CleanHTML(f.Value)
		fields = append(fields, Field{
			Name:       tview.Escape(f.Name),
			Value:      tview.Escape(value),
			VerifiedAt: f.VerifiedAt,
		})
		urls = append(urls, fu...)
	}
	u.Fields = fields

	var controlItems []string
	if tut.Client.Me.ID != user.ID {
		if relation.Following {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserFollow, false))
		} else {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserFollow, true))
		}
		if relation.Blocking {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserBlock, false))
		} else {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserBlock, true))
		}
		if relation.Muting {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserMute, false))
		} else {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserMute, true))
		}
		if len(urls) > 0 {
			controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserLinks, true))
		}
	}
	if showUserControl {
		controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserUser, true))
	}
	controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserAvatar, true))
	controlItems = append(controlItems, config.ColorFromKey(tut.Config, tut.Config.Input.UserYank, true))
	controlsS = strings.Join(controlItems, " ")

	ud := DisplayUserData{
		User:  u,
		Style: tut.Config.Style,
	}
	var output bytes.Buffer
	err := tut.Config.Templates.User.ExecuteTemplate(&output, "user.tmpl", ud)
	if err != nil {
		panic(err)
	}

	if main != nil {
		if additional != "" {
			additional = fmt.Sprintf("%s\n\n", config.SublteText(tut.Config, additional))
		}
		main.SetText(additional + output.String())
	}
	controls.SetText(controlsS)
}
