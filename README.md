# TTFM Bot

A bot to bring some fun and utils on [turntable.fm](https://turntable.fm) rooms.

It's based on [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi), a Golang library to build bots for Turntable.fm.

## Features

Features are inspired by [chillybot](https://github.com/jaycammarano/chillybot), with some differences here and there.

- [x] Easy to install and use: download the binary for your platform and run it!
- [x] Queue moderation
  - [x] Auto-enabled when stage is full. Auto-disabled when stage has free slots and queue is empty.
  - [ ] Max songs per DJ: how many songs before forwarding queue (eg: default 1 song)
- [ ] Max songs per DJ when queue is disabled (0 means unlimited)
  - [ ] Rest time for DJ who has reached max songs limit and has been escorted (eg: default 5 mins)
- [x] Enforce song length limit (10 minutes by default)
- [x] Show song stats when song finishes (only when enabled, and bot is moderator)
- [x] Room greeting (only when enabled, and bot is moderator)
  - [ ] Greeting message configurable
- [x] Playlists
  - [x] Add/Remove songs to your bot's playlists
  - [x] Create/Remove playlists
  - [x] Switch to another playlist
- [x] Auto SNAG every song into current playlist
- [x] Auto BOP
- [x] Favorite rooms
  - [x] Manage the list of favorite rooms
  - [x] Join another room
- [x] Auto DJ when there are `1` (default) or less djs
- [x] GIF reactions
  - [x] Use as many reactions you want
  - [x] Add new reactions or new GIFs to an existing one
  - [ ] Remove reactions or GIFs
- [x] Escorting: a DJ can ask to be escorted immediately or after the current song has been played
- [x] Logging
  - [x] More details for some events/actions
  - [x] Logs rotation
- [ ] Configure a path where to write bot's saved data and logs
- [ ] afk limit
- [ ] afk audience limit(separate from afk limit, both can be toggled on and off)
- [x] DJ stats (shown when dj goes off the stage)
- [x] Print room rules (by command and/or when user joins)
- [x] Room music current theme (default "free play")
- [ ] Bot info (by command: print version, uptime, ...)
- [x] custom prefix for commands (actual is `!`)
- [x] kill switch command (kills/disconnects bot, useful when it turns unresponsive/misbehaved)

## Installation and setup

- Ensure to put the executable where you can write files (bot's saved data and logs)
- Download [latest release](https://github.com/andreapavoni/ttfm_bot/releases/latest) binary for your platform (recommended) or source code (you'll need Go v1.18 or higher to build this project)
- Set environment variables (see below) on the host you want to run the bot
- Run with `./ttfm_bot` (better if you run it from `screen` or `tmux` session)

### Configuration settings

These environment variables are required to make the bot work

- `TTFM_API_AUTH`: the API key to connect to Turntables.fm
- `TTFM_API_USER_ID`: User ID for the user
- `TTFM_API_ROOM_ID`: ID of the room to join
- `TTFM_MAIN_ADMIN_ID`: the ID of the user that will be the first and main admin

## Commands

Each command can be either issued on the chat room or by private message. Bot _might_ reply (or not) in the proper place.

A command might require a certain user (and sometimes bot's) role to execute a command.

### Command prefix customization

By default, the bot has `!` as command prefix (eg: `!command`). It's possible to customize it with `!cfg cmdprefix <new_prefix>`. You can use whatever ASCII character or string.

**NOTE:** when you change prefix, the new one will be the effective prefix to use to run commands. 

### Users

Users are the lowest role, basically anyone who isn't bot's admin or room moderator

- `!props` let the current DJ know you're appreciating the song
- `!help <cmd>` shows description of a command. Without `cmd` shows the list of commands available for the role of the user that issued the command
- `!q [add|rm]` add or remove yourself from the queue. Without args shows the current line in queue
- `!qadd` adds user into queue
- `!qrm` removes user from queue
- `!r <reaction>` shows a funny gif reaction. Without `reaction` shows available ones
- `!r add <reaction> <url>` shows a funny gif reaction. Without `reaction` shows available ones

### DJs

- `!escortme` will escort the DJ off the stage after they played the last song. Requires the bot to be a moderator in the current room.

### Bot admins

Admins are users which can run commands on bot

- `!dj` tells the bot to jump on the stage and starts playing songs, or jump off if it's already djing
- `!snag` tells the bot to snag the current playing song
- `!bop` tells the bot to bop for the current playing song
- `!fan <user_name>` and `!unfan <user_name>` respectively fan/unfan the specified `user_name`
- `!p [[add | rm | switch] <playlist_name> | list | rmsong]` handles playlists
- `!say <something>` say something in the room
- `!cfg <config_key> [<config_value>]` sets config key and value. Without `config_value`, it replies with current configuration for `<config_key>`
- `!room [list | <sub_command> <room_slug>]` handles favorite rooms
- `!die` kills the bot (useful if/when it becomes unresponsive due to some bug)
- `!admin [add | remove <user_name>]`. Without args shows the current admins

These ones require the bot to be moderator of the room to be executed

- `!skip` tells the bot to skip the current playing song
- `!escort <user_name>` tells the bot to escort the specified `user_name` off the stage
- `!boot <user_name>` tells the bot to kick the specified `user_name` off the room

### Upgrading to newer versions

Except if it has been specified differently in a new release notes, it should suffice to:

- Download the new release
- Stop the currently running bot
- Replace old binary with the one from the new release
- Restart the bot



## Credits

- [turntable.fm](https://turntable.fm) for the awesome platform
- [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi) because without it I should have to hack a lot more to get here
- [nugget/cowgod](https://github.com/nugget/cowgod) another Golang bot for Turntable.fm, I did a look at it to learn
- [jaycammarano/chillybot](https://github.com/jaycammarano/chillybot) a JS bot, maybe the most used, I got some inspiration for features from there
- [I ❤️ The 80's](https://turntable.fm/i_the_80s) the best room I've found, full of friendly people

## Rooms

An instance of this bot is called `Mrs.Beats` and can be found in one of these channells. Author username is `pavonz` both on turntable.fm and on Discord servers.

- [Disco Clubbing](https://turntable.fm/disco_clubbing) (its main room, here it's moderator)
- [I ❤️ The 80's](https://turntable.fm/i_the_80s)
- [Aunt Jackie](https://turntable.fm/aunt_jackie)
