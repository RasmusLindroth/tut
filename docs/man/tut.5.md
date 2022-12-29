% tut(5) tut 1.0.30
% Rasmus Lindroth
% 2022-12-29

# NAME
tut - configuration for tut(1)

# DESCRIPTION
The configuration format for tut.

You find the configuration file in *$XDG_CONFIG_HOME/tut/config.ini* which usually equals to *~/.config/tut/config.ini*.

# CONFIGURATION
The configuration file is divided in seven sections named general, media, open-custom, open-pattern, desktop-notification, style and input.

Under each section there is the name of the configuration option. The last line under each options shows the default value. 

# GENERAL
This section is \[general\] in your configuration file

## confirmation
Shows a confirmation view before actions such as favorite, delete toot, boost etc.  
**confirmation**=*true*

## mouse-support
Enable support for using the mouse in tut to select items.  
**mouse-support**=*false*

## timelines
Timelines adds windows of feeds. You can customize the number of feeds, what they should show and the key to activate them.  
  
Available timelines: home, direct, local, federated, special, bookmarks, saved, favorited, notifications, lists, mentions, tag  
  
The one named special are the home timeline with only boosts and/or replies.  
  
Tag is special as you need to add the tag after, see the example below.  
  
The syntax is:  
timelines=feed,[name],[keys...],[showBoosts],[showReplies]  
  
Tha values in brackets are optional. You can see the syntax for keys under the [input] section.  
  
showBoosts and showReplies must be formated as bools. So either true or false. They always defaults to true.  
  
Some examples:  
  
home timeline with the name Home  
timelines=home,Home  
  
local timeline with the name Local and it gets focus when you press 2. It will also hide boosts in the timeline, but show toots that are replies.  
timelines=local,Local,\'2\',false,true  
  
notification timeline with the name [N]otifications and it gets focus when you press n or N  
timelines=notifications,[N]otifications,\'n\',\'N\'  
  
tag timeline for \#linux with the name Linux and it gets focus when you press  
timelines=tag linux,Linux,\"F2\"  
  
  
If you don\'t set any timelines it will default to this:  
timelines=home  
timelines=notifications,[N]otifications,\'n\',\'N\'  
  


## date-format
The date format to be used. See https://godoc.org/time\#Time.Format  
**date-format**=*2006-01-02 15:04*

## date-today-format
Format for dates the same day. See date-format for more info.  
**date-today-format**=*15:04*

## date-relative
This displays relative dates instead for statuses that are one day or older the output is 1y2m1d (1 year 2 months and 1 day)  
  
The value is an integear  
-1     = don\'t use relative dates  
 0     = always use relative dates, except for dates \< 1 day  
 1 - ∞ = number of days to use relative dates  
  
Example: date-relative=28 will display a relative date for toots that are between 1-28 days old. Otherwhise it will use the short or long format.  
**date-relative**=*-1*

## max-width
The max width of text before it wraps when displaying toots.  
0 = no restriction.  
**max-width**=*0*

## list-placement
Where do you want the list of toots to be placed?  
Valid values: left, right, top, bottom.  
**list-placement**=*left*

## list-split
If you have notification-feed set to true you can display it under the main list of toots (row) or place it to the right of the main list of toots (column).  
**list-split**=*row*

## list-proportion
You can change the proportions of the list view in relation to the content view list-proportion=1 and content-proportoin=3 will result in the content taking up 3 times more space.  
Must be n \> 0  
**list-proportion**=*1*

## content-proportion
See list-proportion  
**content-proportion**=*2*

## notifications-to-hide
Hide notifications of this type. If you have multiple you separate them with a comma. Valid types: mention, status, boost, follow, follow_request, favorite, poll, edit.  
**notifications-to-hide**=

## quote-reply
If you always want to quote original message when replying.  
**quote-reply**=*false*

## char-limit
If you\'re on an instance with a custom character limit you can set it here.  
**char-limit**=*500*

## show-icons
If you want to show icons in the list of toots.  
**show-icons**=*true*

## short-hints
If you\'ve learnt all the shortcut keys you can remove the help text and only show the key in tui. So it gets less cluttered.  
**short-hints**=*false*

## show-filter-phrase
If you want to display the filter that filtered a toot.  
**show-filter-phrase**=*true*

## show-help
If you want to show a message in the cmdbar on how to access the help text.  
**show-help**=*true*

## stick-to-top
If you always want tut to jump to the newest post. May ruin your reading experience.  
**stick-to-top**=*false*

## show-boosted-user
If you want to display the username of the person being boosted instead of the person that boosted.  
**show-boosted-user**=*false*

## terminal-title
0 = No terminal title  
1 = Show title in terminal and top bar  
2 = Only show terminal title, and no top bar in tut.  
**terminal-title**=*0*

## redraw-ui
If you don\'t want the whole UI to update, and only the text content you can set this option to true. This will lead to some artifacts being left on the screen when emojis are present. But it will keep the UI from flashing on every single toot in some terminals.  
**redraw-ui**=*true*

## leader-key
The leader is used as a shortcut to run commands as you can do in Vim. By default this is disabled and you enable it by setting a leader-key. It can only consist of one char and I like to use comma as leader key. So to set it you write leader-key=,  
**leader-key**=

## leader-timeout
Number of milliseconds before the leader command resets. So if you tap the leader-key by mistake or are to slow it empties all the input after X milliseconds.  
**leader-timeout**=*1000*

## leader-action
You set actions for the leader-key with one or more leader-action. It consists of two parts first the action then the shortcut. And they\'re separated by a comma.  
  
Available commands: blocking, bookmarks, boosts, clear-notifications, close-window, compose, direct, edit, favorited, favorites, federated, followers, following, history, home, list-placement, list-split, lists, local, mentions, move-window-left, move-window-right, move-window-up, move-window-down, move-window-home, move-window-end, muting, newer, notifications, preferences, profile, proportions, refetch, saved, special-all, special-boosts, special-replies, stick-to-top, switch, tag, tags, window  
  
The ones named special-\* are the home timeline with only boosts and/or replies. All contains both, -boosts only boosts and -replies only replies.  
  
The shortcuts are up to you, but keep them quite short and make sure they don\'t collide. If you have one shortcut that is \"f\" and an other one that is \"fav\", the one with \"f\" will always run and \"fav\" will never run.   
  
Some special leaders:  
tag is special as you need to add the tag after, e.g. tag linux  
window is special as it\'s a shortcut for switching between the timelines you\'ve set under general and they are zero indexed. window 0 = your first timeline, window 1 = your second and so on.  
list-placement as it takes the argument top, right, bottom or left  
list-split as it takes the argument column or row  
proportions takes the arguments [int] [int], where the first integer is the list and the other content, e.g. proportions 1 3. See list-proportion above for more information.  
switch let\'s you go to a timeline if it already exists, if it doesn\'t it will open the timeline in a new window. The syntax is almost the same as in timelines= and is displayed under the examples.  
  
Some examples:  
leader-action=local,lo  
leader-action=lists,li  
leader-action=federated,fed  
leader-action=direct,d  
leader-action=history,h  
leader-action=tag linux,tl  
leader-action=window 0,h  
leader-action=list-placement bottom,b  
leader-action=list-split column,c  
leader-action=proportions 1 3,3  
  
Syntax for switch:  
leader-action=switch feed,shortcut,[name],[showBoosts],[showReplies]  
showBoosts can be either true or false and they are both optional. Here are some examples:  
  
leader-action=switch home,h,false,true  
leader-action=switch tag tut,tt  
  


# MEDIA
This section is \[media\] in your configuration file

## image-viewer
Your image viewer.  
**image-viewer**=*xdg-open*

## image-terminal
Open the image viewer in the same terminal as toot. Only for terminal based viewers.  
**image-terminal**=*false*

## image-single
If images should open one by one e.g. \"imv image.png\" multiple times. If set to false all images will open at the same time like this \"imv image1.png image2.png image3.png\". Not all image viewers support this, so try it first.  
**image-single**=*true*

## image-reverse
If you want to open the images in reverse order. In some image viewers this will display the images in the \"right\" order.  
**image-reverse**=*false*

## video-viewer
Your video viewer.  
**video-viewer**=*xdg-open*

## video-terminal
Open the video viewer in the same terminal as toot. Only for terminal based viewers.  
**video-terminal**=*false*

## video-single
If videos should open one by one. See image-single.  
**video-single**=*true*

## video-reverse
If you want your videos in reverse order. In some video apps this will play the files in the \"right\" order.  
**video-reverse**=*false*

## audio-viewer
Your audio viewer.  
**audio-viewer**=*xdg-open*

## audio-terminal
Open the audio viewer in the same terminal as toot. Only for terminal based viewers.  
**audio-terminal**=*false*

## audio-single
If audio should open one by one. See image-single.  
**audio-single**=*true*

## audio-reverse
If you want to play the audio files in reverse order. In some audio apps this will play the files in the \"right\" order.  
**audio-reverse**=*false*

## link-viewer
Your web browser.  
**link-viewer**=*xdg-open*

## link-terminal
Open the browser in the same terminal as toot. Only for terminal based browsers.  
**link-terminal**=*false*

# OPEN-CUSTOM
This section is \[open-custom\] in your configuration file

This sections allows you to set up to five custom programs to open URLs with. If the url points to an image, you can set c1-name to img and c1-use to imv. If the program runs in a terminal and you want to run it in the same terminal as tut. Set cX-terminal to true. The name will show up in the UI, so keep it short so all five fits.  
  
c1-name=name  
c1-use=program  
c1-terminal=false  
  
c2-name=name  
c2-use=program  
c2-terminal=false  
  
c3-name=name  
c3-use=program  
c3-terminal=false  
  
c4-name=name  
c4-use=program  
c4-terminal=false  
  
c5-name=name  
c5-use=program  
c5-terminal=false  

# OPEN-PATTERN
This section is \[open-pattern\] in your configuration file

Here you can set your own glob patterns for opening matching URLs in the program you want them to open up in. You could for example open Youtube videos in your video player instead of your default browser.  
  
You must name the keys foo-pattern, foo-use and foo-terminal, where use is the program that will open up the URL. To see the syntax for glob pattern you can follow this URL https://github.com/gobwas/glob\#syntax. foo-terminal is if the program runs in the terminal and should open in the same terminal as tut itself.  
  
Example for youtube.com and youtu.be to open up in mpv instead of the browser.  
  
y1-pattern=\*youtube.com/watch\*  
y1-use=mpv  
y1-terminal=false  
  
y2-pattern=\*youtu.be/\*  
y2-use=mpv  
y2-terminal=false  

# DESKTOP-NOTIFICATION
This section is \[desktop-notification\] in your configuration file

## followers
Notification when someone follows you.  
**followers**=*false*

## favorite
Notification when someone favorites one of your toots.  
**favorite**=*false*

## mention
Notification when someone mentions you.  
**mention**=*false*

## update
Notification when someone edits their toot.  
**update**=*false*

## boost
Notification when someone boosts one of your toots.  
**boost**=*false*

## poll
Notification of poll results.  
**poll**=*false*

## posts
Notification when there is new posts in current timeline.  
**posts**=*false*

# STYLE
This section is \[style\] in your configuration file

All styles can be represented in their HEX value like \#ffffff or with their name, so in this case white. The only special value is \"default\" which equals to transparent, so it will be the same color as your terminal.  
  
You can also use xrdb colors like this xrdb:color1 The program will use colors prefixed with an \* first then look for URxvt or XTerm if it can\'t find any color prefixed with an asterisk. If you don\'t want tut to guess the prefix you can set the prefix yourself. If the xrdb color can\'t be found a preset color will be used. You\'ll have to set theme=none for this to work.  

## xrdb-prefix
The xrdb prefix used for colors in .Xresources.  
**xrdb-prefix**=*guess*

## theme
You can use some themes that comes bundled with tut. Check out the themes available on the URL below. If a theme is named \"nord.ini\" you just write theme=nord  
  
https://github.com/RasmusLindroth/tut/tree/master/config/themes  
  
You can also create a theme file in your config directory e.g. ~/.config/tut/themes/foo.ini and then set theme=foo.  
  
If you want to use your own theme but don\'t want to create a new file, set theme=none and then you can create your own theme below.  
**theme**=*default*

## background
The background color used on most elements.  
**background**=

## text
The text color used on most of the text.  
**text**=

## subtle
The color to display subtle elements or subtle text. Like lines and help text.  
**subtle**=

## warning-text
The color for errors or warnings  
**warning-text**=

## text-special-one
This color is used to display username.  
**text-special-one**=

## text-special-two
This color is used to display username and key hints.  
**text-special-two**=

## top-bar-background
The color of the bar at the top  
**top-bar-background**=

## top-bar-text
The color of the text in the bar at the top.  
**top-bar-text**=

## status-bar-background
The color of the bar at the bottom  
**status-bar-background**=

## status-bar-text
The color of the text in the bar at the bottom.  
**status-bar-text**=

## status-bar-view-background
The color of the bar at the bottom in view mode.  
**status-bar-view-background**=

## status-bar-view-text
The color of the text in the bar at the bottom in view mode.  
**status-bar-view-text**=

## command-text
The color of the text in the command bar at the bottom.  
**command-text**=

## list-selected-background
Background of selected list items.  
**list-selected-background**=

## list-selected-text
The text color of selected list items.  
**list-selected-text**=

## list-selected-inactive-background
The background color of selected list items that are out of focus.  
**list-selected-inactive-background**=

## list-selected-inactive-text
The text color of selected list items that are out of focus.  
**list-selected-inactive-text**=

## controls-text
The main color of the text for key hints  
**controls-text**=

## controls-highlight
The highlight color of for key hints  
**controls-highlight**=

## autocomplete-background
The background color in dropdowns and autocompletions  
**autocomplete-background**=

## autocomplete-text
The text color in dropdowns at autocompletions  
**autocomplete-text**=

## autocomplete-selected-background
The background color for selected value in dropdowns and autocompletions  
**autocomplete-selected-background**=

## autocomplete-selected-text
The text color for selected value in dropdowns and autocompletions  
**autocomplete-selected-text**=

## button-color-one
The background color on selected button and the text color of unselected buttons  
**button-color-one**=

## button-color-two
The text color on selected button and the background color of unselected buttons  
**button-color-two**=

## timeline-name-background
The background on named timelines.  
**timeline-name-background**=

## timeline-name-text
The text color on named timelines  
**timeline-name-text**=

# INPUT
This section is \[input\] in your configuration file

You can edit the keys for tut below.  
  
The syntax is a bit weird, but it works. And I\'ll try to explain it as well as I can.  
  
Example:  
status-favorite=\"[F]avorite\",\"Un[F]avorite\",\'f\',\'F\'  
status-delete=\"[D]elete\",\'d\',\'D\'  
  
status-favorite and status-delete differs because favorite can be in two states, so you will have to add two key hints.  
Most keys will only have on key hint. Look at the default value for reference.  
  
Key hints must be in some of the following formats. Remember the quotation marks.  
\"\" = empty  
\"[D]elete\" = Delete with a highlighted D  
\"Un[F]ollow\" = UnFollow with a highlighted F  
\"[Enter]\" = Enter where everything is highlighted  
\"Yan[K]\" = YanK with a highlighted K  
  
After the hint (or hints) you must set the keys. You can do this in two ways, with single quotation marks or double ones.  
  
The single ones are for single chars like \'a\', \'b\', \'c\' and double marks are for special keys like \"Enter\". Remember that they are case sensitive.  
  
To find the names of special keys you have to go to the following site and look for \"var KeyNames = map[Key]string{\"  
  
https://github.com/gdamore/tcell/blob/master/key.go  

## global-down
Keys for moving down  
**global-down**=*\"\",\'j\',\'J\',\"Down\"*

## global-up
Keys for moving up  
**global-up**=*\"\",\'k\',\'K\',\"Up\"*

## global-enter
To select items  
**global-enter**=*\"\",\"Enter\"*

## global-back
To go back  
**global-back**=*\"[Esc]\",\"Esc\"*

## global-exit
To go back and exit Tut  
**global-exit**=*\"[Q]uit\",\'q\',\'Q\'*

## main-home
Move to the top  
**main-home**=*\"\",\'g\',\"Home\"*

## main-end
Move to the bottom  
**main-end**=*\"\",\'G\',\"End\"*

## main-prev-feed
Go to previous feed  
**main-prev-feed**=*\"\",\'h\',\'H\',\"Left\"*

## main-next-feed
Go to next feed  
**main-next-feed**=*\"\",\'l\',\'L\',\"Right\"*

## main-prev-window
Focus on the previous feed window  
**main-prev-window**=*\"\",\"Backtab\"*

## main-next-window
Focus on the next feed window  
**main-next-window**=*\"\",\"Tab\"*

## main-notification-focus
Focus on the notification list  
**main-notification-focus**=*\"[N]otifications\",\'n\',\'N\'*

## main-compose
Compose a new toot  
**main-compose**=*\"\",\'c\',\'C\'*

## status-avatar
Open avatar  
**status-avatar**=*\"[A]vatar\",\'a\',\'A\'*

## status-boost
Boost a toot  
**status-boost**=*\"[B]oost\",\"Un[B]oost\",\'b\',\'B\'*

## status-edit
Edit a toot  
**status-edit**=*\"[E]dit\",\'e\',\'E\'*

## status-delete
Delete a toot  
**status-delete**=*\"[D]elete\",\'d\',\'D\'*

## status-favorite
Favorite a toot  
**status-favorite**=*\"[F]avorite\",\"Un[F]avorite\",\'f\',\'F\'*

## status-media
Open toots media files  
**status-media**=*\"[M]edia\",\'m\',\'M\'*

## status-links
Open links  
**status-links**=*\"[O]pen\",\'o\',\'O\'*

## status-poll
Open poll  
**status-poll**=*\"[P]oll\",\'p\',\'P\'*

## status-reply
Reply to toot  
**status-reply**=*\"[R]eply\",\'r\',\'R\'*

## status-bookmark
Save/bookmark a toot  
**status-bookmark**=*\"[S]ave\",\"Un[S]ave\",\'s\',\'S\'*

## status-thread
View thread  
**status-thread**=*\"[T]hread\",\'t\',\'T\'*

## status-user
Open user profile  
**status-user**=*\"[U]ser\",\'u\',\'U\'*

## status-view-focus
Open the view mode  
**status-view-focus**=*\"[V]iew\",\'v\',\'V\'*

## status-yank
Yank the url of the toot  
**status-yank**=*\"[Y]ank\",\'y\',\'Y\'*

## status-toggle-cw
Show the content in a content warning  
**status-toggle-cw**=*\"Press [Z] to toggle cw\",\'z\',\'Z\'*

## status-show-filtered
Show the content of a filtered toot  
**status-show-filtered**=*\"Press [Z] to view filtered toot\",\'z\',\'Z\'*

## user-avatar
View avatar  
**user-avatar**=*\"[A]vatar\",\'a\',\'A\'*

## user-block
Block the user  
**user-block**=*\"[B]lock\",\"Un[B]lock\",\'b\',\'B\'*

## user-follow
Follow user  
**user-follow**=*\"[F]ollow\",\"Un[F]ollow\",\'f\',\'F\'*

## user-follow-request-decide
Follow user  
**user-follow-request-decide**=*\"Follow [R]equest\",\"Follow [R]equest\",\'r\',\'R\'*

## user-mute
Mute user  
**user-mute**=*\"[M]ute\",\"Un[M]ute\",\'m\',\'M\'*

## user-links
Open links  
**user-links**=*\"[O]pen\",\'o\',\'O\'*

## user-user
View user profile  
**user-user**=*\"[U]ser\",\'u\',\'U\'*

## user-view-focus
Open view mode  
**user-view-focus**=*\"[V]iew\",\'v\',\'V\'*

## user-yank
Yank the user URL  
**user-yank**=*\"[Y]ank\",\'y\',\'Y\'*

## list-open-feed
Open list  
**list-open-feed**=*\"[O]pen\",\'o\',\'O\'*

## list-user-list
List all users in a list  
**list-user-list**=*\"[U]sers\",\'u\',\'U\'*

## list-user-add
Add user to list  
**list-user-add**=*\"[A]dd\",\'a\',\'A\'*

## list-user-delete
Delete user from list  
**list-user-delete**=*\"[D]elete\",\'d\',\'D\'*

## link-open
Open URL  
**link-open**=*\"[O]pen\",\'o\',\'O\'*

## link-yank
Yank the URL  
**link-yank**=*\"[Y]ank\",\'y\',\'Y\'*

## tag-open-feed
Open tag feed  
**tag-open-feed**=*\"[O]pen\",\'o\',\'O\'*

## tag-follow
Toggle follow on tag  
**tag-follow**=*\"[F]ollow\",\"Un[F]ollow\",\'f\',\'F\'*

## compose-edit-cw
Edit content warning text on new toot  
**compose-edit-cw**=*\"[C]W text\",\'c\',\'C\'*

## compose-edit-text
Edit the text on new toot  
**compose-edit-text**=*\"[E]dit text\",\'e\',\'E\'*

## compose-include-quote
Include a quote when replying  
**compose-include-quote**=*\"[I]nclude quote\",\'i\',\'I\'*

## compose-media-focus
Focus on adding media to toot  
**compose-media-focus**=*\"[M]edia\",\'m\',\'M\'*

## compose-post
Post the new toot  
**compose-post**=*\"[P]ost\",\'p\',\'P\'*

## compose-toggle-content-warning
Toggle content warning on toot  
**compose-toggle-content-warning**=*\"[T]oggle CW\",\'t\',\'T\'*

## compose-visibility
Edit the visibility on new toot  
**compose-visibility**=*\"[V]isibility\",\'v\',\'V\'*

## compose-language
Edit the language of a toot  
**compose-language**=*\"[L]ang\",\'l\',\'L\'*

## compose-poll
Switch to creating a poll  
**compose-poll**=*\"P[O]ll\",\'o\',\'O\'*

## media-delete
Delete media file  
**media-delete**=*\"[D]elete\",\'d\',\'D\'*

## media-edit-desc
Edit the description on media file  
**media-edit-desc**=*\"[E]dit desc\",\'e\',\'E\'*

## media-add
Add a new media file  
**media-add**=*\"[A]dd\",\'a\',\'A\'*

## vote-vote
Vote on poll  
**vote-vote**=*\"[V]ote\",\'v\',\'V\'*

## vote-select
Select item to vote on  
**vote-select**=*\"[Enter] to select\",\' \', \"Enter\"*

## poll-add
Add a new poll option  
**poll-add**=*\"[A]dd\",\'a\',\'A\'*

## poll-edit
Edit a poll option  
**poll-edit**=*\"[E]dit\",\'e\',\'E\'*

## poll-delete
Delete a poll option  
**poll-delete**=*\"[D]elete\",\'d\',\'D\'*

## poll-multi-toggle
Toggle voting on multiple options  
**poll-multi-toggle**=*\"Toggle [M]ultiple\",\'m\',\'M\'*

## poll-expiration
Change the expiration of poll  
**poll-expiration**=*\"E[X]pires\",\'x\',\'X\'*

## preference-name
Change display name  
**preference-name**=*\"[N]ame\",\'n\',\'N\'*

## preference-visibility
Change default visibility of toots  
**preference-visibility**=*\"[V]isibility\",\'v\',\'V\'*

## preference-bio
Change bio in profile  
**preference-bio**=*\"[B]io\",\'b\',\'B\'*

## preference-save
Save your preferences  
**preference-save**=*\"[S]ave\",\'s\',\'S\'*

## preference-fields
Edit profile fields  
**preference-fields**=*\"[F]ields\",\'f\',\'F\'*

## preference-fields-add
Add new field  
**preference-fields-add**=*\"[A]dd\",\'a\',\'A\'*

## preference-fields-edit
Edit current field  
**preference-fields-edit**=*\"[E]dit\",\'e\',\'E\'*

## preference-fields-delete
Delete current field  
**preference-fields-delete**=*\"[D]elete\",\'d\',\'D\'*

# SEE ALSO
    tut(1) - flags and commands
    tut(7) - commands and keys inside of tut
