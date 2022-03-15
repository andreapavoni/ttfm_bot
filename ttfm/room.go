package ttfm

import (
	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Room struct {
	Name       string
	Id         string
	Shortcut   string
	Users      *collections.SmartMap[string]
	Moderators *collections.SmartList[string]
	Djs        *collections.SmartList[string]
	MaxDjs     int
	Song       *Song
	escorting  *collections.SmartList[string]
}

type User struct {
	Id   string
	Name string
}

func NewRoom() *Room {
	return &Room{
		Users:      collections.NewSmartMap[string](),
		Moderators: collections.NewSmartList[string](),
		Djs:        collections.NewSmartList[string](),
		Song:       &Song{},
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

	r.UpdateUsers(users)
	r.UpdateModerators(ri.Room.Metadata.ModeratorID)
	r.UpdateDjs(ri.Room.Metadata.Djs)

	return nil
}

func (r *Room) UpdateUsers(users []User) {
	r.Users = collections.NewSmartMap[string]()
	for _, u := range users {
		r.AddUser(u.Id, u.Name)
	}
}

func (r *Room) AddUser(id, name string) {
	r.Users.Set(id, name)
}

func (r *Room) RemoveUser(id string) {
	r.Users.Delete(id)
}

func (r *Room) AddDj(id string) {
	r.Djs.Push(id)
}

func (r *Room) RemoveDj(id string) {
	r.Djs.Remove(id)
}

func (r *Room) UpdateDjs(djs []string) {
	r.Djs = collections.NewSmartListFromSlice(djs)
}

func (r *Room) UpdateModerators(moderators []string) {
	r.Moderators = collections.NewSmartListFromSlice(moderators)
}
