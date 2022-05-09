package config

import (
	"fmt"
	"os"
)

func Load() *Config {
	err := createConfigDir()
	if err != nil {
		fmt.Printf("Couldn't create or access the configuration dir. Error: %v\n", err)
		os.Exit(1)
	}
	path, exists, err := checkConfig("config.ini")
	if err != nil {
		fmt.Printf("Couldn't access config.ini. Error: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		err = CreateDefaultConfig(path)
		if err != nil {
			fmt.Printf("Couldn't create default config. Error: %v\n", err)
			os.Exit(1)
		}
	}
	config, err := parseConfig(path)
	if err != nil {
		fmt.Printf("Couldn't open or parse the config. Error: %v\n", err)
		os.Exit(1)
	}
	return &config
}
