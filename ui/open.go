package ui

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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
		if err != nil {
			log.Fatalln(err)
		}
	})
	return err
}

func openCustom(tv *TutView, program string, args []string, terminal bool, url string) {
	args = append(args, url)
	if terminal {
		openInTerminal(tv, program, args...)
	} else {
		exec.Command(program, args...).Start()
	}
}

func OpenEditor(tv *TutView, content string) (string, error) {
	editor, exists := os.LookupEnv("EDITOR")
	if !exists || editor == "" {
		editor = "vi"
	}
	args := []string{}
	parts := strings.Split(editor, " ")
	if len(parts) > 1 {
		args = append(args, parts[1:]...)
		editor = parts[0]
	}
	f, err := ioutil.TempFile("", "tut")
	if err != nil {
		return "", err
	}
	if content != "" {
		_, err = f.WriteString(content)
		if err != nil {
			return "", err
		}
	}
	args = append(args, f.Name())
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var text []byte
	tv.tut.App.Suspend(func() {
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
		f.Seek(0, 0)
		text, err = ioutil.ReadAll(f)
	})
	f.Close()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(text)), nil
}
