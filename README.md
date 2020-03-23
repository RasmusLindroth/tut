# Tut - a Mastodon TUI

A TUI for Mastodon with vim inspired keys. The program misses some features but they will be added when I get time.

Press `C` to create a new toot.

You can find Linux binaries under [releases](https://github.com/RasmusLindroth/tut/releases).

![Preview](./images/preview.png "Preview")

Currently supported commands
* `:q` `:quit` exit
* `:timeline` home, local, federated, direct

Explanation of the non obvious keys when viewing a toot
* `V` = view. In this mode you can scroll throught the text of the toot if it doesn't fit the screen
* `O` = open. Gives you a list of all URLs in the toot. Opens them in your default browser.
* `M` = media. Opens the media with `xdg-open`.

On my TODO-list:
* Support for config files (theme, default image/video viewer)
* Multiple accounts
* View users profiles
* Support search
* Support tags
* Support lists
* Notifications
* Better error handling (in other words, don't crash the whole program)
