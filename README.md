# TTFM Bot

A bot to bring some fun and utils on [turntable.fm](https://turntable.fm) rooms.

It's based on [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi), a Golang library to build bots for Turntable.fm.

## Installation and setup

**IMPORTANT: Requires Go >= 1.18beta1**

- Install the package with `go install https://github.com/andreapavoni/ttfm_bot/`
- Set environment variables (see below) on the host you want to run the bot
- Run with `$GOPATH/bin/ttfm_bot`

### Configuration settings

#### Mandatory

- `TTFM_API_AUTH`: the API key to connect to Turntables.fm
- `TTFM_API_USER_ID`: User ID for the user
- `TTFM_API_ROOM_ID`: ID of the room to join

#### Optional (with defaults)

- `TTFM_ADMINS` (default: empty): a list of comma-separated usernames
- `TTFM_AUTO_SNAG` (default: false): wether the bot should snag every song played by others
- `TTFM_AUTO_BOP` (default: true): wether the bot should bop every song played by others
- `TTFM_AUTO_DJ` (default: false): if none is playing, then the bot will automatically jump on the stage
- `TTFM_AUTO_QUEUE` (default: false): joins queues managed by others
- `TTFM_AUTO_QUEUE_MSG` (default: empty): react when mentioned to join the stage (works with AutoDJ)
- `TTFM_AUTO_SHOW_SONG_STATS` (default: ): communicate to the room the stats of the last song played
- `TTFM_AUTO_WELCOME` (default: false): welcomes every user that joins the room
- `TTFM_MOD_QUEUE` (default: false): enables queueing when the room is crowded with aspiring DJs
- `TTFM_MOD_SONGS_MAX_DURATION` (default: 0): duration limit of the song in minutes (0 means disabled)
- `TTFM_DEFAULT_PLAYLIST` (default: "default"): which playlist should use the bot (for snag or DJ)
- `TTFM_SET_BOT` (default: false): tells the server that this is a bot

## Commands

**NOTES:**

- Each command can be either issued on the chat room or by private message.
- A command might require a certain user role to execute a command.
- commands which accept `on` or `off` to enable/disable some feature, can be called without arguments to get the current status

### Users

Users are the lowest role, basically are listeners or want to DJ

- `!props` let the current DJ know you're appreciating the song
- `!help` shows the list of commands available for the role of the user that issued the command (_Coming soon_)
- `!q on|off` shows the current queue status
- `!qadd` adds user into queue
- `!qrm` removes user from queue
- `!r reaction` posts a giphy `reaction` in the chat room (_Coming soon_)

### DJs

- `!escortme` will escort the DJ off the stage after they played the last song. Requires the bot to be a moderator in the current room.

### Bot admins

Admins are users which were previously configured on the bot to run commands on it.

- `!dj` tells the bot to jump on the stage and starts playing songs, or jump off if it's already djing
- `!autodj on|off` enables/disables autodj mode
- `!snag` tells the bot to snag the current playing song
- `!autosnag on|off` enables/disables automatic snag
- `!bop` tells the bot to bop for the current playing song
- `!autobop on|off` enables/disables automatic bop
- `!fan username` and `!unfan username` respectively fan/unfan the specified `username`
- `!padd name` creates a new playlist
- `!pdel name` deletes a playlist
- `!pls` lists available playlists
- `!prm` removes the current playing song from the current playlist
- `!pc name` switch playlist

### Moderators

Bot can obey to moderators commands, however it depends by the kind of command issued and by the bot's role in the current room.

- `!skip` tells the bot to skip the current playing song
- `!escort username` tells the bot to escort the specified `username` off the stage
- `!boot username` tells the bot to kick the specified `username` off the room (_Coming soon_)

## Credits

- [turntable.fm](https://turntable.fm) for the awesome platform
- [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi) because without it I should have to hack a lot more to get here
- [nugget/cowgod](https://github.com/nugget/cowgod) another Golang bot for Turntable.fm, I did a look at it to learn

### Rooms

- [Disco Clubbing](https://turntable.fm/disco_clubbing) (here you can find this bot running)
- [Aunt Jackie](https://turntable.fm/aunt_jackie)
- [I ❤️ The 80's](https://turntable.fm/i_the_80s)

## TODO (as of march 10th 2022)

- [x] DJ queues
- [ ] use a struct for command input and output (ex. CommandInput and CommandOutput)
- [ ] room handling (join, list faves, leave)
- [ ] set song max duration
- [ ] set max songs per dj
- [ ] cmd for room stats
- [ ] use command-output struct to determine if message should be sent privately, in room, not at all, or in "/me" form
- [ ] group all string messages into a struct/file
- [ ] add docs for functions
- [ ] shell for commands
