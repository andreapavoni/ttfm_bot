package ttfm

import (
	"errors"
	"fmt"
	"time"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Bot struct {
	api       *ttapi.Bot
	queue     *collections.SmartList[string]
	admins    *collections.SmartList[string]
	config    *Config
	room      *Room
	playlist  *Playlist
	playlists *collections.SmartList[string]
	escorting *collections.SmartList[string]
}

type Config struct {
	ApiAuth             string
	UserId              string
	RoomId              string
	Admins              []string
	AutoSnag            bool
	AutoBop             bool
	AutoDj              bool
	AutoQueue           bool
	AutoQueueMsg        string
	AutoShowSongStats   bool
	ModAutoWelcome      bool
	ModQueue            bool
	ModSongsMaxDuration int64
	ModSongsMaxPerDj    int64
	CurrentPlaylist     string
	SetBot              bool
}

// BOOT
func Start(c Config) {
	bot := Bot{
		api:    ttapi.NewBot(c.ApiAuth, c.UserId, c.RoomId),
		queue:  collections.NewSmartList[string](),
		admins: collections.NewSmartListFromSlice(c.Admins),
		config: &c,
		playlist: &Playlist{
			Name:  c.CurrentPlaylist,
			songs: collections.NewSmartList[SongItem](),
		},
		playlists: collections.NewSmartList[string](),
		room: &Room{
			users:      collections.NewSmartMap[string](),
			moderators: collections.NewSmartList[string](),
			djs:        collections.NewSmartList[string](),
			song:       &Song{},
		},
		escorting: collections.NewSmartList[string](),
	}

	// Commands
	bot.api.OnSpeak(func(e ttapi.SpeakEvt) { handleCommandSpeak(&bot, e) })
	bot.api.OnPmmed(func(e ttapi.PmmedEvt) { handleCommandPm(&bot, e) })

	// Room events
	bot.api.OnReady(func() { onReady(&bot) })
	bot.api.OnRoomChanged(func(e ttapi.RoomInfoRes) { onRoomChanged(&bot, e) })
	bot.api.OnRegistered(func(e ttapi.RegisteredEvt) { onRegistered(&bot, e) })
	bot.api.OnDeregistered(func(e ttapi.DeregisteredEvt) { onDeregistered(&bot, e) })
	bot.api.OnUpdateVotes(func(e ttapi.UpdateVotesEvt) { onUpdateVotes(&bot, e) })
	bot.api.OnSnagged(func(e ttapi.SnaggedEvt) { onSnagged(&bot, e) })

	// DJing
	bot.api.OnRemDJ(func(e ttapi.RemDJEvt) { onRemDj(&bot, e) })
	bot.api.OnAddDJ(func(e ttapi.AddDJEvt) { onAddDj(&bot, e) })
	bot.api.OnNewSong(func(e ttapi.NewSongEvt) { onNewSong(&bot, e) })
	bot.api.Start()
}

// ROOM
func (b *Bot) GetRoomInfo() (ttapi.RoomInfoRes, error) {
	return b.api.RoomInfo()
}

func (b *Bot) AddDjEscorting(userId string) error {
	if b.room.djs.HasElement(userId) {
		return errors.New("You aren't DJing!")
	}

	if !b.escorting.HasElement(userId) {
		b.escorting.Push(userId)
	}

	return nil
}

func (b *Bot) RemoveDjEscorting(userId string) error {
	return b.escorting.Remove(userId)
}

func (b *Bot) EscortDj(userId string) error {
	return b.api.RemDj(userId)
}

// SONG
func (b *Bot) Bop() {
	if b.room.song.djId != b.config.UserId {
		b.api.Bop()
	}
}

func (b *Bot) Downvote() {
	if b.room.song.djId != b.config.UserId {
		b.api.VoteDown()
	}
}
func (b *Bot) SkipSong() {
	b.api.Skip()
}

func (b *Bot) ToggleAutoBop() bool {
	b.config.AutoBop = !b.config.AutoBop
	return b.config.AutoBop
}

func (b *Bot) ShowSongStats() {
	song := b.room.song
	msg := fmt.Sprintf("Stats for `%s` by `%s` played by @%s:", song.title, song.artist, song.djName)
	b.RoomMessage(msg)

	msg = fmt.Sprintf("👍 %d | 👎 %d | ❤️ %d", song.up, song.down, song.snag)

	delay := time.Duration(10) * time.Millisecond
	utils.ExecuteDelayed(delay, func() {
		b.RoomMessage(msg)
	})
}

// AUTO DJ
func (b *Bot) AutoDj() {
	if !b.room.djs.HasElement(b.config.UserId) {
		b.api.AddDj()
	}
}

func (b *Bot) ToggleAutoDj() bool {
	b.config.AutoDj = !b.config.AutoDj
	b.AddDjEscorting(b.config.UserId)
	return b.config.AutoDj
}

// PLAYLISTS
func (b *Bot) Snag(songId string) error {
	if b.room.song.djId == b.config.UserId {
		return errors.New("I'm the current DJ and I already have this song in my playlist...")
	}

	playlist, err := b.api.PlaylistAll(b.config.CurrentPlaylist)

	if err != nil {
	}

	b.api.Snag()
	if err = b.api.PlaylistAdd(songId, b.config.CurrentPlaylist, len(playlist.List)); err != nil {
		return nil
	}

	b.playlist.AddSong(&SongItem{
		id:     b.room.song.id,
		title:  b.room.song.title,
		artist: b.room.song.artist,
		length: b.room.song.length,
	})

	return nil
}

func (b *Bot) ToggleAutoSnag() bool {
	b.config.AutoSnag = !b.config.AutoSnag
	return b.config.AutoSnag
}

func (b *Bot) LoadPlaylist(playlistName string) error {
	playlist, err := b.api.PlaylistAll(b.config.CurrentPlaylist)

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

func (b *Bot) LoadPlaylists() error {
	playlists, err := b.api.PlaylistListAll()
	if err != nil {
		return err
	}

	for _, pl := range playlists.List {
		b.playlists.Push(pl.Name)
	}
	return nil
}

func (b *Bot) AddPlaylist(playlistName string) error {
	if !b.playlists.HasElement(playlistName) {
		if err := b.api.PlaylistCreate(playlistName); err != nil {
			return err
		}
		b.playlists.Push(playlistName)
		return nil
	}

	return errors.New("Playlist not found")
}

func (b *Bot) RemovePlaylist(playlistName string) error {
	if b.playlists.HasElement(playlistName) {
		if err := b.api.PlaylistDelete(playlistName); err != nil {
			return err
		}
		b.playlists.Remove(playlistName)
		return nil
	}

	return errors.New("Playlist not found")
}

func (b *Bot) SwitchPlaylist(playlistName string) error {
	if b.playlists.HasElement(playlistName) {
		if err := b.api.PlaylistSwitch(playlistName); err != nil {
			return err
		}
		b.config.CurrentPlaylist = playlistName
		return b.LoadPlaylist(playlistName)
	}

	return errors.New("Playlist not found")
}

func (b *Bot) PushSongBottomPlaylist() error {
	if err := b.api.PlaylistReorder(b.config.CurrentPlaylist, 0, b.playlist.songs.Size()-1); err == nil {
		currentSong, _ := b.playlist.songs.Shift()
		b.playlist.AddSong(&currentSong)
		return nil
	} else {
		return err
	}
}

func (b *Bot) RemoveSongFromPlaylist(s *SongItem) error {
	idx := b.playlist.songs.IndexOf(*s)

	if idx < 0 {
		return errors.New("Song not found in current playlist")
	}

	if err := b.api.PlaylistRemove(b.config.CurrentPlaylist, idx); err != nil {
		return err
	}

	b.playlist.RemoveSong(s)
	return nil
}

// MESSAGING
func (b *Bot) PrivateMessage(userId, msg string) {
	b.api.PM(userId, msg)
}

func (b *Bot) RoomMessage(msg string) {
	b.api.Speak(msg)
}

// USERS & AUTHORIZATION
func (b *Bot) Fan(userId string) error {
	return b.api.BecomeFan(userId)
}

func (b *Bot) Unfan(userId string) error {
	return b.api.RemoveFan(userId)
}

func (b *Bot) UserFromId(userId string) (*User, error) {
	if userName, ok := b.room.users.Get(userId); ok {
		return &User{Id: userId, Name: userName}, nil
	}
	return &User{}, errors.New("User with ID " + userId + " wasn't found")
}

func (b *Bot) UserFromName(userName string) (*User, error) {
	if id, err := b.api.GetUserID(userName); err == nil {
		return &User{Id: id, Name: userName}, nil
	} else {
		return &User{}, err
	}
}

func (b *Bot) UserIsAdmin(user *User) bool {
	return b.admins.HasElement(user.Name)
}

func (b *Bot) UserIsModerator(user *User) bool {
	return b.room.moderators.HasElement(user.Id)
}
