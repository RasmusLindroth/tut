package main

import (
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

type MediaFocus int

const (
	MediaFocusOverview MediaFocus = iota
	MediaFocusAdd
)

func NewMediaView(app *App) *MediaView {
	m := &MediaView{
		app:        app,
		Flex:       tview.NewFlex(),
		TextTop:    tview.NewTextView(),
		TextBottom: tview.NewTextView(),
		FileList:   tview.NewList(),
		InputField: &MediaInput{app: app, View: tview.NewInputField()},
		Focus:      MediaFocusOverview,
	}
	m.Draw()
	return m
}

type MediaView struct {
	app        *App
	Flex       *tview.Flex
	TextTop    *tview.TextView
	TextBottom *tview.TextView
	FileList   *tview.List
	InputField *MediaInput
	Focus      MediaFocus
	Files      []string
}

func (m *MediaView) Reset() {
	m.Files = nil
	m.FileList.Clear()
	m.Focus = MediaFocusOverview
	m.Draw()
}

func (m *MediaView) AddFile(f string) {
	m.Files = append(m.Files, f)
	m.FileList.AddItem(filepath.Base(f), "", 0, nil)
	m.Draw()
}

func (m *MediaView) Draw() {
	m.TextTop.SetText("List of attached files:")
	m.TextBottom.SetText("[A]dd file [D]elete file [Esc] Done")
}

func (m *MediaView) SetFocus(f MediaFocus) {
	switch f {
	case MediaFocusOverview:
		m.InputField.View.SetText("")
		m.app.App.SetFocus(m.FileList)
	case MediaFocusAdd:
		m.app.App.SetFocus(m.InputField.View)
		pwd, err := os.Getwd()
		if err != nil {
			home, err := os.UserHomeDir()
			if err != nil {
				pwd = ""
			} else {
				pwd = home
			}
		}
		m.InputField.View.SetText(pwd)
	}
	m.Focus = f
}

func (m *MediaView) Prev() {
	index := m.FileList.GetCurrentItem()
	if index-1 >= 0 {
		m.FileList.SetCurrentItem(index - 1)
	}
}

func (m *MediaView) Next() {
	index := m.FileList.GetCurrentItem()
	if index+1 < m.FileList.GetItemCount() {
		m.FileList.SetCurrentItem(index + 1)
	}
}

func (m *MediaView) Delete() {
	index := m.FileList.GetCurrentItem()
	if len(m.Files) == 0 || index > len(m.Files) {
		return
	}
	m.FileList.RemoveItem(index)
	m.Files = append(m.Files[:index], m.Files[index+1:]...)
}

type MediaInput struct {
	app                  *App
	View                 *tview.InputField
	autocompleteIndex    int
	autocompleteList     []string
	originalText         string
	isAutocompleteChange bool
}

func (m *MediaInput) AddRune(r rune) {
	newText := m.View.GetText() + string(r)
	m.View.SetText(newText)
	m.saveAutocompleteState()
}

func (m *MediaInput) HandleChanges(text string) {
	if m.isAutocompleteChange {
		m.isAutocompleteChange = false
		return
	}
	m.saveAutocompleteState()
}

func (m *MediaInput) saveAutocompleteState() {
	text := m.View.GetText()
	m.originalText = text
	m.autocompleteList = FindFiles(text)
	m.autocompleteIndex = 0
}

func (m *MediaInput) AutocompletePrev() {
	if len(m.autocompleteList) == 0 {
		return
	}
	index := m.autocompleteIndex - 1
	if index < 0 {
		index = len(m.autocompleteList) - 1
	}
	m.autocompleteIndex = index
	m.showAutocomplete()
}

func (m *MediaInput) AutocompleteNext() {
	if len(m.autocompleteList) == 0 {
		return
	}
	index := m.autocompleteIndex + 1
	if index >= len(m.autocompleteList) {
		index = 0
	}
	m.autocompleteIndex = index
	m.showAutocomplete()
}

func (m *MediaInput) CheckDone() {
	path := m.View.GetText()
	if IsDir(path) {
		m.saveAutocompleteState()
		return
	}
	m.app.UI.MediaOverlay.AddFile(path)
	m.app.UI.MediaOverlay.SetFocus(MediaFocusOverview)
}

func (m *MediaInput) showAutocomplete() {
	m.isAutocompleteChange = true
	m.View.SetText(m.autocompleteList[m.autocompleteIndex])
	if len(m.autocompleteList) < 3 {
		m.saveAutocompleteState()
	}
}
