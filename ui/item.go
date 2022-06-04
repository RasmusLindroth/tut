package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/feed"
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
		if s.Reblog != nil && cfg.General.ShowIcons {
			acc = fmt.Sprintf("♺ %s", acc)
		}
		d := OutputDate(cfg, s.CreatedAt.Local())
		return fmt.Sprintf("%s %s", d, acc), symbol
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
	default:
		return "", ""
	}
}

func DrawItem(tut *Tut, item api.Item, main *tview.TextView, controls *tview.TextView, ft feed.FeedType) {
	switch item.Type() {
	case api.StatusType:
		drawStatus(tut, item, item.Raw().(*mastodon.Status), main, controls, "")
	case api.UserType, api.ProfileType:
		if ft == feed.FollowRequests {
			drawUser(tut, item.Raw().(*api.User), main, controls, "", true)
		} else {
			drawUser(tut, item.Raw().(*api.User), main, controls, "", false)
		}
	case api.NotificationType:
		drawNotification(tut, item, item.Raw().(*api.NotificationData), main, controls)
	case api.ListsType:
		drawList(tut, item.Raw().(*mastodon.List), main, controls)
	}
}

func DrawItemControls(tut *Tut, item api.Item, controls *tview.TextView, ft feed.FeedType) {
	switch item.Type() {
	case api.StatusType:
		drawStatus(tut, item, item.Raw().(*mastodon.Status), nil, controls, "")
	case api.UserType, api.ProfileType:
		if ft == feed.FollowRequests {
			drawUser(tut, item.Raw().(*api.User), nil, controls, "", true)
		} else {
			drawUser(tut, item.Raw().(*api.User), nil, controls, "", false)
		}
	case api.NotificationType:
		drawNotification(tut, item, item.Raw().(*api.NotificationData), nil, controls)
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
