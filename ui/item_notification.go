package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

func drawNotification(tv *TutView, item api.Item, notification *api.NotificationData, main *tview.TextView, controls *tview.Flex) {
	switch notification.Item.Type {
	case "follow":
		drawUser(tv, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s started following you", util.FormatUsername(notification.Item.Account)), InputUserNormal,
		)
	case "favourite":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			fmt.Sprintf("%s favorited your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "reblog":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			fmt.Sprintf("%s boosted your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "mention":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			fmt.Sprintf("%s mentioned you", util.FormatUsername(notification.Item.Account)),
		)
	case "update":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			fmt.Sprintf("%s updated their toot", util.FormatUsername(notification.Item.Account)),
		)
	case "status":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			fmt.Sprintf("%s posted a new toot", util.FormatUsername(notification.Item.Account)),
		)
	case "poll":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls, false,
			"A poll of yours or one you participated in has ended",
		)
	case "follow_request":
		drawUser(tv, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s  wants to follow you.", util.FormatUsername(notification.Item.Account)),
			InputUserFollowRequest,
		)
	default:
		controls.Clear()
		text := fmt.Sprintf("%s\n", config.SublteText(tv.tut.Config,
			fmt.Sprintf("Notification \"%s\" is not implemented. Open an issue at https://github.com/RasmusLindroth/tut",
				notification.Item.Type),
		))
		main.SetText(text)
	}
}
