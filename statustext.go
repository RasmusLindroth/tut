package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/tview"
)

func NewStatusText(app *App, view *tview.TextView, controls *Controls, lo *LinkOverlay) *StatusText {
	return &StatusText{
		app:         app,
		Index:       0,
		View:        view,
		Controls:    controls,
		LinkOverlay: lo,
	}
}

type StatusText struct {
	app         *App
	Index       int
	View        *tview.TextView
	Controls    *Controls
	LinkOverlay *LinkOverlay
}

func (s *StatusText) ShowToot(index int) {
	s.ShowTootOptions(index, false)
}

func (s *StatusText) ShowTootOptions(index int, showSensitive bool) {
	status, err := s.app.UI.TootList.GetStatus(index)
	if err != nil {
		log.Fatalln(err)
	}

	var line string
	_, _, width, _ := s.View.GetInnerRect()
	for i := 0; i < width; i++ {
		line += "-"
	}
	line += "\n"

	shouldDisplay := !status.Sensitive || showSensitive

	var stripped string
	var urls []URL
	var u []URL
	if status.Sensitive && !showSensitive {
		stripped, u = cleanTootHTML(status.SpoilerText)
		urls = append(urls, u...)
		stripped += "\n" + line
		stripped += "Press [s] to show hidden text"

	} else {
		stripped, u = cleanTootHTML(status.Content)
		urls = append(urls, u...)

		if status.Sensitive {
			sens, u := cleanTootHTML(status.SpoilerText)
			urls = append(urls, u...)
			stripped = sens + "\n\n" + stripped
		}
	}
	s.LinkOverlay.SetURLs(urls)

	var head string
	if status.Reblog != nil {
		if status.Account.DisplayName != "" {
			head += fmt.Sprintf("[gray]%s (%s)\n", status.Account.DisplayName, status.Account.Acct)
		} else {
			head += fmt.Sprintf("[gray]%s\n", status.Account.Acct)
		}
		head += "[gray]Boosted\n"
		head += "[gray]" + line
		status = status.Reblog
	}
	if status.Account.DisplayName != "" {
		head += fmt.Sprintf("[tomato]%s\n", status.Account.DisplayName)
	}
	head += fmt.Sprintf("[yellow]%s\n\n", status.Account.Acct)
	output := head
	content := tview.Escape(stripped)
	if content != "" {
		output += content + "\n\n"
	}

	var poll string
	if status.Poll != nil {
		poll += "[gray]Poll\n"
		poll += "[gray]" + line
		poll += fmt.Sprintf("Number of votes: %d\n\n", status.Poll.VotesCount)
		votes := float64(status.Poll.VotesCount)
		for _, o := range status.Poll.Options {
			res := 0.0
			if votes != 0 {
				res = float64(o.VotesCount) / votes * 100
			}
			poll += fmt.Sprintf("%s - %.2f%% (%d)\n", tview.Escape(o.Title), res, o.VotesCount)
		}
		poll += "\n"
	}

	var media string
	for _, att := range status.MediaAttachments {
		media += "[gray]" + line
		media += fmt.Sprintf("[gray]Attached %s\n", att.Type)
		media += fmt.Sprintf("%s\n", att.URL)
	}

	var card string
	if status.Card != nil {
		card += "[gray]Card type: " + status.Card.Type + "\n"
		card += "[gray]" + line
		if status.Card.Title != "" {
			card += status.Card.Title + "\n\n"
		}
		desc := strings.TrimSpace(status.Card.Description)
		if desc != "" {
			card += desc + "\n\n"
		}
		card += status.Card.URL
	}

	if shouldDisplay {
		output += poll + media + card
	}

	s.View.SetText(output)
	s.View.ScrollToBeginning()

	info := "[B]oost [F]avorite [T]hread [R]eply [V]iew [M]edia [O]ther"
	s.Controls.View.SetText(tview.Escape(info))
}
