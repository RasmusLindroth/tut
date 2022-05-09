package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/util"
	"github.com/gdamore/tcell/v2"
)

func (tv *TutView) Input(event *tcell.EventKey) *tcell.EventKey {
	if tv.PageFocus != LoginFocus {
		switch event.Rune() {
		case ':':
			tv.SetPage(CmdFocus)
		case '?':
			tv.SetPage(HelpFocus)
		}
	}
	switch tv.PageFocus {
	case LoginFocus:
		return tv.InputLoginView(event)
	case MainFocus:
		return tv.InputMainView(event)
	case ViewFocus:
		return tv.InputViewItem(event)
	case ComposeFocus:
		return tv.InputComposeView(event)
	case LinkFocus:
		return tv.InputLinkView(event)
	case CmdFocus:
		return tv.InputCmdView(event)
	case MediaFocus:
		return tv.InputMedia(event)
	case MediaAddFocus:
		return tv.InputMediaAdd(event)
	case VoteFocus:
		return tv.InputVote(event)
	case HelpFocus:
		return tv.InputHelp(event)
	default:
		return event
	}
}

func (tv *TutView) InputLoginView(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.LoginView.Next()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.LoginView.Prev()
		return nil
	}
	if tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.LoginView.Selected()
		return nil
	}
	return event
}

func (tv *TutView) InputMainView(event *tcell.EventKey) *tcell.EventKey {
	switch tv.SubFocus {
	case ListFocus:
		return tv.InputMainViewFeed(event)
	case ContentFocus:
		return tv.InputMainViewContent(event)
	default:
		return event
	}
}

func (tv *TutView) InputMainViewFeed(event *tcell.EventKey) *tcell.EventKey {
	mainFocus := tv.TimelineFocus == FeedFocus

	if tv.tut.Config.Input.MainHome.Match(event.Key(), event.Rune()) {
		tv.Timeline.HomeItemFeed(mainFocus)
		return nil
	}
	if tv.tut.Config.Input.MainEnd.Match(event.Key(), event.Rune()) {
		tv.Timeline.EndItemFeed(mainFocus)
		return nil
	}
	if tv.tut.Config.Input.MainPrevFeed.Match(event.Key(), event.Rune()) {
		if mainFocus {
			tv.Timeline.PrevFeed()
		}
		return nil
	}
	if tv.tut.Config.Input.MainNextFeed.Match(event.Key(), event.Rune()) {
		if mainFocus {
			tv.Timeline.NextFeed()
		}
		return nil
	}
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.Timeline.NextItemFeed(mainFocus)
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.Timeline.PrevItemFeed(mainFocus)
		return nil
	}
	if tv.tut.Config.Input.MainNotificationFocus.Match(event.Key(), event.Rune()) {
		if tv.tut.Config.General.NotificationFeed {
			tv.FocusNotification()
		}
		return nil
	}
	if tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		if mainFocus {
			tv.Timeline.RemoveCurrent(true)
		} else {
			tv.FocusFeed()
		}
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		if mainFocus {
			tv.Timeline.RemoveCurrent(false)
		} else {
			tv.FocusFeed()
		}
		return nil
	}
	return tv.InputItem(event)
}

func (tv *TutView) InputMainViewContent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			tv.Timeline.ScrollDown()
			return nil
		case 'k', 'K':
			tv.Timeline.ScrollUp()
			return nil
		default:
			return event
		}
	}
	return tv.InputItem(event)
}

func (tv *TutView) InputHelp(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		tv.PrevFocus()
		return nil
	}
	switch event.Key() {
	case tcell.KeyEsc:
		tv.PrevFocus()
		return nil
	}
	return event
}

func (tv *TutView) InputViewItem(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		tv.FocusMainNoHistory()
		return nil
	}
	switch event.Key() {
	case tcell.KeyEsc:
		tv.FocusMainNoHistory()
		return nil
	}
	return event
}

func (tv *TutView) InputItem(event *tcell.EventKey) *tcell.EventKey {
	item, err := tv.GetCurrentItem()
	if err != nil {
		return event
	}
	switch event.Rune() {
	case 'c', 'C':
		tv.InitPost(nil)
		return nil
	}
	switch item.Type() {
	case api.StatusType:
		return tv.InputStatus(event, item, item.Raw().(*mastodon.Status))
	case api.UserType, api.ProfileType:
		return tv.InputUser(event, item.Raw().(*api.User))
	case api.NotificationType:
		nd := item.Raw().(*api.NotificationData)
		switch nd.Item.Type {
		case "follow":
			return tv.InputUser(event, nd.User.Raw().(*api.User))
		case "favourite":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status))
		case "reblog":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status))
		case "mention":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status))
		case "status":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status))
		case "poll":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status))
		case "follow_request":
			return tv.InputUser(event, nd.User.Raw().(*api.User))
		}
	case api.ListsType:
		ld := item.Raw().(*mastodon.List)
		return tv.InputList(event, ld)
	}
	return event
}

func (tv *TutView) InputStatus(event *tcell.EventKey, item api.Item, status *mastodon.Status) *tcell.EventKey {
	sr := util.StatusOrReblog(status)

	hasMedia := len(sr.MediaAttachments) > 0
	hasPoll := sr.Poll != nil
	hasSpoiler := sr.Sensitive
	isMine := sr.Account.ID == tv.tut.Client.Me.ID

	boosted := sr.Reblogged
	favorited := sr.Favourited
	bookmarked := sr.Bookmarked

	switch event.Rune() {
	case 'a', 'A':
		openAvatar(tv, sr.Account)
		return nil
	case 'b', 'B':
		txt := "boost"
		if boosted {
			txt = "unboost"
		}
		tv.ModalView.Run(
			fmt.Sprintf("Do you want to %s this toot?", txt), func() {
				ns, err := tv.tut.Client.BoostToggle(status)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't boost toot. Error: %v\n", err),
					)
					return
				}
				*status = *ns
				tv.RedrawControls()
			})
		return nil
	case 'd', 'D':
		if !isMine {
			return nil
		}
		tv.ModalView.Run("Do you want to delete this toot?", func() {
			err := tv.tut.Client.DeleteStatus(sr)
			if err != nil {
				tv.ShowError(
					fmt.Sprintf("Couldn't delete toot. Error: %v\n", err),
				)
				return
			}
			status.Card = nil
			status.Sensitive = false
			status.SpoilerText = ""
			status.Favourited = false
			status.MediaAttachments = nil
			status.Reblogged = false
			status.Content = "Deleted"
			tv.RedrawContent()
		})
		return nil
	case 'f', 'F':
		txt := "favorite"
		if favorited {
			txt = "unfavorite"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s this toot?", txt),
			func() {
				ns, err := tv.tut.Client.FavoriteToogle(status)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't favorite toot. Error: %v\n", err),
					)
					return
				}
				*status = *ns
				tv.RedrawControls()
			})
		return nil
	case 'm', 'M':
		if hasMedia {
			openMedia(tv, sr)
		}
		return nil
	case 'o', 'O':
		tv.SetPage(LinkFocus)
		return nil
	case 'p', 'P':
		if !hasPoll {
			return nil
		}
		tv.VoteView.SetPoll(sr.Poll)
		tv.SetPage(VoteFocus)
		return nil
	case 'r', 'R':
		tv.InitPost(status)
		return nil
	case 's', 'S':
		txt := "save"
		if bookmarked {
			txt = "unsave"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s this toot?", txt),
			func() {
				ns, err := tv.tut.Client.BookmarkToogle(status)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't bookmark toot. Error: %v\n", err),
					)
					return
				}
				*status = *ns
				tv.RedrawControls()
			})
		return nil
	case 't', 'T':
		tv.Timeline.AddFeed(NewThreadFeed(tv, item))
		return nil
	case 'u', 'U':
		user, err := tv.tut.Client.GetUserByID(status.Account.ID)
		if err != nil {
			return nil
		}
		tv.Timeline.AddFeed(NewUserFeed(tv, user))
		return nil
	case 'v', 'V':
		tv.SetPage(ViewFocus)
		return nil
	case 'y', 'Y':
		copyToClipboard(status.URL)
		return nil
	case 'z', 'Z':
		if !hasSpoiler {
			return nil
		}
		if !item.ShowSpoiler() {
			item.ToggleSpoiler()
			tv.RedrawContent()
		}
		return nil
	}

	return event
}

func (tv *TutView) InputUser(event *tcell.EventKey, user *api.User) *tcell.EventKey {
	blocking := user.Relation.Blocking
	muting := user.Relation.Muting
	following := user.Relation.Following
	switch event.Rune() {
	case 'a', 'A':
		openAvatar(tv, *user.Data)
		return nil
	case 'b', 'B':
		txt := "block"
		if blocking {
			txt = "unblock"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s this user?", txt),
			func() {
				rel, err := tv.tut.Client.BlockToggle(user)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't block user. Error: %v\n", err),
					)
					return
				}
				user.Relation = rel
				tv.RedrawControls()
			})
		return nil
	case 'f', 'F':
		txt := "follow"
		if following {
			txt = "unfollow"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s this user?", txt),
			func() {
				rel, err := tv.tut.Client.FollowToggle(user)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't follow user. Error: %v\n", err),
					)
					return
				}
				user.Relation = rel
				tv.RedrawControls()
			})
		return nil
	case 'm', 'M':
		txt := "mute"
		if muting {
			txt = "unmute"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s this user?", txt),
			func() {
				rel, err := tv.tut.Client.MuteToggle(user)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't follow user. Error: %v\n", err),
					)
					return
				}
				user.Relation = rel
				tv.RedrawControls()
			})
		return nil
	case 'o', 'O':
		tv.SetPage(LinkFocus)
		return nil
	case 'u', 'U':
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true)))
		return nil
	case 'v', 'V':
		tv.SetPage(ViewFocus)
		return nil
	case 'y', 'Y':
		copyToClipboard(user.Data.URL)
		return nil
	}
	switch event.Key() {
	case tcell.KeyEnter:
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true)))
		return nil
	}
	return event
}

func (tv *TutView) InputList(event *tcell.EventKey, list *mastodon.List) *tcell.EventKey {
	switch event.Rune() {
	case 'o', 'O':
		tv.Timeline.AddFeed(NewListFeed(tv, list))
		return nil
	}
	switch event.Key() {
	case tcell.KeyEnter:
		tv.Timeline.AddFeed(NewListFeed(tv, list))
		return nil
	}
	return event
}

func (tv *TutView) InputLinkView(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'j', 'J':
			tv.LinkView.Next()
			return nil
		case 'k', 'K':
			tv.LinkView.Prev()
			return nil
		case 'o', 'O':
			tv.LinkView.Open()
			return nil
		case 'y', 'Y':
			tv.LinkView.Yank()
			return nil
		case '1', '2', '3', '4', '5':
			s := string(event.Rune())
			i, _ := strconv.Atoi(s)
			tv.LinkView.OpenCustom(i)
			return nil
		case 'q', 'Q':
			tv.SetPage(MainFocus)
			return nil
		}
	} else {
		switch event.Key() {
		case tcell.KeyEnter:
			tv.LinkView.Open()
			return nil
		case tcell.KeyUp:
			tv.LinkView.Prev()
			return nil
		case tcell.KeyDown:
			tv.LinkView.Next()
			return nil
		case tcell.KeyEsc:
			tv.SetPage(MainFocus)
			return nil
		}
	}
	return event
}

func (tv *TutView) InputComposeView(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'c', 'C':
			tv.ComposeView.EditSpoiler()
			return nil
		case 'e', 'E':
			tv.ComposeView.EditText()
			return nil
		case 'i', 'I':
			tv.ComposeView.IncludeQuote()
			return nil
		case 'm', 'M':
			tv.SetPage(MediaFocus)
			return nil
		case 'p', 'P':
			tv.ComposeView.Post()
			return nil
		case 't', 'T':
			tv.ComposeView.ToggleCW()
			return nil
		case 'v', 'V':
			tv.ComposeView.FocusVisibility()
			return nil
		case 'q', 'Q':
			tv.ModalView.Run(
				"Do you want exit the compose view?", func() {
					tv.FocusMainNoHistory()
				})
			return nil
		}
	} else {
		switch event.Key() {
		case tcell.KeyEsc:
			tv.ModalView.Run(
				"Do you want exit the compose view?", func() {
					tv.FocusMainNoHistory()
				})
			return nil
		}
	}
	return event
}

func (tv *TutView) InputMedia(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j', 'J':
		tv.ComposeView.media.Next()
		return nil
	case 'k', 'K':
		tv.ComposeView.media.Prev()
		return nil
	case 'd', 'D':
		tv.ComposeView.media.Delete()
		return nil
	case 'e', 'E':
		tv.ComposeView.media.EditDesc()
		return nil
	case 'a', 'A':
		tv.SetPage(MediaAddFocus)
		tv.ComposeView.media.SetFocus(false)
		return nil
	case 'q', 'Q':
		tv.SetPage(MediaFocus)
		return nil
	}
	switch event.Key() {
	case tcell.KeyDown:
		tv.ComposeView.media.Next()
		return nil
	case tcell.KeyUp:
		tv.ComposeView.media.Prev()
		return nil
	case tcell.KeyEsc:
		tv.SetPage(ComposeFocus)
		return nil
	}
	return event
}

func (tv *TutView) InputMediaAdd(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		tv.ComposeView.input.AddRune(event.Rune())
		return nil
	}
	switch event.Key() {
	case tcell.KeyTAB:
		tv.ComposeView.input.AutocompleteTab()
		return nil
	case tcell.KeyDown:
		tv.ComposeView.input.AutocompleteNext()
		return nil
	case tcell.KeyBacktab, tcell.KeyUp:
		tv.ComposeView.input.AutocompletePrev()
		return nil
	case tcell.KeyEnter:
		tv.ComposeView.input.CheckDone()
		return nil
	case tcell.KeyEsc:
		tv.SetPage(MediaFocus)
		tv.ComposeView.media.SetFocus(true)
		return nil
	}
	return event
}

func (tv *TutView) InputVote(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j', 'J':
		tv.VoteView.Next()
		return nil
	case 'k', 'K':
		tv.VoteView.Prev()
		return nil
	case 'v', 'V':
		tv.VoteView.Vote()
		return nil
	case ' ':
		tv.VoteView.ToggleSelect()
		return nil
	case 'q', 'Q':
		tv.FocusMainNoHistory()
		return nil
	}
	switch event.Key() {
	case tcell.KeyDown, tcell.KeyTAB:
		tv.VoteView.Next()
		return nil
	case tcell.KeyUp, tcell.KeyBacktab:
		tv.VoteView.Prev()
		return nil
	case tcell.KeyEnter:
		tv.VoteView.ToggleSelect()
		return nil
	case tcell.KeyEsc:
		tv.FocusMainNoHistory()
		return nil
	}
	return event
}

func (tv *TutView) InputCmdView(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		tv.Shared.Bottom.Cmd.DoneFunc(tcell.KeyEnter)
	case tcell.KeyEsc:
		tv.Shared.Bottom.Cmd.Back()
		tv.Shared.Bottom.Cmd.View.Autocomplete()
		return nil
	}
	return event
}
