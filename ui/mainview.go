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
				tv.ShouldSync()
			})
		}
	}()
	return mv
}

func feedList(mv *TutView, fh *FeedHolder) *tview.Flex {
	iw := 3
	if !mv.tut.Config.General.ShowIcons {
		iw = 0
	}
	return tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(fh.GetFeedList().Text, 0, 1, false).
		AddItem(fh.GetFeedList().Symbol, iw, 0, false)
}

func mainViewUI(mv *TutView) *tview.Flex {
	vl := NewVerticalLine(mv.tut.Config)
	hl := NewHorizontalLine(mv.tut.Config)
	lp := mv.tut.Config.General.ListProportion
	cp := mv.tut.Config.General.ContentProportion
	var list *tview.Flex
	if mv.tut.Config.General.ListSplit == config.ListColumn {
		list = tview.NewFlex().SetDirection(tview.FlexColumn)
	} else {
		list = tview.NewFlex().SetDirection(tview.FlexRow)
	}

	if mv.tut.Config.General.ListSplit == config.ListColumn {
		feeds := tview.NewFlex()
		for _, fh := range mv.Timeline.Feeds {
			if mv.tut.Config.General.TimelineName && len(fh.Name) > 0 {
				txt := NewTextView(mv.tut.Config)
				txt.SetText(tview.Escape(fh.Name))
				txt.SetTextColor(mv.tut.Config.Style.Subtle)
				feeds.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(txt, 1, 0, false).
					AddItem(feedList(mv, fh), 0, 1, false), 0, 1, false)
			} else {
				feeds.AddItem(feedList(mv, fh), 0, 1, false)
			}
		}
		list.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(feeds, 0, 1, false), 0, 1, false)
	} else {
		feeds := tview.NewFlex().SetDirection(tview.FlexRow)
		for _, fh := range mv.Timeline.Feeds {
			if mv.tut.Config.General.TimelineName && len(fh.Name) > 0 {
				txt := NewTextView(mv.tut.Config)
				txt.SetText(tview.Escape(fh.Name))
				txt.SetTextColor(mv.tut.Config.Style.Subtle)
				feeds.AddItem(txt, 1, 0, false)
			}
			feeds.AddItem(feedList(mv, fh), 0, 1, false)
		}
		list.AddItem(feeds, 0, 1, false)
	}

	fc := mv.Timeline.GetFeedContent()
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
