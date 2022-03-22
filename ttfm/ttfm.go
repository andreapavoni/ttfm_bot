package ttfm

import (
	"errors"
	"fmt"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Bot struct {
	Config    *Config
	Playlists *collections.SmartList[string]
	Room      *Room
	Queue     *Queue
	api       *ttapi.Bot
	admins    *collections.SmartList[string]
	playlist  *Playlist
	commands  *collections.SmartMap[*Command]
	Reactions *Reactions
	Brain     *Brain
	Actions   *Actions
	Rooms     *Rooms
}

// BOOT
// New bot instance
func New() *Bot {
	SetupLogging()
	brain := NewBrain("./db")
	cfg := NewConfig(brain)
	reactions := NewReactions(brain, "reactions")

	b := Bot{
		Config:    cfg,
		Room:      NewRoom(),
		Queue:     NewQueue(),
		api:       ttapi.NewBot(cfg.ApiAuth, cfg.UserId, cfg.RoomId),
		admins:    collections.NewSmartListFromSlice(cfg.Admins),
		playlist:  NewPlaylist(cfg.CurrentPlaylist),
		Playlists: collections.NewSmartList[string](),
		commands:  collections.NewSmartMap[*Command](),
		Reactions: reactions,
		Brain:     brain,
	}
	b.Actions = &Actions{&b}
	b.Rooms = NewRooms(&b)

	// Commands
	b.api.OnSpeak(func(e ttapi.SpeakEvt) {
		handleCommand(&b, &MessageInput{UserId: e.UserID, Text: e.Text, Source: MessageTypeRoom})
	})
	b.api.OnPmmed(func(e ttapi.PmmedEvt) {
		handleCommand(&b, &MessageInput{UserId: e.SenderID, Text: e.Text, Source: MessageTypePm})
	})

	// Room events
	b.api.OnReady(func() { onReady(&b) })
	b.api.OnRoomChanged(func(e ttapi.RoomInfoRes) { onRoomChanged(&b, e) })
	b.api.OnRegistered(func(e ttapi.RegisteredEvt) { onRegistered(&b, e) })
	b.api.OnDeregistered(func(e ttapi.DeregisteredEvt) { onDeregistered(&b, e) })
	b.api.OnUpdateVotes(func(e ttapi.UpdateVotesEvt) { onUpdateVotes(&b, e) })
	b.api.OnSnagged(func(e ttapi.SnaggedEvt) { onSnagged(&b, e) })

	// DJing
	b.api.OnRemDJ(func(e ttapi.RemDJEvt) { onRemDj(&b, e) })
	b.api.OnAddDJ(func(e ttapi.AddDJEvt) { onAddDj(&b, e) })
	b.api.OnNewSong(func(e ttapi.NewSongEvt) { onNewSong(&b, e) })

	return &b
}

// AddCommand with given alias
func (b *Bot) AddCommand(alias string, cmd *Command) {
	b.commands.Set(alias, cmd)
}

// Start the bot
func (b *Bot) Start() {
	b.api.Start()
}

// COMMANDS
// GetCommand if exists
func (b *Bot) GetCommand(name string) (*Command, error) {
	if cmd, ok := b.commands.Get(name); ok {
		return cmd, nil
	}
	return nil, errors.New("command not found")
}

// ListCommands available
func (b *Bot) ListCommands() (commands []string) {
	for _, k := range b.commands.Keys() {
		commands = append(commands, k)
	}
	return commands
}

// GetRoomInfo from API
func (b *Bot) GetRoomInfo() (ttapi.RoomInfoRes, error) {
	return b.api.RoomInfo()
}

// QUEUE
// ModQueue enabled/disabled
func (b *Bot) ModQueue(status bool) bool {
	b.Config.ModQueue = status
	b.Queue.Empty()
	b.Config.Save()

	return status
}

// AddDjEscorting the dj will be escorted after the current song is played
func (b *Bot) AddDjEscorting(userId string) error {
	if !b.UserIsDj(userId) {
		if userId == b.Config.UserId {
			return errors.New("I'm not on stage!")
		}
		return errors.New("You aren't DJing!")
	}

	if !b.Room.escorting.HasElement(userId) {
		b.Room.escorting.Push(userId)
	}

	return nil
}

// RemoveDjEscorting if dj doesn't want to be escorted anymore
func (b *Bot) RemoveDjEscorting(userId string) error {
	return b.Room.escorting.Remove(userId)
}

// EscortDj
func (b *Bot) EscortDj(userId string) error {
	if !b.UserIsDj(userId) {
		if userId == b.Config.UserId {
			return errors.New("I'm not on stage!")
		}
		return errors.New("user is not on stage")
	}

	return b.api.RemDj(userId)
}

// SONG
// Bop
func (b *Bot) Bop() {
	if b.Room.Song.DjId != b.Config.UserId {
		b.api.Bop()
	}
}

// Downvote
func (b *Bot) Downvote() {
	if b.Room.Song.DjId != b.Config.UserId {
		b.api.VoteDown()
	}
}

// SkipSong current song (must be moderator to skip others songs)
func (b *Bot) SkipSong() {
	b.api.Skip()
}

// AutoBop each song
func (b *Bot) AutoBop(status bool) bool {
	b.Config.AutoBop = status
	b.Config.Save()
	return status
}

// ShowSongStats when song is finished
func (b *Bot) ShowSongStats() (header, data string) {
	song := b.Room.Song
	header = fmt.Sprintf("Stats for `%s` by `%s` played by @%s:", song.Title, song.Artist, song.DjName)
	data = fmt.Sprintf("👍 %d | 👎 %d | ❤️ %d", song.up, song.down, song.snag)

	return header, data
}

// DJ
// Dj
func (b *Bot) Dj() {
	if !b.UserIsDj(b.Config.UserId) {
		b.api.AddDj()
	}
}

// AutoDj enabled/disabled. If bot is djing, it will be escorted when song is finished
func (b *Bot) AutoDj(status bool) bool {
	b.Config.AutoDj = status
	b.Config.Save()

	if status {
		return status
	}

	if !status && b.UserIsDj(b.Config.UserId) {
		if b.UserIsCurrentDj(b.Config.UserId) {
			b.AddDjEscorting(b.Config.UserId)
			return status
		}
		b.api.RemDj("")
	}

	return status
}

// PLAYLISTS
// Snag
func (b *Bot) Snag() error {
	if b.Room.Song.DjId == b.Config.UserId {
		return errors.New("I'm the current DJ and I already have this song in my playlist...")
	}

	playlist, err := b.api.PlaylistAll(b.Config.CurrentPlaylist)
	if err != nil {
		return err
	}

	b.api.Snag()
	if err = b.api.PlaylistAdd(b.Room.Song.Id, b.Config.CurrentPlaylist, len(playlist.List)); err != nil {
		return nil
	}

	b.playlist.AddSong(&SongItem{
		id:     b.Room.Song.Id,
		title:  b.Room.Song.Title,
		artist: b.Room.Song.Artist,
		length: b.Room.Song.Length,
	})

	return nil
}

// AutoSnag each song
func (b *Bot) AutoSnag(status bool) bool {
	b.Config.AutoSnag = status
	b.Config.Save()
	return status
}

// LoadPlaylist
func (b *Bot) LoadPlaylist(playlistName string) error {
	playlist, err := b.api.PlaylistAll(b.Config.CurrentPlaylist)

	if err != nil {
		return err
	}

	for _, s := range playlist.List {
		b.playlist.AddSong(&SongItem{
			id:     s.ID,
			title:  s.Metadata.Song,
			artist: s.Metadata.Artist,
			length: s.Metadata.Length,
		})
	}
	return nil
}

// LoadPlaylists from API and cache them in memory
func (b *Bot) LoadPlaylists() error {
	playlists, err := b.api.PlaylistListAll()
	if err != nil {
		return err
	}

	for _, pl := range playlists.List {
		b.Playlists.Push(pl.Name)
	}
	return nil
}

// AddPlaylist
func (b *Bot) AddPlaylist(playlistName string) error {
	if !b.Playlists.HasElement(playlistName) {
		if err := b.api.PlaylistCreate(playlistName); err != nil {
			return err
		}
		b.Playlists.Push(playlistName)
		return nil
	}

	return errors.New("Playlist not found")
}

// RemovePlaylist
func (b *Bot) RemovePlaylist(playlistName string) error {
	if b.Playlists.HasElement(playlistName) {
		if err := b.api.PlaylistDelete(playlistName); err != nil {
			return err
		}
		b.Playlists.Remove(playlistName)
		return nil
	}

	return errors.New("Playlist not found")
}

// SwitchPlaylist
func (b *Bot) SwitchPlaylist(playlistName string) error {
	if b.Playlists.HasElement(playlistName) {
		if err := b.api.PlaylistSwitch(playlistName); err != nil {
			return err
		}
		b.Config.CurrentPlaylist = playlistName
		b.Config.Save()
		return b.LoadPlaylist(playlistName)
	}

	return errors.New("Playlist not found")
}

// PushSongBottomPlaylist
func (b *Bot) PushSongBottomPlaylist() error {
	if err := b.api.PlaylistReorder(b.Config.CurrentPlaylist, 0, b.playlist.songs.Size()-1); err == nil {
		currentSong, _ := b.playlist.songs.Shift()
		b.playlist.AddSong(&currentSong)
		return nil
	} else {
		return err
	}
}

// RemoveSongFromPlaylist
func (b *Bot) RemoveSongFromPlaylist(songId string) error {
	song, idx, err := b.playlist.GetSongById(songId)

	if err != nil {
		return err
	}

	if err := b.api.PlaylistRemove(b.Config.CurrentPlaylist, idx); err != nil {
		return err
	}

	b.playlist.RemoveSong(song)
	return nil
}

// MESSAGING
// PrivateMessage
func (b *Bot) PrivateMessage(userId, msg string) {
	b.api.PM(userId, msg)
}

// RoomMessage
func (b *Bot) RoomMessage(msg string) {
	b.api.Speak(msg)
}

// USERS & AUTHORIZATION
// BootUser from room
func (b *Bot) BootUser(userId, reason string) error {
	return b.api.BootUser(userId, reason)
}

// Fan another user
func (b *Bot) Fan(userId string) error {
	return b.api.BecomeFan(userId)
}

// Unfan
func (b *Bot) Unfan(userId string) error {
	return b.api.RemoveFan(userId)
}

// UserFromId
func (b *Bot) UserFromId(userId string) (*User, error) {
	if userName, ok := b.Room.Users.Get(userId); ok {
		return &User{Id: userId, Name: userName}, nil
	}
	return &User{}, errors.New("User with ID " + userId + " wasn't found")
}

// UserFromName
func (b *Bot) UserFromName(userName string) (*User, error) {
	if id, err := b.api.GetUserID(userName); err == nil {
		return &User{Id: id, Name: userName}, nil
	} else {
		return &User{}, err
	}
}

// UserIsAdmin
func (b *Bot) UserIsAdmin(user *User) bool {
	return b.admins.HasElement(user.Name)
}

// UserIsDj
func (b *Bot) UserIsDj(userId string) bool {
	return b.Room.Djs.HasElement(userId)
}

// UserIsCurrentDj
func (b *Bot) UserIsCurrentDj(userId string) bool {
	return b.Room.Song.DjId == userId
}

// UserIsModerator
func (b *Bot) UserIsModerator(userId string) bool {
	return b.Room.Moderators.HasElement(userId)
}
