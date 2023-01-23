package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/atotto/clipboard"
)

func downloadFile(url string) (string, error) {
	ext := filepath.Ext(url)
	fname := "tutfile*"
	if len(ext) > 0 {
		fname = fmt.Sprintf("%s%s", fname, ext)
	}
	f, err := os.CreateTemp("", fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		os.Remove(f.Name())
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		os.Remove(f.Name())
		return "", nil
	}

	return f.Name(), nil
}

func openAvatar(tv *TutView, user mastodon.Account) {
	f, err := downloadFile(user.AvatarStatic)
	if err != nil {
		tv.ShowError(
			fmt.Sprintf("Couldn't open avatar. Error: %v\n", err),
		)
		return
	}
	openMediaType(tv, []string{f}, "image")
}

func reverseFiles(filenames []string) []string {
	if len(filenames) == 0 {
		return filenames
	}
	var f []string
	for i := len(filenames) - 1; i >= 0; i-- {
		f = append(f, filenames[i])
	}
	return f
}

type runProgram struct {
	Name      string
	Filenames []string
	Args      []string
	Terminal  bool
}

func newRunProgram(name string, filenames []string, args ...string) runProgram {
	return runProgram{
		Name:      name,
		Filenames: filenames,
		Args:      args,
	}
}

func openMediaType(tv *TutView, filenames []string, mediaType string) {
	terminal := []runProgram{}
	external := []runProgram{}
	mc := tv.tut.Config.Media
	switch mediaType {
	case "image":
		if mc.ImageReverse {
			filenames = reverseFiles(filenames)
		}
		if mc.ImageSingle {
			for _, f := range filenames {
				args := append(mc.ImageArgs, f)
				c := newRunProgram(mc.ImageViewer, []string{f}, args...)
				if mc.ImageTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(mc.ImageArgs, filenames...)
			c := newRunProgram(mc.ImageViewer, filenames, args...)
			if mc.ImageTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	case "video", "gifv":
		if mc.VideoReverse {
			filenames = reverseFiles(filenames)
		}
		if mc.VideoSingle {
			for _, f := range filenames {
				args := append(mc.VideoArgs, f)
				c := newRunProgram(mc.VideoViewer, []string{f}, args...)
				if mc.VideoTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(mc.VideoArgs, filenames...)
			c := newRunProgram(mc.VideoViewer, filenames, args...)
			if mc.VideoTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	case "audio":
		if mc.AudioReverse {
			filenames = reverseFiles(filenames)
		}
		if mc.AudioSingle {
			for _, f := range filenames {
				args := append(mc.AudioArgs, f)
				c := newRunProgram(mc.AudioViewer, []string{f}, args...)
				if mc.AudioTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(mc.AudioArgs, filenames...)
			c := newRunProgram(mc.AudioViewer, filenames, args...)
			if mc.AudioTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	}
	go func() {
		for _, ext := range external {
			exec.Command(ext.Name, ext.Args...).Run()
			deleteFiles(tv, ext.Filenames)
		}
	}()
	for _, term := range terminal {
		openInTerminal(tv, term.Name, term.Args...)
		deleteFiles(tv, term.Filenames)
	}
}

func deleteFiles(tv *TutView, filenames []string) {
	if tv.tut.Config.Media.DeleteTmpFiles {
		for _, filename := range filenames {
			os.Remove(filename)
		}
	} else {
		tv.FileList = append(tv.FileList, filenames...)
	}
}

func openMedia(tv *TutView, status *mastodon.Status) {
	if status.Reblog != nil {
		status = status.Reblog
	}

	if len(status.MediaAttachments) == 0 {
		return
	}

	mediaGroup := make(map[string][]mastodon.Attachment)
	for _, m := range status.MediaAttachments {
		mediaGroup[m.Type] = append(mediaGroup[m.Type], m)
	}

	for key := range mediaGroup {
		var files []string
		for _, m := range mediaGroup[key] {
			//'image', 'video', 'gifv', 'audio' or 'unknown'
			f, err := downloadFile(m.URL)
			if err != nil {
				continue
			}
			files = append(files, f)
		}
		openMediaType(tv, files, key)
		tv.ShouldSync()
	}
}

func copyToClipboard(text string) bool {
	if clipboard.Unsupported {
		return false
	}
	clipboard.WriteAll(text)
	return true
}
