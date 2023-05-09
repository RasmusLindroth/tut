package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func openURL(tv *TutView, url string) {
	for _, m := range tv.tut.Config.OpenPattern.Patterns {
		if m.Compiled.Match(url) {
			args := append(m.Args, url)
			if m.Terminal {
				openInTerminal(tv, m.Program, args...)
			} else {
				exec.Command(m.Program, args...).Start()
			}
			return
		}
	}
	args := append(tv.tut.Config.Media.LinkArgs, url)
	if tv.tut.Config.Media.LinkTerminal {
		openInTerminal(tv, tv.tut.Config.Media.LinkViewer, args...)
	} else {
		exec.Command(tv.tut.Config.Media.LinkViewer, args...).Start()
	}
}

func openInTerminal(tv *TutView, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var err error
	tv.tut.App.Suspend(func() {
		err = cmd.Run()
	})
	if err != nil {
		tv.ShowError(fmt.Sprintf("Eroror while opening: %v", err))
	}
	return nil
}

func openCustom(tv *TutView, program string, args []string, terminal bool, url string) {
	args = append(args, url)
	if terminal {
		openInTerminal(tv, program, args...)
	} else {
		exec.Command(program, args...).Start()
	}
}

func OpenEditorLengthLimit(tv *TutView, s string, limit int) (string, error) {
	text, err := OpenEditor(tv, s)
	if err != nil {
		return text, err
	}
	s = strings.TrimSpace(text)

	if utf8.RuneCountInString(s) > limit {
		ns := ""
		i := 0
		for _, r := range s {
			if i >= limit {
				break
			}
			ns += string(r)
			i++
		}
		s = ns
	}
	return s, nil
}

func OpenEditor(tv *TutView, content string) (string, error) {
	var editor string
	var exists bool
	if tv.tut.Config.General.Editor == strings.TrimSpace("$EDITOR") {
		editor, exists = os.LookupEnv("EDITOR")
		if !exists || editor == "" {
			editor = "vi"
		}
	} else {
		editor = strings.TrimSpace(tv.tut.Config.General.Editor)
	}
	f, err := os.CreateTemp("", "tut")
	if err != nil {
		return "", err
	}
	if content != "" {
		_, err = f.WriteString(content)
		if err != nil {
			return "", err
		}
	}
	fname := f.Name()
	f.Close()
	cmd := exec.Command("/bin/sh", "-c", editor + " " + f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var text []byte
	tv.tut.App.Suspend(func() {
		cmd.Run()
		text, err = os.ReadFile(fname)
	})
	os.Remove(fname)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(text)), nil
}
