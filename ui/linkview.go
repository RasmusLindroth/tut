package ui

import (
	"fmt"
	"strings"

	"github.com/RasmusLindroth/tut/config"
	"github.com/rivo/tview"
)

type LinkView struct {
	tutView  *TutView
	shared   *Shared
	View     *tview.Flex
	list     *tview.List
	controls *tview.TextView
}

func NewLinkView(tv *TutView) *LinkView {
	l := NewList(tv.tut.Config)
	txt := NewTextView(tv.tut.Config)
	lv := &LinkView{
		tutView:  tv,
		shared:   tv.Shared,
		list:     l,
		controls: txt,
	}
	lv.View = linkViewUI(lv)
	return lv
}

func linkViewUI(lv *LinkView) *tview.Flex {
	lv.controls.SetBorderPadding(0, 0, 1, 1)
	items := []string{
		config.ColorFromKey(lv.tutView.tut.Config, lv.tutView.tut.Config.Input.LinkOpen, true),
		config.ColorFromKey(lv.tutView.tut.Config, lv.tutView.tut.Config.Input.LinkYank, true),
	}
	for _, cust := range lv.tutView.tut.Config.OpenCustom.OpenCustoms {
		items = append(items, config.ColorKey(lv.tutView.tut.Config, "", fmt.Sprintf("%d", cust.Index), cust.Name))
	}
	res := strings.Join(items, " ")
	lv.controls.SetText(res)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(lv.shared.Top.View, 1, 0, false).
		AddItem(lv.list, 0, 1, false).
		AddItem(lv.controls, 1, 0, false).
		AddItem(lv.shared.Bottom.View, 2, 0, false)
}

func (lv *LinkView) SetLinks() {
	item, err := lv.tutView.GetCurrentItem()
	if err != nil {
		lv.list.Clear()
		return
	}
	lv.list.Clear()
	urls, mentions, tags, _ := item.URLs()

	for _, url := range urls {
		lv.list.AddItem(url.Text, "", 0, nil)
	}
	for _, mention := range mentions {
		lv.list.AddItem(mention.Acct, "", 0, nil)
	}
	for _, tag := range tags {
		lv.list.AddItem("#"+tag.Name, "", 0, nil)
	}
}

func (lv *LinkView) Next() {
	listNext(lv.list)
}

func (lv *LinkView) Prev() {
	listPrev(lv.list)
}

func (lv *LinkView) Open() {
	item, err := lv.tutView.GetCurrentItem()
	if err != nil {
		return
	}
	urls, mentions, tags, total := item.URLs()
	index := lv.list.GetCurrentItem()

	if total == 0 || index >= total {
		return
	}
	if index < len(urls) {
		openURL(lv.tutView, urls[index].URL)
		return
	}
	mIndex := index - len(urls)
	if mIndex < len(mentions) {
		u, err := lv.tutView.tut.Client.GetUserByID(mentions[mIndex].ID)
		if err != nil {
			lv.tutView.ShowError(
				fmt.Sprintf("Couldn't load user. Error:%v\n", err),
			)
			return
		}
		lv.tutView.Timeline.AddFeed(
			NewUserFeed(lv.tutView, u),
		)
		lv.tutView.FocusMainNoHistory()
		return
	}
	tIndex := index - len(mentions) - len(urls)
	if tIndex < len(tags) {
		lv.tutView.Timeline.AddFeed(
			NewTagFeed(lv.tutView, tags[tIndex].Name),
		)
		lv.tutView.FocusMainNoHistory()
		return
	}
}

func (lv *LinkView) getURL() string {
	item, err := lv.tutView.GetCurrentItem()
	if err != nil {
		return ""
	}
	urls, mentions, tags, total := item.URLs()
	index := lv.list.GetCurrentItem()

	if total == 0 || index >= total {
		return ""
	}
	if index < len(urls) {
		return urls[index].URL
	}
	mIndex := index - len(urls)
	if mIndex < len(mentions) {
		return mentions[mIndex].URL
	}
	tIndex := index - len(mentions) - len(urls)
	if tIndex < len(tags) {
		return tags[tIndex].URL
	}
	return ""
}

func (lv *LinkView) Yank() {
	url := lv.getURL()
	if url == "" {
		return
	}
	copyToClipboard(url)
}

func (lv *LinkView) OpenCustom(index int) {
	url := lv.getURL()
	if url == "" {
		return
	}
	customs := lv.tutView.tut.Config.OpenCustom.OpenCustoms
	for _, c := range customs {
		if c.Index != index {
			continue
		}
		openCustom(lv.tutView, c.Program, c.Args, c.Terminal, url)
		return
	}
}
