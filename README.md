# TTFM Bot

A bot to bring some fun and utils on [turntable.fm](https://turntable.fm) rooms.

It's based on [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi), a Golang library to build bots for Turntable.fm.

## Installation and setup

**IMPORTANT: Requires Go >= 1.18beta1**

- Install the package with `go get https://github.com/andreapavoni/ttfm_bot/`
- Set environment varables (see below)
- Run with `ttfm_bot`

### Configuration

`TTFM_API_AUTH`: the API key to connect to Turntables.fm
`TTFM_API_USER_ID`: User ID for the user
`TTFM_API_ROOM_ID`: lorem ipsum dolor sit

`TTFM_ADMINS` (default: empty): a list of comma-separated usernames
`TTFM_AUTO_SNAG` (default: false): wether the bot should snag every song played by others
`TTFM_AUTO_BOP` (default: true): wether the bot should bop every song played by others
`TTFM_AUTO_DJ` (default: false): if none is playing, then the bot will automatically jump on the stage
`TTFM_AUTO_QUEUE` (default: false): joins queues managed by others
`TTFM_AUTO_QUEUE_MSG` (default: empty): react when mentioned to join the stage (works with AutoDJ)
`TTFM_AUTO_SHOW_SONG_STATS` (default: ): communicate to the room the stats of the last song played
`TTFM_AUTO_WELCOME` (default: false): welcomes every user that joins the room
`TTFM_MOD_QUEUE` (default: false): enables queueing when the room is crowded with aspiring DJs
`TTFM_MOD_SONGS_MAX_DURATION` (default: 0): duration limit of the song in minutes (0 means disabled)
`TTFM_DEFAULT_PLAYLIST` (default: "default"): which playlist should use the bot (for snag or DJ)
`TTFM_SET_BOT` (default: false): tells the server that this is a bot

## Commands

**NOTES:**

- Each command can be either issued on the chat room or by private message.
- A command might require a certain user role to execute a command.

### Users

Users are the lowest role, basically are listeners or want to DJ

- `!help` shows the list of commands available for the role of the user that issued the command (_Coming soon_)
- `!q+` and `!q-`respectively adds/removes user from the DJ queue (_Coming soon_)
- `!r reaction` posts a giphy `reaction` in the chat room (_Coming soon_)

### DJs

- `!escortme` will escort the DJ off the stage after they played the last song. Requires the bot to be a moderator in the current room.

### Bot admins

Admins are users which were previously configured on the bot to run commands on it.

- `!autodj` if autodj is enabled, it tells the bot to jump on the stage and starts playing songs
- `!autodj+` and - `!autodj-` respectively enables/disables autodj mode
- `!snag` tells the bot to snag the current playing song
- `!autosnag+` and - `!autosnag-` respectively enables/disables automatic snag
- `!bop` tells the bot to bop for the current playing song
- `!autobop+` and - `!autobop-` respectively enables/disables automatic bop
- `!fan username` and - `!unfan username` respectively fan/unfan the specified `username`

### Moderators

Bot can obey to moderators commands, however it depends by the kind of command issued and by the bot's role in the current room.

- `!skip` bot to skip the current playing song (_Coming soon_)
- `!escort username` tells the bot to escort the specified `username` off the stage (_Coming soon_)
- `!kick username` tells the bot to kick the specified `username` off the room (_Coming soon_)

## Credits

- [turntable.fm](https://turntable.fm) for the awesome platform
- [alaingilbert/ttapi](https://github.com/alaingilbert/ttapi) because without it I should have to hack a lot more to get here
- [nugget/cowgod](https://github.com/nugget/cowgod) another Golang bot for Turntable.fm, I did a look at it to learn

### Rooms

- [Disco Clubbing](https://turntable.fm/disco_clubbing) (here you can find this bot running)
- [Aunt Jackie](https://turntable.fm/aunt_jackie)
- [I ❤️ The 80's](https://turntable.fm/i_the_80s)
