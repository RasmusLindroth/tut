# Tut - a Mastodon TUI

A TUI for Mastodon with vim inspired keys. The program misses some features but they will be added when I get time.

Press `C` to create a new toot and `N` to focus on your notifications.

You can find Linux binaries under [releases](https://github.com/RasmusLindroth/tut/releases).

![Preview](./images/preview.png "Preview")

## Currently supported commands
* `:q` `:quit` exit
* `:timeline` home, local, federated, direct, notifications
  * `:tl` h, l, f, d, n (shorter form)
* `:blocking` lists users that you have blocked
* `:boosts` lists users that boosted the toot
* `:bookmarks` lists all your bookmarks
* `:compose` compose a new toot
* `:favorites` lists users that favorited the toot
* `:muting`  lists users that you have muted
* `:profile` go to your profile
* `:saved` alias for bookmarks
* `:tag` followed by the hashtag e.g. `:tag linux`
* `:user` followed by a username e.g. `:user rasmus` to narrow a search include 
the instance like this `:user rasmus@mastodon.acc.sunet.se`.

Keys without description in tut
* `c` = Compose a new toot
* `hjkl` = navigation
* `arrow keys` = navigation
* `q` = go back and quit
* `ESC` =  go back

Explanation of the non obvious keys when viewing a toot
* `V` = view. In this mode you can scroll throught the text of the toot if it doesn't fit the screen
* `O` = open. Gives you a list of all URLs in the toot. Opens them in your default browser, if it's
an user or tag they will be opened in tut.
* `M` = media. Opens the media with `xdg-open`.

## Configuration
Tut is configurable, so you can change things like the colors, the default timeline, 
what image viewer to use and some more. Check out the configuration file to see 
all the options.

You find it in `XDG_CONFIG_HOME/tut/config.ini` which usally equals to `~/.config/tut/config.ini`.

You can find an updated configuration file in this repo named `config.example.ini`.
If there are any new configurations options you can copy them frome that file.

## Install instructions
### Binary releases
Head over to https://github.com/RasmusLindroth/tut/releases

### Arch or Manjaro?

You can find it in the Arch User Repository (AUR). I'm the maintainer there.

https://aur.archlinux.org/packages/tut/

### Debian

http://packages.azlux.fr/ (I'm not the maintainer)

### FreeBSD

https://www.freshports.org/net-im/tut (I'm not the maintainer)


## Build it yourself
If you don't use the binary that you find under releases
you will need Go. Use a newer one that supports modules.

```bash
# Fetches and installs tut. Usally /home/user/go/bin
go get -u github.com/RasmusLindroth/tut

# You can also clone the repo if you like
# First clone this repository
git clone https://github.com/RasmusLindroth/tut.git

# Go to that folder
cd tut

# Build or install

# Install (usally /home/user/go/bin)
go install

# Build (same directory i.e. ./ )
go build
```

If you choose to install and want to be able to just run `tut` 
you will have to add `go/bin` to your `$PATH`.

## Flags and commands
```
Commands:
    example-config - creates the default configuration file in the current directory and names it ./config.example.ini

Flags:
    --help -h - prints this message
    --version -v - prints the version
    --new-user -n - add one more user to tut
    --user <name> -u <name> - login directly to user named <name>
        Don't use a = between --user and the <name> 
        If two users are named the same. Use full name like tut@fosstodon.org
```

## Thanks to
* [mattn/go-mastodon](https://github.com/mattn/go-mastodon) - used to make calls to the Mastodon API
* [rivo/tview](https://github.com/rivo/tview) - used to make the TUI
* [gdamore/tcell](https://github.com/gdamore/tcell) - used by tview under the hood
* [microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday) - used to remove HTML-tags
