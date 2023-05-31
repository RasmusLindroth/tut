package config

var conftext = `# Configuration file for tut

[general]
# What editor to use. TUT_USE_INTERNAL will use the editor that comes with tut.
# If you want you can set this to $EDITOR to use your environment variable or
# vim if you want to specify the program directly.
# default="TUT_USE_INTERNAL"
editor="TUT_USE_INTERNAL"

# You need to press yes in a confirmation dialog before favoriting, boosting,
# etc.
# default=true
confirmation=true

# Enable mouse support in tut.
# default=false
mouse-support=false

# The date format to be used. See https://pkg.go.dev/time#pkg-constants
# default="2006-01-02 15:04"
date-format="2006-01-02 15:04"

# Format for dates the same day. See date-format for more info.
# default="15:04"
date-today-format="15:04"

# This displays relative dates instead for statuses that are one day or older
# the output is 1y2m1d (1 year 2 months and 1 day)
# 
# The value is an integer
# -1     = don't use relative dates
#  0     = always use relative dates, except for dates < 1 day
#  1 - ∞ = number of days to use relative dates
# 				
# Value: 28 will display a relative date for toots that are between 1-28 days
# old. Otherwise it will use the short or long format.
# default=-1
date-relative=-1

# The max with of text before it wraps when displaying a toot.
# default=0
max-width=0

# The placement of your panes.
# valid: left, right, top, bottom
# default="left"
list-placement="left"

# How should panes be split?
# valid: row, column
# default="row"
list-split="row"

# The proportion of panes vs. content. 1 on this and 3 on content below results
# in content taking up 3 times more space.
# default=1
list-proportion=1

# See previous.
# default=2
content-proportion=2

# Hide notifications of this type in your notification timelines.
# valid: mention, status, boost, follow, follow_request, favorite, poll, edit
# default=[]
notifications-to-hide=[]

# Always include a quote of the message you're replying to.
# default=false
quote-reply=false

# If you want to show icons in timelines.
# default=true
show-icons=true

# If you only want to you the letter of keys instead of the full hint.
# default=false
short-hints=false

# If you want to display the filter that filtered a toot.
# default=true
show-filter-phrase=true

# Display a message in the commandbar on how to access the help text.
# default=true
show-help=true

# Always jump to the newest post. May ruin your reading experience.
# default=false
stick-to-top=false

# Display the username of the person being boosted instead of the person that
# boosted.
# default=false
show-boosted-user=false

# Open a new pane when you run a command like :timeline home.
# default=true
commands-in-new-pane=true

# Set a default name for the timeline if the name is empty. So if you run :tag
# linux the title of the pane will be set to #linux
# default=true
dynamic-timeline-name=true

# 0 = No terminal title
# 1 = Show title in terminal and top bar
# 2 = Only show terminal title, and no top bar in tut
# 3 = No terminal title and no top bar in tut.
# valid: 0, 1, 2, 3
# default=0
terminal-title=0

# If you don't want the whole UI to update, and only update the text content you
# can disable this. This will lead to some artifacts being left on the screen
# when emojis are present.
# default=true
redraw-ui=true

# The leader is used as a shortcut to run commands as you can do in Vim. By
# default this is disabled and you enable it by setting a key here. It can only
# consist of one char, so set it to something like a comma.
# default=""
leader-key=""

# Number of milliseconds before the leader command resets. So if you tap the
# leader-key by mistake or are to slow it empties all the input after X
# milliseconds.
# default=1000
leader-timeout=1000

# [[general.timelines]]
# Timelines adds panes of feeds. You can customize the number of feeds, what
# they should show and the key to activate them.

# --- START OF EXAMPLE ---
# [[general.timelines]]
# name="home"
# type="home"
# hide-boosts=false
# hide-replies=false
# 
# [[general.timelines]]
# name="Notifications"
# type="notifications"
# keys=["n", "N"]
# closed=true
# on-creation-closed="new-pane"
# on-focus="focus-self"
# --- END OF EXAMPLE ---

# The name to display above the timeline
# default=""
# name=""

# The type of the timeline
# valid: home, direct, local, federated, bookmarks, saved, favorited, notifications,
# lists, mentions, tag
# default=""
# type=""

# Used for the tag type, so here you set the tag. If you have multiple you
# separate them with a space.
# default=""
# data=""

# A list of keys to give this timeline focus. See under the input section to
# learn more about keys.
# default=[]
# keys=[]

# A list of special-keys to give this timeline focus. See under the input
# section to learn more about special-keys.
# default=[]
# special-keys=[]

# A shortcut to give this timeline focus with your leader-key + this shortcut.
# default=""
# shortcut=""

# Hide boosts in this timeline.
# default="false"
# hide-boosts="false"

# Hide replies in this timeline.
# default="false"
# hide-replies="false"

# Don't open this timeline when you start tut. Use your keys or shortcut to open
# it.
# default="false"
# closed="false"

# Don't open this timeline when you start tut. Use your keys or shortcut to open
# it.
# valid: new-pane, current-pane
# default="new-pane"
# on-creation-closed="new-pane"

# Don't open this timeline when you start tut. Use your keys or shortcut to open
# it.
# valid: focus-pane, focus-self
# default="focus-pane"
# on-focus="focus-pane"

# [[general.leader-actions]]
# You set actions leader-key with one or more leader-actions.
# 
# The shortcuts are up to you, but keep them quite short and make sure they
# don't collide. If you have one shortcut that is "f" and an other one that is
# "fav", the one with "f" will always run and "fav" will never run. 
# 
# Some special actions that requires data to be set:
# pane is special as it's a shortcut for switching between the panes you've set
# under general and they are zero indexed. pane 0 = your first timeline, pane 1
# = your second and so on.
# list-placement as it takes the argument top, right, bottom or left
# list-split as it takes the argument column or row
# proportions takes the arguments [int] [int], where the first integer is the
# list and the other content, e.g. proportions 1 3. See list-proportion above
# for more information.

# --- START OF EXAMPLE ---
# [[general.leader-actions]]
# type="close-pane"
# shortcut="q"
# 
# [[general.leader-actions]]
# type="list-split"
# data="row"
# shortcut="r"
# 
# [[general.leader-actions]]
# type="list-split"
# data="column"
# shortcut="c"
# --- END OF EXAMPLE ---

# The action you want to run.
# valid: blocking, boosts, clear-notifications, close-pane, compose, edit, favorited,
# favorites, followers, following, history, list-placement, list-split, lists,
# move-pane-left, move-pane-right, move-pane-up, move-pane-down, move-pane-home,
# move-pane-end, muting, newer, pane, preferences, profile, proportions,
# refetch, stick-to-top, tags
# default=""
# type=""

# Data to pass to the action.
# default=""
# data=""

# A shortcut to run this action with your leader-key + this shortcut.
# default=""
# shortcut=""

[media]
# Media files will be removed directly after they've been opened. Some programs
# doesn't like this, so if your media doesn't open, try set this to false. Tut
# will remove all files once you close the program.
# default=true
delete-temp-files=true

[media.image]
# The program to open images. TUT_OS_DEFAULT equals xdg-open on Linux, open on
# MacOS and start on Windows.
# default="TUT_OS_DEFAULT"
program="TUT_OS_DEFAULT"

# Arguments to pass to the program.
# default=""
args=""

# If the program runs in the terminal set this to true.
# default=false
terminal=false

# If the program should be called multiple times when there is multiple files.
# If set to false all files will be passed as an argument, but not all programs
# support this.
# default=true
single=true

# If the files should be passed in reverse order. This will make some programs
# display the files in the correct order.
# default=false
reverse=false

[media.video]
# The program to open videos. TUT_OS_DEFAULT equals xdg-open on Linux, open on
# MacOS and start on Windows.
# default="TUT_OS_DEFAULT"
program="TUT_OS_DEFAULT"

# Arguments to pass to the program.
# default=""
args=""

# If the program runs in the terminal set this to true.
# default=false
terminal=false

# If the program should be called multiple times when there is multiple files.
# If set to false all files will be passed as an argument, but not all programs
# support this.
# default=true
single=true

# If the files should be passed in reverse order. This will make some programs
# display the files in the correct order.
# default=false
reverse=false

[media.audio]
# The program to open audio. TUT_OS_DEFAULT equals xdg-open on Linux, open on
# MacOS and start on Windows.
# default="TUT_OS_DEFAULT"
program="TUT_OS_DEFAULT"

# Arguments to pass to the program.
# default=""
args=""

# If the program runs in the terminal set this to true.
# default=false
terminal=false

# If the program should be called multiple times when there is multiple files.
# If set to false all files will be passed as an argument, but not all programs
# support this.
# default=true
single=true

# If the files should be passed in reverse order. This will make some programs
# display the files in the correct order.
# default=false
reverse=false

[media.link]
# The program to open links. TUT_OS_DEFAULT equals xdg-open on Linux, open on
# MacOS and start on Windows.
# default="TUT_OS_DEFAULT"
program="TUT_OS_DEFAULT"

# Arguments to pass to the program.
# default=""
args=""

# If the program runs in the terminal set this to true.
# default=false
terminal=false

[desktop-notification]
# Enable notifications when someone follows you.
# default=false
followers=false

# Enable notifications when one of your toots gets favorited.
# default=false
favorite=false

# Enable notifications  when someone mentions you.
# default=false
mention=false

# Enable notifications when a post you have interacted with gets edited.
# default=false
update=false

# Enable notifications when one of your toots gets boosted.
# default=false
boost=false

# Enable notifications when a poll ends.
# default=false
poll=false

# Enable notifications for new posts.
# default=false
posts=false

[open-custom]
# --- START OF EXAMPLE ---
# [[open-custom.programs]]
# program = 'chromium'
# terminal = false
# hint = "[C]hrome"
# keys = ["c", "C"]
# 		
# [[open-custom.programs]]
# program = 'imv'
# terminal = false
# hint = "[I]mv"
# keys = ["i", "I"]"
# --- END OF EXAMPLE ---

# [[open-custom.programs]]
# The program to open the file with.
# default=""
# program=""

# Arguments to pass to the program.
# default=""
# args=""

# If the program runs in the terminal set this to true.
# default=false
# terminal=false

# What should the key hint in tut be for this program. See under the input
# section to learn more about hint.
# default=""
# hint=""

# A list of keys to to open files with this program. See under the input section
# to learn more about keys.
# default=[]
# keys=[]

# A list of special-keys to open files with this program. See under the input
# section to learn more about special-keys.
# default=[]
# special-keys=[]

[open-pattern]
# [[open-pattern.programs]]
# Here you can set your own glob patterns for opening matching URLs in the
# program you want them to open up in. You could for example open Youtube videos
# in your video player instead of your default browser. To see the syntax for
# glob pattern you can follow this URL https://github.com/gobwas/glob#syntax.
# default=""
# matching=""

# The program to open the file with.
# default=""
# program=""

# Arguments to pass to the program.
# default=""
# args=""

# If the program runs in the terminal set this to true.
# default=false
# terminal=false

[style]
# All styles can be represented in their HEX value like #ffffff or with their
# name, so in this case white. The only special value is "default" which equals
# to transparent, so it will be the same color as your terminal.
# You can also use xrdb colors like this xrdb:color1 The program will use colors
# prefixed with an * first then look for URxvt or XTerm if it can't find any
# color prefixed with an asterisk. If you don't want tut to guess the prefix you
# can set the prefix yourself. If the xrdb color can't be found a preset color
# will be used. You'll have to set theme="none" for this to work.

# The theme to use. You can use some themes that comes bundled with tut. Check
# out the themes available on the URL below. If a theme is named nord.toml you
# just write theme="nord".
# 
# https://github.com/RasmusLindroth/tut/tree/master/config/themes
# 
# You can also create a theme file in your config directory e.g.
# ~/.config/tut/themes/foo.toml and then set theme=foo.
# 
# If you want to use your own theme but don't want to create a new file, set
# theme="none" and then you can create your own theme below.
# 
# default="default"
theme="default"

# The xrdb prefix used for colors in .Xresources.
# default="guess"
xrdb-prefix="guess"

# The background color used on most elements.
# default=""
background=""

# The text color used on most of the text.
# default=""
text=""

# The color to display subtle elements or subtle text. Like lines and help text.
# default=""
subtle=""

# The color for errors or warnings
# default=""
warning-text=""

# This color is used to display username.
# default=""
text-special-one=""

# This color is used to display username and key hints.
# default=""
text-special-two=""

# The color of the bar at the top
# default=""
top-bar-background=""

# The color of the text in the bar at the top.
# default=""
top-bar-text=""

# The color of the bar at the bottom
# default=""
status-bar-background=""

# The color of the text in the bar at the bottom.
# default=""
status-bar-text=""

# The color of the bar at the bottom in view mode.
# default=""
status-bar-view-background=""

# The color of the text in the bar at the bottom in view mode.
# default=""
status-bar-view-text=""

# The color of the text in the command bar at the bottom.
# default=""
command-text=""

# Background of selected list items.
# default=""
list-selected-background=""

# The text color of selected list items.
# default=""
list-selected-text=""

# The background color of selected list items that are out of focus.
# default=""
list-selected-inactive-background=""

# The text color of selected list items that are out of focus.
# default=""
list-selected-inactive-text=""

# The main color of the text for key hints
# default=""
controls-text=""

# The highlight color of for key hints
# default=""
controls-highlight=""

# The background color in drop-downs and autocompletions
# default=""
autocomplete-background=""

# The text color in drop-downs at autocompletions
# default=""
autocomplete-text=""

# The background color for selected value in drop-downs and autocompletions
# default=""
autocomplete-selected-background=""

# The text color for selected value in drop-downs and autocompletions
# default=""
autocomplete-selected-text=""

# The background color on selected button and the text color of unselected
# buttons
# default=""
button-color-one=""

# The text color on selected button and the background color of unselected
# buttons
# default=""
button-color-two=""

# The background on named timelines.
# default=""
timeline-name-background=""

# The text color on named timelines
# default=""
timeline-name-text=""

[input]
# In this section you set the keys to be used in tut.
# 		
# The hint option lets you set which part of the hint that will be highlighted
# in tut. E.g. [F]avorite results in a highlighted F and the rest of the text is
# displayed normally.
# Some of the options can be in two states, like favorites, so there you can set
# the hint-alt option to something like Un[F]avorite.
# 
# Examples:
# "[D]elete" = Delete with a highlighted D
# "Un[F]ollow" = UnFollow with a highlighted F
# "[Enter]" = Enter where everything is highlighted
# "Yan[K]" = YanK with a highlighted K
# 
# The keys option lets you define what key that should be pressed. This is
# limited to on character only and they are case sensitive.
# Example:
# keys=["j","J"]
# 
# You can also set special-keys and they're for keys like Escape and Enter. To
# find the names of special keys you have to go to the following site and look
# for "var KeyNames = map[Key]string{"
# 
# https://github.com/gdamore/tcell/blob/master/key.go

[input.global-down]
# Keys for moving down

# default=["j", "J"]
keys=["j","J"]

# default=["Down"]
special-keys=["Down"]

[input.global-up]
# Keys for moving down

# default=["k", "K"]
keys=["k","K"]

# default=["Up"]
special-keys=["Up"]

[input.global-enter]
# To select items

# default=["Enter"]
special-keys=["Enter"]

[input.global-back]
# To go back

# default="[Esc]"
hint="[Esc]"

# default=["Esc"]
special-keys=["Esc"]

[input.global-exit]
# To go back or exit

# default="[Q]uit"
hint="[Q]uit"

# default=["q", "Q"]
keys=["q","Q"]

[input.main-home]
# Move to the top

# default=["g"]
keys=["g"]

# default=["Home"]
special-keys=["Home"]

[input.main-end]
# Move to the bottom

# default=["G"]
keys=["G"]

# default=["End"]
special-keys=["End"]

[input.main-prev-feed]
# Go to previous feed

# default=["h", "H"]
keys=["h","H"]

# default=["Left"]
special-keys=["Left"]

[input.main-next-feed]
# Go to next feed

# default=["l", "L"]
keys=["l","L"]

# default=["Right"]
special-keys=["Right"]

[input.main-prev-pane]
# Focus on the previous feed pane

# default=["Backtab"]
special-keys=["Backtab"]

[input.main-next-pane]
# Focus on the next feed pane

# default=["Tab"]
special-keys=["Tab"]

[input.main-next-account]
# Focus on the next account

# default=["Ctrl-N"]
special-keys=["Ctrl-N"]

[input.main-prev-account]
# Focus on the previous account

# default=["Ctrl-P"]
special-keys=["Ctrl-P"]

[input.main-compose]
# Compose a new toot

# default=["c", "C"]
keys=["c","C"]

[input.status-avatar]
# Open avatar

# default="[A]vatar"
hint="[A]vatar"

# default=["a", "A"]
keys=["a","A"]

[input.status-boost]
# Boost a toot

# default="[B]oost"
hint="[B]oost"

# default=["b", "B"]
keys=["b","B"]

[input.status-edit]
# Edit a toot

# default="[E]dit"
hint="[E]dit"

# default=["e", "E"]
keys=["e","E"]

[input.status-delete]
# Delete a toot

# default="[D]elete"
hint="[D]elete"

# default=["d", "D"]
keys=["d","D"]

[input.status-favorite]
# Favorite a toot

# default="[F]avorite"
hint="[F]avorite"

# default=["f", "F"]
keys=["f","F"]

[input.status-media]
# Open toots media files

# default="[M]edia"
hint="[M]edia"

# default=["m", "M"]
keys=["m","M"]

[input.status-links]
# Open links

# default="[O]pen"
hint="[O]pen"

# default=["o", "O"]
keys=["o","O"]

[input.status-poll]
# Open poll

# default="[P]oll"
hint="[P]oll"

# default=["p", "P"]
keys=["p","P"]

[input.status-reply]
# Reply to toot

# default="[R]eply"
hint="[R]eply"

# default=["r", "R"]
keys=["r","R"]

[input.status-bookmark]
# Save/bookmark a toot

# default="[S]ave"
hint="[S]ave"

# default="Un[S]ave"
hint-alt="Un[S]ave"

# default=["s", "S"]
keys=["s","S"]

[input.status-thread]
# View thread

# default="[T]hread"
hint="[T]hread"

# default=["t", "T"]
keys=["t","T"]

[input.status-user]
# Open user profile

# default="[U]ser"
hint="[U]ser"

# default=["u", "U"]
keys=["u","U"]

[input.status-view-focus]
# Open the view mode

# default="[V]iew"
hint="[V]iew"

# default=["v", "V"]
keys=["v","V"]

[input.status-yank]
# Yank the url of the toot

# default="[Y]ank"
hint="[Y]ank"

# default=["y", "Y"]
keys=["y","Y"]

[input.status-toggle-cw]
# Show the content in a content warning

# default="Press [Z] to toggle cw"
hint="Press [Z] to toggle cw"

# default=["z", "Z"]
keys=["z","Z"]

[input.status-show-filtered]
# Show the content of a filtered toot

# default="Press [Z] to view filtered toot"
hint="Press [Z] to view filtered toot"

# default=["z", "Z"]
keys=["z","Z"]

[input.user-avatar]
# View avatar

# default="[A]vatar"
hint="[A]vatar"

# default=["a", "A"]
keys=["a","A"]

[input.user-block]
# Block the user

# default="[B]lock"
hint="[B]lock"

# default="Un[B]lock"
hint-alt="Un[B]lock"

# default=["b", "B"]
keys=["b","B"]

[input.user-follow]
# Follow user

# default="[F]ollow"
hint="[F]ollow"

# default="Un[F]ollow"
hint-alt="Un[F]ollow"

# default=["f", "F"]
keys=["f","F"]

[input.user-follow-request-decide]
# Follow user

# default="Follow [R]equest"
hint="Follow [R]equest"

# default="Follow [R]equest"
hint-alt="Follow [R]equest"

# default=["r", "R"]
keys=["r","R"]

[input.user-mute]
# Mute user

# default="[M]ute"
hint="[M]ute"

# default="Un[M]ute"
hint-alt="Un[M]ute"

# default=["m", "M"]
keys=["m","M"]

[input.user-links]
# Open links

# default="[O]pen"
hint="[O]pen"

# default=["o", "O"]
keys=["o","O"]

[input.user-user]
# View user profile

# default="[U]ser"
hint="[U]ser"

# default=["u", "U"]
keys=["u","U"]

[input.user-view-focus]
# Open view mode

# default="[V]iew"
hint="[V]iew"

# default=["v", "V"]
keys=["v","V"]

[input.user-yank]
# Yank the user URL

# default="[Y]ank"
hint="[Y]ank"

# default=["y", "Y"]
keys=["y","Y"]

[input.list-open-feed]
# Open list

# default="[O]pen"
hint="[O]pen"

# default=["o", "O"]
keys=["o","O"]

[input.list-user-list]
# List all users in a list

# default="[U]sers"
hint="[U]sers"

# default=["u", "U"]
keys=["u","U"]

[input.list-user-add]
# Add user to list

# default="[A]dd"
hint="[A]dd"

# default=["a", "A"]
keys=["a","A"]

[input.list-user-delete]
# Delete user from list

# default="[D]elete"
hint="[D]elete"

# default=["d", "D"]
keys=["d","D"]

[input.link-open]
# Open URL

# default="[O]pen"
hint="[O]pen"

# default=["o", "O"]
keys=["o","O"]

[input.link-yank]
# Yank the URL

# default="[Y]ank"
hint="[Y]ank"

# default=["y", "Y"]
keys=["y","Y"]

[input.tag-open-feed]
# Open tag feed

# default="[O]pen"
hint="[O]pen"

# default=["o", "O"]
keys=["o","O"]

[input.tag-follow]
# Toggle follow on tag

# default="[F]ollow"
hint="[F]ollow"

# default="Un[F]ollow"
hint-alt="Un[F]ollow"

# default=["f", "F"]
keys=["f","F"]

[input.compose-edit-cw]
# Edit content warning text on new toot

# default="[C]W text"
hint="[C]W text"

# default=["c", "C"]
keys=["c","C"]

[input.compose-edit-text]
# Edit the text on new toot

# default="[E]dit text"
hint="[E]dit text"

# default=["e", "E"]
keys=["e","E"]

[input.compose-include-quote]
# Include a quote when replying

# default="[I]nclude quote"
hint="[I]nclude quote"

# default=["i", "I"]
keys=["i","I"]

[input.compose-media-focus]
# Focus on adding media to toot

# default="[M]edia"
hint="[M]edia"

# default=["m", "M"]
keys=["m","M"]

[input.compose-post]
# Post the new toot

# default="[P]ost"
hint="[P]ost"

# default=["p", "P"]
keys=["p","P"]

[input.compose-toggle-content-warning]
# Toggle content warning on toot

# default="[T]oggle CW"
hint="[T]oggle CW"

# default=["t", "T"]
keys=["t","T"]

[input.compose-visibility]
# Edit the visibility on new toot

# default="[V]isibility"
hint="[V]isibility"

# default=["v", "V"]
keys=["v","V"]

[input.compose-language]
# Edit the language of a toot

# default="[L]ang"
hint="[L]ang"

# default=["l", "L"]
keys=["l","L"]

[input.compose-poll]
# Switch to creating a poll

# default="P[O]ll"
hint="P[O]ll"

# default=["o", "O"]
keys=["o","O"]

[input.media-delete]
# Delete media file

# default="[D]elete"
hint="[D]elete"

# default=["d", "D"]
keys=["d","D"]

[input.media-edit-desc]
# Edit the description on media file

# default="[E]dit desc"
hint="[E]dit desc"

# default=["e", "E"]
keys=["e","E"]

[input.media-add]
# Add a new media file

# default="[A]dd"
hint="[A]dd"

# default=["a", "A"]
keys=["a","A"]

[input.vote-vote]
# Vote on poll

# default="[V]ote"
hint="[V]ote"

# default=["v", "V"]
keys=["v","V"]

[input.vote-select]
# Select item to vote on

# default="[Enter] to select"
hint="[Enter] to select"

# default=["Enter"]
special-keys=["Enter"]

[input.poll-add]
# Add a new poll option

# default="[A]dd"
hint="[A]dd"

# default=["a", "A"]
keys=["a","A"]

[input.poll-edit]
# Edit a poll option

# default="[E]dit"
hint="[E]dit"

# default=["e", "E"]
keys=["e","E"]

[input.poll-delete]
# Delete a poll option

# default="[D]elete"
hint="[D]elete"

# default=["d", "D"]
keys=["d","D"]

[input.poll-multi-toggle]
# Toggle voting on multiple options

# default="Toggle [M]ultiple"
hint="Toggle [M]ultiple"

# default=["m", "M"]
keys=["m","M"]

[input.poll-expiration]
# Change the expiration of poll

# default="E[X]pires"
hint="E[X]pires"

# default=["x", "X"]
keys=["x","X"]

[input.preference-name]
# Change display name

# default="[N]ame"
hint="[N]ame"

# default=["n", "N"]
keys=["n","N"]

[input.preference-visibility]
# Change default visibility of toots

# default="[V]isibility"
hint="[V]isibility"

# default=["v", "V"]
keys=["v","V"]

[input.preference-bio]
# Change bio in profile

# default="[B]io"
hint="[B]io"

# default=["b", "B"]
keys=["b","B"]

[input.preference-save]
# Save your preferences

# default="[S]ave"
hint="[S]ave"

# default=["s", "S"]
keys=["s","S"]

[input.preference-fields]
# Edit profile fields

# default="[F]ields"
hint="[F]ields"

# default=["f", "F"]
keys=["f","F"]

[input.preference-fields-add]
# Add new field

# default="[A]dd"
hint="[A]dd"

# default=["a", "A"]
keys=["a","A"]

[input.preference-fields-edit]
# Edit current field

# default="[E]dit"
hint="[E]dit"

# default=["e", "E"]
keys=["e","E"]

[input.preference-fields-delete]
# Delete current field

# default="[D]elete"
hint="[D]elete"

# default=["d", "D"]
keys=["d","D"]

[input.editor-exit]
# Exit the editor

# default="[Esc] when done"
hint="[Esc] when done"

# default=["Esc"]
special-keys=["Esc"]
`
