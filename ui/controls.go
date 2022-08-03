package ui

import "github.com/RasmusLindroth/tut/config"

type Control struct {
	Label string
	Len   int
}

func NewControl(c *config.Config, k config.Key, first bool) Control {
	label, length := config.ColorFromKey(c, k, first)
	return Control{
		Label: label,
		Len:   length,
	}
}
