package util

import (
	"os/exec"
	"runtime"
)

func GetDefaultForOS() (program string, args []string) {
	switch runtime.GOOS {
	case "windows":
		program = "start"
		args = []string{"/wait"}
	case "darwin":
		program = "open"
		args = []string{"-W"}
	default:
		program = "xdg-open"
	}
	return program, args
}

func OpenURL(url string) {
	program, args := GetDefaultForOS()
	args = append(args, url)
	exec.Command(program, args...).Start()
}
