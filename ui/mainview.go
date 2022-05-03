package ui

import (
	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type MainView struct {
	View *tview.Flex
}

func NewMainView(tv *TutView, update chan bool) *MainView {
	mv := &MainView{
		View: mainViewUI(tv),
	}
	go func() {
		for range update {
			tv.tut.App.QueueUpdateDraw(func() {
				*tv.MainView.View = *mainViewUI(tv)
			})
		}
	}()
	return mv
}

func feedList(mv *TutView) *tview.Flex {
	iw := 3
	if !mv.tut.Config.General.ShowIcons {
		iw = 0
	}
	return tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(mv.Timeline.GetFeedList().Text, 0, 1, false).
		AddItem(mv.Timeline.GetFeedList().Symbol, iw, 0, false) //fix so you can hide
}
func notificationList(mv *TutView) *tview.Flex {
	iw := 3
	if !mv.tut.Config.General.ShowIcons {
		iw = 0
	}
	return tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(mv.Timeline.Notifications.List.Text, 0, 1, false).
		AddItem(mv.Timeline.Notifications.List.Symbol, iw, 0, false) //fix so you can hide
}

func mainViewUI(mv *TutView) *tview.Flex {
	showMain := mv.TimelineFocus == FeedFocus
	vl := NewVerticalLine(mv.tut.Config)
	hl := NewHorizontalLine(mv.tut.Config)
	nt := NewTextView(mv.tut.Config)
	lp := mv.tut.Config.General.ListProportion
	cp := mv.tut.Config.General.ContentProportion
	nt.SetTextColor(mv.tut.Config.Style.Subtle)
	nt.SetDynamicColors(false)
	nt.SetText("[N]otifications")

	var list *tview.Flex
	if mv.tut.Config.General.ListSplit == config.ListColumn {
		list = tview.NewFlex().SetDirection(tview.FlexColumn)
	} else {
		list = tview.NewFlex().SetDirection(tview.FlexRow)
	}

	if mv.tut.Config.General.NotificationFeed && !mv.tut.Config.General.HideNotificationText {
		if mv.tut.Config.General.ListSplit == config.ListColumn {
			list.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox(), 1, 0, false).
				AddItem(feedList(mv), 0, 1, false), 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(nt, 1, 0, false).
					AddItem(notificationList(mv), 0, 1, false), 0, 1, false)
		} else {
			list.AddItem(feedList(mv), 0, 1, false).
				AddItem(nt, 1, 0, false).
				AddItem(notificationList(mv), 0, 1, false)
		}

	} else if mv.tut.Config.General.NotificationFeed && mv.tut.Config.General.HideNotificationText {
		if mv.tut.Config.General.ListSplit == config.ListColumn {
			list.AddItem(feedList(mv), 0, 1, false).
				AddItem(notificationList(mv), 0, 1, false)

		} else {
			list.AddItem(feedList(mv), 0, 1, false).
				AddItem(notificationList(mv), 0, 1, false)
		}
	}
	fc := mv.Timeline.GetFeedContent(showMain)
	content := fc.Main
	controls := fc.Controls
	r := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mv.Shared.Top.View, 1, 0, false)
	if mv.tut.Config.General.ListPlacement == config.ListPlacementTop {
		r.AddItem(list, 0, lp, false).
			AddItem(hl, 1, 0, false).
			AddItem(content, 0, cp, false).
			AddItem(controls, 1, 0, false).
			AddItem(mv.Shared.Bottom.View, 2, 0, false)
	} else if mv.tut.Config.General.ListPlacement == config.ListPlacementBottom {
		r.AddItem(content, 0, cp, false).
			AddItem(controls, 1, 0, false).
			AddItem(hl, 1, 0, false).
			AddItem(list, 0, lp, false).
			AddItem(mv.Shared.Bottom.View, 2, 0, false)
	} else if mv.tut.Config.General.ListPlacement == config.ListPlacementLeft {
		r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, lp, false).
			AddItem(vl, 1, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(content, 0, 1, false).
				AddItem(controls, 1, 0, false), 0, cp, false), 0, 1, false).
			AddItem(mv.Shared.Bottom.View, 2, 0, false)
	} else if mv.tut.Config.General.ListPlacement == config.ListPlacementRight {
		r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(content, 0, 1, false).
				AddItem(controls, 1, 0, false), 0, cp, false).
			AddItem(vl, 1, 0, false).
			AddItem(list, 0, lp, false), 0, 1, false).
			AddItem(mv.Shared.Bottom.View, 2, 0, false)
	}
	return r
}
