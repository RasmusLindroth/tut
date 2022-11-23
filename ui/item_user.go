package ui

import (
	"bytes"
	"fmt"
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

func drawUser(tv *TutView, data *api.User, main *tview.TextView, controls *tview.Flex, additional string, ut InputUserType) {
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

	var controlItems []Control
	if ut == InputUserFollowRequest {
		controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserFollowRequestDecide, false))
	}
	if tv.tut.Client.Me.ID != user.ID {
		if relation.Following {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserFollow, false))
		} else {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserFollow, true))
		}
		if relation.Blocking {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserBlock, false))
		} else {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserBlock, true))
		}
		if relation.Muting {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserMute, false))
		} else {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserMute, true))
		}
		if len(urls) > 0 {
			controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserLinks, true))
		}
	}
	if showUserControl {
		controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserUser, true))
	}
	controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserAvatar, true))
	controlItems = append(controlItems, NewControl(tv.tut.Config, tv.tut.Config.Input.UserYank, true))

	// Clear controls and only have add and delete for lists.
	if ut == InputUserListAdd {
		controlItems = []Control{NewControl(tv.tut.Config, tv.tut.Config.Input.ListUserAdd, true)}
	} else if ut == InputUserListDelete {
		controlItems = []Control{NewControl(tv.tut.Config, tv.tut.Config.Input.ListUserDelete, true)}
	}

	controls.Clear()
	for i, item := range controlItems {
		if i < len(controlItems)-1 {
			controls.AddItem(NewControlButton(tv, item), item.Len+1, 0, false)
		} else {
			controls.AddItem(NewControlButton(tv, item), item.Len, 0, false)
		}
	}

	ud := DisplayUserData{
		User:  u,
		Style: tv.tut.Config.Style,
	}
	var output bytes.Buffer
	err := tv.tut.Config.Templates.User.ExecuteTemplate(&output, "user.tmpl", ud)
	if err != nil {
		panic(err)
	}

	if main != nil {
		if additional != "" {
			additional = fmt.Sprintf("%s\n\n", config.SublteText(tv.tut.Config, additional))
		}
		main.SetText(additional + output.String())
	}
}
