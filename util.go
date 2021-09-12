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
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/gen2brain/beeep"
	"github.com/icza/gox/timex"
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

//Runs commands prefixed !CMD!
func CmdToString(cmd string) (string, error) {
	cmd = strings.TrimPrefix(cmd, "!CMD!")
	parts := strings.Split(cmd, " ")
	s, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	return string(s), err
}

func getURLs(text string) []URL {
	doc := html.NewTokenizer(strings.NewReader(text))
	var urls []URL

	for {
		n := doc.Next()
		switch n {
		case html.ErrorToken:
			return urls

		case html.StartTagToken:
			token := doc.Token()
			if token.Data == "a" {
				url := URL{}
				var appendUrl = true
				for _, a := range token.Attr {
					switch a.Key {
					case "href":
						url.URL = a.Val
						url.Text = a.Val
					case "class":
						url.Classes = strings.Split(a.Val, " ")
						if strings.Contains(a.Val, "hashtag") ||
							strings.Contains(a.Val, "mention") {
							appendUrl = false
						}
					}
				}
				if appendUrl {
					urls = append(urls, url)
				}
			}
		}
	}
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
	cmd.Stderr = os.Stderr
	var text []byte
	app.Suspend(func() {
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

func copyToClipboard(text string) bool {
	if clipboard.Unsupported {
		return false
	}
	clipboard.WriteAll(text)
	return true
}

func openCustom(program string, args []string, url string) {
	args = append(args, url)
	exec.Command(program, args...).Start()
}

func openURL(conf MediaConfig, pc OpenPatternConfig, url string) {
	for _, m := range pc.Patterns {
		if m.Compiled.Match(url) {
			args := append(m.Args, url)
			exec.Command(m.Program, args...).Start()
			return
		}
	}
	args := append(conf.LinkArgs, url)
	exec.Command(conf.LinkViewer, args...).Start()
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

func openMediaType(conf MediaConfig, filenames []string, mediaType string) {
	switch mediaType {
	case "image":
		if conf.ImageReverse {
			filenames = reverseFiles(filenames)
		}
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
		if conf.VideoReverse {
			filenames = reverseFiles(filenames)
		}
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
		if conf.AudioReverse {
			filenames = reverseFiles(filenames)
		}
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

func ColorKey(c *Config, pre, key, end string) string {
	color := ColorMark(c.Style.TextSpecial2)
	normal := ColorMark(c.Style.Text)
	key = tview.Escape("[" + key + "]")
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

func FloorDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func OutputDate(status time.Time, today time.Time, long, short string, relativeDate int) string {
	ty, tm, td := today.Date()
	sy, sm, sd := status.Date()

	format := long
	sameDay := false
	displayRelative := false

	if ty == sy && tm == sm && td == sd {
		format = short
		sameDay = true
	}

	todayFloor := FloorDate(today)
	statusFloor := FloorDate(status)

	if relativeDate > -1 && !sameDay {
		days := int(todayFloor.Sub(statusFloor).Hours() / 24)
		if relativeDate == 0 || days <= relativeDate {
			displayRelative = true
		}
	}
	var dateOutput string
	if displayRelative {
		y, m, d, _, _, _ := timex.Diff(statusFloor, todayFloor)
		if y > 0 {
			dateOutput = fmt.Sprintf("%s%dy", dateOutput, y)
		}
		if dateOutput != "" || m > 0 {
			dateOutput = fmt.Sprintf("%s%dm", dateOutput, m)
		}
		if dateOutput != "" || d > 0 {
			dateOutput = fmt.Sprintf("%s%dd", dateOutput, d)
		}
	} else {
		dateOutput = status.Format(format)
	}
	return dateOutput
}

func Notify(nc NotificationConfig, t NotificationType, title string, body string) {
	switch t {
	case NotificationFollower:
		if !nc.NotificationFollower {
			return
		}
	case NotificationFavorite:
		if !nc.NotificationFavorite {
			return
		}
	case NotificationMention:
		if !nc.NotificationMention {
			return
		}
	case NotificationBoost:
		if !nc.NotificationBoost {
			return
		}
	case NotificationPoll:
		if !nc.NotificationPoll {
			return
		}
	case NotificationPost:
		if !nc.NotificationPost {
			return
		}
	default:
		return
	}

	beeep.Notify(title, body, "")
}
