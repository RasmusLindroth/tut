package ui

import (
	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
)

type PageFocusAt uint

const (
	LoginFocus PageFocusAt = iota
	MainFocus
	ModalFocus
	LinkFocus
	ComposeFocus
	MediaFocus
	MediaAddFocus
	CmdFocus
)

func (tv *TutView) GetCurrentItem() (api.Item, error) {
	foc := tv.TimelineFocus
	var f *Feed
	if foc == FeedFocus {
		f = tv.Timeline.Feeds[tv.Timeline.FeedIndex]
	} else {
		f = tv.Timeline.Notifications
	}
	return f.Data.Item(f.List.Text.GetCurrentItem())
}

func (tv *TutView) RedrawContent() {
	foc := tv.TimelineFocus
	var f *Feed
	if foc == FeedFocus {
		f = tv.Timeline.Feeds[tv.Timeline.FeedIndex]
	} else {
		f = tv.Timeline.Notifications
	}
	item, err := f.Data.Item(f.List.Text.GetCurrentItem())
	if err != nil {
		return
	}
	DrawItem(tv.tut, item, f.Content.Main, f.Content.Controls)
}
func (tv *TutView) RedrawControls() {
	foc := tv.TimelineFocus
	var f *Feed
	if foc == FeedFocus {
		f = tv.Timeline.Feeds[tv.Timeline.FeedIndex]
	} else {
		f = tv.Timeline.Notifications
	}
	item, err := f.Data.Item(f.List.Text.GetCurrentItem())
	if err != nil {
		return
	}
	DrawItemControls(tv.tut, item, f.Content.Controls)
}

func (tv *TutView) SetPage(f PageFocusAt) {
	tv.PrevPageFocus = tv.PageFocus
	if tv.PrevPageFocus == LoginFocus {
		tv.PrevPageFocus = MainFocus
	}
	switch f {
	case LoginFocus:
		tv.PageFocus = LoginFocus
		tv.View.SwitchToPage("login")
		tv.Shared.Bottom.StatusBar.SetMode(UserMode)
		tv.Shared.Top.SetText("select accouth with <Enter>")
		tv.tut.App.SetFocus(tv.View)
	case MainFocus:
		tv.PageFocus = MainFocus
		tv.View.SwitchToPage("main")
		tv.Shared.Bottom.StatusBar.SetMode(ListMode)
		tv.Shared.Top.SetText(tv.Timeline.GetTitle())
		tv.tut.App.SetFocus(tv.View)
	case LinkFocus:
		tv.PageFocus = LinkFocus
		tv.View.SwitchToPage("link")
		tv.Shared.Bottom.StatusBar.SetMode(ListMode)
		tv.Shared.Top.SetText("select link with <Enter>")
		tv.tut.App.SetFocus(tv.View)
	case ComposeFocus:
		tv.PageFocus = ComposeFocus
		tv.View.SwitchToPage("compose")
		tv.Shared.Bottom.StatusBar.SetMode(ComposeMode)
		tv.Shared.Top.SetText("write a toot")
		tv.ComposeView.SetControls(ComposeNormal)
		tv.tut.App.SetFocus(tv.ComposeView.content)
	case MediaFocus:
		tv.PageFocus = MediaFocus
		tv.ComposeView.SetControls(ComposeMedia)
		tv.tut.App.SetFocus(tv.View)
	case MediaAddFocus:
		tv.PageFocus = MediaAddFocus
		tv.tut.App.SetFocus(tv.ComposeView.input.View)
	case CmdFocus:
		tv.PageFocus = CmdFocus
		tv.tut.App.SetFocus(tv.Shared.Bottom.Cmd.View)
		tv.Shared.Bottom.StatusBar.SetMode(CmdMode)
		tv.Shared.Bottom.Cmd.ClearInput()
	case ModalFocus:
		tv.PageFocus = ModalFocus
		tv.View.SwitchToPage("modal")
		tv.tut.App.SetFocus(tv.ModalView.View)

	}
	tv.ShouldSync()
}

func (tv *TutView) FocusMainNoHistory() {
	tv.SetPage(MainFocus)
	tv.PrevPageFocus = MainFocus
}

func (tv *TutView) PrevFocus() {
	tv.SetPage(tv.PrevPageFocus)
	tv.PrevPageFocus = MainFocus
}

func (tv *TutView) InitPost(status *mastodon.Status) {
	tv.ComposeView.SetStatus(status)
	tv.SetPage(ComposeFocus)
}

func (tv *TutView) ShouldSync() {
	if !tv.tut.Config.General.RedrawUI {
		return
	}
	tv.tut.App.Sync()
}
