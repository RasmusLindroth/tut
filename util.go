package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
		urls = append(urls, url)
	}
	return urls
}

func cleanTootHTML(content string) (string, []URL) {
	stripped := bluemonday.NewPolicy().AllowElements("p", "br").AllowAttrs("href", "class").OnElements("a").Sanitize(content)
	urls := getURLs(stripped)
	stripped = bluemonday.NewPolicy().AllowElements("p", "br").Sanitize(content)
	stripped = strings.ReplaceAll(stripped, "<br>", "\n")
	stripped = strings.ReplaceAll(stripped, "<p>", "")
	stripped = strings.ReplaceAll(stripped, "</p>", "\n\n")
	stripped = strings.TrimSpace(stripped)
	stripped = html.UnescapeString(stripped)
	return stripped, urls
}

func openEditor(app *tview.Application, content string) (string, error) {
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
	cmd := exec.Command("nvim", f.Name())
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

func openURL(url string) {
	exec.Command("xdg-open", url).Start()
}

func openMedia(filenames []string) {
	for _, f := range filenames {
		exec.Command("xdg-open", f).Run()
	}
	for _, f := range filenames {
		os.Remove(f)
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
