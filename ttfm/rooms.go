package ttfm

import (
	"errors"
	"fmt"

	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Rooms struct {
	bot *Bot
	*collections.SmartMap[string]
}

func NewRooms(b *Bot) *Rooms {
	r := Rooms{bot: b, SmartMap: collections.NewSmartMap[string]()}
	r.LoadRoomsFromDb()
	return &r
}

// AddFavoriteRoom
func (r *Rooms) AddFavorite(roomSlug, roomId string) error {
	r.Set(roomSlug, roomId)
	return r.storeRoomsOnDb()
}

// ListFavoriteRooms
func (r *Rooms) ListFavorites() (rooms []string) {
	for _, k := range r.Keys() {
		rooms = append(rooms, k)
	}
	return rooms
}

// JoinRoom
func (r *Rooms) Join(roomSlug string) error {
	roomId, ok := r.Get(roomSlug)

	if !ok {
		return fmt.Errorf("room `%s` hasn't been found in my brain", roomSlug)
	}

	if roomId == r.bot.Room.Id {
		return fmt.Errorf("I'm already in `%s`", roomSlug)
	}

	if err := r.bot.api.RoomRegister(roomId); err != nil {
		return err
	}

	return nil
}

func (r *Rooms) LoadRoomsFromDb() error {
	rooms := map[string]string{}
	if err := r.bot.Brain.Get("rooms", &rooms); err != nil {
		if r.bot.Room.Shortcut != "" {
			// set first one if none is found
			r.Set(r.bot.Room.Shortcut, r.bot.Room.Id)
			rooms[r.bot.Room.Shortcut] = r.bot.Room.Id
			r.bot.Brain.Put("rooms", &rooms)
			return nil
		}
		return errors.New("can't load rooms")
	}

	for k, v := range rooms {
		r.Set(k, v)
	}

	return nil
}

func (r *Rooms) storeRoomsOnDb() error {
	rooms := map[string]string{}
	for i := range r.SmartMap.Iter() {
		rooms[i.Key] = i.Value
	}
	return r.bot.Brain.Put("rooms", &rooms)
}
