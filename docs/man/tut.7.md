% tut(7) tut 1.0.34
% Rasmus Lindroth
% 2023-01-01

# NAME
tut - keys and commands inside of tut(1)

# DESCRIPTION
This page contains information for some of the keys and all the commands you can use in tut(1).

To change the keys look at tut(5) under the *INPUT* section. 

# KEYS
## Keys without description in tut
**c** = Compose a new toot  
**j** or **Down arrow** = Navigate down in feed list or toot  
**k** or **Up arrow** = Navigate up in feed list or toot  
**h** or **Left arrow** = Cycle back in open timelines  
**l** or **Right arrow** = Cycle forward in open timelines  
**g** or **Home** = Go to top in feed list or toot  
**G** or **End** = Go to bottom in feed list or toot  
**?** = View help  
**q** = Go back or quit  
**Esc** = Go back

## Explanation of the non obvious keys when viewing a toot
**v** = view. In this mode you can scroll throught the text of the toot if it doesn\'t fit the screen  
**o** = open. Gives you a list of all URLs in the toot. Opens them in your default browser, if it\'s an user or tag they will be opened in tut  
**m** = media. Opens the media with xdg-open

# Commands
**:quit**
: Exit tut

**:q**
: Shorter form of former command

**:timeline** *home|local|federated|direct|notifications|mentions|favorited|special-all|special-boosts|special-replies*
: Open selected timeline

**:tl** *h|l|f|d|n|m|fav|sa|sb|sr*
: Shorter form of former command

**:blocking**
: Lists users that you have blocked

**:boosts**
: Lists users that have boosted the toot

**:bookmarks**
: List all your bookmarks

**:clear-notifications**
: Remove all of your notifications

**:close-window**
: Closes the current window, including all the timelines in said window

**:compose**
: Compose a new toot

**:edit**
: Edit one of your toots

**:favorited**
: Lists toots  you\'ve favorited

**:favorites**
: Lists users that favorited the toot

**:follow-tag** *\<tag\>*
: Follow a hashtag named \<tag\>

**:followers**
: List of people the account are following. It only works on profiles

**:following**
: List of people follwing the account. It only works on profiles

**:help**
: Show help for how to use tut

**:h**
: Shorter form of former command

**:history**
: Show edits of a toot

**:lists**
: Show a list of your lists

**:list-placement** *top|right|bottom|left*
: Place the list in choosen placement

**:list-split** *row|column*
: Split the timelines in window by row or column

**:move-window** *left|right|up|down|home|end*
: Moves the window in choosen direction

**:mv** *l|r|u|d|h|e*
: Shorter form of former command

**:muting**
: Lists users that you\'ve muted

**:newer**
: Force load newer toots in current timeline

**:preferences**
: Update your profile and some other settings

**:profile**
: Go to your profile

**:proportions** *[int] [int]*
: Sets the proportions of the windows and the content. The first integer is your windows and the other for content, e.g. :proportions 1 3

**:refetch**
: Refetches the current item that you\'re viewing. Can be used to update poll results.

**:saved**
: Alias for bookmarks

**:stick-to-top**
: Toggle the stick-to-top setting that always shows the latest toot in all timelines

**:tag** *\<tag\>*
: Shows toots tagged with \<tag\>, e.g. :tag linux. You can input multiple tags if you want to show them in the same timeline

**:tags**
: List of tags that you\'re following

**:unfollow-tag** *\<tag\>*
: Unfollow the hashtag named \<tag\>, e.g. :unfollow-tag tut

**:user** *\<username\>*
: Search for users named \<username\>, e.g. :user rasmus. To narrow a search include the instance like this :user rasmus@mastodon.acc.sunet.se

**:window** *\<int\>*
: Switch window by index (zero indexed) e.g. :window 0 for the left/top window

# SEE ALSO
    tut(1) - flags and commands
    tut(5) - configuration format
