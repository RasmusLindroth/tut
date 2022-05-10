package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CmdBar struct {
	tutView *TutView
	View    *tview.InputField
}

func NewCmdBar(tv *TutView) *CmdBar {
	c := &CmdBar{
		tutView: tv,
		View:    NewInputField(tv.tut.Config),
	}
	c.View.SetAutocompleteFunc(c.Autocomplete)
	c.View.SetDoneFunc(c.DoneFunc)

	if tv.tut.Config.General.ShowHelp {
		c.ShowMsg("Press ? or :help to learn how tut functions")
	}

	return c
}

func (c *CmdBar) GetInput() string {
	return strings.TrimSpace(c.View.GetText())
}

func (c *CmdBar) ShowError(s string) {
	c.View.SetFieldTextColor(c.tutView.tut.Config.Style.WarningText)
	c.View.SetText(s)
}

func (c *CmdBar) ShowMsg(s string) {
	c.View.SetFieldTextColor(c.tutView.tut.Config.Style.StatusBarText)
	c.View.SetText(s)
}

func (c *CmdBar) ClearInput() {
	c.View.SetFieldTextColor(c.tutView.tut.Config.Style.StatusBarText)
	c.View.SetText("")
}

func (c *CmdBar) Back() {
	c.ClearInput()
	c.View.Autocomplete()
	c.tutView.PrevFocus()
}

func (c *CmdBar) DoneFunc(key tcell.Key) {
	if key == tcell.KeyTAB {
		return
	}
	input := c.GetInput()
	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case ":q":
		fallthrough
	case ":quit":
		c.tutView.tut.App.Stop()
	case ":compose":
		c.tutView.ComposeCommand()
		c.ClearInput()
		c.View.Autocomplete()
	case ":blocking":
		c.tutView.BlockingCommand()
		c.Back()
	case ":bookmarks", ":saved":
		c.tutView.BookmarksCommand()
		c.Back()
	case ":favorited":
		c.tutView.FavoritedCommand()
		c.Back()
	case ":boosts":
		c.tutView.BoostsCommand()
		c.Back()
	case ":favorites":
		c.tutView.FavoritesCommand()
		c.Back()
	case ":following":
		c.tutView.FollowingCommand()
		c.Back()
	case ":followers":
		c.tutView.FollowersCommand()
		c.Back()
	case ":muting":
		c.tutView.MutingCommand()
		c.Back()
	case ":profile":
		c.tutView.ProfileCommand()
		c.Back()
	case ":timeline", ":tl":
		if len(parts) < 2 {
			break
		}
		switch parts[1] {
		case "local", "l":
			c.tutView.LocalCommand()
			c.Back()
		case "federated", "f":
			c.tutView.FederatedCommand()
			c.Back()
		case "direct", "d":
			c.tutView.DirectCommand()
			c.Back()
		case "home", "h":
			c.tutView.HomeCommand()
			c.Back()
		case "notifications", "n":
			c.tutView.NotificationsCommand()
			c.Back()
		case "favorited", "fav":
			c.tutView.FavoritedCommand()
			c.Back()
		}
		c.ClearInput()
	case ":tag":
		if len(parts) < 2 {
			break
		}
		tag := strings.TrimSpace(strings.TrimPrefix(parts[1], "#"))
		if len(tag) == 0 {
			break
		}
		c.tutView.Timeline.AddFeed(
			NewTagFeed(c.tutView, tag),
		)
		c.Back()
	case ":user":
		if len(parts) < 2 {
			break
		}
		user := strings.TrimSpace(parts[1])
		if len(user) == 0 {
			break
		}
		c.tutView.Timeline.AddFeed(
			NewUserSearchFeed(c.tutView, user),
		)
		c.Back()
	case ":lists":
		c.tutView.ListsCommand()
		c.Back()
	case ":help", ":h":
		c.tutView.PageFocus = c.tutView.PrevPageFocus
		c.tutView.SetPage(HelpFocus)
		c.ClearInput()
		c.View.Autocomplete()
	}
}

func (c *CmdBar) Autocomplete(curr string) []string {
	var entries []string
	words := strings.Split(":blocking,:boosts,:bookmarks,:compose,:favorites,:favorited,:followers,:following,:help,:h,:lists,:muting,:profile,:saved,:tag,:timeline,:tl,:user,:quit,:q", ",")
	if curr == "" {
		return entries
	}

	if len(curr) > 2 && curr[:3] == ":tl" {
		words = strings.Split(":tl home,:tl notifications,:tl local,:tl federated,:tl direct,:tl favorited", ",")
	}
	if len(curr) > 8 && curr[:9] == ":timeline" {
		words = strings.Split(":timeline home,:timeline notifications,:timeline local,:timeline federated,:timeline direct,:timeline favorited", ",")
	}

	for _, word := range words {
		if strings.HasPrefix(strings.ToLower(word), strings.ToLower(curr)) {
			entries = append(entries, word)
		}
	}
	if len(entries) < 1 {
		entries = nil
	}
	return entries
}
