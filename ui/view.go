package ui

import (
	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
)

type PageFocusAt uint

const (
	LoginFocus PageFocusAt = iota
	MainFocus
	ViewFocus
	ModalFocus
	LinkFocus
	ComposeFocus
	MediaFocus
	MediaAddFocus
	CmdFocus
	VoteFocus
	HelpFocus
)

func (tv *TutView) GetCurrentFeed() *Feed {
	foc := tv.TimelineFocus
	if foc == FeedFocus {
		return tv.Timeline.Feeds[tv.Timeline.FeedIndex]
	}
	return tv.Timeline.Notifications
}

func (tv *TutView) GetCurrentItem() (api.Item, error) {
	f := tv.GetCurrentFeed()
	return f.Data.Item(f.List.Text.GetCurrentItem())
}

func (tv *TutView) RedrawContent() {
	f := tv.GetCurrentFeed()
	item, err := f.Data.Item(f.List.Text.GetCurrentItem())
	if err != nil {
		f.Content.Main.SetText("")
		f.Content.Controls.SetText("")
		return
	}
	DrawItem(tv.tut, item, f.Content.Main, f.Content.Controls, f.Data.Type())
}
func (tv *TutView) RedrawPoll(poll *mastodon.Poll) {
	f := tv.GetCurrentFeed()
	item, err := f.Data.Item(f.List.Text.GetCurrentItem())
	if err != nil {
		return
	}
	if item.Type() != api.StatusType {
		tv.RedrawContent()
		return
	}
	so := item.Raw().(*mastodon.Status)
	if so.Reblog != nil {
		so.Reblog.Poll = poll
	} else {
		so.Poll = poll
	}
	DrawItem(tv.tut, item, f.Content.Main, f.Content.Controls, f.Data.Type())
}
func (tv *TutView) RedrawControls() {
	f := tv.GetCurrentFeed()
	item, err := f.Data.Item(f.List.Text.GetCurrentItem())
	if err != nil {
		return
	}
	DrawItemControls(tv.tut, item, f.Content.Controls, f.Data.Type())
}

func (tv *TutView) SetPage(f PageFocusAt) {
	if f == tv.PageFocus {
		return
	}
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
	case ViewFocus:
		f := tv.GetCurrentFeed()
		tv.PageFocus = ViewFocus
		tv.Shared.Bottom.StatusBar.SetMode(ScrollMode)
		tv.tut.App.SetFocus(f.Content.Main)
	case LinkFocus:
		tv.LinkView.SetLinks()
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
	case VoteFocus:
		tv.PageFocus = VoteFocus
		tv.View.SwitchToPage("vote")
		tv.tut.App.SetFocus(tv.View)
		tv.Shared.Bottom.StatusBar.SetMode(VoteMode)
		tv.Shared.Top.SetText("vote on poll")
	case HelpFocus:
		tv.PageFocus = HelpFocus
		tv.View.SwitchToPage("help")
		tv.Shared.Bottom.StatusBar.SetMode(HelpMode)
		tv.tut.App.SetFocus(tv.HelpView.content)
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

func (tv *TutView) ShowError(s string) {
	tv.Shared.Bottom.Cmd.ShowError(s)
}

func (tv *TutView) ShouldSync() {
	if !tv.tut.Config.General.RedrawUI {
		return
	}
	tv.tut.App.Sync()
}
