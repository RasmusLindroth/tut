package config

import (
	"fmt"
	"os"
)

func Load(cnfPath string, cnfDir string) *Config {
	err := createConfigDir()
	if err != nil {
		fmt.Printf("Couldn't create or access the configuration dir. Error: %v\n", err)
		os.Exit(1)
	}
	path, exists, err := checkConfig("config.toml", cnfPath, cnfDir)
	if err != nil {
		fmt.Printf("Couldn't access config.toml. Error: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		err = CreateDefaultConfig(path)
		if err != nil {
			fmt.Printf("Couldn't create default config. Error: %v\n", err)
			os.Exit(1)
		}
	}
	config, err := parseConfig(path, cnfPath, cnfDir)
	if err != nil {
		fmt.Printf("Couldn't open or parse the config. Error: %v\n", err)
		os.Exit(1)
	}
	return &config
}
