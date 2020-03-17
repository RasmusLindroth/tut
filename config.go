package main

import (
	"os"

	"github.com/kyoh86/xdg"
)

type Config struct {
}

func CreateConfigDir() error {
	path := xdg.ConfigHome() + "/tut"
	return os.MkdirAll(path, os.ModePerm)
}

func CheckConfig(filename string) (path string, exists bool, err error) {
	dir := xdg.ConfigHome() + "/tut/"
	path = dir + filename
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, false, nil
	} else if err != nil {
		return path, true, err
	}
	return path, true, err
}
