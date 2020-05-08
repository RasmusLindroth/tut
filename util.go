package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-mastodon"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

type URL struct {
	Text    string
	URL     string
	Classes []string
}

func getURLs(text string) []URL {
	urlReg := regexp.MustCompile(`<a\s+?(.+?)>(.+?)<\/a>`)
	attrReg := regexp.MustCompile(`(\w+?)="(.+?)"`)
	matches := urlReg.FindAllStringSubmatch(text, -1)

	var urls []URL
	for _, m := range matches {
		url := URL{
			Text: m[2],
		}
		attrs := attrReg.FindAllStringSubmatch(m[1], -1)
		if attrs == nil {
			continue
		}
		for _, a := range attrs {
			switch a[1] {
			case "href":
				url.URL = a[2]
			case "class":
				url.Classes = strings.Split(a[2], " ")
			}
		}
		if len(url.Classes) == 0 {
			urls = append(urls, url)
		}
	}
	return urls
}

func cleanTootHTML(content string) (string, []URL) {
	stripped := bluemonday.NewPolicy().AllowElements("p", "br").AllowAttrs("href", "class").OnElements("a").Sanitize(content)
	urls := getURLs(stripped)
	stripped = bluemonday.NewPolicy().AllowElements("p", "br").Sanitize(content)
	stripped = strings.ReplaceAll(stripped, "<br>", "\n")
	stripped = strings.ReplaceAll(stripped, "<br/>", "\n")
	stripped = strings.ReplaceAll(stripped, "<p>", "")
	stripped = strings.ReplaceAll(stripped, "</p>", "\n\n")
	stripped = strings.TrimSpace(stripped)
	stripped = html.UnescapeString(stripped)
	return stripped, urls
}

func openEditor(app *tview.Application, content string) (string, error) {
	editor, exists := os.LookupEnv("EDITOR")
	if !exists || editor == "" {
		editor = "vi"
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
	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	var text []byte
	app.Suspend(func() {
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
		f.Seek(0, 0)
		text, err = ioutil.ReadAll(f)
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(text)), nil
}

func openURL(conf MediaConfig, url string) {
	args := append(conf.LinkArgs, url)
	exec.Command(conf.LinkViewer, args...).Start()
}

func openMediaType(conf MediaConfig, filenames []string, mediaType string) {
	switch mediaType {
	case "image":
		if conf.ImageSingle {
			for _, f := range filenames {
				args := append(conf.ImageArgs, f)
				exec.Command(conf.ImageViewer, args...).Run()
			}
		} else {
			args := append(conf.ImageArgs, filenames...)
			exec.Command(conf.ImageViewer, args...).Run()
		}
	case "video", "gifv":
		if conf.VideoSingle {
			for _, f := range filenames {
				args := append(conf.VideoArgs, f)
				exec.Command(conf.VideoViewer, args...).Run()
			}
		} else {
			args := append(conf.VideoArgs, filenames...)
			exec.Command(conf.VideoViewer, args...).Run()
		}
	case "audio":
		if conf.AudioSingle {
			for _, f := range filenames {
				args := append(conf.AudioArgs, f)
				exec.Command(conf.AudioViewer, args...).Run()
			}
		} else {
			args := append(conf.AudioArgs, filenames...)
			exec.Command(conf.AudioViewer, args...).Run()
		}
	}
}

func downloadFile(url string) (string, error) {
	f, err := ioutil.TempFile("", "tutfile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", nil
	}

	return f.Name(), nil
}

func getConfigDir() string {
	home, _ := os.LookupEnv("HOME")
	xdgConfig, exists := os.LookupEnv("XDG_CONFIG_HOME")
	if !exists {
		xdgConfig = home + "/.config"
	}
	xdgConfig += "/tut"
	return xdgConfig
}

func testConfigPath(name string) (string, error) {
	xdgConfig := getConfigDir()
	path := xdgConfig + "/" + name
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return path, nil
}

func GetAccountsPath() (string, error) {
	return testConfigPath("accounts.yaml")
}

func GetConfigPath() (string, error) {
	return testConfigPath("config.yaml")
}

func CheckPath(input string, inclHidden bool) (string, bool) {
	info, err := os.Stat(input)
	if err != nil {
		return "", false
	}
	if !inclHidden && strings.HasPrefix(info.Name(), ".") {
		return "", false
	}

	if info.IsDir() {
		if input == "/" {
			return input, true
		}
		return input + "/", true
	}
	return input, true
}

func IsDir(input string) bool {
	info, err := os.Stat(input)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FindFiles(s string) []string {
	input := filepath.Clean(s)
	if len(s) > 2 && s[len(s)-2:] == "/." {
		input += "/."
	}
	var files []string
	path, exists := CheckPath(input, true)
	if exists {
		files = append(files, path)
	}

	base := filepath.Base(input)
	inclHidden := strings.HasPrefix(base, ".") || (len(input) > 1 && input[len(input)-2:] == "/.")
	matches, _ := filepath.Glob(input + "*")
	if strings.HasSuffix(path, "/") {
		matchesDir, _ := filepath.Glob(path + "*")
		matches = append(matches, matchesDir...)
	}
	for _, f := range matches {
		p, exists := CheckPath(f, inclHidden)
		if exists && p != path {
			files = append(files, p)
		}
	}
	return files
}

func ColorKey(style StyleConfig, pre, key, end string) string {
	color := ColorMark(style.TextSpecial2)
	normal := ColorMark(style.Text)
	key = tview.Escape("[" + key + "]")
	text := fmt.Sprintf("%s%s%s%s%s", pre, color, key, normal, end)
	return text
}

func ColorMark(color tcell.Color) string {
	return fmt.Sprintf("[#%x]", color.Hex())
}

func FormatUsername(a mastodon.Account) string {
	if a.DisplayName != "" {
		return fmt.Sprintf("%s (%s)", a.DisplayName, a.Acct)
	}
	return a.Acct
}

func SublteText(style StyleConfig, text string) string {
	subtle := ColorMark(style.Subtle)
	return fmt.Sprintf("%s%s", subtle, text)
}
