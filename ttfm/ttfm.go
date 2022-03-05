package ttfm

import (
	"errors"
	"fmt"
	"time"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils"
)

type Bot struct {
	api       *ttapi.Bot
	queue     *SmartList
	admins    *SmartList
	config    *Config
	room      *Room
	escorting *SmartList
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
	DefaultPlaylist     string
	SetBot              bool
}

// BOOT
func Start(c Config) {
	bot := Bot{
		api:    ttapi.NewBot(c.ApiAuth, c.UserId, c.RoomId),
		queue:  NewSmartList(),
		admins: NewSmartListFromSlice(c.Admins),
		config: &c,
		room: &Room{
			users:      NewSmartMap(),
			moderators: NewSmartList(),
			djs:        NewSmartList(),
			song:       &Song{},
		},
		escorting: NewSmartList(),
	}

	// Commands
	bot.api.OnSpeak(func(e ttapi.SpeakEvt) { handleCommandSpeak(bot, e) })
	bot.api.OnPmmed(func(e ttapi.PmmedEvt) { handleCommandPm(bot, e) })

	// Room events
	bot.api.OnReady(func() { onReady(bot) })
	bot.api.OnRoomChanged(func(e ttapi.RoomInfoRes) { onRoomChanged(bot, e) })
	bot.api.OnRegistered(func(e ttapi.RegisteredEvt) { onRegistered(bot, e) })
	bot.api.OnDeregistered(func(e ttapi.DeregisteredEvt) { onDeregistered(bot, e) })
	bot.api.OnUpdateVotes(func(e ttapi.UpdateVotesEvt) { onUpdateVotes(bot, e) })
	bot.api.OnSnagged(func(e ttapi.SnaggedEvt) { onSnagged(bot, e) })

	// DJing
	bot.api.OnRemDJ(func(e ttapi.RemDJEvt) { onRemDj(bot, e) })
	bot.api.OnAddDJ(func(e ttapi.AddDJEvt) { onAddDj(bot, e) })
	bot.api.OnNewSong(func(e ttapi.NewSongEvt) { onNewSong(bot, e) })
	bot.api.Start()
}

// ROOM STUFF

func (b *Bot) GetRoomInfo() (ttapi.RoomInfoRes, error) {
	return b.api.RoomInfo()
}

// STAGE / QUEUE MANAGEMENT
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

// func (b *Bot) QueueList() []string {
// 	var list = []string{}
// 	for _, u := range ListElements(b.queue, list) {
// 		userName := b.room.UserNameFromId(u)
// 		list = append(list, userName)
// 	}

// 	return list
// }

// func (b *Bot) QueueAdd(userId string) error {
// 	if b.queue.HasElement(userId) {
// 		b.queue.Push(userId)
// 	} else {
// 		return errors.New("DJ already in queue")
// 	}
// 	return nil
// }

// func (b *Bot) QueueRemove(userId string) error {
// 	return b.queue.Remove(userId)
// }

// AUTO DJ
func (b *Bot) AutoDj() {
	if !b.room.djs.HasElement(b.config.UserId) {
		b.api.AddDj()
	}
}

func (b *Bot) ToggleAutoDj() bool {
	b.config.AutoDj = !b.config.AutoDj
	return b.config.AutoDj
}

// SONGS
func (b *Bot) Snag(songId string) error {
	if b.room.song.djId != b.config.UserId {
		if playlist, err := b.api.PlaylistAll(b.config.DefaultPlaylist); err == nil {
			b.api.Snag()
			return b.api.PlaylistAdd(songId, b.config.DefaultPlaylist, len(playlist.List))
		} else {
			return err
		}
	} else {
		return errors.New("I'm the current DJ and I already have this song in my playlist...")
	}
}

func (b *Bot) PushSongBottom() error {
	if b.room.song.djId == b.config.UserId {
		if playlist, err := b.api.PlaylistAll(b.config.DefaultPlaylist); err == nil {
			return b.api.PlaylistReorder(b.config.DefaultPlaylist, 0, len(playlist.List)-1)
		} else {
			return err
		}
	} else {
		return errors.New("I'm not the current DJ")
	}
}

func (b *Bot) ToggleAutoSnag() bool {
	b.config.AutoSnag = !b.config.AutoSnag
	return b.config.AutoSnag
}

func (b *Bot) Bop() {
	if b.room.song.djId != b.config.UserId {
		b.api.Bop()
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

// MESSAGING
func (b *Bot) PrivateMessage(userId, msg string) {
	b.api.PM(userId, msg)
}

func (b *Bot) RoomMessage(msg string) {
	b.api.Speak(msg)
}

// AUTHORIZATION
func (b *Bot) UserIsAdmin(userId string) bool {
	if profile, err := b.api.GetProfile(userId); err == nil {
		return b.admins.HasElement(profile.Name)
	}
	return false
}
