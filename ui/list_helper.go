package ui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

func GetCurrentID(l *tview.List) uint {
	if l.GetItemCount() == 0 {
		return 0
	}
	i := l.GetCurrentItem()
	_, sec := l.GetItemText(i)
	id, err := strconv.ParseUint(sec, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

func SetByID(id uint, l *tview.List) {
	if l.GetItemCount() == 0 {
		return
	}
	s := fmt.Sprintf("%d", id)
	items := l.FindItems("", s, false, false)
	for _, i := range items {
		_, sec := l.GetItemText(i)
		if sec == s {
			l.SetCurrentItem(i)
			break
		}
	}
}
