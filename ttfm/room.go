package ttfm

import (
	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Room struct {
	name       string
	id         string
	shortcut   string
	users      *collections.SmartMap[string]
	admins     *collections.SmartList[string]
	moderators *collections.SmartList[string]
	djs        *collections.SmartList[string]
	Song       *Song
}

type User struct {
	Id   string
	Name string
}

func (r *Room) Update(ri ttapi.RoomInfoRes) error {
	r.name = ri.Room.Name
	r.id = ri.Room.Roomid
	r.shortcut = ri.Room.Shortcut

	song := ri.Room.Metadata.CurrentSong
	r.Song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)
	r.UpdateModerators(ri.Room.Metadata.ModeratorID)

	users := []User{}
	for _, u := range ri.Users {
		users = append(users, User{Id: u.ID, Name: u.Name})
	}
	r.UpdateUsers(users)

	r.UpdateDjs(ri.Room.Metadata.Djs)

	return nil
}

func (r *Room) UpdateUsers(users []User) {
	r.users = collections.NewSmartMap[string]()
	for _, u := range users {
		r.AddUser(u.Id, u.Name)
	}
}

func (r *Room) AddUser(id, name string) {
	r.users.Set(id, name)
}

func (r *Room) RemoveUser(id string) {
	r.users.Delete(id)
}

func (r *Room) AddDj(id string) {
	r.djs.Push(id)
}

func (r *Room) RemoveDj(id string) {
	r.djs.Remove(id)
}

func (r *Room) UpdateDjs(djs []string) {
	r.djs = collections.NewSmartListFromSlice(djs)
}

func (r *Room) UpdateModerators(moderators []string) {
	r.moderators = collections.NewSmartListFromSlice(moderators)
}
