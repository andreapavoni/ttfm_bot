package ttfm

import (
	"errors"
	"fmt"

	"github.com/alaingilbert/ttapi"
	"github.com/sirupsen/logrus"

	"github.com/andreapavoni/ttfm_bot/utils"
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
	escorting *collections.SmartList[string]
	commands  *collections.SmartMap[CommandHandler]
}

type Config struct {
	ApiAuth                string
	UserId                 string
	RoomId                 string
	Admins                 []string
	AutoSnag               bool
	AutoBop                bool
	AutoDj                 bool
	AutoQueue              bool
	AutoQueueMsg           string
	AutoShowSongStats      bool
	ModAutoWelcome         bool
	ModQueue               bool
	ModQueueInviteDuration int64
	ModSongsMaxDuration    int64
	ModSongsMaxPerDj       int64
	ModDjRestDuration      int64
	CurrentPlaylist        string
	SetBot                 bool
}

// BOOT
func New() *Bot {
	logrus.SetFormatter(&LogFormatter{})
	// f, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0777)
	// logrus.SetOutput(f)

	c := loadConfig()

	b := Bot{
		Config:    c,
		Room:      NewRoom(),
		Queue:     NewQueue(),
		api:       ttapi.NewBot(c.ApiAuth, c.UserId, c.RoomId),
		admins:    collections.NewSmartListFromSlice(c.Admins),
		playlist:  NewPlaylist(c.CurrentPlaylist),
		Playlists: collections.NewSmartList[string](),
		escorting: collections.NewSmartList[string](),
		commands:  collections.NewSmartMap[CommandHandler](),
	}

	// Commands
	b.api.OnSpeak(func(e ttapi.SpeakEvt) { handleCommandSpeak(&b, e.UserID, e.Text) })
	b.api.OnPmmed(func(e ttapi.PmmedEvt) { handleCommandPm(&b, e.SenderID, e.Text) })

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
func loadConfig() *Config {
	return &Config{
		ApiAuth:                utils.GetEnvOrPanic("TTFM_API_AUTH"),
		UserId:                 utils.GetEnvOrPanic("TTFM_API_USER_ID"),
		Admins:                 utils.StringToSlice(utils.GetEnvOrDefault("TTFM_ADMINS", "pavonz"), ","),
		RoomId:                 utils.GetEnvOrPanic("TTFM_API_ROOM_ID"),
		AutoSnag:               utils.GetEnvBoolOrDefault("TTFM_AUTO_SNAG", true),
		AutoBop:                utils.GetEnvBoolOrDefault("TTFM_AUTO_BOP", true),
		AutoDj:                 utils.GetEnvBoolOrDefault("TTFM_AUTO_DJ", false),
		AutoQueue:              utils.GetEnvBoolOrDefault("TTFM_AUTO_QUEUE", false),
		AutoQueueMsg:           utils.GetEnvOrDefault("TTFM_AUTO_QUEUE_MSG", "Next in line is: @JumpingMonkey. Time to claim your spot!"),
		AutoShowSongStats:      utils.GetEnvBoolOrDefault("TTFM_AUTO_SHOW_SONG_STATS", false),
		ModAutoWelcome:         utils.GetEnvBoolOrDefault("TTFM_AUTO_WELCOME", false),
		ModQueue:               utils.GetEnvBoolOrDefault("TTFM_MOD_QUEUE", false),
		ModQueueInviteDuration: utils.GetEnvIntOrDefault("TTFM_MOD_QUEUE_INVITE_DURATION", 1),
		ModSongsMaxDuration:    utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_DURATION", 10),
		ModSongsMaxPerDj:       utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_DURATION", 0),
		ModDjRestDuration:      utils.GetEnvIntOrDefault("TTFM_MOD_DJ_REST_DURATION", 0),
		CurrentPlaylist:        utils.GetEnvOrDefault("TTFM_DEFAULT_PLAYLIST", "default"),
		SetBot:                 utils.GetEnvBoolOrDefault("TTFM_SET_BOT", false),
	}
}
func (b *Bot) AddCommand(trigger string, h CommandHandler) {
	b.commands.Set(trigger, h)
}
func (b *Bot) Start() {
	b.api.Start()
}

// QUEUE
func (b *Bot) ToggleModQueue() bool {
	b.Config.ModQueue = !b.Config.ModQueue
	return b.Config.ModQueue
}

// ROOM
func (b *Bot) GetRoomInfo() (ttapi.RoomInfoRes, error) {
	return b.api.RoomInfo()
}
func (b *Bot) AddDjEscorting(userId string) error {
	if !b.UserIsDj(userId) {
		if userId == b.Config.UserId {
			return errors.New("I'm not on stage!")
		}
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
	if b.Room.Song.djId != b.Config.UserId {
		b.api.Bop()
	}
}
func (b *Bot) Downvote() {
	if b.Room.Song.djId != b.Config.UserId {
		b.api.VoteDown()
	}
}
func (b *Bot) SkipSong() {
	b.api.Skip()
}
func (b *Bot) ToggleAutoBop() bool {
	b.Config.AutoBop = !b.Config.AutoBop
	return b.Config.AutoBop
}
func (b *Bot) ShowSongStats() (header, data string) {
	song := b.Room.Song
	header = fmt.Sprintf("Stats for `%s` by `%s` played by @%s:", song.Title, song.Artist, song.DjName)
	data = fmt.Sprintf("👍 %d | 👎 %d | ❤️ %d", song.up, song.down, song.snag)

	return header, data
}

// AUTO DJ
func (b *Bot) AutoDj() {
	if !b.UserIsDj(b.Config.UserId) {
		b.api.AddDj()
	}
}
func (b *Bot) ToggleAutoDj() bool {
	b.Config.AutoDj = !b.Config.AutoDj
	b.AddDjEscorting(b.Config.UserId)
	return b.Config.AutoDj
}

// PLAYLISTS
func (b *Bot) Snag() error {
	if b.Room.Song.djId == b.Config.UserId {
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
func (b *Bot) ToggleAutoSnag() bool {
	b.Config.AutoSnag = !b.Config.AutoSnag
	return b.Config.AutoSnag
}
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
func (b *Bot) SwitchPlaylist(playlistName string) error {
	if b.Playlists.HasElement(playlistName) {
		if err := b.api.PlaylistSwitch(playlistName); err != nil {
			return err
		}
		b.Config.CurrentPlaylist = playlistName
		return b.LoadPlaylist(playlistName)
	}

	return errors.New("Playlist not found")
}
func (b *Bot) PushSongBottomPlaylist() error {
	if err := b.api.PlaylistReorder(b.Config.CurrentPlaylist, 0, b.playlist.songs.Size()-1); err == nil {
		currentSong, _ := b.playlist.songs.Shift()
		b.playlist.AddSong(&currentSong)
		return nil
	} else {
		return err
	}
}
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
func (b *Bot) PrivateMessage(userId, msg string) {
	b.api.PM(userId, msg)
}
func (b *Bot) RoomMessage(msg string) {
	b.api.Speak(msg)
}

// USERS & AUTHORIZATION
func (b *Bot) BootUser(userId, reason string) error {
	return b.api.BootUser(userId, reason)
}

func (b *Bot) Fan(userId string) error {
	return b.api.BecomeFan(userId)
}
func (b *Bot) Unfan(userId string) error {
	return b.api.RemoveFan(userId)
}
func (b *Bot) UserFromId(userId string) (*User, error) {
	if userName, ok := b.Room.users.Get(userId); ok {
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
func (b *Bot) UserIsDj(userId string) bool {
	return b.Room.djs.HasElement(userId)
}
func (b *Bot) UserIsCurrentDj(userId string) bool {
	return b.Room.Song.djId == userId
}
func (b *Bot) UserIsModerator(userId string) bool {
	return b.Room.moderators.HasElement(userId)
}
