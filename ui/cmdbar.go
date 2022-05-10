package ui

import (
	"fmt"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/util"
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
	item, itemErr := c.tutView.GetCurrentItem()
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case ":q":
		fallthrough
	case ":quit":
		c.tutView.tut.App.Stop()
	case ":compose":
		c.tutView.InitPost(nil)
		c.ClearInput()
		c.View.Autocomplete()
	case ":blocking":
		c.tutView.Timeline.AddFeed(
			NewBlocking(c.tutView),
		)
		c.Back()
	case ":bookmarks", ":saved":
		c.tutView.Timeline.AddFeed(
			NewBookmarksFeed(c.tutView),
		)
		c.Back()
	case ":favorited":
		c.tutView.Timeline.AddFeed(
			NewFavoritedFeed(c.tutView),
		)
		c.Back()
	case ":boosts":
		if itemErr != nil {
			c.Back()
			return
		}
		if item.Type() != api.StatusType {
			c.Back()
			return
		}
		s := item.Raw().(*mastodon.Status)
		s = util.StatusOrReblog(s)
		c.tutView.Timeline.AddFeed(
			NewBoosts(c.tutView, s.ID),
		)
		c.Back()
	case ":favorites":
		if itemErr != nil {
			c.Back()
			return
		}
		if item.Type() != api.StatusType {
			c.Back()
			return
		}
		s := item.Raw().(*mastodon.Status)
		s = util.StatusOrReblog(s)
		c.tutView.Timeline.AddFeed(
			NewFavoritesStatus(c.tutView, s.ID),
		)
		c.Back()
	case ":following":
		if itemErr != nil {
			c.Back()
			return
		}
		if item.Type() != api.UserType && item.Type() != api.ProfileType {
			c.Back()
			return
		}
		s := item.Raw().(*api.User)
		c.tutView.Timeline.AddFeed(
			NewFollowing(c.tutView, s.Data.ID),
		)
		c.Back()
	case ":followers":
		if itemErr != nil {
			c.Back()
			return
		}
		if item.Type() != api.UserType && item.Type() != api.ProfileType {
			c.Back()
			return
		}
		s := item.Raw().(*api.User)
		c.tutView.Timeline.AddFeed(
			NewFollowers(c.tutView, s.Data.ID),
		)
		c.Back()
	case ":muting":
		c.tutView.Timeline.AddFeed(
			NewMuting(c.tutView),
		)
		c.Back()
	case ":profile":
		item, err := c.tutView.tut.Client.GetUserByID(c.tutView.tut.Client.Me.ID)
		if err != nil {
			c.ShowError(fmt.Sprintf("Couldn't load user. Error: %v\n", err))
			c.Back()
		}
		c.tutView.Timeline.AddFeed(
			NewUserFeed(c.tutView, item),
		)
		c.Back()
	case ":timeline", ":tl":
		if len(parts) < 2 {
			break
		}
		switch parts[1] {
		case "local", "l":
			c.tutView.Timeline.AddFeed(
				NewLocalFeed(c.tutView),
			)
			c.Back()
		case "federated", "f":
			c.tutView.Timeline.AddFeed(
				NewFederatedFeed(c.tutView),
			)
			c.Back()
		case "direct", "d":
			c.tutView.Timeline.AddFeed(
				NewConversationsFeed(c.tutView),
			)
			c.Back()
		case "home", "h":
			c.tutView.Timeline.AddFeed(
				NewHomeFeed(c.tutView),
			)
			c.Back()
		case "notifications", "n":
			c.tutView.Timeline.AddFeed(
				NewNotificationFeed(c.tutView),
			)
			c.Back()
		case "favorited", "fav":
			c.tutView.Timeline.AddFeed(
				NewFavoritedFeed(c.tutView),
			)
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
		c.tutView.Timeline.AddFeed(
			NewListsFeed(c.tutView),
		)
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
