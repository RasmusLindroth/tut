package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/rivo/tview"
)

type Toot struct {
	Visibility         string
	Boosted            bool
	BoostedDisplayName string
	BoostedAcct        string
	Bookmarked         bool
	AccountDisplayName string
	Account            string
	Spoiler            bool
	SpoilerText        string
	ShowSpoiler        bool
	ContentText        string
	Width              int
	HasExtra           bool
	Poll               Poll
	Media              []Media
	Card               Card
	Replies            int
	Boosts             int
	Favorites          int
	Controls           string
}

type Poll struct {
	ID         string
	ExpiresAt  time.Time
	Expired    bool
	Multiple   bool
	VotesCount int64
	Options    []PollOption
	Voted      bool
}

type PollOption struct {
	Title      string
	VotesCount int64
	Percent    string
}

type Media struct {
	Type        string
	Description string
	URL         string
}

type Card struct {
	Type        string
	Title       string
	Description string
	URL         string
}

type DisplayTootData struct {
	Toot  Toot
	Style config.Style
}

func drawStatus(tut *Tut, item api.Item, status *mastodon.Status, main *tview.TextView, controls *tview.TextView, additional string) {
	showSensitive := item.ShowSpoiler()

	var strippedContent string
	var strippedSpoiler string

	so := status
	if status.Reblog != nil {
		status = status.Reblog
	}

	strippedContent, _ = util.CleanHTML(status.Content)
	strippedContent = tview.Escape(strippedContent)

	width := 0
	if main != nil {
		_, _, width, _ = main.GetInnerRect()
	}
	toot := Toot{
		Width:              width,
		ContentText:        strippedContent,
		Boosted:            so.Reblog != nil,
		BoostedDisplayName: tview.Escape(so.Account.DisplayName),
		BoostedAcct:        tview.Escape(so.Account.Acct),
		ShowSpoiler:        showSensitive,
	}

	toot.AccountDisplayName = tview.Escape(status.Account.DisplayName)
	toot.Account = tview.Escape(status.Account.Acct)
	toot.Bookmarked = status.Bookmarked
	toot.Visibility = status.Visibility
	toot.Spoiler = status.Sensitive

	if status.Poll != nil {
		p := *status.Poll
		toot.Poll = Poll{
			ID:         string(p.ID),
			ExpiresAt:  p.ExpiresAt,
			Expired:    p.Expired,
			Multiple:   p.Multiple,
			VotesCount: p.VotesCount,
			Voted:      p.Voted,
			Options:    []PollOption{},
		}
		for _, item := range p.Options {
			percent := 0.0
			if p.VotesCount > 0 {
				percent = float64(item.VotesCount) / float64(p.VotesCount) * 100
			}

			o := PollOption{
				Title:      tview.Escape(item.Title),
				VotesCount: item.VotesCount,
				Percent:    fmt.Sprintf("%.2f", percent),
			}
			toot.Poll.Options = append(toot.Poll.Options, o)
		}

	} else {
		toot.Poll = Poll{}
	}

	if status.Sensitive {
		strippedSpoiler, _ = util.CleanHTML(status.SpoilerText)
		strippedSpoiler = tview.Escape(strippedSpoiler)
	}

	toot.SpoilerText = strippedSpoiler

	media := []Media{}
	for _, att := range status.MediaAttachments {
		m := Media{
			Type:        att.Type,
			Description: tview.Escape(att.Description),
			URL:         att.URL,
		}
		media = append(media, m)
	}
	toot.Media = media

	if status.Card != nil {
		toot.Card = Card{
			Type:        status.Card.Type,
			Title:       tview.Escape(strings.TrimSpace(status.Card.Title)),
			Description: tview.Escape(strings.TrimSpace(status.Card.Description)),
			URL:         status.Card.URL,
		}
	} else {
		toot.Card = Card{}
	}

	toot.HasExtra = len(status.MediaAttachments) > 0 || status.Card != nil || status.Poll != nil
	toot.Replies = int(status.RepliesCount)
	toot.Boosts = int(status.ReblogsCount)
	toot.Favorites = int(status.FavouritesCount)

	if main != nil {
		main.ScrollToBeginning()
	}

	var info []string
	if status.Favourited {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusFavorite, false))
	} else {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusFavorite, true))
	}
	if status.Reblogged {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusBoost, false))
	} else {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusBoost, true))
	}
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusThread, true))
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusReply, true))
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusViewFocus, true))
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusUser, true))
	if len(status.MediaAttachments) > 0 {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusMedia, true))
	}
	_, _, _, length := item.URLs()
	if length > 0 {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusLinks, true))
	}
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusAvatar, true))
	if status.Account.ID == tut.Client.Me.ID {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusDelete, true))
	}

	if !status.Bookmarked {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusBookmark, true))
	} else {
		info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusBookmark, false))
	}
	info = append(info, config.ColorFromKey(tut.Config, tut.Config.Input.StatusYank, true))

	controlsS := strings.Join(info, " ")

	td := DisplayTootData{
		Toot:  toot,
		Style: tut.Config.Style,
	}
	var output bytes.Buffer
	err := tut.Config.Templates.Toot.ExecuteTemplate(&output, "toot.tmpl", td)
	if err != nil {
		panic(err)
	}

	if main != nil {
		if additional != "" {
			additional = fmt.Sprintf("%s\n\n", config.SublteText(tut.Config, additional))
		}
		main.SetText(additional + output.String())
	}
	controls.SetText(controlsS)
}
