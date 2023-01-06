package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/util"
	"github.com/spf13/pflag"
)

func CliView(version string) (newUser bool, selectedUser string, confPath string, confDir string) {
	showHelp := pflag.BoolP("help", "h", false, "config path")
	showVersion := pflag.BoolP("version", "v", false, "config path")
	nu := pflag.BoolP("new-user", "n", false, "add one more user to tut")
	user := pflag.StringP("user", "u", "", "login directly to user named `<name>`")
	cnf := pflag.StringP("config", "c", "", "load config.toml from `<path>`")
	cnfDir := pflag.StringP("config-dir", "d", "", "load all config from `<path>`")
	pflag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "example-config":
			config.CreateDefaultConfig("./config.example.toml")
			os.Exit(0)
		}
	}
	if nu != nil && *nu {
		newUser = true
	}
	if user != nil && *user != "" {
		selectedUser = strings.TrimSpace(*user)
	}
	if cnf != nil && *cnf != "" {
		cp := strings.TrimSpace(*cnf)
		abs, err := util.GetAbsPath(cp)
		if err != nil {
			log.Fatalln(err)
		}
		confPath = abs
	} else if os.Getenv("TUT_CONF") != "" {
		cp := os.Getenv("TUT_CONF")
		abs, err := util.GetAbsPath(cp)
		if err != nil {
			log.Fatalln(err)
		}
		confPath = abs
	}
	if cnfDir != nil && *cnfDir != "" {
		cd := strings.TrimSpace(*cnfDir)
		abs, err := util.GetAbsPath(cd)
		if err != nil {
			log.Fatalln(err)
		}
		confDir = abs
	} else if os.Getenv("TUT_CONF_DIR") != "" {
		cd := os.Getenv("TUT_CONF_DIR")
		abs, err := util.GetAbsPath(cd)
		if err != nil {
			log.Fatalln(err)
		}
		confDir = abs
	}
	if showHelp != nil && *showHelp {
		fmt.Print("tut - a TUI for Mastodon with vim inspired keys.\n\n")
		fmt.Print("Usage:\n")
		fmt.Print("\tTo run the program you just have to write tut\n\n")

		fmt.Print("Commands:\n")
		fmt.Print("\texample-config - creates the default configuration file in the current directory and names it ./config.example.toml\n\n")

		fmt.Print("Flags:\n")
		fmt.Print("\t-h  --help             prints this message\n")
		fmt.Print("\t-v  --version          prints the version\n")
		fmt.Print("\t-n  --new-user         add one more user to tut\n")
		fmt.Print("\t-c  --config <path>    load config.toml from <path>\n")
		fmt.Print("\t-d --config-dir <path> load all config from <path>\n")
		fmt.Print("\t-u  --user <name>      login directly to user named <name>\n")
		fmt.Print("\t\tIf two users are named the same. Use full name like tut@fosstodon.org\n\n")

		fmt.Print("Configuration:\n")
		fmt.Printf("\tThe config is located in XDG_CONFIG_HOME/tut/config.toml which usually equals to ~/.config/tut/config.toml.\n")
		fmt.Printf("\tThe program will generate the file the first time you run tut. The file has comments which exmplains what each configuration option does.\n\n")

		fmt.Print("Contact info for issues or questions:\n")
		fmt.Printf("\t@tut@fosstodon.org\n\t@rasmus@mastodon.acc.sunet.se\n\trasmus@lindroth.xyz\n")
		fmt.Printf("\thttps://github.com/RasmusLindroth/tut\n")
		os.Exit(0)
	}
	if showVersion != nil && *showVersion {
		fmt.Printf("tut version %s\n", version)
		fmt.Printf("https://github.com/RasmusLindroth/tut\n")
		os.Exit(0)

	}
	return newUser, selectedUser, confPath, confDir
}
