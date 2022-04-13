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

func (d *Dj) UpdateStats(up, down, snag, play int) {
	d.up += up
	d.down += down
	d.snag += snag
	d.plays += play
}

type Room struct {
	Name        string
	Id          string
	Shortcut    string
	moderators  *collections.SmartList[string]
	Djs         *collections.SmartMap[*Dj]
	MaxDjs      int
	CurrentSong *Song
	CurrentDj   *Dj
	escorting   *collections.SmartList[string]
	bot         *Bot
}

func NewRoom(bot *Bot) *Room {
	return &Room{
		bot:         bot,
		moderators:  collections.NewSmartList[string](),
		Djs:         collections.NewSmartMap[*Dj](),
		CurrentSong: NewSong(bot),
		CurrentDj:   &Dj{"", 0, 0, 0, 0},
		escorting:   collections.NewSmartList[string](),
	}
}

func (r *Room) Update(ri ttapi.RoomInfoRes) error {
	r.Name = ri.Room.Name
	r.Id = ri.Room.Roomid
	r.Shortcut = ri.Room.Shortcut
	r.MaxDjs = ri.Room.Metadata.MaxDjs

	song := ri.Room.Metadata.CurrentSong
	r.CurrentSong.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)

	users := []User{}
	for _, u := range ri.Users {
		users = append(users, User{Id: u.ID, Name: u.Name})
	}

	r.bot.Users.Update(users)
	r.UpdateModerators(ri.Room.Metadata.ModeratorID)
	r.UpdateDjs(ri.Room.Metadata.Djs)
	if err := r.ResetDj(); err != nil {
		return err
	}

	return nil
}

func (r *Room) ResetDj() error {
	d, ok := r.Djs.Get(r.CurrentSong.DjId)
	if !ok {
		return errors.New("Can't reset current Dj")
	}
	r.CurrentDj = d
	return nil
}

func (r *Room) AddDj(id string) {
	r.Djs.Set(id, NewDj(id))
}

func (r *Room) RemoveDj(id string) {
	r.Djs.Delete(id)
}

func (r *Room) UpdateDjs(djs []string) {
	newDjs := collections.NewSmartMap[*Dj]()
	for _, djId := range djs {
		if dj, ok := r.Djs.Get(djId); ok {
			newDjs.Set(djId, dj)
		} else {
			newDjs.Set(djId, NewDj(djId))
		}
	}
	r.Djs = newDjs
}

func (r *Room) UpdateModerators(moderators []string) {
	r.moderators = collections.NewSmartListFromSlice(moderators)
}

func (r *Room) HasModerator(userId string) bool {
	return r.moderators.HasElement(userId)
}

func (r *Room) UpdateDataFromApi() error {
	if roomInfo, err := r.bot.api.RoomInfo(); err != nil {
		return err
	} else {
		if err := r.Update(roomInfo); err != nil {
			return err
		}
	}
	return nil
}

// SongStats
func (r *Room) SongStats() (data string) {
	song := r.CurrentSong
	data = fmt.Sprintf("Stats for `%s` by `%s` played by @%s: 👍 %d | 👎 %d | ❤️ %d", song.Title, song.Artist, song.DjName, song.up, song.down, song.snag)
	return data
}

// DjStats
func (r *Room) DjStats(userId string) (data string, err error) {
	dj, ok := r.Djs.Get(userId)
	user, err := r.bot.Users.UserFromId(userId)

	if !ok || err != nil {
		return "", errors.New("Dj not found")
	}

	data = fmt.Sprintf("Stats for @%s: 👍 %d | 👎 %d | ❤️ %d | 🎧 %d", user.Name, dj.up, dj.down, dj.snag, dj.plays)
	return data, nil
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
