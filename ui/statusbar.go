package ui

import "github.com/rivo/tview"

type StatusBar struct {
	tutView *TutView
	View    *tview.TextView
}

func NewStatusBar(tv *TutView) *StatusBar {
	sb := &StatusBar{
		tutView: tv,
		View:    NewTextView(tv.tut.Config),
	}
	sb.View.SetBackgroundColor(tv.tut.Config.Style.StatusBarBackground)
	sb.View.SetTextColor(tv.tut.Config.Style.StatusBarText)
	return sb
}

type ViewMode uint

const (
	CmdMode ViewMode = iota
	ComposeMode
	HelpMode
	EditorMode
	LinkMode
	ListMode
	MediaMode
	NotificationsMode
	ScrollMode
	UserMode
	VoteMode
	PollMode
	PreferenceMode
)

func (sb *StatusBar) SetMode(m ViewMode) {
	sb.View.SetBackgroundColor(sb.tutView.tut.Config.Style.StatusBarBackground)
	sb.View.SetTextColor(sb.tutView.tut.Config.Style.StatusBarText)
	switch m {
	case CmdMode:
		sb.View.SetText("-- CMD --")
	case ComposeMode:
		sb.View.SetText("-- COMPOSE --")
	case HelpMode:
		sb.View.SetText("-- HELP --")
	case LinkMode:
		sb.View.SetText("-- LINK --")
	case ListMode:
		sb.View.SetText("-- LIST --")
	case EditorMode:
		sb.View.SetText("-- EDITOR --")
	case MediaMode:
		sb.View.SetText("-- MEDIA --")
	case NotificationsMode:
		sb.View.SetText("-- NOTIFICATIONS --")
	case VoteMode:
		sb.View.SetText("-- VOTE --")
	case ScrollMode:
		sb.View.SetBackgroundColor(sb.tutView.tut.Config.Style.StatusBarViewBackground)
		sb.View.SetTextColor(sb.tutView.tut.Config.Style.StatusBarViewText)
		sb.View.SetText("-- VIEW --")
	case UserMode:
		sb.View.SetText("-- SELECT USER --")
	case PollMode:
		sb.View.SetText("-- CREATE POLL --")
	case PreferenceMode:
		sb.View.SetText("-- PREFERENCES --")
	}
}
