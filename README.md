# Tut - a Mastodon TUI

A TUI for Mastodon with vim inspired keys. The program misses some features but they will be added when I get time.

Press `C` to create a new toot.

You can find Linux binaries under [releases](https://github.com/RasmusLindroth/tut/releases).

![Preview](./images/preview.png "Preview")

### Currently supported commands
* `:q` `:quit` exit
* `:timeline` home, local, federated, direct

Explanation of the non obvious keys when viewing a toot
* `V` = view. In this mode you can scroll throught the text of the toot if it doesn't fit the screen
* `O` = open. Gives you a list of all URLs in the toot. Opens them in your default browser.
* `M` = media. Opens the media with `xdg-open`.

### Install instructions
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



### On my TODO-list:
* Support for config files (theme, default image/video viewer)
* Multiple accounts
* View users profiles
* Support search
* Support tags
* Support lists
* Notifications
* Better error handling (in other words, don't crash the whole program)

### Thanks to
* [mattn/go-mastodon](https://github.com/mattn/go-mastodon) - used to make calls to the Mastodon API
* [rivo/tview](https://github.com/rivo/tview) - used to make the TUI
* [gdamore/tcell](https://github.com/gdamore/tcell) - used by tview under the hood
* [microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday) - used to remove HTML-tags
