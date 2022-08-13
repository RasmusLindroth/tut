package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

func drawNotification(tv *TutView, item api.Item, notification *api.NotificationData, main *tview.TextView, controls *tview.Flex) {
	switch notification.Item.Type {
	case "follow":
		drawUser(tv, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s started following you", util.FormatUsername(notification.Item.Account)), false,
		)
	case "favourite":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s favorited your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "reblog":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s boosted your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "mention":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s mentioned you", util.FormatUsername(notification.Item.Account)),
		)
	case "status":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s posted a new toot", util.FormatUsername(notification.Item.Account)),
		)
	case "poll":
		drawStatus(tv, notification.Status, notification.Item.Status, main, controls,
			"A poll of yours or one you participated in has ended",
		)
	case "follow_request":
		drawUser(tv, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s  wants to follow you.", util.FormatUsername(notification.Item.Account)),
			true,
		)
	}
}
