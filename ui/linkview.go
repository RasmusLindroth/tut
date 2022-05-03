package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/tut/api"
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
		config.ColorKey(lv.tutView.tut.Config, "", "O", "pen"),
		config.ColorKey(lv.tutView.tut.Config, "", "Y", "ank"),
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

func (lv *LinkView) SetLinks(item api.Item) {
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
			//l.app.UI.CmdBar.ShowError(fmt.Sprintf("Couldn't load user. Error: %v\n", err))
			fmt.Printf("Couldn't load user. Error:%v\n", err)
			os.Exit(1)
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

func (lv *LinkView) Yank() {
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
		copyToClipboard(urls[index].URL)
		return
	}
	mIndex := index - len(urls)
	if mIndex < len(mentions) {
		copyToClipboard(mentions[mIndex].URL)
		return
	}
	tIndex := index - len(mentions) - len(urls)
	if tIndex < len(tags) {
		copyToClipboard(tags[tIndex].URL)
		return
	}
}
