package util

import (
	"bufio"
	"strings"
)

func ReadLine(r *bufio.Reader) (string, error) {
	text, err := r.ReadString('\n')
	if err != nil {
		return text, err
	}
	text = strings.TrimSpace(text)
	return text, err
}
