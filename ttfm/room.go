package ttfm

import "github.com/alaingilbert/ttapi"

type Room struct {
	name       string
	id         string
	shortcut   string
	users      *SmartMap
	admins     *SmartList
	moderators *SmartList
	djs        *SmartList
	song       *Song
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
	r.song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)
	r.UpdateModerators(ri.Room.Metadata.ModeratorID)

	users := []User{}
	for _, u := range ri.Users {
		users = append(users, User{Id: u.ID, Name: u.Name})
	}
	r.UpdateUsers(users)

	return nil
}

func (r *Room) UpdateUsers(users []User) {
	r.users = NewSmartMap()
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

func (r *Room) UpdateModerators(moderators []string) {
	r.moderators = NewSmartListFromSlice(moderators)
}

func (r *Room) UserIsModerator(userId string) bool {
	return r.moderators.HasElement(userId)
}

func (r *Room) UserNameFromId(id string) string {
	if userName, ok := r.users.Get(id); ok {
		return userName.(string)
	}
	return ""
}
