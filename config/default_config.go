package config

var conftext = `# Configuration file for tut

[general]
# Shows a confirmation view before actions such as favorite, delete toot, boost
# etc.
# default=true
confirmation=true

# Enable support for using the mouse in tut to select items.
# default=false
mouse-support=false

# Timelines adds windows of feeds. You can customize the number of feeds, what
# they should show and the key to activate them.
# 
# Available timelines: home, direct, local, federated, special, bookmarks,
# saved, favorited, notifications, lists, tag
# 
# The one named special are the home timeline with only boosts and/or replies.
# 
# Tag is special as you need to add the tag after, see the example below.
# 
# The syntax is:
# timelines=feed,[name],[keys...],[showBoosts],[showReplies]
# 
# Tha values in brackets are optional. You can see the syntax for keys under the
# [input] section.
# 
# showBoosts and showReplies must be formated as bools. So either true or false.
# They always defaults to true.
# 
# Some examples:
# 
# home timeline with the name Home
# timelines=home,Home
# 
# local timeline with the name Local and it gets focus when you press 2. It will
# also hide boosts in the timeline, but show toots that are replies.
# timelines=local,Local,'2',false,true
# 
# notification timeline with the name [N]otifications and it gets focus when you
# press n or N
# timelines=notifications,[N]otifications,'n','N'
# 
# tag timeline for #linux with the name Linux and it gets focus when you press
# timelines=tag linux,Linux,"F2"
# 
# 
# If you don't set any timelines it will default to this:
# timelines=home
# timelines=notifications,[N]otifications,'n','N'
# 


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

# The max width of text before it wraps when displaying toots.
# 0 = no restriction.
# default=0
max-width=0

# Where do you want the list of toots to be placed?
# Valid values: left, right, top, bottom.
# default=left
list-placement=left

# If you have notification-feed set to true you can display it under the main
# list of toots (row) or place it to the right of the main list of toots
# (column).
# default=row
list-split=row

# You can change the proportions of the list view in relation to the content
# view list-proportion=1 and content-proportoin=3 will result in the content
# taking up 3 times more space.
# Must be n > 0
# default=1
list-proportion=1

# See list-proportion
# default=2
content-proportion=2

# Hide notifications of this type. If you have multiple you separate them with a
# comma. Valid types: mention, status, boost, follow, follow_request, favorite,
# poll, edit.
# default=
notifications-to-hide=

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

# If you want to display the filter that filtered a toot.
# default=true
show-filter-phrase=true

# If you want to show a message in the cmdbar on how to access the help text.
# default=true
show-help=true

# If you always want tut to jump to the newest post. May ruin your reading
# experience.
# default=false
stick-to-top=false

# If you want to display the username of the person being boosted instead of the
# person that boosted.
# default=false
show-boosted-user=false

# 0 = No terminal title
# 1 = Show title in terminal and top bar
# 2 = Only show terminal title, and no top bar in tut.
# default=0
terminal-title=0

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
# of two parts first the action then the shortcut. And they're separated by a
# comma.
# 
# Available commands: home, direct, local, federated, special-all,
# special-boosts, special-replies, clear-notifications, compose, edit, history,
# blocking, bookmarks, refetch, saved, favorited, boosts, favorites, following,
# followers, muting, newer, preferences, profile, notifications, lists,
# stick-to-top, tag, tags, window, list-placement, list-split, proportions
# 
# The ones named special-* are the home timeline with only boosts and/or
# replies. All contains both, -boosts only boosts and -replies only replies.
# 
# The shortcuts are up to you, but keep them quite short and make sure they
# don't collide. If you have one shortcut that is "f" and an other one that is
# "fav", the one with "f" will always run and "fav" will never run. 
# 
# Some special leaders:
# tag is special as you need to add the tag after, e.g. tag linux
# window is special as it's a shortcut for switching between the timelines
# you've set under general and they are zero indexed. window 0 = your first
# timeline, window 1 = your second and so on.
# list-placement as it takes the argument top, right, bottom or left
# list-split as it takes the argument column or row
# proportions takes the arguments [int] [int], where the first integer is the
# list and the other content, e.g. proportions 1 3. See list-proportion above
# for more information.
# 
# Some examples:
# leader-action=local,lo
# leader-action=lists,li
# leader-action=federated,fed
# leader-action=direct,d
# leader-action=history,h
# leader-action=tag linux,tl
# leader-action=window 0,h
# leader-action=list-placement bottom,b
# leader-action=list-split column,c
# leader-action=proportions 1 3,3
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
# This sections allows you to set up to five custom programs to open URLs with.
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

# Notification when someone edits their toot.
# default=false
update=false

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
# color prefixed with an asterisk. If you don't want tut to guess the prefix you
# can set the prefix yourself. If the xrdb color can't be found a preset color
# will be used. You'll have to set theme=none for this to work.

# The xrdb prefix used for colors in .Xresources.
# default=guess
xrdb-prefix=guess

# You can use some themes that comes bundled with tut. Check out the themes
# available on the URL below. If a theme is named "nord.ini" you just write
# theme=nord
# 
# https://github.com/RasmusLindroth/tut/tree/master/config/themes
# 
# You can also create a theme file in your config directory e.g.
# ~/.config/tut/themes/foo.ini and then set theme=foo.
# 
# If you want to use your own theme but don't want to create a new file, set
# theme=none and then you can create your own theme below.
# default=default
theme=default

# The background color used on most elements.
# default=
background=

# The text color used on most of the text.
# default=
text=



# The color to display subtle elements or subtle text. Like lines and help text.
# default=
subtle=

# The color for errors or warnings
# default=
warning-text=

# This color is used to display username.
# default=
text-special-one=

# This color is used to display username and key hints.
# default=
text-special-two=

# The color of the bar at the top
# default=
top-bar-background=

# The color of the text in the bar at the top.
# default=
top-bar-text=

# The color of the bar at the bottom
# default=
status-bar-background=

# The color of the text in the bar at the bottom.
# default=
status-bar-text=

# The color of the bar at the bottom in view mode.
# default=
status-bar-view-background=

# The color of the text in the bar at the bottom in view mode.
# default=
status-bar-view-text=

# The color of the text in the command bar at the bottom.
# default=
command-text=

# Background of selected list items.
# default=
list-selected-background=

# The text color of selected list items.
# default=
list-selected-text=

# The background color of selected list items that are out of focus.
# default=
list-selected-inactive-background=

# The text color of selected list items that are out of focus.
# default=
list-selected-inactive-text=

# The main color of the text for key hints
# default=
controls-text=

# The highlight color of for key hints
# default=
controls-highlight=

# The background color in dropdowns and autocompletions
# default=
autocomplete-background=

# The text color in dropdowns at autocompletions
# default=
autocomplete-text=

# The background color for selected value in dropdowns and autocompletions
# default=
autocomplete-selected-background=

# The text color for selected value in dropdowns and autocompletions
# default=
autocomplete-selected-text=

# The background color on selected button and the text color of unselected
# buttons
# default=
button-color-one=

# The text color on selected button and the background color of unselected
# buttons
# default=
button-color-two=

# The background on named timelines.
# default=
timeline-name-background=

# The text color on named timelines
# default=
timeline-name-text=

# The text color used for date/time in the timeline view 
# defaults to the same color as subtle
datetime-text=

# The text color used for boosts in the timeline view 
# defaults to the primary text color
boost-text=

# The text color used for toots with media or cards in the timeline view 
# defaults to the primary text color
media-text=

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
# for special keys like "Enter". Remember that they are case sensitive.
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

# Focus on the previous feed window
# default="","Backtab"
main-prev-window="","Backtab"

# Focus on the next feed window
# default="","Tab"
main-next-window="","Tab"

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

# Edit a toot
# default="[E]dit",'e','E'
status-edit="[E]dit",'e','E'

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

# Show the content in a content warning
# default="Press [Z] to toggle cw",'z','Z'
status-toggle-cw="Press [Z] to toggle cw",'z','Z'

# Show the content of a filtered toot
# default="Press [Z] to view filtered toot",'z','Z'
status-show-filtered="Press [Z] to view filtered toot",'z','Z'

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

# List all users in a list
# default="[U]sers",'u','U'
list-user-list="[U]sers",'u','U'

# Add user to list
# default="[A]dd",'a','A'
list-user-add="[A]dd",'a','A'

# Delete user from list
# default="[D]elete",'d','D'
list-user-delete="[D]elete",'d','D'

# Open URL
# default="[O]pen",'o','O'
link-open="[O]pen",'o','O'

# Yank the URL
# default="[Y]ank",'y','Y'
link-yank="[Y]ank",'y','Y'

# Open tag feed
# default="[O]pen",'o','O'
tag-open-feed="[O]pen",'o','O'

# Toggle follow on tag
# default="[F]ollow","Un[F]ollow",'f','F'
tag-follow="[F]ollow","Un[F]ollow",'f','F'

# Edit content warning text on new toot
# default="[C]W text",'c','C'
compose-edit-cw="[C]W text",'c','C'

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

# Edit the language of a toot
# default="[L]ang",'l','L'
compose-language="[L]ang",'l','L'

# Switch to creating a poll
# default="P[O]ll",'o','O'
compose-poll="P[O]ll",'o','O'

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

# Add a new poll option
# default="[A]dd",'a','A'
poll-add="[A]dd",'a','A'

# Edit a poll option
# default="[E]dit",'e','E'
poll-edit="[E]dit",'e','E'

# Delete a poll option
# default="[D]elete",'d','D'
poll-delete="[D]elete",'d','D'

# Toggle voting on multiple options
# default="Toggle [M]ultiple",'m','M'
poll-multi-toggle="Toggle [M]ultiple",'m','M'

# Change the expiration of poll
# default="E[X]pires",'x','X'
poll-expiration="E[X]pires",'x','X'

# Change display name
# default="[N]ame",'n','N'
preference-name="[N]ame",'n','N'

# Change default visibility of toots
# default="[V]isibility",'v','V'
preference-visibility="[V]isibility",'v','V'

# Change bio in profile
# default="[B]io",'b','B'
preference-bio="[B]io",'b','B'

# Save your preferences
# default="[S]ave",'s','S'
preference-save="[S]ave",'s','S'

# Edit profile fields
# default="[F]ields",'f','F'
preference-fields="[F]ields",'f','F'

# Add new field
# default="[A]dd",'a','A'
preference-fields-add="[A]dd",'a','A'

# Edit current field
# default="[E]dit",'e','E'
preference-fields-edit="[E]dit",'e','E'

# Delete current field
# default="[D]elete",'d','D'
preference-fields-delete="[D]elete",'d','D'
`
