package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/tview"
)

func NewTootView(app *App) *TootView {
	t := &TootView{
		app:      app,
		Index:    0,
		Text:     tview.NewTextView(),
		Controls: tview.NewTextView(),
	}

	t.Text.SetWordWrap(true).SetDynamicColors(true)
	t.Text.SetBackgroundColor(app.Config.Style.Background)
	t.Text.SetTextColor(app.Config.Style.Text)
	t.Controls.SetDynamicColors(true)
	t.Controls.SetBackgroundColor(app.Config.Style.Background)

	return t
}

type TootView struct {
	app      *App
	Index    int
	Text     *tview.TextView
	Controls *tview.TextView
}

func (s *TootView) ShowToot(index int) {
	s.ShowTootOptions(index, false)
}

func (s *TootView) ShowTootOptions(index int, showSensitive bool) {
	status, err := s.app.UI.TootList.GetStatus(index)
	if err != nil {
		log.Fatalln(err)
	}

	var line string
	_, _, width, _ := s.Text.GetInnerRect()
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
	s.app.UI.LinkOverlay.SetURLs(urls)

	subtleColor := fmt.Sprintf("[#%x]", s.app.Config.Style.Subtle.Hex())
	special1 := fmt.Sprintf("[#%x]", s.app.Config.Style.TextSpecial1.Hex())
	special2 := fmt.Sprintf("[#%x]", s.app.Config.Style.TextSpecial2.Hex())
	var head string
	if status.Reblog != nil {
		if status.Account.DisplayName != "" {
			head += fmt.Sprintf(subtleColor+"%s (%s)\n", status.Account.DisplayName, status.Account.Acct)
		} else {
			head += fmt.Sprintf(subtleColor+"%s\n", status.Account.Acct)
		}
		head += subtleColor + "Boosted\n"
		head += subtleColor + line
		status = status.Reblog
	}

	if status.Account.DisplayName != "" {
		head += fmt.Sprintf(special2+"%s\n", status.Account.DisplayName)
	}
	head += fmt.Sprintf(special1+"%s\n\n", status.Account.Acct)
	output := head
	content := tview.Escape(stripped)
	if content != "" {
		output += content + "\n\n"
	}

	var poll string
	if status.Poll != nil {
		poll += subtleColor + "Poll\n"
		poll += subtleColor + line
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
		media += subtleColor + line
		media += fmt.Sprintf(subtleColor+"Attached %s\n", att.Type)
		media += fmt.Sprintf("%s\n", att.URL)
	}

	var card string
	if status.Card != nil {
		card += subtleColor + "Card type: " + status.Card.Type + "\n"
		card += subtleColor + line
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

	s.Text.SetText(output)
	s.Text.ScrollToBeginning()
	var info []string
	if status.Favourited == true {
		info = append(info, "Un[F]avorite")
	} else {
		info = append(info, "[F]avorite")
	}
	if status.Reblogged == true {
		info = append(info, "Un[B]oost")
	} else {
		info = append(info, "[B]oost")
	}
	info = append(info, "[T]hread", "[R]eply", "[V]iew")
	if len(status.MediaAttachments) > 0 {
		info = append(info, "[M]edia")
	}
	if len(urls) > 0 {
		info = append(info, "[O]pen")
	}

	if status.Account.ID == s.app.Me.ID {
		info = append(info, "[D]elete")
	}

	s.Controls.SetText(tview.Escape(strings.Join(info, " ")))
}
