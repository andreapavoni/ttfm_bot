package ttfm

import (
	"github.com/alaingilbert/ttapi"
)

type Bot struct {
	Brain  *Brain
	Config *Config
	api    *ttapi.Bot

	Actions         *Actions
	Admins          *Admins
	CurrentPlaylist *Playlist
	Identity        *User

	Users     *Users
	Playlists *Playlists
	Room      *Room
	Queue     *Queue
	Commands  *Commands
	Reactions *Reactions
	FavRooms  *FavRooms
}

// BOOT
// New bot instance
func New() *Bot {
	SetupLogging()
	brain := NewBrain("./db")
	cfg := NewConfig(brain)

	b := Bot{
		Config:    cfg,
		Queue:     NewQueue(),
		api:       ttapi.NewBot(cfg.ApiAuth, cfg.UserId, cfg.RoomId),
		Admins:    NewAdmins(brain),
		Commands:  NewCommands(),
		Reactions: NewReactions(brain),
		Brain:     brain,
		Identity:  &User{Id: cfg.UserId},
	}

	b.CurrentPlaylist = NewPlaylist(&b, cfg.CurrentPlaylist)
	b.Playlists = NewPlaylists(&b)
	b.Users = NewUsers(&b)
	b.Room = NewRoom(&b)
	b.Actions = NewActions(&b)
	b.FavRooms = NewFavRooms(&b)

	// Commands
	b.api.OnSpeak(func(e ttapi.SpeakEvt) {
		mi := &MessageInput{UserId: e.UserID, Text: e.Text, Source: MessageTypeRoom}
		mi.HandleCommand(&b)
	})
	b.api.OnPmmed(func(e ttapi.PmmedEvt) {
		mi := &MessageInput{UserId: e.SenderID, Text: e.Text, Source: MessageTypePm}
		mi.HandleCommand(&b)
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

// Start the bot
func (b *Bot) Start() {
	b.api.Start()
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
