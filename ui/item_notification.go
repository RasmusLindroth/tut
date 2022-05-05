package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

func drawNotification(tut *Tut, item api.Item, notification *api.NotificationData, main *tview.TextView, controls *tview.TextView) {
	switch notification.Item.Type {
	case "follow":
		drawUser(tut, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s started following you", util.FormatUsername(notification.Item.Account)),
		)
	case "favourite":
		drawStatus(tut, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s favorited your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "reblog":
		drawStatus(tut, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s boosted your toot", util.FormatUsername(notification.Item.Account)),
		)
	case "mention":
		drawStatus(tut, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s mentioned you", util.FormatUsername(notification.Item.Account)),
		)
	case "status":
		drawStatus(tut, notification.Status, notification.Item.Status, main, controls,
			fmt.Sprintf("%s posted a new toot", util.FormatUsername(notification.Item.Account)),
		)
	case "poll":
		drawStatus(tut, notification.Status, notification.Item.Status, main, controls,
			"A poll of yours or one you participated in has ended",
		)
	case "follow_request":
		drawUser(tut, notification.User.Raw().(*api.User), main, controls,
			fmt.Sprintf("%s  wants to follow you. This is currently not implemented, so use another app to accept or reject the request.", util.FormatUsername(notification.Item.Account)),
		)
	}
}
