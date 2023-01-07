package ui

import (
	"fmt"

	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type MainView struct {
	View    *tview.Flex
	accView *tview.Flex
	update  chan bool
}

func NewMainView(tv *TutView, update chan bool) *MainView {
	mv := &MainView{
		update:  update,
		accView: NewControlView(tv.tut.Config),
	}
	mv.View = mv.mainViewUI(tv)
	go func() {
		for range mv.update {
			tv.tut.App.QueueUpdateDraw(func() {
				*tv.MainView.View = *mv.mainViewUI(tv)
				tv.ShouldSync()
			})
		}
	}()
	return mv
}

func (mv *MainView) ForceUpdate() {
	mv.update <- true
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

func (mv *MainView) mainViewUI(tv *TutView) *tview.Flex {
	vl := NewVerticalLine(tv.tut.Config)
	hl := NewHorizontalLine(tv.tut.Config)
	lp := tv.tut.Config.General.ListProportion
	cp := tv.tut.Config.General.ContentProportion
	var list *tview.Flex
	if tv.tut.Config.General.ListSplit == config.ListColumn {
		list = tview.NewFlex().SetDirection(tview.FlexColumn)
	} else {
		list = tview.NewFlex().SetDirection(tview.FlexRow)
	}

	if tv.tut.Config.General.ListSplit == config.ListColumn {
		feeds := tview.NewFlex()
		for _, fh := range tv.Timeline.Feeds {
			fTitle := fh.GetTitle()
			if len(fTitle) > 0 {
				txt := NewTextView(tv.tut.Config)
				txt.SetText(tview.Escape(fTitle))
				txt.SetBackgroundColor(tv.tut.Config.Style.TimelineNameBackground)
				txt.SetTextColor(tv.tut.Config.Style.TimelineNameText)
				feeds.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(txt, 1, 0, false).
					AddItem(feedList(tv, fh), 0, 1, false), 0, 1, false)
			} else {
				feeds.AddItem(feedList(tv, fh), 0, 1, false)
			}
		}
		list.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(feeds, 0, 1, false), 0, 1, false)
	} else {
		feeds := tview.NewFlex().SetDirection(tview.FlexRow)
		for _, fh := range tv.Timeline.Feeds {
			fTitle := fh.GetTitle()
			if len(fTitle) > 0 {
				txt := NewTextView(tv.tut.Config)
				txt.SetText(tview.Escape(fTitle))
				txt.SetBackgroundColor(tv.tut.Config.Style.TimelineNameBackground)
				txt.SetTextColor(tv.tut.Config.Style.TimelineNameText)
				feeds.AddItem(txt, 1, 0, false)
			}
			feeds.AddItem(feedList(tv, fh), 0, 1, false)
		}
		list.AddItem(feeds, 0, 1, false)
	}

	fc := tv.Timeline.GetFeedContent()
	content := fc.Main
	controls := fc.Controls

	mv.accView.Clear()
	for i, t := range TutViews.Views {
		acct := t.tut.Client.Me.Acct
		acct = fmt.Sprintf("%s ", acct)
		if i > 0 {
			acct = fmt.Sprintf(" %s", acct)
		}
		item := NewAccButton(tv, tv.tut.Config, acct, i, i == TutViews.Current)
		mv.accView.AddItem(item, len(acct), 0, false)
	}

	r := tview.NewFlex().SetDirection(tview.FlexRow)
	if tv.tut.Config.General.TerminalTitle < 2 {
		r.AddItem(tv.Shared.Top.View, 1, 0, false)
	}
	if tv.tut.Config.General.ListPlacement == config.ListPlacementTop {
		r.AddItem(list, 0, lp, false).
			AddItem(hl, 1, 0, false).
			AddItem(content, 0, cp, false).
			AddItem(controls, 1, 0, false).
			AddItem(tv.Shared.Bottom.View, 2, 0, false)
	} else if tv.tut.Config.General.ListPlacement == config.ListPlacementBottom {
		r.AddItem(content, 0, cp, false).
			AddItem(controls, 1, 0, false).
			AddItem(hl, 1, 0, false).
			AddItem(list, 0, lp, false).
			AddItem(tv.Shared.Bottom.View, 2, 0, false)
	} else if tv.tut.Config.General.ListPlacement == config.ListPlacementLeft {
		r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, lp, false).
			AddItem(vl, 1, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(content, 0, 1, false).
				AddItem(controls, 1, 0, false), 0, cp, false), 0, 1, false).
			AddItem(tv.Shared.Bottom.View, 2, 0, false)
	} else if tv.tut.Config.General.ListPlacement == config.ListPlacementRight {
		r.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(content, 0, 1, false).
				AddItem(controls, 1, 0, false), 0, cp, false).
			AddItem(vl, 1, 0, false).
			AddItem(list, 0, lp, false), 0, 1, false).
			AddItem(tv.Shared.Bottom.View, 2, 0, false)
	}
	if len(TutViews.Views) > 1 {
		r.AddItem(mv.accView, 1, 0, false)
	}
	return r
}
