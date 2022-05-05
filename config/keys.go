package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func ColorKey(c *Config, pre, key, end string) string {
	color := ColorMark(c.Style.TextSpecial2)
	normal := ColorMark(c.Style.Text)
	key = TextFlags("b") + key + TextFlags("-")
	if c.General.ShortHints {
		pre = ""
		end = ""
	}
	text := fmt.Sprintf("%s%s%s%s%s%s", normal, pre, color, key, normal, end)
	return text
}

func ColorMark(color tcell.Color) string {
	return fmt.Sprintf("[#%06x]", color.Hex())
}

func TextFlags(s string) string {
	return fmt.Sprintf("[::%s]", s)
}

func SublteText(c *Config, text string) string {
	subtle := ColorMark(c.Style.Subtle)
	return fmt.Sprintf("%s%s", subtle, text)
}
