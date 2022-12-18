package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/icza/gox/timex"
	"github.com/rivo/tview"
)

func DrawListItem(cfg *config.Config, item api.Item) (string, string) {
	switch item.Type() {
	case api.StatusType:
		s := item.Raw().(*mastodon.Status)
		symbol := ""
		status := s
		if s.Reblog != nil {
			status = s.Reblog
		}
		if status.RepliesCount > 0 {
			symbol = " ⤶ "
		}
		if item.Pinned() {
			symbol = " ! "
		}
		acc := strings.TrimSpace(s.Account.Acct)
		if cfg.General.ShowBoostedUser && s.Reblog != nil {
			acc = strings.TrimSpace(s.Reblog.Account.Acct)
		}
		if s.Reblog != nil && cfg.General.ShowIcons {
			acc = fmt.Sprintf("♺ %s", acc)
		}
		d := OutputDate(cfg, s.CreatedAt.Local())
		return fmt.Sprintf("%s %s", d, acc), symbol
	case api.StatusHistoryType:
		s := item.Raw().(*mastodon.StatusHistory)
		acc := strings.TrimSpace(s.Account.Acct)
		d := OutputDate(cfg, s.CreatedAt.Local())
		return fmt.Sprintf("%s %s", d, acc), ""
	case api.UserType:
		a := item.Raw().(*api.User)
		return strings.TrimSpace(a.Data.Acct), ""
	case api.ProfileType:
		return "Profile", ""
	case api.NotificationType:
		a := item.Raw().(*api.NotificationData)
		symbol := ""
		switch a.Item.Type {
		case "follow", "follow_request":
			symbol += " + "
		case "favourite":
			symbol = " ★ "
		case "reblog":
			symbol = " ♺ "
		case "mention":
			symbol = " ⤶ "
		case "update":
			symbol = " ☢ "
		case "poll":
			symbol = " = "
		case "status":
			symbol = " ⤶ "
		}
		d := OutputDate(cfg, a.Item.CreatedAt.Local())
		return fmt.Sprintf("%s %s", d, strings.TrimSpace(a.Item.Account.Acct)), symbol
	case api.ListsType:
		a := item.Raw().(*mastodon.List)
		return tview.Escape(a.Title), ""
	case api.TagType:
		a := item.Raw().(*mastodon.Tag)
		return tview.Escape("#" + a.Name), ""
	default:
		return "", ""
	}
}

func DrawItem(tv *TutView, item api.Item, main *tview.TextView, controls *tview.Flex, ft config.FeedType) {
	switch item.Type() {
	case api.StatusType:
		drawStatus(tv, item, item.Raw().(*mastodon.Status), main, controls, ft, false, "")
	case api.StatusHistoryType:
		s := item.Raw().(*mastodon.StatusHistory)
		status := mastodon.Status{
			Content:          s.Content,
			SpoilerText:      s.SpoilerText,
			Account:          s.Account,
			Sensitive:        s.Sensitive,
			CreatedAt:        s.CreatedAt,
			Emojis:           s.Emojis,
			MediaAttachments: s.MediaAttachments,
			Visibility:       mastodon.VisibilityPublic,
		}
		drawStatus(tv, item, &status, main, controls, ft, true, "")
	case api.UserType, api.ProfileType:
		switch ft {
		case config.FollowRequests:
			drawUser(tv, item.Raw().(*api.User), main, controls, "", InputUserFollowRequest)
		case config.ListUsersAdd:
			drawUser(tv, item.Raw().(*api.User), main, controls, "", InputUserListAdd)
		case config.ListUsersIn:
			drawUser(tv, item.Raw().(*api.User), main, controls, "", InputUserListDelete)
		default:
			drawUser(tv, item.Raw().(*api.User), main, controls, "", InputUserNormal)
		}
	case api.NotificationType:
		drawNotification(tv, item, item.Raw().(*api.NotificationData), main, controls)
	case api.ListsType:
		drawList(tv, item.Raw().(*mastodon.List), main, controls)
	case api.TagType:
		drawTag(tv, item.Raw().(*mastodon.Tag), main, controls)
	}
}

func DrawItemControls(tv *TutView, item api.Item, controls *tview.Flex, ft config.FeedType) {
	switch item.Type() {
	case api.StatusType:
		drawStatus(tv, item, item.Raw().(*mastodon.Status), nil, controls, ft, false, "")
	case api.StatusHistoryType:
		s := item.Raw().(*mastodon.StatusHistory)
		status := mastodon.Status{
			Content:          s.Content,
			SpoilerText:      s.SpoilerText,
			Account:          s.Account,
			Sensitive:        s.Sensitive,
			CreatedAt:        s.CreatedAt,
			Emojis:           s.Emojis,
			MediaAttachments: s.MediaAttachments,
			Visibility:       mastodon.VisibilityPublic,
		}
		drawStatus(tv, item, &status, nil, controls, ft, true, "")
	case api.UserType, api.ProfileType:
		if ft == config.FollowRequests {
			drawUser(tv, item.Raw().(*api.User), nil, controls, "", InputUserFollowRequest)
		} else {
			drawUser(tv, item.Raw().(*api.User), nil, controls, "", InputUserNormal)
		}
	case api.NotificationType:
		drawNotification(tv, item, item.Raw().(*api.NotificationData), nil, controls)
	case api.ListsType:
		drawList(tv, item.Raw().(*mastodon.List), nil, controls)
	case api.TagType:
		drawTag(tv, item.Raw().(*mastodon.Tag), nil, controls)
	}

}

func OutputDate(cfg *config.Config, status time.Time) string {
	today := time.Now()
	ty, tm, td := today.Date()
	sy, sm, sd := status.Date()

	format := cfg.General.DateFormat
	sameDay := false
	displayRelative := false

	if ty == sy && tm == sm && td == sd {
		format = cfg.General.DateTodayFormat
		sameDay = true
	}

	todayFloor := FloorDate(today)
	statusFloor := FloorDate(status)

	if cfg.General.DateRelative > -1 && !sameDay {
		days := int(todayFloor.Sub(statusFloor).Hours() / 24)
		if cfg.General.DateRelative == 0 || days <= cfg.General.DateRelative {
			displayRelative = true
		}
	}
	var dateOutput string
	if displayRelative {
		y, m, d, _, _, _ := timex.Diff(statusFloor, todayFloor)
		if y > 0 {
			dateOutput = fmt.Sprintf("%s%dy", dateOutput, y)
		}
		if dateOutput != "" || m > 0 {
			dateOutput = fmt.Sprintf("%s%dm", dateOutput, m)
		}
		if dateOutput != "" || d > 0 {
			dateOutput = fmt.Sprintf("%s%dd", dateOutput, d)
		}
	} else {
		dateOutput = status.Format(format)
	}
	return dateOutput
}

func FloorDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
