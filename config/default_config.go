package config

var conftext = `# Configuration file for tut

[general]
# Shows a confirmation view before actions such as favorite, delete toot, boost
# etc.
# default=true
confirmation=true

# The date format to be used. See https://godoc.org/time#Time.Format
# default=2006-01-02 15:04
date-format=2006-01-02 15:04

# Format for dates the same day. See date-format for more info.
# default=15:04
date-today-format=15:04

# This displays relative dates instead for statuses that are one day or older
# the output is 1y2m1d (1 year 2 months and 1 day)
# 
# The value is an integear
# -1     = don't use relative dates
#  0     = always use relative dates, except for dates < 1 day
#  1 - ∞ = number of days to use relative dates
# 
# Example: date-relative=28 will display a relative date for toots that are
# between 1-28 days old. Otherwhise it will use the short or long format.
# default=-1
date-relative=-1

# The timeline that opens up when you start tut.
# Valid values: home, direct, local, federated
# default=home
timeline=home

# The max width of text before it wraps when displaying toots.
# 0 = no restriction.
# default=0
max-width=0

# If you want to display a list of notifications under your timeline feed.
# default=true
notification-feed=true

# Where do you want the list of toots to be placed?
# Valid values: left, right, top, bottom.
# default=left
list-placement=left

# If you have notification-feed set to true you can display it under the main
# list of toots (row) or place it to the right of the main list of toots
# (column).
# default=row
list-split=row

# Hide notification text above list in column split. It's displayed as
# [N]otifications.
# default=false
hide-notification-text=false

# You can change the proportions of the list view in relation to the content
# view list-proportion=1 and content-proportoin=3 will result in the content
# taking up 3 times more space.
# Must be n > 0
# default=1
list-proportion=1

# See list-proportion
# default=2
content-proportion=2

# If you always want to quote original message when replying.
# default=false
quote-reply=false

# If you're on an instance with a custom character limit you can set it here.
# default=500
char-limit=500

# If you want to show icons in the list of toots.
# default=true
show-icons=true

# If you've learnt all the shortcut keys you can remove the help text and only
# show the key in tui. So it gets less cluttered.
# default=false
short-hints=false

# If you want to show a message in the cmdbar on how to access the help text.
# default=true
show-help=true

# If you don't want the whole UI to update, and only the text content you can
# set this option to true. This will lead to some artifacts being left on the
# screen when emojis are present. But it will keep the UI from flashing on every
# single toot in some terminals.
# default=true
redraw-ui=true

# The leader is used as a shortcut to run commands as you can do in Vim. By
# default this is disabled and you enable it by setting a leader-key. It can
# only consist of one char and I like to use comma as leader key. So to set it
# you write leader-key=,
# default=
leader-key=

# Number of milliseconds before the leader command resets. So if you tap the
# leader-key by mistake or are to slow it empties all the input after X
# milliseconds.
# default=1000
leader-timeout=1000

# You set actions for the leader-key with one or more leader-action. It consists
# of two parts first the action then the shortcut. And they're seperated by a
# comma.
# 
# Available commands: home, direct, local, federated, compose, blocking,
# bookmarks, saved, favorited, boosts, favorites, following, followers, muting,
# profile, notifications, lists, tag
# 
# The shortcuts are up to you, but keep them quite short and make sure they
# don't collide. If you have one shortcut that is "f" and an other one that is
# "fav", the one with "f" will always run and "fav" will never run. Tag is
# special as you need to add the tag after, see the example below.
# 
# Some examples:
# leader-action=local,lo
# leader-action=lists,li
# leader-action=federated,fed
# leader-action=direct,d
# leader-action=tag linux,tl
# 


[media]
# Your image viewer.
# default=xdg-open
image-viewer=xdg-open

# Open the image viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
image-terminal=false

# If images should open one by one e.g. "imv image.png" multiple times. If set
# to false all images will open at the same time like this "imv image1.png
# image2.png image3.png". Not all image viewers support this, so try it first.
# default=true
image-single=true

# If you want to open the images in reverse order. In some image viewers this
# will display the images in the "right" order.
# default=false
image-reverse=false

# Your video viewer.
# default=xdg-open
video-viewer=xdg-open

# Open the video viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
video-terminal=false

# If videos should open one by one. See image-single.
# default=true
video-single=true

# If you want your videos in reverse order. In some video apps this will play
# the files in the "right" order.
# default=false
video-reverse=false

# Your audio viewer.
# default=xdg-open
audio-viewer=xdg-open

# Open the audio viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
audio-terminal=false

# If audio should open one by one. See image-single.
# default=true
audio-single=true

# If you want to play the audio files in reverse order. In some audio apps this
# will play the files in the "right" order.
# default=false
audio-reverse=false

# Your web browser.
# default=xdg-open
link-viewer=xdg-open

# Open the browser in the same terminal as toot. Only for terminal based
# browsers.
# default=false
link-terminal=false

[open-custom]
# This sections allows you to set up to five custom programs to upen URLs with.
# If the url points to an image, you can set c1-name to img and c1-use to imv.
# If the program runs in a terminal and you want to run it in the same terminal
# as tut. Set cX-terminal to true. The name will show up in the UI, so keep it
# short so all five fits.
# 
# c1-name=name
# c1-use=program
# c1-terminal=false
# 
# c2-name=name
# c2-use=program
# c2-terminal=false
# 
# c3-name=name
# c3-use=program
# c3-terminal=false
# 
# c4-name=name
# c4-use=program
# c4-terminal=false
# 
# c5-name=name
# c5-use=program
# c5-terminal=false

[open-pattern]
# Here you can set your own glob patterns for opening matching URLs in the
# program you want them to open up in. You could for example open Youtube videos
# in your video player instead of your default browser.
# 
# You must name the keys foo-pattern, foo-use and foo-terminal, where use is the
# program that will open up the URL. To see the syntax for glob pattern you can
# follow this URL https://github.com/gobwas/glob#syntax. foo-terminal is if the
# program runs in the terminal and should open in the same terminal as tut
# itself.
# 
# Example for youtube.com and youtu.be to open up in mpv instead of the browser.
# 
# y1-pattern=*youtube.com/watch*
# y1-use=mpv
# y1-terminal=false
# 
# y2-pattern=*youtu.be/*
# y2-use=mpv
# y2-terminal=false

[desktop-notification]
# Notification when someone follows you.
# default=false
followers=false

# Notification when someone favorites one of your toots.
# default=false
favorite=false

# Notification when someone mentions you.
# default=false
mention=false

# Notification when someone boosts one of your toots.
# default=false
boost=false

# Notification of poll results.
# default=false
poll=false

# Notification when there is new posts in current timeline.
# default=false
posts=false

[style]
# All styles can be represented in their HEX value like #ffffff or with their
# name, so in this case white. The only special value is "default" which equals
# to transparent, so it will be the same color as your terminal.
# 
# You can also use xrdb colors like this xrdb:color1 The program will use colors
# prefixed with an * first then look for URxvt or XTerm if it can't find any
# color prefixed with an asterik. If you don't want tut to guess the prefix you
# can set the prefix yourself. If the xrdb color can't be found a preset color
# will be used. You'll have to set theme=none for this to work.

# The xrdb prefix used for colors in .Xresources.
# default=guess
xrdb-prefix=guess

# You can use some themes that comes bundled with tut check out the themes
# available on the URL below. If a theme is named "nord.ini" you just write
# theme=nord
# 
# https://github.com/RasmusLindroth/tut/tree/master/themes
# 
# If you want to use your own theme set theme to none then you can create your
# own theme below
# default=default
theme=default

# The background color used on most elements.
# default=xrdb:background
background=xrdb:background

# The text color used on most of the text.
# default=xrdb:foreground
text=xrdb:foreground

# The color to display sublte elements or subtle text. Like lines and help text.
# default=xrdb:color14
subtle=xrdb:color14

# The color for errors or warnings
# default=xrdb:color1
warning-text=xrdb:color1

# This color is used to display username.
# default=xrdb:color5
text-special-one=xrdb:color5

# This color is used to display username and key hints.
# default=xrdb:color2
text-special-two=xrdb:color2

# The color of the bar at the top
# default=xrdb:color5
top-bar-background=xrdb:color5

# The color of the text in the bar at the top.
# default=xrdb:background
top-bar-text=xrdb:background

# The color of the bar at the bottom
# default=xrdb:color5
status-bar-background=xrdb:color5

# The color of the text in the bar at the bottom.
# default=xrdb:foreground
status-bar-text=xrdb:foreground

# The color of the bar at the bottom in view mode.
# default=xrdb:color4
status-bar-view-background=xrdb:color4

# The color of the text in the bar at the bottom in view mode.
# default=xrdb:foreground
status-bar-view-text=xrdb:foreground

# Background of selected list items.
# default=xrdb:color5
list-selected-background=xrdb:color5

# The text color of selected list items.
# default=xrdb:background
list-selected-text=xrdb:background

[input]
# You can edit the keys for tut below.
# 
# The syntax is a bit weird, but it works. And I'll try to explain it as well as
# I can.
# 
# Example:
# status-favorite="[F]avorite","Un[F]avorite",'f','F'
# status-delete="[D]elete",'d','D'
# 
# status-favorite and status-delete differs because favorite can be in two
# states, so you will have to add two key hints.
# Most keys will only have on key hint. Look at the default value for reference.
# 
# Key hints must be in some of the following formats. Remember the quotation
# marks.
# "" = empty
# "[D]elete" = Delete with a highlighted D
# "Un[F]ollow" = UnFollow with a highlighted F
# "[Enter]" = Enter where everything is highlighted
# "Yan[K]" = YanK with a highlighted K
# 
# After the hint (or hints) you must set the keys. You can do this in two ways,
# with single quotation marks or double ones.
# 
# The single ones are for single chars like 'a', 'b', 'c' and double marks are
# for special keys like "Enter". Remember that they are case sensetive.
# 
# To find the names of special keys you have to go to the following site and
# look for "var KeyNames = map[Key]string{"
# 
# https://github.com/gdamore/tcell/blob/master/key.go

# Keys for moving down
# default="",'j','J',"Down"
global-down="",'j','J',"Down"

# Keys for moving up
# default="",'k','K',"Up"
global-up="",'k','K',"Up"

# To select items
# default="","Enter"
global-enter="","Enter"

# To go back
# default="[Esc]","Esc"
global-back="[Esc]","Esc"

# To go back and exit Tut
# default="[Q]uit",'q','Q'
global-exit="[Q]uit",'q','Q'

# Move to the top
# default="",'g',"Home"
main-home="",'g',"Home"

# Move to the bottom
# default="",'G',"End"
main-end="",'G',"End"

# Go to previous feed
# default="",'h','H',"Left"
main-prev-feed="",'h','H',"Left"

# Go to next feed
# default="",'l','L',"Right"
main-next-feed="",'l','L',"Right"

# Focus on the notification list
# default="[N]otifications",'n','N'
main-notification-focus="[N]otifications",'n','N'

# Compose a new toot
# default="",'c','C'
main-compose="",'c','C'

# Open avatar
# default="[A]vatar",'a','A'
status-avatar="[A]vatar",'a','A'

# Boost a toot
# default="[B]oost","Un[B]oost",'b','B'
status-boost="[B]oost","Un[B]oost",'b','B'

# Delete a toot
# default="[D]elete",'d','D'
status-delete="[D]elete",'d','D'

# Favorite a toot
# default="[F]avorite","Un[F]avorite",'f','F'
status-favorite="[F]avorite","Un[F]avorite",'f','F'

# Open toots media files
# default="[M]edia",'m','M'
status-media="[M]edia",'m','M'

# Open links
# default="[O]pen",'o','O'
status-links="[O]pen",'o','O'

# Open poll
# default="[P]oll",'p','P'
status-poll="[P]oll",'p','P'

# Reply to toot
# default="[R]eply",'r','R'
status-reply="[R]eply",'r','R'

# Save/bookmark a toot
# default="[S]ave","Un[S]ave",'s','S'
status-bookmark="[S]ave","Un[S]ave",'s','S'

# View thread
# default="[T]hread",'t','T'
status-thread="[T]hread",'t','T'

# Open user profile
# default="[U]ser",'u','U'
status-user="[U]ser",'u','U'

# Open the view mode
# default="[V]iew",'v','V'
status-view-focus="[V]iew",'v','V'

# Yank the url of the toot
# default="[Y]ank",'y','Y'
status-yank="[Y]ank",'y','Y'

# Remove the spoiler
# default="Press [Z] to toggle spoiler",'z','Z'
status-toggle-spoiler="Press [Z] to toggle spoiler",'z','Z'

# View avatar
# default="[A]vatar",'a','A'
user-avatar="[A]vatar",'a','A'

# Block the user
# default="[B]lock","Un[B]lock",'b','B'
user-block="[B]lock","Un[B]lock",'b','B'

# Follow user
# default="[F]ollow","Un[F]ollow",'f','F'
user-follow="[F]ollow","Un[F]ollow",'f','F'

# Follow user
# default="Follow [R]equest","Follow [R]equest",'r','R'
user-follow-request-decide="Follow [R]equest","Follow [R]equest",'r','R'

# Mute user
# default="[M]ute","Un[M]ute",'m','M'
user-mute="[M]ute","Un[M]ute",'m','M'

# Open links
# default="[O]pen",'o','O'
user-links="[O]pen",'o','O'

# View user profile
# default="[U]ser",'u','U'
user-user="[U]ser",'u','U'

# Open view mode
# default="[V]iew",'v','V'
user-view-focus="[V]iew",'v','V'

# Yank the user URL
# default="[Y]ank",'y','Y'
user-yank="[Y]ank",'y','Y'

# Open list
# default="[O]pen",'o','O'
list-open-feed="[O]pen",'o','O'

# Open URL
# default="[O]pen",'o','O'
link-open="[O]pen",'o','O'

# Yank the URL
# default="[Y]ank",'y','Y'
link-yank="[Y]ank",'y','Y'

# Edit spoiler text on new toot
# default="[C]W text",'c','C'
compose-edit-spoiler="[C]W text",'c','C'

# Edit the text on new toot
# default="[E]dit text",'e','E'
compose-edit-text="[E]dit text",'e','E'

# Include a quote when replying
# default="[I]nclude quote",'i','I'
compose-include-quote="[I]nclude quote",'i','I'

# Focus on adding media to toot
# default="[M]edia",'m','M'
compose-media-focus="[M]edia",'m','M'

# Post the new toot
# default="[P]ost",'p','P'
compose-post="[P]ost",'p','P'

# Toggle content warning on toot
# default="[T]oggle CW",'t','T'
compose-toggle-content-warning="[T]oggle CW",'t','T'

# Edit the visibility on new toot
# default="[V]isibility",'v','V'
compose-visibility="[V]isibility",'v','V'

# Delete media file
# default="[D]elete",'d','D'
media-delete="[D]elete",'d','D'

# Edit the description on media file
# default="[E]dit desc",'e','E'
media-edit-desc="[E]dit desc",'e','E'

# Add a new media file
# default="[A]dd",'a','A'
media-add="[A]dd",'a','A'

# Vote on poll
# default="[V]ote",'v','V'
vote-vote="[V]ote",'v','V'

# Select item to vote on
# default="[Enter] to select",' ', "Enter"
vote-select="[Enter] to select",' ', "Enter"
`
