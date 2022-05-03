package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RasmusLindroth/tut/config"
)

func CliView(version string) (newUser bool, selectedUser string) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "example-config":
			config.CreateDefaultConfig("./config.example.ini")
			os.Exit(0)
		case "--new-user", "-n":
			newUser = true
		case "--user", "-u":
			if len(os.Args) > 2 {
				name := os.Args[2]
				selectedUser = strings.TrimSpace(name)
			} else {
				log.Fatalln("--user/-u must be followed by a user name. Like -u tut")
			}
		case "--help", "-h":
			fmt.Print("tut - a TUI for Mastodon with vim inspired keys.\n\n")
			fmt.Print("Usage:\n\n")
			fmt.Print("\tTo run the program you just have to write tut\n\n")

			fmt.Print("Commands:\n\n")
			fmt.Print("\texample-config - creates the default configuration file in the current directory and names it ./config.example.ini\n\n")

			fmt.Print("Flags:\n\n")
			fmt.Print("\t--help -h - prints this message\n")
			fmt.Print("\t--version -v - prints the version\n")
			fmt.Print("\t--new-user -n - add one more user to tut\n")
			fmt.Print("\t--user <name> -u <name> - login directly to user named <name>\n")
			fmt.Print("\t\tDon't use a = between --user and the <name>\n")
			fmt.Print("\t\tIf two users are named the same. Use full name like tut@fosstodon.org\n\n")

			fmt.Print("Configuration:\n\n")
			fmt.Printf("\tThe config is located in XDG_CONFIG_HOME/tut/config.ini which usally equals to ~/.config/tut/config.ini.\n")
			fmt.Printf("\tThe program will generate the file the first time you run tut. The file has comments which exmplains what each configuration option does.\n\n")

			fmt.Print("Contact info for issues or questions:\n\n")
			fmt.Printf("\t@rasmus@mastodon.acc.sunet.se\n\trasmus@lindroth.xyz\n")
			fmt.Printf("\thttps://github.com/RasmusLindroth/tut\n")
			os.Exit(0)
		case "--version", "-v":
			fmt.Printf("tut version %s\n\n", version)
			fmt.Printf("https://github.com/RasmusLindroth/tut\n")
			os.Exit(0)
		}
	}
	return newUser, selectedUser
}
