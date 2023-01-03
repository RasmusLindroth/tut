package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	case PollFocus:
		return tv.InputPollView(event)
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
	case PreferenceFocus:
		return tv.InputPreference(event)
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
	if tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.tut.App.Stop()
		tv.CleanExit(0)
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
		case config.LeaderSpecialAll:
			tv.SpecialCommand(true, true)
		case config.LeaderSpecialBoosts:
			tv.SpecialCommand(false, true)
		case config.LeaderSpecialReplies:
			tv.SpecialCommand(true, false)
		case config.LeaderClearNotifications:
			tv.ClearNotificationsCommand()
		case config.LeaderCompose:
			tv.ComposeCommand()
		case config.LeaderEdit:
			tv.EditCommand()
		case config.LeaderBlocking:
			tv.BlockingCommand()
		case config.LeaderBookmarks, config.LeaderSaved:
			tv.BookmarksCommand()
		case config.LeaderFavorited:
			tv.FavoritedCommand()
		case config.LeaderHistory:
			tv.HistoryCommand()
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
		case config.LeaderPreferences:
			tv.PreferencesCommand()
		case config.LeaderProfile:
			tv.ProfileCommand()
		case config.LeaderNotifications:
			tv.NotificationsCommand()
		case config.LeaderMentions:
			tv.MentionsCommand()
		case config.LeaderLoadNewer:
			tv.LoadNewerCommand()
		case config.LeaderLists:
			tv.ListsCommand()
		case config.LeaderStickToTop:
			tv.ToggleStickToTop()
		case config.LeaderRefetch:
			tv.RefetchCommand()
		case config.LeaderTag:
			tv.TagCommand(subaction)
		case config.LeaderTags:
			tv.TagsCommand()
		case config.LeaderWindow:
			tv.WindowCommand(subaction)
		case config.LeaderCloseWindow:
			tv.CloseWindowCommand()
		case config.LeaderMoveWindowLeft:
			tv.MoveWindowLeft()
		case config.LeaderMoveWindowRight:
			tv.MoveWindowRight()
		case config.LeaderMoveWindowHome:
			tv.MoveWindowHome()
		case config.LeaderMoveWindowEnd:
			tv.MoveWindowEnd()
		case config.LeaderSwitch:
			tv.SwitchCommand(subaction)
		case config.LeaderListPlacement:
			switch subaction {
			case "top":
				tv.ListPlacementCommand(config.ListPlacementTop)
			case "right":
				tv.ListPlacementCommand(config.ListPlacementRight)
			case "bottom":
				tv.ListPlacementCommand(config.ListPlacementBottom)
			case "left":
				tv.ListPlacementCommand(config.ListPlacementLeft)
			}
		case config.LeaderListSplit:
			switch subaction {
			case "row":
				tv.ListSplitCommand(config.ListRow)
			case "column":
				tv.ListSplitCommand(config.ListColumn)
			}
		case config.LeaderProportions:
			parts := strings.Split(subaction, " ")
			if len(parts) == 2 {
				tv.ProportionsCommand(parts[0], parts[1])
			}
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
		tv.PrevFeed()
		return nil
	}
	if tv.tut.Config.Input.MainNextWindow.Match(event.Key(), event.Rune()) {
		tv.NextFeed()
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
	if tv.tut.Config.Input.MainCompose.Match(event.Key(), event.Rune()) {
		tv.InitPost(nil, nil)
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
	if tv.tut.Config.Input.MainCompose.Match(event.Key(), event.Rune()) {
		tv.InitPost(nil, nil)
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
	switch item.Type() {
	case api.StatusType:
		return tv.InputStatus(event, item, item.Raw().(*mastodon.Status), nil, fd.Data.Type())
	case api.StatusHistoryType:
		return tv.InputStatusHistory(event, item, item.Raw().(*mastodon.StatusHistory), nil)
	case api.UserType, api.ProfileType:
		switch ft {
		case config.FollowRequests:
			return tv.InputUser(event, item.Raw().(*api.User), InputUserFollowRequest)
		case config.ListUsersAdd:
			return tv.InputUser(event, item.Raw().(*api.User), InputUserListAdd)
		case config.ListUsersIn:
			return tv.InputUser(event, item.Raw().(*api.User), InputUserListDelete)
		default:
			return tv.InputUser(event, item.Raw().(*api.User), InputUserNormal)
		}
	case api.NotificationType:
		nd := item.Raw().(*api.NotificationData)
		switch nd.Item.Type {
		case "follow":
			return tv.InputUser(event, nd.User.Raw().(*api.User), InputUserNormal)
		case "favourite":
			user := nd.User.Raw().(*api.User)
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), user.Data, config.Notifications)
		case "reblog":
			user := nd.User.Raw().(*api.User)
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), user.Data, config.Notifications)
		case "mention":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), nil, config.Notifications)
		case "update":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), nil, config.Notifications)
		case "status":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), nil, config.Notifications)
		case "poll":
			return tv.InputStatus(event, nd.Status, nd.Status.Raw().(*mastodon.Status), nil, config.Notifications)
		case "follow_request":
			return tv.InputUser(event, nd.User.Raw().(*api.User), InputUserFollowRequest)
		}
	case api.ListsType:
		ld := item.Raw().(*mastodon.List)
		return tv.InputList(event, ld)
	case api.TagType:
		tag := item.Raw().(*mastodon.Tag)
		return tv.InputTag(event, tag)
	}
	return event
}

func (tv *TutView) InputStatus(event *tcell.EventKey, item api.Item, status *mastodon.Status, nAcc *mastodon.Account, fd config.FeedType) *tcell.EventKey {
	sr := util.StatusOrReblog(status)

	hasMedia := len(sr.MediaAttachments) > 0
	hasPoll := sr.Poll != nil
	hasSpoiler := sr.Sensitive
	isMine := sr.Account.ID == tv.tut.Client.Me.ID

	boosted, favorited, bookmarked := false, false, false
	if sr.Reblogged != nil {
		boosted = sr.Reblogged.(bool)
	}
	if sr.Favourited != nil {
		favorited = sr.Favourited.(bool)
	}
	if sr.Bookmarked != nil {
		bookmarked = sr.Bookmarked.(bool)
	}

	if tv.tut.Config.Input.StatusAvatar.Match(event.Key(), event.Rune()) {
		if nAcc != nil {
			openAvatar(tv, *nAcc)
		} else {
			openAvatar(tv, sr.Account)
		}
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
	if tv.tut.Config.Input.StatusEdit.Match(event.Key(), event.Rune()) {
		tv.EditCommand()
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
		tv.InitPost(status, nil)
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
		tv.Timeline.AddFeed(NewThreadFeed(tv, item, &config.Timeline{
			FeedType: config.Thread,
		}))
		return nil
	}
	if tv.tut.Config.Input.StatusUser.Match(event.Key(), event.Rune()) {
		id := sr.Account.ID
		if nAcc != nil {
			id = nAcc.ID
		}
		user, err := tv.tut.Client.GetUserByID(id)
		if err != nil {
			return nil
		}
		tv.Timeline.AddFeed(NewUserFeed(tv, user, &config.Timeline{
			FeedType: config.User,
		}))
		return nil
	}
	if tv.tut.Config.Input.StatusViewFocus.Match(event.Key(), event.Rune()) {
		tv.SetPage(ViewFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusYank.Match(event.Key(), event.Rune()) {
		copyToClipboard(sr.URL)
		return nil
	}
	if tv.tut.Config.Input.StatusToggleCW.Match(event.Key(), event.Rune()) {
		filtered, _, _, forceView := item.Filtered(fd)
		if filtered && !forceView {
			item.ForceViewFilter()
			tv.RedrawContent()
			return nil
		}
		if !hasSpoiler {
			return nil
		}
		if !item.ShowCW() {
			item.ToggleCW()
			tv.RedrawContent()
		}
		return nil
	}

	return event
}

func (tv *TutView) InputStatusHistory(event *tcell.EventKey, item api.Item, sr *mastodon.StatusHistory, nAcc *mastodon.Account) *tcell.EventKey {
	hasMedia := len(sr.MediaAttachments) > 0
	hasSpoiler := sr.Sensitive

	status := &mastodon.Status{
		Content:          sr.Content,
		SpoilerText:      sr.SpoilerText,
		Account:          sr.Account,
		Sensitive:        sr.Sensitive,
		CreatedAt:        sr.CreatedAt,
		Emojis:           sr.Emojis,
		MediaAttachments: sr.MediaAttachments,
	}

	if tv.tut.Config.Input.StatusAvatar.Match(event.Key(), event.Rune()) {
		if nAcc != nil {
			openAvatar(tv, *nAcc)
		} else {
			openAvatar(tv, sr.Account)
		}
		return nil
	}
	if tv.tut.Config.Input.StatusMedia.Match(event.Key(), event.Rune()) {
		if hasMedia {
			openMedia(tv, status)
		}
		return nil
	}
	if tv.tut.Config.Input.StatusLinks.Match(event.Key(), event.Rune()) {
		tv.SetPage(LinkFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusViewFocus.Match(event.Key(), event.Rune()) {
		tv.SetPage(ViewFocus)
		return nil
	}
	if tv.tut.Config.Input.StatusToggleCW.Match(event.Key(), event.Rune()) {
		if !hasSpoiler {
			return nil
		}
		if !item.ShowCW() {
			item.ToggleCW()
			tv.RedrawContent()
		}
		return nil
	}

	return event
}

type InputUserType uint

const (
	InputUserNormal = iota
	InputUserFollowRequest
	InputUserListAdd
	InputUserListDelete
)

func (tv *TutView) InputUser(event *tcell.EventKey, user *api.User, ut InputUserType) *tcell.EventKey {
	blocking := user.Relation.Blocking
	muting := user.Relation.Muting
	following := user.Relation.Following

	if ut == InputUserListAdd {
		if tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) ||
			tv.tut.Config.Input.ListUserAdd.Match(event.Key(), event.Rune()) {
			ad := user.AdditionalData
			switch ad.(type) {
			case *mastodon.List:
				l := user.AdditionalData.(*mastodon.List)
				err := tv.tut.Client.AddUserToList(user.Data, l)
				if err != nil {
					tv.ShowError(fmt.Sprintf("Couldn't add user to list. Error: %v", err))
				}
				return nil
			default:
				return event
			}

		}
		return event
	}

	if ut == InputUserListDelete {
		if tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) ||
			tv.tut.Config.Input.ListUserDelete.Match(event.Key(), event.Rune()) {
			ad := user.AdditionalData
			switch ad.(type) {
			case *mastodon.List:
				l := user.AdditionalData.(*mastodon.List)
				err := tv.tut.Client.DeleteUserFromList(user.Data, l)
				if err != nil {
					tv.ShowError(fmt.Sprintf("Couldn't remove user from list. Error: %v", err))
				}
				return nil
			default:
				return event
			}
		}
		return event
	}

	if ut == InputUserFollowRequest && tv.tut.Config.Input.UserFollowRequestDecide.Match(event.Key(), event.Rune()) {
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
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true), &config.Timeline{
			FeedType: config.User,
		}))
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
		tv.Timeline.AddFeed(NewUserFeed(tv, api.NewUserItem(user, true), &config.Timeline{
			FeedType: config.User,
		}))
		return nil
	}
	return event
}

func (tv *TutView) InputList(event *tcell.EventKey, list *mastodon.List) *tcell.EventKey {
	if tv.tut.Config.Input.ListOpenFeed.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewListFeed(tv, list, &config.Timeline{
			FeedType: config.List,
		}))
		return nil
	}
	if tv.tut.Config.Input.ListUserList.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewUsersInListFeed(tv, list, &config.Timeline{
			FeedType: config.ListUsersIn,
		}))
		return nil
	}
	if tv.tut.Config.Input.ListUserAdd.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewUsersAddListFeed(tv, list, &config.Timeline{
			FeedType: config.ListUsersAdd,
		}))
		return nil
	}
	return event
}

func (tv *TutView) InputTag(event *tcell.EventKey, tag *mastodon.Tag) *tcell.EventKey {
	if tv.tut.Config.Input.TagOpenFeed.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalEnter.Match(event.Key(), event.Rune()) {
		tv.Timeline.AddFeed(NewTagFeed(tv, &config.Timeline{
			FeedType:  config.Tag,
			Subaction: tag.Name,
		}))
		return nil
	}
	if tv.tut.Config.Input.TagFollow.Match(event.Key(), event.Rune()) {
		txt := "follow"
		if tag.Following != nil && tag.Following == true {
			txt = "unfollow"
		}
		tv.ModalView.Run(fmt.Sprintf("Do you want to %s #%s?", txt, tag.Name),
			func() {
				nt, err := tv.tut.Client.TagToggleFollow(tag)
				if err != nil {
					tv.ShowError(
						fmt.Sprintf("Couldn't %s tag. Error: %v\n", txt, err),
					)
					return
				}
				*tag = *nt
				tv.RedrawControls()
			})
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
	for _, oc := range tv.tut.Config.OpenCustom.OpenCustoms {
		if oc.Key.Match(event.Key(), event.Rune()) {
			tv.LinkView.OpenCustom(oc)
			return nil
		}
	}
	return event
}

func (tv *TutView) InputComposeView(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.ComposeEditCW.Match(event.Key(), event.Rune()) {
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
		if tv.PollView.HasPoll() {
			tv.ShowError("Can't add media when you have a poll")
			return nil
		}
		tv.SetPage(MediaFocus)
		return nil
	}
	if tv.tut.Config.Input.ComposePoll.Match(event.Key(), event.Rune()) {
		if tv.ComposeView.HasMedia() {
			tv.ShowError("Can't add poll when you have added media")
			return nil
		}
		tv.SetPage(PollFocus)
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
	if tv.tut.Config.Input.ComposeLanguage.Match(event.Key(), event.Rune()) {
		tv.ComposeView.FocusLang()
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

func (tv *TutView) InputPollView(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.PollAdd.Match(event.Key(), event.Rune()) {
		tv.PollView.Add()
		return nil
	}
	if tv.tut.Config.Input.PollEdit.Match(event.Key(), event.Rune()) {
		tv.PollView.Edit()
		return nil
	}
	if tv.tut.Config.Input.PollDelete.Match(event.Key(), event.Rune()) {
		tv.PollView.Delete()
		return nil
	}
	if tv.tut.Config.Input.PollMultiToggle.Match(event.Key(), event.Rune()) {
		tv.PollView.ToggleMultiple()
		return nil
	}
	if tv.tut.Config.Input.PollExpiration.Match(event.Key(), event.Rune()) {
		tv.PollView.FocusExpiration()
		return nil
	}
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.PollView.Next()
		return nil
	}
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.PollView.Prev()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.SetPage(ComposeFocus)
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

func (tv *TutView) InputPreference(event *tcell.EventKey) *tcell.EventKey {
	if tv.PreferenceView.HasFieldFocus() {
		return tv.InputPreferenceFields(event)
	}
	if tv.tut.Config.Input.PreferenceFields.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.FieldFocus()
		return nil
	}
	if tv.tut.Config.Input.PreferenceName.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.EditDisplayname()
		return nil
	}
	if tv.tut.Config.Input.PreferenceVisibility.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.FocusVisibility()
		return nil
	}
	if tv.tut.Config.Input.PreferenceBio.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.EditBio()
		return nil
	}
	if tv.tut.Config.Input.PreferenceSave.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.Save()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.ModalView.Run(
			"Do you want exit the preference view?", func() {
				tv.FocusMainNoHistory()
			})
		return nil
	}
	return event
}
func (tv *TutView) InputPreferenceFields(event *tcell.EventKey) *tcell.EventKey {
	if tv.tut.Config.Input.GlobalUp.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.PrevField()
		return nil
	}
	if tv.tut.Config.Input.GlobalDown.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.NextField()
		return nil
	}
	if tv.tut.Config.Input.PreferenceFieldsAdd.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.AddField()
		return nil
	}
	if tv.tut.Config.Input.PreferenceFieldsEdit.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.EditField()
		return nil
	}
	if tv.tut.Config.Input.PreferenceFieldsDelete.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.DeleteField()
		return nil
	}
	if tv.tut.Config.Input.GlobalBack.Match(event.Key(), event.Rune()) ||
		tv.tut.Config.Input.GlobalExit.Match(event.Key(), event.Rune()) {
		tv.PreferenceView.MainFocus()
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

func (tv *TutView) MouseInput(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	if event == nil {
		return nil, action
	}
	switch action {
	case tview.MouseLeftUp, tview.MouseMiddleUp, tview.MouseRightUp:
		return event, action
	}

	switch tv.PageFocus {
	case ViewFocus, MainFocus:
		return tv.MouseInputMainView(event, action)
	case LoginFocus:
		tv.MouseInputLoginView(event, action)
	case LinkFocus:
		return tv.MouseInputLinkView(event, action)
	case MediaFocus:
		return tv.MouseInputMediaView(event, action)
	case VoteFocus:
		return tv.MouseInputVoteView(event, action)
	case ModalFocus:
		tv.MouseInputModalView(event, action)
	case PollFocus:
		return tv.MouseInputPollView(event, action)
	case HelpFocus:
		return tv.MouseInputHelpView(event, action)
	case ComposeFocus:
		return tv.MouseInputComposeView(event, action)
	case PreferenceFocus:
		return tv.MouseInputPreferenceView(event, action)
	}

	return nil, action
}

func (tv *TutView) feedListMouse(list *tview.List, i int, event *tcell.EventMouse, action tview.MouseAction) {
	tv.SetPage(MainFocus)
	tv.FocusFeed(i)
	mh := list.MouseHandler()
	if mh == nil {
		return
	}
	lastIndex := list.GetCurrentItem()
	mh(action, event, func(p tview.Primitive) {})
	newIndex := list.GetCurrentItem()
	if lastIndex != newIndex {
		tv.Timeline.SetItemFeedIndex(newIndex)
	}
}

var scrollSleepTime time.Duration = 150

type scrollSleep struct {
	mux  sync.Mutex
	last time.Time
	next func()
	prev func()
}

func NewScrollSleep(next func(), prev func()) *scrollSleep {
	return &scrollSleep{
		next: next,
		prev: prev,
	}
}

func (sc *scrollSleep) Action(list *tview.List, action tview.MouseAction) {
	mh := list.MouseHandler()
	if mh == nil {
		return
	}
	lock := sc.mux.TryLock()
	if !lock {
		return
	}
	if time.Since(sc.last) < (scrollSleepTime * time.Millisecond) {
		sc.mux.Unlock()
		return
	}
	if action == tview.MouseScrollDown {
		sc.next()
	}
	if action == tview.MouseScrollUp {
		sc.prev()
	}
	sc.last = time.Now()
	sc.mux.Unlock()
}

func (tv *TutView) MouseInputMainView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseScrollDown, tview.MouseScrollUp:
		f := tv.GetCurrentFeed()
		if f.Content.Main.InRect(x, y) {
			if action == tview.MouseScrollDown {
				tv.Timeline.ScrollDown()
				return nil, action
			}
			if action == tview.MouseScrollUp {
				tv.Timeline.ScrollUp()
				return nil, action
			}
		}
		for _, tl := range tv.Timeline.Feeds {
			fl := tl.GetFeedList()
			if fl.Text.InRect(x, y) {
				tv.Timeline.scrollSleep.Action(fl.Text, action)
				return nil, action
			}
			if fl.Symbol.InRect(x, y) {
				tv.Timeline.scrollSleep.Action(fl.Symbol, action)
				return nil, action
			}
		}
	case tview.MouseLeftClick:
		f := tv.GetCurrentFeed()
		if f.Content.Main.InRect(x, y) {
			tv.SetPage(ViewFocus)
			return nil, action
		}
		if f.Content.Controls.InRect(x, y) {
			return event, action
		}
		for i, tl := range tv.Timeline.Feeds {
			fl := tl.GetFeedList()
			if fl.Text.InRect(x, y) {
				tv.feedListMouse(fl.Text, i, event, action)
				return nil, action
			}
			if fl.Symbol.InRect(x, y) {
				tv.feedListMouse(fl.Symbol, i, event, action)
				return nil, action
			}
		}
	}
	return nil, action
}

func (tv *TutView) MouseInputLoginView(event *tcell.EventMouse, action tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		list := tv.LoginView.list
		if !list.InRect(x, y) {
			return
		}
		mh := list.MouseHandler()
		if mh == nil {
			return
		}
		mh(action, event, func(p tview.Primitive) {})
		tv.LoginView.Selected()
	case tview.MouseScrollDown, tview.MouseScrollUp:
		tv.LoginView.scrollSleep.Action(tv.LoginView.list, action)
	}
}

func (tv *TutView) MouseInputLinkView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.LinkView.controls.InRect(x, y) {
			return event, action
		}
		list := tv.LinkView.list
		if !list.InRect(x, y) {
			return nil, action
		}
		mh := list.MouseHandler()
		if mh == nil {
			return nil, action
		}
		mh(action, event, func(p tview.Primitive) {})
		tv.LinkView.Open()
	case tview.MouseScrollDown, tview.MouseScrollUp:
		tv.LinkView.scrollSleep.Action(tv.LinkView.list, action)
	}
	return nil, action
}

func (tv *TutView) MouseInputMediaView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.ComposeView.controls.InRect(x, y) {
			return event, action
		}
		list := tv.ComposeView.media.list
		if !list.InRect(x, y) {
			return nil, action
		}
		mh := list.MouseHandler()
		if mh == nil {
			return nil, action
		}
		mh(action, event, func(p tview.Primitive) {})
	case tview.MouseScrollDown, tview.MouseScrollUp:
		tv.ComposeView.media.scrollSleep.Action(tv.ComposeView.media.list, action)
	}
	return nil, action
}

func (tv *TutView) MouseInputVoteView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.VoteView.controls.InRect(x, y) {
			return event, action
		}
		list := tv.VoteView.list
		if !list.InRect(x, y) {
			return nil, action
		}
		mh := list.MouseHandler()
		if mh == nil {
			return nil, action
		}
		mh(action, event, func(p tview.Primitive) {})
		tv.VoteView.ToggleSelect()
	case tview.MouseScrollDown, tview.MouseScrollUp:
		tv.VoteView.scrollSleep.Action(tv.VoteView.list, action)
	}
	return nil, action
}

func (tv *TutView) MouseInputModalView(event *tcell.EventMouse, action tview.MouseAction) {
	switch action {
	case tview.MouseLeftClick:
		modal := tv.ModalView.View
		mh := modal.MouseHandler()
		if mh == nil {
			return
		}
		mh(action, event, func(p tview.Primitive) {})
	}
}

func (tv *TutView) MouseInputPollView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.PollView.controls.InRect(x, y) {
			return event, action
		}
		list := tv.PollView.list
		if !list.InRect(x, y) {
			return nil, action
		}
		mh := list.MouseHandler()
		if mh == nil {
			return nil, action
		}
		mh(action, event, func(p tview.Primitive) {})
	case tview.MouseScrollDown, tview.MouseScrollUp:
		tv.PollView.scrollSleep.Action(tv.PollView.list, action)
	}
	return nil, action
}

func (tv *TutView) MouseInputHelpView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.HelpView.controls.InRect(x, y) {
			return event, action
		}
	}
	return nil, action
}

func (tv *TutView) MouseInputPreferenceView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.PreferenceView.controls.InRect(x, y) {
			return event, action
		}
	}
	return nil, action
}

func (tv *TutView) MouseInputComposeView(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	x, y := event.Position()
	switch action {
	case tview.MouseLeftClick:
		if tv.ComposeView.controls.InRect(x, y) {
			return event, action
		}
	}
	return nil, action
}
