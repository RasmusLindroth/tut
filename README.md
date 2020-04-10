# Tut - a Mastodon TUI

A TUI for Mastodon with vim inspired keys. The program misses some features but they will be added when I get time.

Press `C` to create a new toot.

You can find Linux binaries under [releases](https://github.com/RasmusLindroth/tut/releases).

![Preview](./images/preview.png "Preview")

## Currently supported commands
* `:q` `:quit` exit
* `:timeline` home, local, federated, direct, notifications
* `:tl` h, l, f, d, n (a shorter form of the former)
* `:tag` followed by the hashtag e.g. `:tag linux`
* `:user` followed by a username e.g. `:user rasmus` to narrow a search include 
the instance like this `:user rasmus@mastodon.acc.sunet.se`.

Explanation of the non obvious keys when viewing a toot
* `V` = view. In this mode you can scroll throught the text of the toot if it doesn't fit the screen
* `O` = open. Gives you a list of all URLs in the toot. Opens them in your default browser, if it's
an user or tag they will be opened in tut.
* `M` = media. Opens the media with `xdg-open`.

## Configuration
Tut if configurable, so you can change things like the colors, the default timeline, 
what image viewer to use and some more. Check out the configuration file to see 
all the options.

You find it in `XDG_CONFIG_HOME/tut/config.ini` which usally equals to `~/.config/tut/config.ini`.

## Install instructions
### Using Arch?

You can find it in the Arch User Repository (AUR).

https://aur.archlinux.org/packages/tut/

## Build it yourself
If you don't use the binary that you find under releases
you will need Go. Use a newer one that supports modules.

```bash
# First clone this repository
git clone https://github.com/RasmusLindroth/tut.git

# Go to that folder
cd tut

# Build or install

# Install (usally /home/user/go/bin)
go install

#Build (same directory i.e. ./ )
go build
```

If you choose to install and want to be able to just run `tut` 
you will have to add `go/bin` to your `$PATH`.



## On my TODO-list:
* Multiple accounts
* Support search
* Support lists
* Better error handling (in other words, don't crash the whole program)

## Thanks to
* [mattn/go-mastodon](https://github.com/mattn/go-mastodon) - used to make calls to the Mastodon API
* [rivo/tview](https://github.com/rivo/tview) - used to make the TUI
* [gdamore/tcell](https://github.com/gdamore/tcell) - used by tview under the hood
* [microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday) - used to remove HTML-tags
