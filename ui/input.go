package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/feed"
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
	if tv.PageFocus != LoginFocus && tv.PageFocus != CmdFocus {
		event = tv.InputLeaderKey(event)
		if event == nil {
			return nil
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

func (tv *TutView) InputLeaderKey(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.General.LeaderKey == rune(0) {
		return event
	}
	if event.Rune() == tv.tut.Config.General.LeaderKey {
		tv.Leader.Reset()
		return nil
	} else if tv.Leader.IsActive() {
		if event.Rune() != rune(0) {
			tv.Leader.AddRune(event.Rune())
		}
		action := config.LeaderNone
		var subaction string
		content := tv.Leader.Content()
		for _, la := range tv.tut.Config.General.LeaderActions {
			if la.Shortcut == content {
				action = la.Command
				subaction = la.Subaction
				break
			}
		}
		if action == config.LeaderNone {
			return nil
		}
		switch action {
		case config.LeaderHome:
			tv.HomeCommand()
		case config.LeaderDirect:
			tv.DirectCommand()
		case config.LeaderLocal:
			tv.LocalCommand()
		case config.LeaderFederated:
			tv.FederatedCommand()
		case config.LeaderCompose:
			tv.ComposeCommand()
		case config.LeaderBlocking:
			tv.BlockingCommand()
		case config.LeaderBookmarks, config.LeaderSaved:
			tv.BookmarksCommand()
		case config.LeaderFavorited:
			tv.FavoritedCommand()
		case config.LeaderBoosts:
			tv.BoostsCommand()
		case config.LeaderFavorites:
			tv.FavoritesCommand()
		case config.LeaderFollowing:
			tv.FollowingCommand()
		case config.LeaderFollowers:
			tv.FollowersCommand()
		case config.LeaderMuting:
			tv.MutingCommand()
		case config.LeaderProfile:
			tv.ProfileCommand()
		case config.LeaderNotifications:
			tv.NotificationsCommand()
		case config.LeaderLists:
			tv.ListsCommand()
		case config.LeaderTag:
			tv.TagCommand(subaction)
		case config.LeaderWindow:
			tv.WindowCommand(subaction)
		}
		tv.Leader.ResetInactive()
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
	if tv.tut.Config.Input.MainHome.Match(event.Key(), event.Rune()) {
		tv.Timeline.HomeItemFeed()
		return nil
	}
	if tv.tut.Config.Input.MainEnd.Match(event.Key(), event.Rune()) {
		tv.Timeline.EndItemFeed()
		return nil
	}
	if tv.tut.Config.Input.MainPrevFeed.Match(event.Key(), event.Rune()) {
		tv.Timeline.PrevFeed()
		return nil
	}
	if tv.tut.Config.Input.MainNextFeed.Match(event.Key(), event.Rune()) {
		tv.Timeline.NextFeed()
		return nil
	}
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.Timeline.NextItemFeed()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.Timeline.PrevItemFeed()
		return nil
	}
	if tv.tut.Config.Input.MainPrevWindow.Match(event.Key(), event.Rune()) {
		if tv.tut.Config.General.NotificationFeed {
			tv.PrevFeed()
		}
		return nil
	}
	if tv.tut.Config.Input.MainNextWindow.Match(event.Key(), event.Rune()) {
		if tv.tut.Config.General.NotificationFeed {
			tv.NextFeed()
		}
		return nil
	}
	if tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		exiting := tv.Timeline.RemoveCurrent(false)
		if exiting && tv.Timeline.FeedFocusIndex == 0 {
			tv.ModalView.Run("Do you want to exit tut?",
				func() {
					tv.Timeline.RemoveCurrent(true)
				})
			return nil
		} else if exiting && tv.Timeline.FeedFocusIndex != 0 {
			tv.FocusFeed(0)
		}
		return nil
	}
	for i, tl := range tv.tut.Config.General.Timelines {
		if tl.Key.Match(event.Key(), event.Rune()) {
			tv.FocusFeed(i)
		}
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) {
		exiting := tv.Timeline.RemoveCurrent(false)
		if exiting && tv.Timeline.FeedFocusIndex != 0 {
			tv.FocusFeed(0)
		}
		return nil
	}
	return tv.InputItem(event)
}

func (tv *TutView) InputMainViewContent(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.Timeline.ScrollDown()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.Timeline.ScrollDown()
		return nil
	}
	return tv.InputItem(event)
}

func (tv *TutView) InputHelp(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.PrevFocus()
		return nil
	}
	return event
}

func (tv *TutView) InputViewItem(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.FocusMainNoHistory()
		return nil
	}
	return event
}

func (tv *TutView) InputItem(event *tcell.EventKey) *tcell.EventKey {
	fd := tv.GetCurrentFeed()
	ft := fd.Data.Type()
	item, err := tv.GetCurrentItem()
	if err != nil {
		return event
	}
	if tv.tut.Config.Input.MainCompose.Match(event.Key(), event.Rune()) {
		tv.InitPost(nil)
		return nil
	}
	switch item.Type() {
	case api.StatusType:
		return tv.InputStatus(event, item, item.Raw().(*mastodon.Status))
	case api.UserType, api.ProfileType:
		if ft == feed.FollowRequests {
			return tv.InputUser(event, item.Raw().(*api.User), true)
		} else {
			return tv.InputUser(event, item.Raw().(*api.User), false)
		}
	case api.NotificationType:
		nd := item.Raw().(*api.NotificationData)
		switch nd.Item.Type {
		case "follow":
			return tv.InputUser(event, nd.User.Raw().(*api.User), false)
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
			return tv.InputUser(event, nd.User.Raw().(*api.User), true)
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

	if tv.tut.Config.Input.StatusAvatar.Match(event.Key(), event.Rune()) {
		openAvatar(tv, sr.Account)
		return nil
	}
	if tv.tut.Config.Input.StatusBoost.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.StatusDelete.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.StatusFavorite.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.StatusMedia.Match(event.Key(), event.Rune()) {
		if hasMedia {
			openMedia(tv, sr)
		}
		return nil
	}
	if tv.tut.Config.Input.StatusLinks.Match(event.Key(), event.Rune()) {
		tv.SetPage(LinkFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusPoll.Match(event.Key(), event.Rune()) {
		if !hasPoll {
			return nil
		}
		tv.VoteView.SetPoll(sr.Poll)
		tv.SetPage(VoteFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusReply.Match(event.Key(), event.Rune()) {
		tv.InitPost(status)
		return nil
	}
	if tv.tut.Config.Input.StatusBookmark.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.StatusThread.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewThreadFeed(tv, item))
		return nil
	}
	if tv.tut.Config.Input.StatusUser.Match(event.Key(), event.Rune()) {
		user, err := tv.tut.Client.GetUserByID(status.Account.ID)
		if err != nil {
			return nil
		}
		tv.Timeline.AddFeed(NewUserFeed(tv, user))
		return nil
	}
	if tv.tut.Config.Input.StatusViewFocus.Match(event.Key(), event.Rune()) {
		tv.SetPage(ViewFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusYank.Match(event.Key(), event.Rune()) {
		copyToClipboard(status.URL)
		return nil
	}
	if tv.tut.Config.Input.StatusToggleSpoiler.Match(event.Key(), event.Rune()) {
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

func (tv *TutView) InputUser(event *tcell.EventKey, user *api.User, fr bool) *tcell.EventKey {
	blocking := user.Relation.Blocking
	muting := user.Relation.Muting
	following := user.Relation.Following

	if tv.tut.Config.Input.UserFollowRequestDecide.Match(event.Key(), event.Rune()) {
		tv.ModalView.RunDecide("Do you want accept the follow request?",
			func() {
				err := tv.tut.Client.FollowRequestAccept(user.Data)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't accept follow request. Error: %v\n", err),
					)
					return
				}
				f := tv.GetCurrentFeed()
				f.Delete()
				tv.RedrawContent()
			},
			func() {
				err := tv.tut.Client.FollowRequestDeny(user.Data)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't deny follow request. Error: %v\n", err),
					)
					return
				}
				f := tv.GetCurrentFeed()
				f.Delete()
				tv.RedrawContent()
			})
		return nil
	}
	if tv.tut.Config.Input.UserAvatar.Match(event.Key(), event.Rune()) {
		openAvatar(tv, *user.Data)
		return nil
	}
	if tv.tut.Config.Input.UserBlock.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.UserFollow.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.UserMute.Match(event.Key(), event.Rune()) {
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
	}
	if tv.tut.Config.Input.UserLinks.Match(event.Key(), event.Rune()) {
		tv.SetPage(LinkFocus)
		return nil
	}
	if tv.tut.Config.Input.UserUser.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true)))
		return nil
	}
	if tv.tut.Config.Input.UserViewFocus.Match(event.Key(), event.Rune()) {
		tv.SetPage(ViewFocus)
		return nil
	}
	if tv.tut.Config.Input.UserYank.Match(event.Key(), event.Rune()) {
		copyToClipboard(user.Data.URL)
		return nil
	}
	if tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true)))
		return nil
	}
	return event
}

func (tv *TutView) InputList(event *tcell.EventKey, list *mastodon.List) *tcell.EventKey {
	if tv.tut.Config.Input.ListOpenFeed.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewListFeed(tv, list))
		return nil
	}
	return event
}

func (tv *TutView) InputLinkView(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.LinkView.Next()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.LinkView.Prev()
		return nil
	}
	if tv.tut.Config.Input.LinkOpen.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.LinkView.Open()
		return nil
	}
	if tv.tut.Config.Input.LinkYank.Match(event.Key(), event.Rune()) {
		tv.LinkView.Yank()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.SetPage(MainFocus)
		return nil
	}
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case '1', '2', '3', '4', '5':
			s := string(event.Rune())
			i, _ := strconv.Atoi(s)
			tv.LinkView.OpenCustom(i)
			return nil
		}
	}
	return event
}

func (tv *TutView) InputComposeView(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.ComposeEditSpoiler.Match(event.Key(), event.Rune()) {
		tv.ComposeView.EditSpoiler()
		return nil
	}
	if tv.tut.Config.Input.ComposeEditText.Match(event.Key(), event.Rune()) {
		tv.ComposeView.EditText()
		return nil
	}
	if tv.tut.Config.Input.ComposeIncludeQuote.Match(event.Key(), event.Rune()) {
		tv.ComposeView.IncludeQuote()
		return nil
	}
	if tv.tut.Config.Input.ComposeMediaFocus.Match(event.Key(), event.Rune()) {
		tv.SetPage(MediaFocus)
		return nil
	}
	if tv.tut.Config.Input.ComposePost.Match(event.Key(), event.Rune()) {
		tv.ComposeView.Post()
		return nil
	}
	if tv.tut.Config.Input.ComposeToggleContentWarning.Match(event.Key(), event.Rune()) {
		tv.ComposeView.ToggleCW()
		return nil
	}
	if tv.tut.Config.Input.ComposeVisibility.Match(event.Key(), event.Rune()) {
		tv.ComposeView.FocusVisibility()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.ModalView.Run(
			"Do you want exit the compose view?", func() {
				tv.FocusMainNoHistory()
			})
		return nil
	}
	return event
}

func (tv *TutView) InputMedia(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.ComposeView.media.Next()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.ComposeView.media.Prev()
		return nil
	}
	if tv.tut.Config.Input.MediaDelete.Match(event.Key(), event.Rune()) {
		tv.ComposeView.media.Delete()
		return nil
	}
	if tv.tut.Config.Input.MediaEditDesc.Match(event.Key(), event.Rune()) {
		tv.ComposeView.media.EditDesc()
		return nil
	}
	if tv.tut.Config.Input.MediaAdd.Match(event.Key(), event.Rune()) {
		tv.SetPage(MediaAddFocus)
		tv.ComposeView.media.SetFocus(false)
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
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
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.VoteView.Next()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.VoteView.Prev()
		return nil
	}
	if tv.tut.Config.Input.VoteVote.Match(event.Key(), event.Rune()) {
		tv.VoteView.Vote()
		return nil
	}
	if tv.tut.Config.Input.VoteSelect.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.VoteView.ToggleSelect()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
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
