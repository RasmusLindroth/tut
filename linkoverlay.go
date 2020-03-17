package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/tview"
)

func NewLinkOverlay(app *App, view *tview.TextView, controls *Controls) *LinkOverlay {
	return &LinkOverlay{
		app:      app,
		View:     view,
		Controls: controls,
		hints:    generateHints(),
		Scroll:   false,
	}
}

type LinkOverlay struct {
	app      *App
	View     *tview.TextView
	Controls *Controls
	URLs     []URL
	hints    []string
	input    string
	Scroll   bool
}

func generateHints() []string {
	//TODO: REMOVE T
	chars := strings.Split("asdfghjkl", "")
	if len(chars) == 0 {
		chars = strings.Split("asdfghjkl", "")
	}

	var one []string
	var two []string
	var three []string

	for _, a := range chars {
		one = append(one, a)
		for _, b := range chars {
			if b == a {
				continue
			}
			two = append(two, a+b)
			for _, c := range chars {
				if c == b {
					continue
				}
				three = append(three, a+b+c)
			}
		}
	}
	return append(one, append(two, three...)...)
}

func (l *LinkOverlay) SetURLs(urls []URL) {
	l.DisableScroll()
	l.input = ""
	l.URLs = urls
}

func (l *LinkOverlay) Draw() {
	if len(l.hints) < len(l.URLs) {
		log.Fatalln("No hints")
		return
	}
	var output string
	for i, url := range l.URLs {
		hint := tview.Escape(
			fmt.Sprintf("[%s]", l.hints[i]),
		)
		output += fmt.Sprintf("%s %s\n", hint,
			tview.Escape(url.URL))
	}

	l.View.SetText(output)
	if l.Scroll {
		l.Controls.View.SetText(tview.Escape("\n--SCROLL-- [t]oggle scroll"))
	} else {
		l.Controls.View.SetText(tview.Escape("\n--OPEN-- [t]oggle scroll"))
	}
}

func (l *LinkOverlay) AddRune(r rune) {
	l.input += string(r)
	l.Controls.View.SetText(tview.Escape(l.input + "\n--OPEN-- [t]oggle scroll"))
	for i, key := range l.hints {
		if key == l.input && i < len(l.URLs) {
			openURL(l.URLs[i].URL)
			l.Clear()
			l.app.UI.SetFocus(LeftPaneFocus)
		}
	}
}

func (l *LinkOverlay) ActivateScroll() {
	l.Scroll = true
	l.Controls.View.SetText(tview.Escape("\n--SCROLL-- [t]oggle scroll"))
}

func (l *LinkOverlay) DisableScroll() {
	l.Scroll = false
	l.input = ""
	l.Controls.View.SetText(tview.Escape("\n--OPEN-- [t]oggle scroll"))
}

func (l *LinkOverlay) HasInput() bool {
	return l.input != ""
}

func (l *LinkOverlay) Clear() {
	l.input = ""
	l.Controls.View.SetText(tview.Escape(l.input + "\n--OPEN-- [t]oggle scroll"))
}
