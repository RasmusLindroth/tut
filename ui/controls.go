package ui

import (
	"github.com/RasmusLindroth/tut/config"
	"github.com/gdamore/tcell/v2"
)

type Control struct {
	key   config.Key
	Label string
	Len   int
}

func NewControl(c *config.Config, k config.Key, first bool) Control {
	label, length := config.ColorFromKey(c, k, first)
	return Control{
		key:   k,
		Label: label,
		Len:   length,
	}
}

func (c Control) Click() *tcell.EventKey {
	for _, k := range c.key.Keys {
		return tcell.NewEventKey(k, 0, tcell.ModNone)
	}
	for _, r := range c.key.Runes {
		return tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
	}
	return tcell.NewEventKey(tcell.KeyRune, 0, tcell.ModNone)
}
