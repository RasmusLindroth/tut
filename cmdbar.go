package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewCmdBar(app *App) *CmdBar {
	c := &CmdBar{
		app:   app,
		Input: tview.NewInputField(),
	}

	c.Input.SetBackgroundColor(app.Config.Style.Background)
	c.Input.SetFieldBackgroundColor(app.Config.Style.Background)
	c.Input.SetFieldTextColor(app.Config.Style.Text)
	c.Input.SetDoneFunc(c.DoneFunc)

	return c
}

type CmdBar struct {
	app   *App
	Input *tview.InputField
}

func (c *CmdBar) GetInput() string {
	return strings.TrimSpace(c.Input.GetText())
}

func (c *CmdBar) ShowError(s string) {
	c.Input.SetFieldTextColor(c.app.Config.Style.WarningText)
	c.Input.SetText(s)
}

func (c *CmdBar) ShowMsg(s string) {
	c.Input.SetFieldTextColor(c.app.Config.Style.StatusBarText)
	c.Input.SetText(s)
}

func (c *CmdBar) ClearInput() {
	c.Input.SetFieldTextColor(c.app.Config.Style.Text)
	c.Input.SetText("")
	c.Input.Autocomplete()
}

func (c *CmdBar) DoneFunc(key tcell.Key) {
	input := c.GetInput()
	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case ":q":
		fallthrough
	case ":quit":
		c.app.UI.Root.Stop()
	case ":compose":
		c.app.UI.NewToot()
		c.app.UI.CmdBar.ClearInput()
	case ":blocking":
		c.app.UI.StatusView.AddFeed(NewUserListFeed(c.app, UserListBlocking, ""))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":bookmarks", ":saved":
		c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineBookmarked, nil))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":favorited":
		c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineFavorited, nil))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":boosts":
		c.app.UI.CmdBar.ClearInput()
		status := c.app.UI.StatusView.GetCurrentStatus()
		if status == nil {
			return
		}

		if status.Reblog != nil {
			status = status.Reblog
		}
		c.app.UI.StatusView.AddFeed(NewUserListFeed(c.app, UserListBoosts, string(status.ID)))
		c.app.UI.SetFocus(LeftPaneFocus)
	case ":favorites":
		c.app.UI.CmdBar.ClearInput()
		status := c.app.UI.StatusView.GetCurrentStatus()
		if status == nil {
			return
		}
		if status.Reblog != nil {
			status = status.Reblog
		}
		c.app.UI.StatusView.AddFeed(NewUserListFeed(c.app, UserListFavorites, string(status.ID)))
		c.app.UI.SetFocus(LeftPaneFocus)
	/*
		case ":followers":
			app.UI.CmdBar.ClearInput()
			user := app.UI.StatusView.GetCurrentUser()
			if user == nil {
				return
			}
			app.UI.StatusView.AddFeed(NewUserListFeed(app, UserListFollowers, string(user.ID)))
			app.UI.SetFocus(LeftPaneFocus)
		case ":following":
			app.UI.CmdBar.ClearInput()
			user := app.UI.StatusView.GetCurrentUser()
			if user == nil {
				return
			}
			app.UI.StatusView.AddFeed(NewUserListFeed(app, UserListFollowing, string(user.ID)))
			app.UI.SetFocus(LeftPaneFocus)
	*/
	case ":muting":
		c.app.UI.StatusView.AddFeed(NewUserListFeed(c.app, UserListMuting, ""))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":profile":
		c.app.UI.CmdBar.ClearInput()
		if c.app.Me == nil {
			return
		}
		c.app.UI.StatusView.AddFeed(NewUserFeed(c.app, *c.app.Me))
		c.app.UI.SetFocus(LeftPaneFocus)
	case ":timeline", ":tl":
		if len(parts) < 2 {
			break
		}
		switch parts[1] {
		case "local", "l":
			c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineLocal, nil))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		case "federated", "f":
			c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineFederated, nil))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		case "direct", "d":
			c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineDirect, nil))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		case "home", "h":
			c.app.UI.StatusView.AddFeed(NewTimelineFeed(c.app, TimelineHome, nil))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		case "notifications", "n":
			c.app.UI.StatusView.AddFeed(NewNotificationFeed(c.app, false))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		case "favrotied", "fav":
			c.app.UI.StatusView.AddFeed(NewNotificationFeed(c.app, false))
			c.app.UI.SetFocus(LeftPaneFocus)
			c.app.UI.CmdBar.ClearInput()
		}
	case ":tag":
		if len(parts) < 2 {
			break
		}
		tag := strings.TrimSpace(strings.TrimPrefix(parts[1], "#"))
		if len(tag) == 0 {
			break
		}
		c.app.UI.StatusView.AddFeed(NewTagFeed(c.app, tag))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":user":
		if len(parts) < 2 {
			break
		}
		user := strings.TrimSpace(parts[1])
		if len(user) == 0 {
			break
		}
		c.app.UI.StatusView.AddFeed(NewUserListFeed(c.app, UserListSearch, user))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	case ":lists":
		c.app.UI.StatusView.AddFeed(NewListFeed(c.app))
		c.app.UI.SetFocus(LeftPaneFocus)
		c.app.UI.CmdBar.ClearInput()
	}
}
