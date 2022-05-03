package util

import "os/exec"

func OpenURL(url string) {
	exec.Command("xdg-open", url).Start()
}
