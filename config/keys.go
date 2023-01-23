package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func ColorFromKey(c *Config, k Key, first bool) (string, int) {
	if len(k.Hint) == 0 {
		return "", 0
	}
	parts := k.Hint[0]
	if !first && len(k.Hint) > 1 {
		parts = k.Hint[1]
	}
	if len(parts) != 3 {
		return "", 0
	}
	return ColorKey(c, parts[0], parts[1], parts[2]), len(fmt.Sprintf("%s%s%s", parts[0], parts[1], parts[2]))
}

func ColorKey(c *Config, pre, key, end string) string {
	color := ColorMark(c.Style.ControlsHighlight)
	normal := ColorMark(c.Style.ControlsText)
	key = TextFlags("b") + key + TextFlags("-")
	if c.General.ShortHints {
		pre = ""
		end = ""
	}
	text := fmt.Sprintf("%s%s%s%s%s%s", normal, pre, color, key, normal, end)
	return text
}

func TextFlags(s string) string {
	return fmt.Sprintf("[::%s]", s)
}

func SublteText(c *Config, text string) string {
	subtle := ColorMark(c.Style.Subtle)
	return fmt.Sprintf("%s%s", subtle, text)
}

func ColorMark(color tcell.Color) string {
	return fmt.Sprintf("[#%06x]", color.Hex())
}
