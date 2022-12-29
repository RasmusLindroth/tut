% tut(1) tut 1.0.30
% Rasmus Lindroth
% 2022-12-29

# NAME
tut - a Mastodon TUI

# SYNOPSIS
**tut** [command] [options...]

# DESCRIPTION
A TUI for Mastodon with vim inspired keys. The program has most of the features you can find in the web client.
To see keys and commands you can use inside of tut check tut(7).

# OPTIONS

**-h**, **\--help**
: Show help message

**-v**, **\--version**
: Show the version number

**-n**, **\--new-user**
: Add one more user to tut

**-c**,  **\--config** \<path\>
: Load config.ini from *\<path\>*

**-d**,  **\--config-dir** \<path\>
: Load all config from *\<path\>*

**-u**,  **\--user** \<name\>
: Login directly to user named *\<name\>*.
: If two users are named the same, use full name like *tut@fosstodon.org*

# COMMANDS

**no command**
: Runs the TUI

**example-config**
: Generates the default configuration file in the current directory and names it ./config.example.ini

# CONFIGURATION
Tut is configurable, so you can change things like the colors, the default timeline, what image viewer to use and some more. Check out tut(5) or the configuration file to see all the options.

You find it in *$XDG_CONFIG_HOME/tut/config.ini* on Linux which usually equals to *~/.config/tut/config.ini*.
If you don't run Linux it will use the path of the Go funcdtion os.UserConfigDir().
But if you move the tut folder to *XDG_CONFIG_HOME/tut/* and have set the environment variable *XDG_CONFIG_HOME*
it will look there instead of the standard place.

You can generate an example configuration file with *tut example-config*. It will be updated with potential new features.

# SEE ALSO
    tut(5) - configuration format
    tut(7) - commands and keys inside of tut
