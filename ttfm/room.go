package ttfm

import (
	"errors"
	"fmt"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Dj struct {
	userId string
	up     int
	down   int
	snag   int
	plays  int
}

func NewDj(userId string) *Dj {
	return &Dj{userId, 0, 0, 0, 0}
}

type Room struct {
	Name       string
	Id         string
	Shortcut   string
	moderators *collections.SmartList[string]
	Djs        *collections.SmartMap[*Dj]
	MaxDjs     int
	Song       *Song
	escorting  *collections.SmartList[string]
	bot        *Bot
}

func NewRoom(bot *Bot) *Room {
	return &Room{
		bot:        bot,
		moderators: collections.NewSmartList[string](),
		Djs:        collections.NewSmartMap[*Dj](),
		Song:       NewSong(bot),
		escorting:  collections.NewSmartList[string](),
	}
}

func (r *Room) Update(ri ttapi.RoomInfoRes) error {
	r.Name = ri.Room.Name
	r.Id = ri.Room.Roomid
	r.Shortcut = ri.Room.Shortcut
	r.MaxDjs = ri.Room.Metadata.MaxDjs

	song := ri.Room.Metadata.CurrentSong
	r.Song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)

	users := []User{}
	for _, u := range ri.Users {
		users = append(users, User{Id: u.ID, Name: u.Name})
	}

	r.bot.Users.Update(users)
	r.UpdateModerators(ri.Room.Metadata.ModeratorID)
	r.UpdateDjs(ri.Room.Metadata.Djs)

	return nil
}

func (r *Room) AddDj(id string) {
	r.Djs.Set(id, NewDj(id))
}

func (r *Room) RemoveDj(id string) {
	r.Djs.Delete(id)
}

func (r *Room) UpdateDjs(djs []string) {
	r.Djs = collections.NewSmartMap[*Dj]()
	for _, djId := range djs {
		r.AddDj(djId)
	}
}

func (r *Room) UpdateDjStats(userId string, up, down, snag int) error {
	d, ok := r.Djs.Get(userId)

	if !ok {
		return errors.New("Dj not found")
	}

	d.up += up
	d.down += down
	d.snag += snag
	d.plays += 1
	return nil
}

func (r *Room) UpdateModerators(moderators []string) {
	r.moderators = collections.NewSmartListFromSlice(moderators)
}

func (r *Room) HasModerator(userId string) bool {
	return r.moderators.HasElement(userId)
}

func (r *Room) UpdateDataFromApi() {
	if roomInfo, err := r.bot.api.RoomInfo(); err == nil {
		r.Update(roomInfo)
	} else {
		panic(err)
	}
}

// SongStats
func (r *Room) SongStats() (header, data string) {
	song := r.Song
	header = fmt.Sprintf("Stats for `%s` by `%s` played by @%s:", song.Title, song.Artist, song.DjName)
	data = fmt.Sprintf("👍 %d | 👎 %d | ❤️ %d", song.up, song.down, song.snag)
	return header, data
}

// DjStats
func (r *Room) DjStats(userId string) (header, data string, err error) {
	dj, ok := r.Djs.Get(userId)

	if !ok {
		return "", "", errors.New("Dj not found")
	}

	header = fmt.Sprintf("Stats for @%s", r.Song.DjName)
	data = fmt.Sprintf("👍 %d | 👎 %d | ❤️ %d | 🎧 %d", dj.up, dj.down, dj.snag, dj.plays)
	return header, data, nil
}

// AddDjEscorting the dj will be escorted after the current song is played
func (r *Room) AddDjEscorting(userId string) error {
	if !r.bot.Users.UserIsDj(userId) {
		if userId == r.bot.Identity.Id {
			return errors.New("I'm not on stage!")
		}
		return errors.New("You aren't DJing!")
	}
	if !r.escorting.HasElement(userId) {
		r.escorting.Push(userId)
	}
	return nil
}

// RemoveDjEscorting if dj doesn't want to be escorted anymore
func (r *Room) RemoveDjEscorting(userId string) error {
	return r.escorting.Remove(userId)
}

// BootUser from room
func (r *Room) BootUser(userId, reason string) error {
	return r.bot.api.BootUser(userId, reason)
}

// EscortDj
func (r *Room) EscortDj(userId string) error {
	if !r.bot.Users.UserIsDj(userId) {
		if userId == r.bot.Identity.Id {
			return errors.New("I'm not on stage!")
		}
		return errors.New("user is not on stage")
	}
	return r.bot.api.RemDj(userId)
}
