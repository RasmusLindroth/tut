% tut(5) tut 1.0.35
% Rasmus Lindroth
% 2023-01-07

# NAME
tut - configuration for tut(1)

# DESCRIPTION
The configuration format for tut.

You find it in *$XDG_CONFIG_HOME/tut/config.ini* on Linux which usually equals to *~/.config/tut/config.ini*.
If you don't run Linux it will use the path of the Go funcdtion os.UserConfigDir().
But if you move the tut folder to *XDG_CONFIG_HOME/tut/* and have set the environment variable *XDG_CONFIG_HOME*
it will look there instead of the standard place.

# CONFIGURATION
The configuration file is divided in seven sections named general, media, open-custom, open-pattern, desktop-notification, style and input.

Under each section there is the name of the configuration option. The last line under each options shows the default value. 

# GENERAL
This section is \[general\] in your configuration file

## confirmation
You need to press yes in a confirmation dialog before favoriting, boosting, etc.  

valid: true, false

**confirmation**=*true*

## mouse-support
Enable mouse support in tut.  

valid: true, false

**mouse-support**=*false*

## date-format
The date format to be used. See https://pkg.go.dev/time\#pkg-constants  
**date-format**=*"2006-01-02 15:04"*

## date-tody-format
Format for dates the same day. See date-format for more info.  
**date-tody-format**=*"15:04"*

## date-relative
This displays relative dates instead for statuses that are one day or older the output is 1y2m1d (1 year 2 months and 1 day)  
  
The value is an integear  
-1     = don\'t use relative dates  
 0     = always use relative dates, except for dates \< 1 day  
 1 - ∞ = number of days to use relative dates  
				  
Value: 28 will display a relative date for toots that are between 1-28 days old. Otherwhise it will use the short or long format.  
**date-relative**=*-1*

## max-width
The max with of text before it wraps when displaying a toot.  
**max-width**=*0*

## list-placement
The placement of your panes.  

valid: left, right, top, bottom

**list-placement**=*"left"*

## list-split
How should panes be split?  

valid: row, column

**list-split**=*"row"*

## list-proportion
The proportion of panes vs. content. 1 on this and 3 on content below results in content taking up 3 times more space.  
**list-proportion**=*1*

## content-proportion
See previous.  
**content-proportion**=*2*

## notifications-to-hide
Hide notifications of this type in your notification timelines.  

valid: mention, status, boost, follow, follow_request, favorite, poll, edit

**notifications-to-hide**=*[]*

## quote-reply
Always include a quote of the message you\'re replying to.  

valid: true, false

**quote-reply**=*false*

## show-icons
If you want to show icons in timelines.  

valid: true, false

**show-icons**=*true*

## short-hints
If you only want to you the letter of keys instead of the full hint.  

valid: true, false

**short-hints**=*false*

## show-filter-phrase
If you want to display the filter that filtered a toot.  

valid: true, false

**show-filter-phrase**=*true*

## show-help
Display a message in the commandbar on how to access the help text.  

valid: true, false

**show-help**=*true*

## stick-to-top
Always jump to the newest post. May ruin your reading experience.  

valid: true, false

**stick-to-top**=*false*

## show-boosted-user
Display the username of the person being boosted insted of the person that boosted.  

valid: true, false

**show-boosted-user**=*false*

## terminal-title
0 = No terminal title  
1 = Show title in terminal and top bar  
2 = Only show terminal title, and no top bar in tut.  

valid: 0, 1, 2

**terminal-title**=*0*

## redraw-ui
If you don\'t want the whole UI to update, and only update the text content you can disable this. This will lead to some artifacts being left on the screen when emojis are present.  

valid: true, false

**redraw-ui**=*true*

## leader-key
The leader is used as a shortcut to run commands as you can do in Vim. By default this is disabled and you enable it by setting a key here. It can only consist of one char, so set it to something like a comma.  
**leader-key**=*""*

## leader-timeout
Number of milliseconds before the leader command resets. So if you tap the leader-key by mistake or are to slow it empties all the input after X milliseconds.  
**leader-timeout**=*1000*

# GENERAL.TIMELINES
This section is \[\[general.timelines\]\] in your configuration file. You can have multiple of them.

Example:

[[general.timelines]]  
name=\"home\"  
type=\"home\"  
show-boosts=true  
show-replies=true  
  
[[general.timelines]]  
name = \"Notifications\"  
type = \"notifications\"  
keys = [\"n\", \"N\"]  
closed = true  
on-creation-closed = \"new-pane\"  
on-focus=\"focus-self\"  

## name
The name to display above the timeline  
**name**=*""*

## type
The type of the timeline  

valid: home, direct, local, federated, bookmarks, saved, favorited, notifications, lists, mentions, tag

**type**=*""*

## data
Used for the tag type, so here you set the tag.  
**data**=*""*

## keys
A list of keys to give this timeline focus. See under the input section to learn more about keys.  
**keys**=*[]*

## special-keys
A list of special-keys to give this timeline focus. See under the input section to learn more about special-keys.  
**special-keys**=*[]*

## shortcut
A shortcut to give this timeline focus with your leader-key + this shortcut.  
**shortcut**=*""*

## hide-boosts
Hide boosts in this timeline.  

valid: true, false

**hide-boosts**=*"false"*

## hide-replies
Hide replies in this timeline.  

valid: true, false

**hide-replies**=*"false"*

## closed
Don\'t open this timeline when you start tut. Use your keys or shortcut to open it.  

valid: true, false

**closed**=*"false"*

## on-creation-closed
Don\'t open this timeline when you start tut. Use your keys or shortcut to open it.  

valid: new-pane, current-pane

**on-creation-closed**=*"new-pane"*

## on-focus
Don\'t open this timeline when you start tut. Use your keys or shortcut to open it.  

valid: focus-pane, focus-self

**on-focus**=*"focus-pane"*

# GENERAL.LEADER-ACTIONS
This section is \[\[general.leader-actions\]\] in your configuration file. You can have multiple of them.

## type
The action you want to run.  

valid: blocking, boosts, clear-notifications, close-pane, compose, edit, favorited, favorites, followers, following, history, list-placement, list-split, lists, move-pane-left, move-pane-right, move-pane-up, move-pane-down, move-pane-home, move-pane-end, muting, newer, pane, preferences, profile, proportions, refetch, stick-to-top, tags

**type**=*""*

## data
Data to pass to the action.  
**data**=*""*

## shortcut
A shortcut to run this action with your leader-key + this shortcut.  
**shortcut**=*""*

# MEDIA
This section is \[media\] in your configuration file

# MEDIA.IMAGE
This section is \[media.image\] in your configuration file

## program
The program to open images. TUT_OS_DEFAULT equals xdg-open on Linux, open on MacOS and start on Windows.  
**program**=*"TUT_OS_DEFAULT"*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

## single
If the program should be called multiple times when there is multiple files. If set to false all files will be passed as an argument, but not all programs support this.  

valid: true, false

**single**=*true*

## reverse
If the files should be passed in reverse order. This will make some programs display the files in the correct order.  

valid: true, false

**reverse**=*false*

# MEDIA.VIDEO
This section is \[media.video\] in your configuration file

## program
The program to open videos. TUT_OS_DEFAULT equals xdg-open on Linux, open on MacOS and start on Windows.  
**program**=*"TUT_OS_DEFAULT"*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

## single
If the program should be called multiple times when there is multiple files. If set to false all files will be passed as an argument, but not all programs support this.  

valid: true, false

**single**=*true*

## reverse
If the files should be passed in reverse order. This will make some programs display the files in the correct order.  

valid: true, false

**reverse**=*false*

# MEDIA.AUDIO
This section is \[media.audio\] in your configuration file

## program
The program to open audio. TUT_OS_DEFAULT equals xdg-open on Linux, open on MacOS and start on Windows.  
**program**=*"TUT_OS_DEFAULT"*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

## single
If the program should be called multiple times when there is multiple files. If set to false all files will be passed as an argument, but not all programs support this.  

valid: true, false

**single**=*true*

## reverse
If the files should be passed in reverse order. This will make some programs display the files in the correct order.  

valid: true, false

**reverse**=*false*

# MEDIA.LINK
This section is \[media.link\] in your configuration file

## program
The program to open links. TUT_OS_DEFAULT equals xdg-open on Linux, open on MacOS and start on Windows.  
**program**=*"TUT_OS_DEFAULT"*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

# DESKTOP-NOTIFICATION
This section is \[desktop-notification\] in your configuration file

## followers
Enable notifications when someone follows you.  

valid: true, false

**followers**=*false*

## favorite
Enable notifications when one of your toots gets favorited.  

valid: true, false

**favorite**=*false*

## mention
Enable notifications  when someone mentions you.  

valid: true, false

**mention**=*false*

## update
Enable notifications when a post you have interacted with gets edited.  

valid: true, false

**update**=*false*

## boost
Enable notifications when one of your toots gets boosted.  

valid: true, false

**boost**=*false*

## poll
Enable notifications when a poll ends.  

valid: true, false

**poll**=*false*

## posts
Enable notifications for new posts.  

valid: true, false

**posts**=*false*

# OPEN-CUSTOM
This section is \[open-custom\] in your configuration file

# OPEN-CUSTOM.PROGRAMS
This section is \[\[open-custom.programs\]\] in your configuration file. You can have multiple of them.

## program
The program to open the file with.  
**program**=*""*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

## hint
What should the key hint in tut be for this program. See under the input section to learn more about hint.  
**hint**=*""*

## keys
A list of keys to to open files with this program. See under the input section to learn more about keys.  
**keys**=*[]*

## special-keys
A list of special-keys to open files with this program. See under the input section to learn more about special-keys.  
**special-keys**=*[]*

# OPEN-PATTERN
This section is \[open-pattern\] in your configuration file

# OPEN-PATTERN.PROGRAMS
This section is \[\[open-pattern.programs\]\] in your configuration file. You can have multiple of them.

## matching
Here you can set your own glob patterns for opening matching URLs in the program you want them to open up in. You could for example open Youtube videos in your video player instead of your default browser. To see the syntax for glob pattern you can follow this URL https://github.com/gobwas/glob\#syntax.  
**matching**=*""*

## program
The program to open the file with.  
**program**=*""*

## args
Arguments to pass to the program.  
**args**=*""*

## terminal
If the program runs in the terminal set this to true.  

valid: true, false

**terminal**=*false*

# STYLE
This section is \[style\] in your configuration file

All styles can be represented in their HEX value like \#ffffff or with their name, so in this case white. The only special value is \"default\" which equals to transparent, so it will be the same color as your terminal.  
You can also use xrdb colors like this xrdb:color1 The program will use colors prefixed with an \* first then look for URxvt or XTerm if it can\'t find any color prefixed with an asterisk. If you don\'t want tut to guess the prefix you can set the prefix yourself. If the xrdb color can\'t be found a preset color will be used. You\'ll have to set theme=\"none\" for this to work.  

## theme
The theme to use. You can use some themes that comes bundled with tut. Check out the themes available on the URL below. If a theme is named nord.ini you just write theme=\"nord\".  
  
https://github.com/RasmusLindroth/tut/tree/master/config/themes  
  
You can also create a theme file in your config directory e.g. ~/.config/tut/themes/foo.ini and then set theme=foo.  
  
If you want to use your own theme but don\'t want to create a new file, set theme=\"none\" and then you can create your own theme below.  
  
**theme**=*"default"*

## xrdb-prefix
The xrdb prefix used for colors in .Xresources.  
**xrdb-prefix**=*"guess"*

## background
The background color used on most elements.  
**background**=*""*

## text
The text color used on most of the text.  
**text**=*""*

## subtle
The color to display subtle elements or subtle text. Like lines and help text.  
**subtle**=*""*

## warning-text
The color for errors or warnings  
**warning-text**=*""*

## text-special-one
This color is used to display username.  
**text-special-one**=*""*

## text-special-two
This color is used to display username and key hints.  
**text-special-two**=*""*

## top-bar-background
The color of the bar at the top  
**top-bar-background**=*""*

## top-bar-text
The color of the text in the bar at the top.  
**top-bar-text**=*""*

## status-bar-background
The color of the bar at the bottom  
**status-bar-background**=*""*

## status-bar-text
The color of the text in the bar at the bottom.  
**status-bar-text**=*""*

## status-bar-view-background
The color of the bar at the bottom in view mode.  
**status-bar-view-background**=*""*

## status-bar-view-text
The color of the text in the bar at the bottom in view mode.  
**status-bar-view-text**=*""*

## command-text
The color of the text in the command bar at the bottom.  
**command-text**=*""*

## list-selected-background
Background of selected list items.  
**list-selected-background**=*""*

## list-selected-text
The text color of selected list items.  
**list-selected-text**=*""*

## list-selected-inactive-background
The background color of selected list items that are out of focus.  
**list-selected-inactive-background**=*""*

## list-selected-inactive-text
The text color of selected list items that are out of focus.  
**list-selected-inactive-text**=*""*

## controls-text
The main color of the text for key hints  
**controls-text**=*""*

## controls-highlight
The highlight color of for key hints  
**controls-highlight**=*""*

## autocomplete-background
The background color in dropdowns and autocompletions  
**autocomplete-background**=*""*

## autocomplete-text
The text color in dropdowns at autocompletions  
**autocomplete-text**=*""*

## autocomplete-selected-background
The background color for selected value in dropdowns and autocompletions  
**autocomplete-selected-background**=*""*

## autocomplete-selected-text
The text color for selected value in dropdowns and autocompletions  
**autocomplete-selected-text**=*""*

## button-color-one
The background color on selected button and the text color of unselected buttons  
**button-color-one**=*""*

## button-color-two
The text color on selected button and the background color of unselected buttons  
**button-color-two**=*""*

## timeline-name-background
The background on named timelines.  
**timeline-name-background**=*""*

## timeline-name-text
The text color on named timelines  
**timeline-name-text**=*""*

# INPUT
This section is \[input\] in your configuration file

In this section you set the keys to be used in tut.  
		  
The hint option lets you set which part of the hint that will be highlighted in tut. E.g. [F]avorite results in a highlighted F and the rest of the text is displayed normaly.  
Some of the options can be in two states, like favorites, so there you can set the hint-alt option to something like Un[F]avorite.  
  
Examples:  
\"[D]elete\" = Delete with a highlighted D  
\"Un[F]ollow\" = UnFollow with a highlighted F  
\"[Enter]\" = Enter where everything is highlighted  
\"Yan[K]\" = YanK with a highlighted K  
  
The keys option lets you define what key that should be pressed. This is limited to on character only and they are case sensetive.  
Example:  
keys=[\"j\",\"J\"]  
  
You can also set special-keys and they\'re for keys like Escape and Enter. To find the names of special keys you have to go to the following site and look for \"var KeyNames = map[Key]string{\"  
  
https://github.com/gdamore/tcell/blob/master/key.go  

# INPUT.GLOBAL-DOWN
This section is \[input.global-down\] in your configuration file

Keys for moving down  

## keys
**keys**=*["j","J"]*

## special-keys
**special-keys**=*["Down"]*

# INPUT.GLOBAL-UP
This section is \[input.global-up\] in your configuration file

Keys for moving down  

## keys
**keys**=*["k","K"]*

## special-keys
**special-keys**=*["Up"]*

# INPUT.GLOBAL-ENTER
This section is \[input.global-enter\] in your configuration file

To select items  

## special-keys
**special-keys**=*["Enter"]*

# INPUT.GLOBAL-BACK
This section is \[input.global-back\] in your configuration file

To go back  

## hint
**hint**=*"[Esc]"*

## special-keys
**special-keys**=*["Esc"]*

# INPUT.GLOBAL-EXIT
This section is \[input.global-exit\] in your configuration file

To go back or exit  

## hint
**hint**=*"[Q]uit"*

## keys
**keys**=*["q","Q"]*

# INPUT.MAIN-HOME
This section is \[input.main-home\] in your configuration file

Move to the top  

## keys
**keys**=*["g"]*

## special-keys
**special-keys**=*["Home"]*

# INPUT.MAIN-END
This section is \[input.main-end\] in your configuration file

Move to the bottom  

## keys
**keys**=*["G"]*

## special-keys
**special-keys**=*["End"]*

# INPUT.MAIN-PREV-FEED
This section is \[input.main-prev-feed\] in your configuration file

Go to previous feed  

## keys
**keys**=*["h","H"]*

## special-keys
**special-keys**=*["Left"]*

# INPUT.MAIN-NEXT-FEED
This section is \[input.main-next-feed\] in your configuration file

Go to next feed  

## keys
**keys**=*["l","L"]*

## specialkeys
**specialkeys**=*["Right"]*

# INPUT.MAIN-PREV-PANE
This section is \[input.main-prev-pane\] in your configuration file

Focus on the previous feed pane  

## special-keys
**special-keys**=*["Backtab"]*

# INPUT.MAIN-NEXT-PANE
This section is \[input.main-next-pane\] in your configuration file

Focus on the next feed pane  

## special-keys
**special-keys**=*["Tab"]*

# INPUT.MAIN-NEXT-ACCOUNT
This section is \[input.main-next-account\] in your configuration file

Focus on the next account  

## special-keys
**special-keys**=*["Ctrl-N"]*

# INPUT.MAIN-PREV-ACCOUNT
This section is \[input.main-prev-account\] in your configuration file

Focus on the previous account  

## special-keys
**special-keys**=*["Ctrl-P"]*

# INPUT.MAIN-COMPOSE
This section is \[input.main-compose\] in your configuration file

Compose a new toot  

## keys
**keys**=*["c","C"]*

# INPUT.STATUS-AVATAR
This section is \[input.status-avatar\] in your configuration file

Open avatar  

## hint
**hint**=*"[A]vatar"*

## keys
**keys**=*["a","A"]*

# INPUT.STATUS-BOOST
This section is \[input.status-boost\] in your configuration file

Boost a toot  

## hint
**hint**=*"[B]oost"*

## keys
**keys**=*["b","B"]*

# INPUT.STATUS-EDIT
This section is \[input.status-edit\] in your configuration file

Edit a toot  

## hint
**hint**=*"[E]dit"*

## keys
**keys**=*["e","E"]*

# INPUT.STATUS-DELETE
This section is \[input.status-delete\] in your configuration file

Delete a toot  

## hint
**hint**=*"[D]elete"*

## keys
**keys**=*["d","D"]*

# INPUT.STATUS-FAVORITE
This section is \[input.status-favorite\] in your configuration file

Favorite a toot  

## hint
**hint**=*"[F]avorite"*

## keys
**keys**=*["f","F"]*

# INPUT.STATUS-MEDIA
This section is \[input.status-media\] in your configuration file

Open toots media files  

## hint
**hint**=*"[M]edia"*

## keys
**keys**=*["m","M"]*

# INPUT.STATUS-LINKS
This section is \[input.status-links\] in your configuration file

Open links  

## hint
**hint**=*"[O]pen"*

## keys
**keys**=*["o","O"]*

# INPUT.STATUS-POLL
This section is \[input.status-poll\] in your configuration file

Open poll  

## hint
**hint**=*"[P]oll"*

## keys
**keys**=*["p","P"]*

# INPUT.STATUS-REPLY
This section is \[input.status-reply\] in your configuration file

Reply to toot  

## hint
**hint**=*"[R]eply"*

## keys
**keys**=*["r","R"]*

# INPUT.STATUS-BOOKMARK
This section is \[input.status-bookmark\] in your configuration file

Save/bookmark a toot  

## hint
**hint**=*"[S]ave"*

## hint-alt
**hint-alt**=*"Un[S]ave"*

## keys
**keys**=*["s","S"]*

# INPUT.STATUS-THREAD
This section is \[input.status-thread\] in your configuration file

View thread  

## hint
**hint**=*"[T]hread"*

## keys
**keys**=*["t","T"]*

# INPUT.STATUS-USER
This section is \[input.status-user\] in your configuration file

Open user profile  

## hint
**hint**=*"[U]ser"*

## keys
**keys**=*["u","U"]*

# INPUT.STATUS-VIEW-FOCUS
This section is \[input.status-view-focus\] in your configuration file

Open the view mode  

## hint
**hint**=*"[V]iew"*

## keys
**keys**=*["v","V"]*

# INPUT.STATUS-YANK
This section is \[input.status-yank\] in your configuration file

Yank the url of the toot  

## hint
**hint**=*"[Y]ank"*

## keys
**keys**=*["y","Y"]*

# INPUT.STATUS-TOGGLE-CW
This section is \[input.status-toggle-cw\] in your configuration file

Show the content in a content warning  

## hint
**hint**=*"Press [Z] to toggle cw"*

## keys
**keys**=*["z","Z"]*

# INPUT.STATUS-SHOW-FILTERED
This section is \[input.status-show-filtered\] in your configuration file

Show the content of a filtered toot  

## hint
**hint**=*"Press [Z] to view filtered toot"*

## keys
**keys**=*["z","Z"]*

# INPUT.USER-AVATAR
This section is \[input.user-avatar\] in your configuration file

View avatar  

## hint
**hint**=*"[A]vatar"*

## keys
**keys**=*["a","A"]*

# INPUT.USER-BLOCK
This section is \[input.user-block\] in your configuration file

Block the user  

## hint
**hint**=*"[B]lock"*

## hint-alt
**hint-alt**=*"Un[B]lock"*

## keys
**keys**=*["b","B"]*

# INPUT.USER-FOLLOW
This section is \[input.user-follow\] in your configuration file

Follow user  

## hint
**hint**=*"[F]ollow"*

## hint-alt
**hint-alt**=*"Un[F]ollow"*

## keys
**keys**=*["f","F"]*

# INPUT.USER-FOLLOW-REQUEST-DECIDE
This section is \[input.user-follow-request-decide\] in your configuration file

Follow user  

## hint
**hint**=*"Follow [R]equest"*

## hint-alt
**hint-alt**=*"Follow [R]equest"*

## keys
**keys**=*["r","R"]*

# INPUT.USER-MUTE
This section is \[input.user-mute\] in your configuration file

Mute user  

## hint
**hint**=*"[M]ute"*

## hint-alt
**hint-alt**=*"Un[M]ute"*

## keys
**keys**=*["m","M"]*

# INPUT.USER-LINKS
This section is \[input.user-links\] in your configuration file

Open links  

## hint
**hint**=*"[O]pen"*

## keys
**keys**=*["o","O"]*

# INPUT.USER-USER
This section is \[input.user-user\] in your configuration file

View user profile  

## hint
**hint**=*"[U]ser"*

## keys
**keys**=*["u","U"]*

# INPUT.USER-VIEW-FOCUS
This section is \[input.user-view-focus\] in your configuration file

Open view mode  

## hint
**hint**=*"[V]iew"*

## keys
**keys**=*["v","V"]*

# INPUT.USER-YANK
This section is \[input.user-yank\] in your configuration file

Yank the user URL  

## hint
**hint**=*"[Y]ank"*

## keys
**keys**=*["y","Y"]*

# INPUT.LIST-OPEN-FEED
This section is \[input.list-open-feed\] in your configuration file

Open list  

## hint
**hint**=*"[O]pen"*

## keys
**keys**=*["o","O"]*

# INPUT.LIST-USER-LIST
This section is \[input.list-user-list\] in your configuration file

List all users in a list  

## hint
**hint**=*"[U]sers"*

## keys
**keys**=*["u","U"]*

# INPUT.LIST-USER-ADD
This section is \[input.list-user-add\] in your configuration file

Add user to list  

## hint
**hint**=*"[A]dd"*

## keys
**keys**=*["a","A"]*

# INPUT.LIST-USER-DELETE
This section is \[input.list-user-delete\] in your configuration file

Delete user from list  

## hint
**hint**=*"[D]elete"*

## keys
**keys**=*["d","D"]*

# INPUT.LINK-OPEN
This section is \[input.link-open\] in your configuration file

Open URL  

## hint
**hint**=*"[O]pen"*

## keys
**keys**=*["o","O"]*

# INPUT.LINK-YANK
This section is \[input.link-yank\] in your configuration file

Yank the URL  

## hint
**hint**=*"[Y]ank"*

## keys
**keys**=*["y","Y"]*

# INPUT.TAG-OPEN-FEED
This section is \[input.tag-open-feed\] in your configuration file

Open tag feed  

## hint
**hint**=*"[O]pen"*

## keys
**keys**=*["o","O"]*

# INPUT.TAG-FOLLOW
This section is \[input.tag-follow\] in your configuration file

Toggle follow on tag  

## hint
**hint**=*"[F]ollow"*

## hint-alt
**hint-alt**=*"Un[F]ollow"*

## keys
**keys**=*["f","F"]*

# INPUT.COMPOSE-EDIT-CW
This section is \[input.compose-edit-cw\] in your configuration file

Edit content warning text on new toot  

## hint
**hint**=*"[C]W text"*

## keys
**keys**=*["c","C"]*

# INPUT.COMPOSE-EDIT-TEXT
This section is \[input.compose-edit-text\] in your configuration file

Edit the text on new toot  

## hint
**hint**=*"[E]dit text"*

## keys
**keys**=*["e","E"]*

# INPUT.COMPOSE-INCLUDE-QUOTE
This section is \[input.compose-include-quote\] in your configuration file

Include a quote when replying  

## hint
**hint**=*"[I]nclude quote"*

## keys
**keys**=*["i","I"]*

# INPUT.COMPOSE-MEDIA-FOCUS
This section is \[input.compose-media-focus\] in your configuration file

Focus on adding media to toot  

## hint
**hint**=*"[M]edia"*

## keys
**keys**=*["m","M"]*

# INPUT.COMPOSE-POST
This section is \[input.compose-post\] in your configuration file

Post the new toot  

## hint
**hint**=*"[P]ost"*

## keys
**keys**=*["p","P"]*

# INPUT.COMPOSE-TOGGLE-CONTENT-WARNING
This section is \[input.compose-toggle-content-warning\] in your configuration file

Toggle content warning on toot  

## hint
**hint**=*"[T]oggle CW"*

## keys
**keys**=*["t","T"]*

# INPUT.COMPOSE-VISIBILITY
This section is \[input.compose-visibility\] in your configuration file

Edit the visibility on new toot  

## hint
**hint**=*"[V]isibility"*

## keys
**keys**=*["v","V"]*

# INPUT.COMPOSE-LANGUAGE
This section is \[input.compose-language\] in your configuration file

Edit the language of a toot  

## hint
**hint**=*"[L]ang"*

## keys
**keys**=*["l","L"]*

# INPUT.COMPOSE-POLL
This section is \[input.compose-poll\] in your configuration file

Switch to creating a poll  

## hint
**hint**=*"P[O]ll"*

## keys
**keys**=*["o","O"]*

# INPUT.MEDIA-DELETE
This section is \[input.media-delete\] in your configuration file

Delete media file  

## hint
**hint**=*"[D]elete"*

## keys
**keys**=*["d","D"]*

# INPUT.MEDIA-EDIT-DESC
This section is \[input.media-edit-desc\] in your configuration file

Edit the description on media file  

## hint
**hint**=*"[E]dit desc"*

## keys
**keys**=*["e","E"]*

# INPUT.MEDIA-ADD
This section is \[input.media-add\] in your configuration file

Add a new media file  

## hint
**hint**=*"[A]dd"*

## keys
**keys**=*["a","A"]*

# INPUT.VOTE-VOTE
This section is \[input.vote-vote\] in your configuration file

Vote on poll  

## hint
**hint**=*"[V]ote"*

## keys
**keys**=*["v","V"]*

# INPUT.VOTE-SELECT
This section is \[input.vote-select\] in your configuration file

Select item to vote on  

## hint
**hint**=*"[Enter] to select"*

## special-keys
**special-keys**=*["Enter"]*

# INPUT.POLL-ADD
This section is \[input.poll-add\] in your configuration file

Add a new poll option  

## hint
**hint**=*"[A]dd"*

## keys
**keys**=*["a","A"]*

# INPUT.POLL-EDIT
This section is \[input.poll-edit\] in your configuration file

Edit a poll option  

## hint
**hint**=*"[E]dit"*

## keys
**keys**=*["e","E"]*

# INPUT.POLL-DELETE
This section is \[input.poll-delete\] in your configuration file

Delete a poll option  

## hint
**hint**=*"[D]elete"*

## keys
**keys**=*["d","D"]*

# INPUT.POLL-MULTI-TOGGLE
This section is \[input.poll-multi-toggle\] in your configuration file

Toggle voting on multiple options  

## hint
**hint**=*"Toggle [M]ultiple"*

## keys
**keys**=*["m","M"]*

# INPUT.POLL-EXPIRATION
This section is \[input.poll-expiration\] in your configuration file

Change the expiration of poll  

## hint
**hint**=*"E[X]pires"*

## keys
**keys**=*["x","X"]*

# INPUT.PREFERENCE-NAME
This section is \[input.preference-name\] in your configuration file

Change display name  

## hint
**hint**=*"[N]ame"*

## keys
**keys**=*["n","N"]*

# INPUT.PREFERENCE-VISIBILITY
This section is \[input.preference-visibility\] in your configuration file

Change default visibility of toots  

## hint
**hint**=*"[V]isibility"*

## keys
**keys**=*["v","V"]*

# INPUT.PREFERENCE-BIO
This section is \[input.preference-bio\] in your configuration file

Change bio in profile  

## hint
**hint**=*"[B]io"*

## keys
**keys**=*["b","B"]*

# INPUT.PREFERENCE-SAVE
This section is \[input.preference-save\] in your configuration file

Save your preferences  

## hint
**hint**=*"[S]ave"*

## keys
**keys**=*["s","S"]*

# INPUT.PREFERENCE-FIELDS
This section is \[input.preference-fields\] in your configuration file

Edit profile fields  

## hint
**hint**=*"[F]ields"*

## keys
**keys**=*["f","F"]*

# INPUT.PREFERENCE-FIELDS-ADD
This section is \[input.preference-fields-add\] in your configuration file

Add new field  

## hint
**hint**=*"[A]dd"*

## keys
**keys**=*["a","A"]*

# INPUT.PREFERENCE-FIELDS-EDIT
This section is \[input.preference-fields-edit\] in your configuration file

Edit current field  

## hint
**hint**=*"[E]dit"*

## keys
**keys**=*["e","E"]*

# INPUT.PREFERENCE-FIELDS-DELETE
This section is \[input.preference-fields-delete\] in your configuration file

Delete current field  

## hint
**hint**=*"[D]elete"*

## keys
**keys**=*["d","D"]*

# INPUT.EDITOR-EXIT
This section is \[input.editor-exit\] in your configuration file

Exit the editor  

## hint
**hint**=*"[Esc] when done"*

## special-keys
**special-keys**=*["Esc"]*

# SEE ALSO
    tut(1) - flags and commands
    tut(7) - commands and keys inside of tut
