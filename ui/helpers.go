package ui

import "github.com/rivo/tview"

func listNext(l *tview.List) (loadOlder bool) {
	ni := l.GetCurrentItem() + 1
	if ni >= l.GetItemCount() {
		ni = l.GetItemCount() - 1
		if ni < 0 {
			ni = 0
		}
	}
	l.SetCurrentItem(ni)
	return l.GetItemCount()-(ni+1) < 5
}

func listPrev(l *tview.List) (loadNewer bool) {
	ni := l.GetCurrentItem() - 1
	if ni < 0 {
		ni = 0
	}
	l.SetCurrentItem(ni)
	return ni < 4
}
