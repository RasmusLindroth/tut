package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

type URL struct {
	Text    string
	URL     string
	Classes []string
}

func CleanHTML(content string) (string, []URL) {
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

						if strings.Contains(a.Val, "hashtag") {
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

func CmdToString(cmd string) (string, error) {
	cmd = strings.TrimPrefix(cmd, "!CMD!")
	parts := strings.Split(cmd, " ")
	s, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	return strings.TrimSpace(string(s)), err
}

func MakeDirs() {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Printf("couldn't find $HOME. Error: %v\n", err)
		os.Exit(1)
	}
	dir := cd + "/tut"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Printf("couldn't create dirs. Error: %v\n", err)
		os.Exit(1)
	}
}

func CheckConfig(filename string) (path string, exists bool, err error) {
	cd, err := os.UserConfigDir()
	if err != nil {
		log.Printf("couldn't find $HOME. Error: %v\n", err)
		os.Exit(1)
	}
	dir := cd + "/tut/"
	path = dir + filename
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, false, nil
	} else if err != nil {
		return path, true, err
	}
	return path, true, err
}

func FormatUsername(a mastodon.Account) string {
	if a.DisplayName != "" {
		return fmt.Sprintf("%s (%s)", a.DisplayName, a.Acct)
	}
	return a.Acct
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

func StatusOrReblog(s *mastodon.Status) *mastodon.Status {
	if s.Reblog != nil {
		return s.Reblog
	}
	return s
}

func SetTerminalTitle(s string) {
	fmt.Printf("\033]0;%s\a", s)
}
