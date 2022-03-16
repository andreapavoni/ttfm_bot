# TTFM Bot

A bot to bring some fun and utils on [turntable.fm](https://turntable.fm) rooms.

It's based on [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi), a Golang library to build bots for Turntable.fm.

## Installation and setup

**IMPORTANT: Requires Go >= 1.18beta1 because it uses generics!**

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
- `!help <cmd>` shows description of a command. Without `cmd` shows the list of commands available for the role of the user that issued the command
- `!q [on|off]` enables/disables queue. Without `on` or `off` replies with current status
- `!qadd` adds user into queue
- `!qrm` removes user from queue
- `!r <reaction>` shows a funny gif reaction. Without `reaction` shows available ones.

### DJs

- `!escortme` will escort the DJ off the stage after they played the last song. Requires the bot to be a moderator in the current room.

### Bot admins

Admins are users which were previously configured on the bot to run commands on it.

- `!dj` tells the bot to jump on the stage and starts playing songs, or jump off if it's already djing
- `!autodj [on|off]` enables/disables autodj mode. Without `on` or `off` replies with current status
- `!snag` tells the bot to snag the current playing song
- `!autosnag [on|off]` enables/disables automatic snag. Without `on` or `off` replies with current status
- `!bop` tells the bot to bop for the current playing song
- `!autobop [on|off]` enables/disables automatic bop. Without args replies with current status
- `!fan <user_name>` and `!unfan <user_name>` respectively fan/unfan the specified `user_name`
- `!padd <playlist_name>` creates a new playlist
- `!pdel <playlist_name>` deletes a playlist
- `!pls` lists available playlists
- `!prm` removes the current playing song from the current playlist
- `!pc <playlist_name>` switch to `playlist_name` playlist
- `!say <something>` say something in the room
- `!cfg <config_key> [<config_value>]` sets config key and value. Without `config_value`, it replies with current configuration for `<config_key>`

### Moderators

Bot can obey to moderators commands, however it depends by the kind of command issued and by the bot's role in the current room.

- `!skip` tells the bot to skip the current playing song
- `!escort <user_name>` tells the bot to escort the specified `user_name` off the stage
- `!boot <user_name>` tells the bot to kick the specified `user_name` off the room

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
- [x] use a struct for command input and output (ex. CommandInput and CommandOutput)
  - [x] input
  - [x] output
  - [x] use Enum types for source/reply fields
- [x] add say cmd
- [x] song max duration
  - [x] enforce on new song
  - [x] update value at runtime
- [x] add help cmd
- [ ] multiple config source (env, toml) and load before init bot
- [x] log to file
- [x] set max songs per dj at runtime
- [ ] cmd for room stats?
- [ ] group all string messages into a struct/file?
- [ ] add docs for functions
- [ ] shell for commands
